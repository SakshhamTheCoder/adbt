package screens

import (
	"strings"
	"time"

	"github.com/SakshhamTheCoder/adbt/internal/adb"
	"github.com/SakshhamTheCoder/adbt/internal/state"
	"github.com/SakshhamTheCoder/adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type intentMode int

const (
	intentModeActivity intentMode = iota
	intentModeBroadcast
)

var intentModeNames = []string{"Activity", "Broadcast"}

var intentActionSuggestions = []string{
	"android.intent.action.VIEW",
	"android.intent.action.SEND",
	"android.intent.action.SENDTO",
	"android.intent.action.DIAL",
	"android.intent.action.MAIN",
	"android.intent.action.PICK",
	"android.intent.action.GET_CONTENT",
	"android.settings.SETTINGS",
	"android.settings.WIFI_SETTINGS",
	"android.settings.BLUETOOTH_SETTINGS",
	"android.settings.APPLICATION_SETTINGS",
	"android.settings.APPLICATION_DETAILS_SETTINGS",
}

type Intents struct {
	state *state.AppState

	form components.FormModal
	mode intentMode

	lastOutput string
	toast      components.Toast
}

func NewIntents(state *state.AppState) *Intents {
	i := &Intents{
		state: state,
	}
	i.showForm()
	return i
}

func (i *Intents) showForm() {
	title := "Start Activity"
	if i.mode == intentModeBroadcast {
		title = "Send Broadcast"
	}
	i.form.Show(title, []components.FormField{
		{
			Label:       "Action",
			Value:       "android.intent.action.VIEW",
			Type:        components.FormFieldAutocomplete,
			Suggestions: intentActionSuggestions,
		},
		{
			Label:       "Data URI",
			Placeholder: "https://...  geo:...  tel:...  mailto:...",
		},
		{Label: "Extras (k=v;k2=v2)"},
	})
}

func (i *Intents) Init() tea.Cmd {
	return nil
}

func (i *Intents) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	i.toast.Update(msg)

	if i.form.Visible {
		switch msg := msg.(type) {
		case components.FormSubmitMsg:
			values := msg.Values
			i.form.Hide()

			action := ""
			dataURI := ""
			extras := ""
			if len(values) > 0 {
				action = values[0]
			}
			if len(values) > 1 {
				dataURI = values[1]
			}
			if len(values) > 2 {
				extras = values[2]
			}

			if action == "" {
				var cmd tea.Cmd
				i.toast, cmd = components.ShowToast(
					"Action is required",
					true,
					2*time.Second,
				)
				return i, cmd
			}

			serial := i.state.DeviceSerial()
			var toastCmd tea.Cmd
			i.toast, toastCmd = components.ShowToast(
				"Sending intent...",
				false,
				2*time.Second,
			)

			var intentCmd tea.Cmd
			if i.mode == intentModeBroadcast {
				intentCmd = adb.SendBroadcastCmd(serial, action, extras)
			} else {
				intentCmd = adb.SendIntentCmd(serial, action, dataURI, extras)
			}

			return i, tea.Batch(toastCmd, intentCmd)

		case components.FormCancelMsg:
			i.form.Hide()
			return i, nil
		}

		return i, i.form.Update(msg)
	}

	switch msg := msg.(type) {
	case adb.IntentResultMsg:
		i.lastOutput = msg.Output
		var cmd tea.Cmd
		i.toast, cmd = components.ShowToast(
			"Intent sent successfully",
			false,
			2*time.Second,
		)
		return i, cmd

	case adb.IntentErrorMsg:
		i.lastOutput = msg.Error.Error()
		var cmd tea.Cmd
		i.toast, cmd = components.ShowToast(
			"Intent failed",
			true,
			3*time.Second,
		)
		return i, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			i.showForm()
		case "left", "right":
			i.mode = (i.mode + 1) % 2
			i.showForm()
		}
	}

	return i, nil
}

func (i *Intents) View() string {
	if !i.state.HasDevice() {
		return components.RenderNoDevice(i.state, "Intents")
	}

	var staticContent strings.Builder
	staticContent.WriteString(components.TitleStyle.Render("Intent Tester") + "\n")

	staticContent.WriteString("  ")
	for idx, name := range intentModeNames {
		if intentMode(idx) == i.mode {
			staticContent.WriteString(components.HelpKeyStyle.Render(name))
		} else {
			staticContent.WriteString(components.StatusMuted.Render(name))
		}
		if idx < len(intentModeNames)-1 {
			staticContent.WriteString(components.StatusMuted.Render(" / "))
		}
	}
	staticContent.WriteString("\n\n")

	var scrollableContent strings.Builder
	if i.lastOutput != "" {
		scrollableContent.WriteString(components.StatusMuted.Render("Last Result:") + "\n")
		scrollableContent.WriteString(i.lastOutput + "\n")
	} else {
		scrollableContent.WriteString(components.StatusMuted.Render("Press [n] to send an intent") + "\n")
	}

	footer := components.Help("n", "new intent") + "  " +
		components.Help("←/→", "mode") + "  " +
		components.Help("esc", "back")

	rendered := components.RenderLayoutWithScrollableSection(i.state, components.LayoutWithScrollProps{
		Title:             "Intents",
		StaticContent:     staticContent.String(),
		ScrollableContent: scrollableContent.String(),
		Footer:            footer,
	})

	if i.form.Visible {
		rendered = components.RenderFormOverlay(rendered, i.form, i.state)
	}

	if i.toast.Visible {
		rendered = components.RenderOverlay(rendered, i.toast.View(), i.state)
	}

	return rendered
}
