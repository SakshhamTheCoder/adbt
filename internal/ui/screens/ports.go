package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/SakshhamTheCoder/adbt/internal/adb"
	"github.com/SakshhamTheCoder/adbt/internal/state"
	"github.com/SakshhamTheCoder/adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type portMode int

const (
	portModeForward portMode = iota
	portModeReverse
)

var portModeNames = []string{"Forward", "Reverse"}

type Ports struct {
	state   *state.AppState
	ports   []adb.PortForward
	loading bool
	cursor  int

	mode    portMode
	form    components.FormModal
	confirm components.ConfirmPrompt
	toast   components.Toast
}

func NewPorts(state *state.AppState) *Ports {
	return &Ports{
		state: state,
	}
}

func (p *Ports) Init() tea.Cmd {
	if !p.state.HasDevice() {
		return nil
	}
	p.loading = true
	return adb.ListPortForwardsCmd(p.state.DeviceSerial())
}

func (p *Ports) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p.toast.Update(msg)

	if p.form.Visible {
		switch msg := msg.(type) {
		case components.FormSubmitMsg:
			values := msg.Values
			p.form.Hide()
			local := ""
			remote := ""
			if len(values) > 0 {
				local = values[0]
			}
			if len(values) > 1 {
				remote = values[1]
			}

			if local == "" || remote == "" {
				var cmd tea.Cmd
				p.toast, cmd = components.ShowToast(
					"Both local and remote are required",
					true,
					2*time.Second,
				)
				return p, cmd
			}

			serial := p.state.DeviceSerial()
			var toastCmd tea.Cmd

			if p.mode == portModeReverse {
				p.toast, toastCmd = components.ShowToast(
					fmt.Sprintf("Reversing %s → %s...", remote, local),
					false,
					2*time.Second,
				)
				return p, tea.Batch(
					toastCmd,
					adb.ReversePortCmd(serial, remote, local),
				)
			}

			p.toast, toastCmd = components.ShowToast(
				fmt.Sprintf("Forwarding %s → %s...", local, remote),
				false,
				2*time.Second,
			)
			return p, tea.Batch(
				toastCmd,
				adb.ForwardPortCmd(serial, local, remote),
			)

		case components.FormCancelMsg:
			p.form.Hide()
			return p, nil
		}
		return p, p.form.Update(msg)
	}

	if p.confirm.Visible {
		switch msg.(type) {
		case components.ConfirmYesMsg:
			p.confirm.Hide()
			if p.cursor < len(p.ports) {
				port := p.ports[p.cursor]
				return p, adb.RemoveForwardCmd(p.state.DeviceSerial(), port.Local)
			}
			return p, nil
		case components.ConfirmNoMsg:
			p.confirm.Hide()
			return p, nil
		}
		return p, p.confirm.Update(msg)
	}

	switch msg := msg.(type) {
	case adb.PortsLoadedMsg:
		p.loading = false
		p.ports = msg.Ports
		if p.cursor >= len(p.ports) {
			p.cursor = 0
		}
		components.ViewportGotoTop("Ports")

	case adb.PortsLoadErrorMsg:
		p.loading = false
		var cmd tea.Cmd
		p.toast, cmd = components.ShowToast(
			"Failed to load port forwards",
			true,
			3*time.Second,
		)
		return p, cmd

	case adb.PortActionResultMsg:
		var cmd tea.Cmd
		p.toast, cmd = components.ShowToast(
			msg.Action+" successful",
			false,
			2*time.Second,
		)
		return p, tea.Batch(
			cmd,
			adb.ListPortForwardsCmd(p.state.DeviceSerial()),
		)

	case adb.PortActionErrorMsg:
		var cmd tea.Cmd
		p.toast, cmd = components.ShowToast(
			msg.Action+" failed: "+msg.Error.Error(),
			true,
			3*time.Second,
		)
		return p, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
				components.ViewportEnsureVisible("Ports", p.cursor)
			}
		case "down", "j":
			if p.cursor < len(p.ports)-1 {
				p.cursor++
				components.ViewportEnsureVisible("Ports", p.cursor)
			}
		case "a":
			title := "Add Forward Rule"
			if p.mode == portModeReverse {
				title = "Add Reverse Rule"
			}
			p.form.Show(title, []components.FormField{
				{Label: "Local", Value: "tcp:"},
				{Label: "Remote", Value: "tcp:"},
			})
		case "d":
			if len(p.ports) > 0 && p.cursor < len(p.ports) {
				port := p.ports[p.cursor]
				p.confirm.Show("Remove forward:\n" + port.Local + " → " + port.Remote)
			}
		case "t":
			p.mode = (p.mode + 1) % 2
		case "r":
			p.loading = true
			p.cursor = 0
			components.ViewportGotoTop("Ports")
			return p, adb.ListPortForwardsCmd(p.state.DeviceSerial())
		default:
			return p, components.UpdateViewport("Ports", msg)
		}
	}

	return p, nil
}

func (p *Ports) View() string {
	if !p.state.HasDevice() {
		return components.RenderNoDevice(p.state, "Ports")
	}

	var staticContent strings.Builder
	staticContent.WriteString(components.TitleStyle.Render("Port Forwarding") + "\n")

	staticContent.WriteString("  ")
	for idx, name := range portModeNames {
		if portMode(idx) == p.mode {
			staticContent.WriteString(components.HelpKeyStyle.Render(name))
		} else {
			staticContent.WriteString(components.StatusMuted.Render(name))
		}
		if idx < len(portModeNames)-1 {
			staticContent.WriteString(components.StatusMuted.Render(" / "))
		}
	}
	staticContent.WriteString("\n")

	maxWidth := p.state.Width - 8
	if maxWidth < 20 {
		maxWidth = 20
	}
	truncStyle := lipgloss.NewStyle().MaxWidth(maxWidth)

	var scrollableContent strings.Builder

	if p.loading {
		scrollableContent.WriteString(components.StatusMuted.Render("Loading port forwards..."))
	} else if len(p.ports) == 0 {
		scrollableContent.WriteString(components.StatusMuted.Render("No active port forwards") + "\n")
		scrollableContent.WriteString(components.StatusMuted.Render("Press [a] to add a new rule"))
	} else {
		for i, port := range p.ports {
			prefix := "  "
			if i == p.cursor {
				prefix = "› "
			}

			label := fmt.Sprintf("%s → %s", port.Local, port.Remote)

			var line string
			if i == p.cursor {
				line = prefix + components.ListItemSelectedStyle.Render(label)
			} else {
				line = prefix + components.ListItemStyle.Render(label)
			}
			scrollableContent.WriteString(truncStyle.Render(line) + "\n")
		}
	}

	footer := components.Help("↑/↓", "navigate") + "  " +
		components.Help("a", "add") + "  " +
		components.Help("d", "remove") + "  " +
		components.Help("t", "toggle mode") + "  " +
		components.Help("r", "refresh") + "  " +
		components.Help("esc", "back")

	rendered := components.RenderLayoutWithScrollableSection(p.state, components.LayoutWithScrollProps{
		Title:             "Ports",
		StaticContent:     staticContent.String(),
		ScrollableContent: scrollableContent.String(),
		Footer:            footer,
	})

	if p.form.Visible {
		rendered = components.RenderOverlay(rendered, p.form.View(), p.state)
	}

	if p.confirm.Visible {
		rendered = components.RenderOverlay(rendered, p.confirm.View(), p.state)
	}

	if p.toast.Visible {
		rendered = components.RenderOverlay(rendered, p.toast.View(), p.state)
	}

	return rendered
}
