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
	toast   components.Toast
}

func NewDeviceInfo(state *state.AppState) *DeviceInfo {
	return &DeviceInfo{
		state: state,
		actions: []deviceAction{
			{"s", "Start scrcpy", adb.StartScrcpyCmd},
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
		for _, a := range d.actions {
			if msg.String() == a.key {
				return d, a.cmd(d.state.DeviceSerial())
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

	body.WriteString("\n\n")
	body.WriteString(components.TitleStyle.Render("Actions") + "\n\n")

	for _, a := range d.actions {
		body.WriteString(
			components.HelpKeyStyle.Render("["+a.key+"]") +
				" " +
				components.ListItemStyle.Render(a.label) +
				"\n",
		)
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
