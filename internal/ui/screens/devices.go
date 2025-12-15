package screens

import (
	"strings"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type Devices struct {
	state   *state.AppState
	cursor  int
	loading bool
}

func NewDevices(appState *state.AppState) *Devices {
	return &Devices{
		state: appState,
	}
}

func (d *Devices) Init() tea.Cmd {
	d.loading = true
	return adb.ListDevices()
}

func (d *Devices) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if d.cursor > 0 {
				d.cursor--
			}
		case "down", "j":
			if d.cursor < len(d.state.Devices)-1 {
				d.cursor++
			}
		case "enter":
			if len(d.state.Devices) > 0 && d.cursor < len(d.state.Devices) {
				device := &d.state.Devices[d.cursor]
				if device.IsConnected() {
					d.state.SelectDevice(device)
					return d, func() tea.Msg {
						return SwitchScreenMsg{Screen: "dashboard"}
					}
				}
			}
		case "r":
			d.loading = true
			return d, adb.ListDevices()
		}

	case adb.DevicesLoadedMsg:
		d.loading = false
		if msg.Error == nil {
			d.state.Devices = msg.Devices

			if d.cursor >= len(d.state.Devices) {
				d.cursor = len(d.state.Devices) - 1
			}
			if d.cursor < 0 {
				d.cursor = 0
			}
		}
	}

	return d, nil
}

func (d *Devices) View() string {
	var b strings.Builder

	b.WriteString(components.RenderHeader(d.state, "Device Selection") + "\n")

	if d.loading {
		content := components.ContentStyle.Width(d.state.Width - 4).Render("Loading devices...")
		b.WriteString(content + "\n")
	} else {
		deviceList := components.DeviceList(d.state.Devices, d.cursor)
		content := components.ContentStyle.Width(d.state.Width - 4).Render(deviceList)
		b.WriteString(content + "\n")
	}

	footer := components.Help("↑/k", "up") + "  " +
		components.Help("↓/j", "down") + "  " +
		components.Help("enter", "select") + "  " +
		components.Help("r", "refresh") + "  " +
		components.Help("esc", "back") + "  " +
		components.Help("q", "quit")
	b.WriteString(components.FooterStyle.Render(footer))

	return b.String()
}
