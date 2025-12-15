package screens

import (
	"strings"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	key           string
	label         string
	description   string
	action        Action
	requireDevice bool
}

type Dashboard struct {
	state     *state.AppState
	loading   bool
	menuItems []menuItem
	cursor    int
}

func NewDashboard(appState *state.AppState) *Dashboard {
	return &Dashboard{
		state: appState,
		menuItems: []menuItem{
			{"d", "Devices", "View and select connected devices", ActionDevices, false},
			{"i", "Device Info", "View device details and controls", ActionDeviceInfo, true},
			{"l", "Logcat", "View live device logs", ActionLogcat, true},
			{"s", "Shell", "Interactive ADB shell", ActionShell, true},
			{"a", "Apps", "Manage installed applications", ActionApps, true},
			{"f", "Files", "Browse device storage", ActionFiles, true},
		},
	}
}

func (d *Dashboard) Init() tea.Cmd {
	d.loading = true
	return adb.ListDevices()
}

func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case adb.DevicesLoadedMsg:
		d.loading = false
		if msg.Error == nil {
			d.state.Devices = msg.Devices
			if len(msg.Devices) == 1 && msg.Devices[0].IsConnected() {
				d.state.SelectDevice(&msg.Devices[0])
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if d.cursor > 0 {
				d.cursor--
			}
		case "down", "j":
			if d.cursor < len(d.menuItems)-1 {
				d.cursor++
			}
		case "enter":
			item := d.menuItems[d.cursor]
			return d, ResolveAction(item.action, d.state)
		default:
			for _, item := range d.menuItems {
				if msg.String() == item.key {
					return d, ResolveAction(item.action, d.state)
				}
			}
		}
	}
	return d, nil
}

func (d *Dashboard) View() string {
	if d.loading {
		return components.RenderLayout(d.state, components.LayoutProps{
			Title: "Dashboard",
			Body:  "Loading devices...",
		})
	}

	var body strings.Builder

	body.WriteString(components.TitleStyle.Render("Device") + "\n")

	if d.state.HasDevice() {
		dev := d.state.SelectedDevice
		body.WriteString(
			components.KeyValueList([]components.KeyValueRow{
				{Key: "Status:", Value: components.StatusConnected.Render("● Connected")},
				{Key: "Model:", Value: dev.Model},
				{Key: "Serial:", Value: dev.Serial},
				{Key: "Android:", Value: dev.Android},
				{Key: "State:", Value: dev.State},
			}),
		)
	} else {
		body.WriteString(components.StatusDisconnected.Render("● No device connected\n"))
		body.WriteString(components.StatusMuted.Render(
			"Connect a device with USB debugging enabled.\nPress D to open device manager.",
		))
	}

	body.WriteString("\n")
	body.WriteString(components.TitleStyle.Render("Quick Actions") + "\n")

	for i, item := range d.menuItems {
		line := "  "
		if i == d.cursor {
			line = "› "
		}

		disabled := item.requireDevice && !d.state.HasDevice()

		if i == d.cursor {
			line += components.HelpKeyStyle.Render("[" + item.key + "]")
			line += " " + components.ListItemSelectedStyle.Render(item.label)
		} else {
			line += components.StatusMuted.Render("[" + item.key + "] ")
			line += components.ListItemStyle.Render(item.label)
		}

		line += " " + components.StatusMuted.Render("- "+item.description)

		if disabled {
			line += " " + components.ErrorStyle.Render("(requires device)")
		}

		body.WriteString(line + "\n")
	}

	return components.RenderLayout(d.state, components.LayoutProps{
		Title:  "Dashboard",
		Body:   body.String(),
		Footer: components.Help("↑/↓", "navigate") + "  " + components.Help("enter", "select"),
	})
}
