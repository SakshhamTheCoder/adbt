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
	key   string
	label string
	cmd   func(string) tea.Cmd
}

type DeviceInfo struct {
	state   *state.AppState
	actions []deviceAction
	cursor  int
	toast   components.Toast
}

func NewDeviceInfo(state *state.AppState) *DeviceInfo {
	return &DeviceInfo{
		state: state,
		actions: []deviceAction{
			{"c", "Start scrcpy", adb.StartScrcpyCmd},
			{"w", "Toggle Wi-Fi", adb.ToggleWifiCmd},
			{"s", "Toggle Screen", adb.ToggleScreenCmd},
			{"r", "Reboot device", adb.RebootCmd},
			{"R", "Reboot to recovery", adb.RebootRecoveryCmd},
			{"b", "Reboot to bootloader", adb.RebootBootloaderCmd},
		},
	}
}

func (d *DeviceInfo) Init() tea.Cmd {
	return nil
}

func (d *DeviceInfo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	d.toast.Update(msg)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if !d.state.HasDevice() {
			return d, nil
		}

		switch msg.String() {

		case "up", "k":
			if d.cursor > 0 {
				d.cursor--
			}

		case "down", "j":
			if d.cursor < len(d.actions)-1 {
				d.cursor++
			}

		case "enter":
			a := d.actions[d.cursor]
			return d, a.cmd(d.state.DeviceSerial())

		default:

			for _, a := range d.actions {
				if msg.String() == a.key {
					return d, a.cmd(d.state.DeviceSerial())
				}
			}
		}

	case adb.DeviceActionResultMsg:
		var cmd tea.Cmd
		d.toast, cmd = components.ShowToast(
			msg.Action+" executed successfully",
			false,
			2*time.Second,
		)
		return d, cmd

	case adb.DeviceActionErrorMsg:
		var cmd tea.Cmd
		d.toast, cmd = components.ShowToast(
			msg.Action+" failed: "+msg.Error.Error(),
			true,
			3*time.Second,
		)
		return d, cmd
	}

	return d, nil
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
	body.WriteString(components.TitleStyle.Render("Actions") + "\n")

	for i, a := range d.actions {
		line := "  "
		if i == d.cursor {
			line = "› "
		}

		if i == d.cursor {
			line += components.HelpKeyStyle.Render("[" + a.key + "]")
			line += " " + components.ListItemSelectedStyle.Render(a.label)
		} else {
			line += components.StatusMuted.Render("[" + a.key + "] ")
			line += components.ListItemStyle.Render(a.label)
		}

		body.WriteString(line + "\n")
	}

	if d.toast.Visible {
		body.WriteString("\n")
		body.WriteString(d.toast.View())
	}

	return components.RenderLayout(d.state, components.LayoutProps{
		Title: "Device Info",
		Body:  body.String(),
		Footer: components.Help("↑/↓", "navigate") + "  " +
			components.Help("enter", "select") + "  " +
			components.Help("esc", "back") + "  ",
	})
}
