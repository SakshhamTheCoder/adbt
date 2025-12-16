package screens

import (
	"strings"
	"time"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type Devices struct {
	state   *state.AppState
	cursor  int
	loading bool

	form        components.FormModal
	formLoading bool
	toast       components.Toast
}

func NewDevices(state *state.AppState) *Devices {
	return &Devices{state: state}
}

func (d *Devices) Init() tea.Cmd {
	d.loading = true
	return adb.ListDevices()
}

func (d *Devices) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	d.toast.Update(msg)

	if d.form.Visible {
		switch m := msg.(type) {

		case components.FormSubmitMsg:
			d.form.Hide()
			d.formLoading = true

			ip := m.Values[0]
			port := m.Values[1]
			pin := m.Values[2]

			return d, adb.PairWirelessCmd(ip, port, pin)

		case components.FormCancelMsg:
			d.form.Hide()
			return d, nil
		}

		if cmd := d.form.Update(msg); cmd != nil {
			return d, cmd
		}
		return d, nil
	}

	switch msg := msg.(type) {

	case adb.PairWirelessResultMsg:
		d.formLoading = false

		if msg.Error != nil {
			var cmd tea.Cmd
			d.toast, cmd = components.ShowToast(
				"Pairing failed: "+msg.Error.Error(),
				true,
				3*time.Second,
			)
			return d, cmd
		}

		d.loading = true
		return d, adb.ListDevices()
	}

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
			if len(d.state.Devices) > 0 {
				d.state.SelectDevice(&d.state.Devices[d.cursor])
				return d, func() tea.Msg {
					return SwitchScreenMsg{Screen: "dashboard"}
				}
			}

		case "w":
			d.form.Show(
				"Pair Device Wirelessly",
				[]components.FormField{
					{Label: "IP Address"},
					{Label: "Port"},
					{Label: "PIN Code"},
				},
			)

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
	var body strings.Builder

	if d.loading {
		body.WriteString("Loading devices...")
	} else {
		body.WriteString(components.DeviceList(d.state.Devices, d.cursor))
	}

	if d.formLoading {
		body.WriteString("\n\n")
		body.WriteString(
			components.StatusMuted.Render("Connecting to device…"),
		)
	}

	if d.form.Visible {
		body.WriteString("\n\n")
		body.WriteString(d.form.View())
	}

	if d.toast.Visible {
		body.WriteString("\n\n")
		body.WriteString(d.toast.View())
	}

	return components.RenderLayout(d.state, components.LayoutProps{
		Title: "Device Selection",
		Body:  body.String(),
		Footer: components.Help("↑/↓", "navigate") + "  " +
			components.Help("enter", "select") + "  " +
			components.Help("w", "wireless pair") + "  " +
			components.Help("r", "refresh") + "  " +
			components.Help("esc", "back"),
	})
}
