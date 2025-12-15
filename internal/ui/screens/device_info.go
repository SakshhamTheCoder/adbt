package screens

import (
	"strings"
	"time"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type deviceAction struct {
	key     string
	label   string
	run     func(string) error
	confirm bool
	prompt  string
}

type DeviceInfo struct {
	state *state.AppState

	actions       []deviceAction
	pendingAction deviceAction

	toast   components.Toast
	confirm components.ConfirmPrompt
}

func NewDeviceInfo(state *state.AppState) *DeviceInfo {
	return &DeviceInfo{
		state: state,
		actions: []deviceAction{
			{"s", "Start scrcpy (screen mirror)", adb.StartScrcpy, true, "Start screen mirroring?"},
			{"r", "Reboot device", adb.Reboot, true, "Reboot the device?"},
			{"R", "Reboot to recovery", adb.RebootRecovery, true, "Reboot into recovery?"},
			{"b", "Reboot to bootloader", adb.RebootBootloader, true, "Reboot into bootloader?"},
		},
	}
}

func (d *DeviceInfo) Init() tea.Cmd {
	return nil
}

func (d *DeviceInfo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	d.toast.Update(msg)

	if d.confirm.Visible {
		if cmd := d.confirm.Update(msg); cmd != nil {
			return d, cmd
		}
	}

	switch msg := msg.(type) {
	case components.ConfirmYesMsg:
		d.confirm.Hide()
		return d, d.runAction(d.pendingAction)

	case components.ConfirmNoMsg:
		d.confirm.Hide()
		return d, nil

	case tea.KeyMsg:
		for _, a := range d.actions {
			if msg.String() != a.key {
				continue
			}

			if a.confirm {
				d.pendingAction = a
				d.confirm.Show(a.prompt)
				return d, nil
			}

			return d, d.runAction(a)
		}
	}

	return d, nil
}

func (d *DeviceInfo) runAction(a deviceAction) tea.Cmd {
	err := a.run(d.state.DeviceSerial())
	if err != nil {
		var cmd tea.Cmd
		d.toast, cmd = components.ShowToast(
			"Error: "+err.Error(),
			true,
			3*time.Second,
		)
		return cmd
	}

	var cmd tea.Cmd
	d.toast, cmd = components.ShowToast(
		"Command executed successfully",
		false,
		2*time.Second,
	)
	return cmd
}

func (d *DeviceInfo) View() string {
	if !d.state.HasDevice() {
		return components.RenderLayout(d.state, components.LayoutProps{
			Title:  "Device Info",
			Body:   components.StatusDisconnected.Render("No device selected"),
			Footer: components.Help("esc", "back"),
		})
	}

	dev := d.state.SelectedDevice
	var body strings.Builder

	body.WriteString(components.TitleStyle.Render("Device Details") + "\n")
	body.WriteString(
		components.KeyValueList([]components.KeyValueRow{
			{Key: "Model:", Value: dev.Model},
			{Key: "Serial:", Value: dev.Serial},
			{Key: "Android:", Value: dev.Android},
			{Key: "State:", Value: dev.State},
		}),
	)

	body.WriteString("\n")

	if d.confirm.Visible {
		body.WriteString(d.confirm.View())
	} else {
		body.WriteString(components.TitleStyle.Render("Actions") + "\n\n")
		for _, a := range d.actions {
			body.WriteString(
				components.HelpKeyStyle.Render("["+a.key+"]") +
					" " +
					components.ListItemStyle.Render(a.label) +
					"\n",
			)
		}
	}

	if d.toast.Visible {
		body.WriteString("\n\n")
		body.WriteString(d.toast.View())
	}

	return components.RenderLayout(d.state, components.LayoutProps{
		Title:  "Device Info",
		Body:   body.String(),
		Footer: components.Help("esc", "back") + "  " + components.Help("q", "quit"),
	})
}
