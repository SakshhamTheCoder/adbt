package screens

import (
	"fmt"
	"strings"
	"time"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type deviceAction struct {
	key         string
	label       string
	cmd         func(string) tea.Cmd
	destructive bool
}

type DeviceInfo struct {
	state   *state.AppState
	actions []deviceAction
	cursor  int

	confirm components.ConfirmPrompt
	pending *deviceAction

	toast   components.Toast
	details *adb.DeviceDetails
	loading bool
}

func NewDeviceInfo(state *state.AppState) *DeviceInfo {
	return &DeviceInfo{
		state: state,
		actions: []deviceAction{
			{"c", "Start scrcpy", adb.StartScrcpyCmd, false},
			{"w", "Toggle Wi-Fi", adb.ToggleWifiCmd, false},
			{"s", "Toggle Screen", adb.ToggleScreenCmd, false},

			{"r", "Reboot device", adb.RebootCmd, true},
			{"R", "Reboot to recovery", adb.RebootRecoveryCmd, true},
			{"b", "Reboot to bootloader", adb.RebootBootloaderCmd, true},
		},
	}
}

func (d *DeviceInfo) Init() tea.Cmd {
	if !d.state.HasDevice() {
		return nil
	}
	d.loading = true
	return adb.FetchDeviceDetailsCmd(d.state.DeviceSerial())
}

func (d *DeviceInfo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	d.toast.Update(msg)

	if d.confirm.Visible {
		switch msg.(type) {

		case components.ConfirmYesMsg:
			a := d.pending
			d.pending = nil
			d.confirm.Hide()
			return d, a.cmd(d.state.DeviceSerial())

		case components.ConfirmNoMsg:
			d.pending = nil
			d.confirm.Hide()
			return d, tea.Batch()
		}

		return d, d.confirm.Update(msg)
	}

	switch msg := msg.(type) {

	case adb.DeviceDetailsMsg:
		d.loading = false
		if msg.Error == nil {
			d.details = &msg.Details
		}

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
			return d, d.triggerAction(d.actions[d.cursor])

		default:
			for i := range d.actions {
				if msg.String() == d.actions[i].key {
					d.cursor = i
					return d, d.triggerAction(d.actions[i])
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

func (d *DeviceInfo) triggerAction(a deviceAction) tea.Cmd {
	if a.destructive {
		d.pending = &a
		d.confirm.Show(a.label + "?")
		return nil
	}

	return a.cmd(d.state.DeviceSerial())
}

func infoCard(label, value string) string {
	if value == "" {
		value = components.StatusMuted.Render("—")
	}
	return fmt.Sprintf(
		"%s  %s",
		components.StatusMuted.Render(label),
		lipgloss.NewStyle().Bold(true).Render(value),
	)
}

func (d *DeviceInfo) View() string {
	if !d.state.HasDevice() {
		return components.RenderNoDevice(d.state, "Device Info")
	}

	dev := d.state.SelectedDevice
	colWidth := (d.state.Width - 12) / 2
	if colWidth < 20 {
		colWidth = 20
	}
	colStyle := lipgloss.NewStyle().Width(colWidth)

	var body strings.Builder

	body.WriteString(components.TitleStyle.Render("Device") + "\n\n")

	body.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		colStyle.Render(infoCard("Model", dev.Model)),
		colStyle.Render(infoCard("Android", dev.Android)),
	) + "\n")

	body.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		colStyle.Render(infoCard("Serial", dev.Serial)),
		colStyle.Render(infoCard("State", dev.State)),
	) + "\n")

	body.WriteString("\n" + components.TitleStyle.Render("Details") + "\n\n")

	if d.loading {
		body.WriteString(components.StatusMuted.Render("  Loading details...") + "\n")
	} else if d.details != nil {
		dt := d.details

		batteryVal := dt.BatteryLevel
		if batteryVal != "" && dt.BatteryStatus != "" {
			batteryVal += " (" + dt.BatteryStatus + ")"
		}

		storageVal := ""
		if dt.StorageUsed != "" && dt.StorageTotal != "" {
			storageVal = dt.StorageUsed + " / " + dt.StorageTotal
		}

		screenVal := dt.ScreenSize
		if screenVal != "" && dt.ScreenDensity != "" {
			screenVal += " @ " + dt.ScreenDensity + "dpi"
		}

		body.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			colStyle.Render(infoCard("Battery", batteryVal)),
			colStyle.Render(infoCard("Storage", storageVal)),
		) + "\n")

		body.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			colStyle.Render(infoCard("Screen", screenVal)),
			colStyle.Render(infoCard("IP", dt.IPAddress)),
		) + "\n")
	} else {
		body.WriteString(components.StatusMuted.Render("  Could not load details") + "\n")
	}

	body.WriteString("\n" + components.TitleStyle.Render("Actions") + "\n")

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

		if a.destructive {
			line += " " + components.ErrorStyle.Render("!")
		}

		body.WriteString(line + "\n")
	}

	rendered := components.RenderLayoutWithScrollableSection(d.state, components.LayoutWithScrollProps{
		Title:             "Device Info",
		ScrollableContent: body.String(),
		Footer: components.Help("↑/↓", "navigate") + "  " +
			components.Help("enter", "select") + "  " +
			components.Help("esc", "back"),
	})

	if d.confirm.Visible {
		rendered = components.RenderOverlay(rendered, d.confirm.View(), d.state)
	}

	if d.toast.Visible {
		rendered = components.RenderOverlay(rendered, d.toast.View(), d.state)
	}

	return rendered
}
