package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"

	"github.com/SakshhamTheCoder/adbt/internal/adb"
	"github.com/SakshhamTheCoder/adbt/internal/state"
	"github.com/SakshhamTheCoder/adbt/internal/ui/components"
	"github.com/SakshhamTheCoder/adbt/internal/ui/navigation"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	key         string
	label       string
	description string
	screen      navigation.Screen
}

type Dashboard struct {
	state     *state.AppState
	loading   bool
	menuItems []menuItem
	cursor    int
	viewport  viewport.Model
}

func NewDashboard(appState *state.AppState) *Dashboard {
	return &Dashboard{
		state:    appState,
		viewport: viewport.New(0, 0),
		menuItems: []menuItem{
			{"d", "Devices", "View and select connected devices", navigation.ScreenDevices},
			{"i", "Device Info", "View device details and controls", navigation.ScreenDeviceInfo},
			{"l", "Logcat", "View live device logs", navigation.ScreenLogcat},
			{"a", "Apps", "Manage installed applications", navigation.ScreenApps},
			{"f", "Files", "Browse device file system", navigation.ScreenFiles},
			{"m", "Monitor", "Performance stats (CPU, RAM, Net)", navigation.ScreenPerfMonitor},
			{"t", "Intent Tester", "Test deep links and intents", navigation.ScreenIntents},
			{"p", "Port Forwarding", "Manage adb port forwarding", navigation.ScreenPorts},
		},
	}
}

func (d *Dashboard) Init() tea.Cmd {
	d.loading = true
	return adb.ListDevicesCmd()
}

func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case adb.DevicesLoadedMsg:
		d.loading = false
		if msg.Error == nil {
			d.state.Devices = msg.Devices
			if len(msg.Devices) == 1 && msg.Devices[0].IsConnected() {
				d.state.SelectDevice(msg.Devices[0].Serial)
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if d.cursor > 0 {
				d.cursor--
				d.ensureCursorVisible()
			}
		case "down", "j":
			if d.cursor < len(d.menuItems)-1 {
				d.cursor++
				d.ensureCursorVisible()
			}
		case "enter":
			item := d.menuItems[d.cursor]
			return d, navigation.Navigate(item.screen, d.state)
		default:
			for _, item := range d.menuItems {
				if msg.String() == item.key {
					return d, navigation.Navigate(item.screen, d.state)
				}
			}
		}
	}
	return d, nil
}

func (d *Dashboard) ensureCursorVisible() {
	ensureViewportLineVisible(&d.viewport, d.cursor)
}

func (d *Dashboard) View() string {
	var staticContent strings.Builder
	var scrollableContent strings.Builder

	if d.loading {
		staticContent.WriteString(components.StatusMuted.Render("Loading devices..."))
	} else {
		staticContent.WriteString(components.SectionTitle("Device") + "\n")

		if dev := d.state.SelectedDevice(); dev != nil {
			staticContent.WriteString(
				components.KeyValueList([]components.KeyValueRow{
					{Key: "Status:", Value: components.StatusConnected.Render("● Connected")},
					{Key: "Model:", Value: dev.Model},
					{Key: "Serial:", Value: dev.Serial},
				}),
			)
		} else {
			staticContent.WriteString(components.StatusDisconnected.Render("● No device connected") + "\n")
			staticContent.WriteString(components.StatusMuted.Render("Connect a device with USB debugging enabled."))
		}

		staticContent.WriteString("\n\n")
		staticContent.WriteString(components.SectionTitle("Quick Actions") + "\n")

		for i, item := range d.menuItems {
			line := "  "
			if i == d.cursor {
				line = "› "
			}

			disabled := navigation.Registry[item.screen].RequireDevice && !d.state.HasDevice()

			paddedLabel := fmt.Sprintf("%-16s", item.label)
			if i == d.cursor {
				line += components.HelpKeyStyle.Render("[" + item.key + "]")
				line += " " + components.ListItemSelectedStyle.Render(paddedLabel)
			} else {
				line += components.StatusMuted.Render("[" + item.key + "] ")
				line += components.ListItemStyle.Render(paddedLabel)
			}

			line += " " + components.StatusMuted.Render("- "+item.description)

			if disabled {
				line += " " + components.ErrorStyle.Render("(requires device)")
			}

			scrollableContent.WriteString(line + "\n")
		}
	}

	return components.RenderLayoutWithScrollableSection(d.state, components.LayoutWithScrollProps{
		Title:             "Dashboard",
		StaticContent:     staticContent.String(),
		ScrollableContent: scrollableContent.String(),
		Viewport:          &d.viewport,
		Footer:            components.JoinHelp([2]string{"↑/↓", "navigate"}, [2]string{"enter", "select"}),
	})
}
