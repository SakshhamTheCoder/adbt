package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/SakshhamTheCoder/adbt/internal/adb"
	"github.com/SakshhamTheCoder/adbt/internal/state"
	"github.com/SakshhamTheCoder/adbt/internal/ui/components"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	FilterAll AppFilter = iota
	FilterUser
	FilterSystem
)

type AppFilter int

var filterNames = []string{"All", "User", "System"}

type AppManager struct {
	state   *state.AppState
	apps    []adb.App
	loading bool
	cursor  int

	search components.SearchState

	filterType AppFilter

	viewport viewport.Model

	confirm components.ConfirmPrompt
	toast   components.Toast
	pending string

	installForm components.FormModal
}

func NewAppManager(state *state.AppState) *AppManager {
	return &AppManager{
		state:    state,
		viewport: viewport.New(0, 0),
	}
}

func (a *AppManager) Init() tea.Cmd {
	if !a.state.HasDevice() {
		return nil
	}

	a.loading = true
	return adb.ListAppsCmd(a.state.DeviceSerial())
}

func (a *AppManager) filteredApps() []adb.App {
	var filtered []adb.App
	lowerSearch := strings.ToLower(a.search.Query)

	for _, app := range a.apps {
		switch a.filterType {
		case FilterUser:
			if app.IsSystem {
				continue
			}
		case FilterSystem:
			if !app.IsSystem {
				continue
			}
		}

		if a.search.Query != "" {
			if !strings.Contains(strings.ToLower(app.PackageName), lowerSearch) {
				continue
			}
		}

		filtered = append(filtered, app)
	}

	return filtered
}

func (a *AppManager) selectedApp() *adb.App {
	filtered := a.filteredApps()
	if len(filtered) == 0 || a.cursor >= len(filtered) {
		return nil
	}
	return &filtered[a.cursor]
}

func (a *AppManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	a.toast.Update(msg)

	if a.installForm.Visible {
		switch msg := msg.(type) {
		case components.FormSubmitMsg:
			values := msg.Values
			a.installForm.Hide()
			if len(values) > 0 && values[0] != "" {
				var toastCmd tea.Cmd
				a.toast, toastCmd = components.ShowToast(
					"Installing APK...",
					false,
					2*time.Second,
				)
				return a, tea.Batch(
					toastCmd,
					adb.InstallAppCmd(a.state.DeviceSerial(), values[0]),
				)
			}
			return a, nil
		case components.FormCancelMsg:
			a.installForm.Hide()
			return a, nil
		}
		return a, a.installForm.Update(msg)
	}

	if a.confirm.Visible {
		switch msg.(type) {

		case components.ConfirmYesMsg:
			a.confirm.Hide()
			app := a.selectedApp()
			if app == nil {
				return a, nil
			}
			serial := a.state.DeviceSerial()
			switch a.pending {
			case "uninstall":
				return a, adb.UninstallAppCmd(serial, app.PackageName)
			case "clear data":
				return a, adb.ClearAppDataCmd(serial, app.PackageName)
			case "force_stop":
				return a, adb.ForceStopAppCmd(serial, app.PackageName)
			}
			return a, nil

		case components.ConfirmNoMsg:
			a.confirm.Hide()
			a.pending = ""
			return a, tea.Batch()
		}

		return a, a.confirm.Update(msg)
	}

	switch msg := msg.(type) {

	case adb.AppsLoadedMsg:
		a.loading = false
		a.apps = msg.Apps
		a.cursor = 0
		a.gotoTop()
		return a, nil

	case adb.AppsLoadErrorMsg:
		a.loading = false
		var cmd tea.Cmd
		a.toast, cmd = components.ShowToast(
			"Failed to load apps",
			true,
			3*time.Second,
		)
		return a, cmd

	case adb.AppActionResultMsg:
		var cmd tea.Cmd
		a.toast, cmd = components.ShowToast(
			msg.Action+" successful",
			false,
			2*time.Second,
		)
		if msg.Action == "uninstall" || msg.Action == "install" {
			return a, tea.Batch(
				cmd,
				adb.ListAppsCmd(a.state.DeviceSerial()),
			)
		}
		return a, cmd

	case adb.AppActionErrorMsg:
		var cmd tea.Cmd
		a.toast, cmd = components.ShowToast(
			msg.Action+" failed: "+msg.Error.Error(),
			true,
			3*time.Second,
		)
		return a, cmd

	case tea.KeyMsg:
		if a.search.Active {
			before := a.search.Query
			a.search.HandleKey(msg)
			if a.search.Query != before {
				a.cursor = 0
				a.gotoTop()
			}
			return a, consumeKeyCmd()
		}

		filtered := a.filteredApps()

		switch msg.String() {
		case "up", "k":
			if a.cursor > 0 {
				a.cursor--
				a.ensureCursorVisible()
			}

		case "down", "j":
			if a.cursor < len(filtered)-1 {
				a.cursor++
				a.ensureCursorVisible()
			}

		case "enter", "l":
			if app := a.selectedApp(); app != nil {
				return a, adb.LaunchAppCmd(
					a.state.DeviceSerial(),
					app.PackageName,
				)
			}

		case "s":
			if app := a.selectedApp(); app != nil {
				a.pending = "force_stop"
				a.confirm.Show("Force stop:\n" + app.PackageName)
			}

		case "u":
			if app := a.selectedApp(); app != nil {
				if app.IsSystem {
					var cmd tea.Cmd
					a.toast, cmd = components.ShowToast(
						"Cannot uninstall system app",
						true,
						2*time.Second,
					)
					return a, cmd
				}
				a.pending = "uninstall"
				a.confirm.Show("Uninstall:\n" + app.PackageName)
			}

		case "x":
			if app := a.selectedApp(); app != nil {
				a.pending = "clear data"
				a.confirm.Show("Clear data:\n" + app.PackageName)
			}

		case "/":
			a.search.Start()

		case "r":
			if a.state.HasDevice() {
				a.loading = true
				a.cursor = 0
				a.gotoTop()
				return a, adb.ListAppsCmd(a.state.DeviceSerial())
			}

		case "i":
			a.installForm.Show("Install APK", []components.FormField{
				{Label: "APK Path", Value: ""},
			})

		case "right":
			a.filterType = (a.filterType + 1) % 3
			a.cursor = 0
			a.gotoTop()

		case "left":
			a.filterType = (a.filterType + 2) % 3
			a.cursor = 0
			a.gotoTop()

		case "esc":
			if a.search.Query != "" {
				a.search.Clear()
				a.cursor = 0
				a.gotoTop()
				return a, consumeKeyCmd()
			}

		default:
			return a, a.updateViewport(msg)
		}
	}

	return a, nil
}

