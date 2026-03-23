package adb

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type PortForward struct {
	Serial string
	Local  string
	Remote string
}

type PortsLoadedMsg struct {
	Ports []PortForward
}

type PortsLoadErrorMsg struct {
	Error error
}

type PortActionResultMsg struct {
	Action string
}

type PortActionErrorMsg struct {
	Action string
	Error  error
}

func ListPortForwardsCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		out, err := ExecuteCommand("", "forward", "--list")
		if err != nil {
			return PortsLoadErrorMsg{Error: err}
		}

		lines := ParseLines(out)
		ports := make([]PortForward, 0, len(lines))

		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}
			// Only show forwardings for the selected device
			if serial != "" && fields[0] != serial {
				continue
			}
			ports = append(ports, PortForward{
				Serial: fields[0],
				Local:  fields[1],
				Remote: fields[2],
			})
		}

		return PortsLoadedMsg{Ports: ports}
	}
}

func ForwardPortCmd(serial, local, remote string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(serial, "forward", local, remote)
		if err != nil {
			return PortActionErrorMsg{
				Action: "forward",
				Error:  fmt.Errorf("failed to forward %s → %s: %w", local, remote, err),
			}
		}
		return PortActionResultMsg{Action: "forward"}
	}
}

func ReversePortCmd(serial, remote, local string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(serial, "reverse", remote, local)
		if err != nil {
			return PortActionErrorMsg{
				Action: "reverse",
				Error:  fmt.Errorf("failed to reverse %s → %s: %w", remote, local, err),
			}
		}
		return PortActionResultMsg{Action: "reverse"}
	}
}

func RemoveForwardCmd(serial, local string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(serial, "forward", "--remove", local)
		if err != nil {
			return PortActionErrorMsg{
				Action: "remove",
				Error:  err,
			}
		}
		return PortActionResultMsg{Action: "remove"}
	}
}