func (a *AppManager) View() string {
	if !a.state.HasDevice() {
		return components.RenderNoDevice(a.state, "Apps")
	}

	var staticContent strings.Builder
	staticContent.WriteString(components.TitleStyle.Render("Installed Applications") + "\n")

	staticContent.WriteString("  ")
	for i, name := range filterNames {
		if AppFilter(i) == a.filterType {
			staticContent.WriteString(components.HelpKeyStyle.Render(name))
		} else {
			staticContent.WriteString(components.StatusMuted.Render(name))
		}
		if i < len(filterNames)-1 {
			staticContent.WriteString(components.StatusMuted.Render(" / "))
		}
	}
	staticContent.WriteString("\n")

	if a.search.Active {
		staticContent.WriteString(
			components.HelpKeyStyle.Render("search: ") + a.search.Query + "▌\n",
		)
	} else if a.search.Query != "" {
		staticContent.WriteString(
			components.StatusMuted.Render("filter: \""+a.search.Query+"\"") + "\n",
		)
	}

	maxWidth := a.state.Width - 8
	if maxWidth < 20 {
		maxWidth = 20
	}
	truncStyle := lipgloss.NewStyle().MaxWidth(maxWidth)

	var scrollableContent strings.Builder

	if a.loading {
		scrollableContent.WriteString(components.StatusMuted.Render("Loading apps..."))
	} else {
		filtered := a.filteredApps()

		if len(filtered) == 0 {
			scrollableContent.WriteString(components.StatusMuted.Render("No apps found"))
		} else {
			for i, app := range filtered {
				prefix := "  "
				if i == a.cursor {
					prefix = "› "
				}

				tag := components.StatusMuted.Render("[U]")
				if app.IsSystem {
					tag = components.StatusMuted.Render("[S]")
				}

				var line string
				if i == a.cursor {
					line = fmt.Sprintf(
						"%s%s %s",
						prefix,
						tag,
						components.ListItemSelectedStyle.Render(app.PackageName),
					)
				} else {
					line = fmt.Sprintf(
						"%s%s %s",
						prefix,
						tag,
						components.ListItemStyle.Render(app.PackageName),
					)
				}

				scrollableContent.WriteString(truncStyle.Render(line) + "\n")
			}
		}
	}

	var footer string
	if a.search.Active {
		footer = components.Help("enter", "apply") + "  " +
			components.Help("esc", "cancel")
	} else if a.search.Query != "" {
		footer = components.Help("↑/↓", "navigate") + "  " +
			components.Help("enter", "launch") + "  " +
			components.Help("/", "search") + "  " +
			components.Help("esc", "clear filter")
	} else {
		footer = components.Help("↑/↓", "navigate") + "  " +
			components.Help("enter", "launch") + "  " +
			components.Help("i", "install") + "  " +
			components.Help("s", "stop") + "  " +
			components.Help("u", "uninstall") + "  " +
			components.Help("x", "clear") + "  " +
			components.Help("←/→", "filter") + "  " +
			components.Help("/", "search") + "  " +
			components.Help("r", "reload") + "  " +
			components.Help("esc", "back")
	}

	rendered := components.RenderLayoutWithScrollableSection(a.state, components.LayoutWithScrollProps{
		Title:             "Apps",
		StaticContent:     staticContent.String(),
		ScrollableContent: scrollableContent.String(),
		Footer:            footer,
		Viewport:          &a.viewport,
	})

	if a.installForm.Visible {
		rendered = components.RenderFormOverlay(rendered, a.installForm, a.state)
	}

	if a.confirm.Visible {
		rendered = components.RenderOverlay(rendered, a.confirm.View(), a.state)
	}

	if a.toast.Visible {
		rendered = components.RenderOverlay(rendered, a.toast.View(), a.state)
	}

	return rendered
}

func (a *AppManager) updateViewport(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	a.viewport, cmd = a.viewport.Update(msg)
	return cmd
}

func (a *AppManager) gotoTop() {
	a.viewport.GotoTop()
}

func (a *AppManager) ensureCursorVisible() {
	ensureViewportLineVisible(&a.viewport, a.cursor)
}

func consumeKeyCmd() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}
