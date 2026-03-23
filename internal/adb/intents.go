package adb

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type IntentResultMsg struct {
	Output string
}

type IntentErrorMsg struct {
	Error error
}

func SendIntentCmd(serial, action, dataURI, extras string) tea.Cmd {
	return func() tea.Msg {
		args := []string{
			"shell", "am", "start",
			"-W",
			"-a", action,
		}

		if dataURI != "" {
			args = append(args, "-d", dataURI)
		}

		if extras != "" {
			for _, extra := range strings.Split(extras, ";") {
				extra = strings.TrimSpace(extra)
				if extra == "" {
					continue
				}
				parts := strings.SplitN(extra, "=", 2)
				if len(parts) == 2 {
					args = append(args, "--es", parts[0], parts[1])
				}
			}
		}

		out, err := ExecuteCommand(serial, args...)
		if err != nil {
			return IntentErrorMsg{Error: fmt.Errorf("%w: %s", err, string(out))}
		}
		return IntentResultMsg{Output: strings.TrimSpace(string(out))}
	}
}

func SendBroadcastCmd(serial, action, extras string) tea.Cmd {
	return func() tea.Msg {
		args := []string{
			"shell", "am", "broadcast",
			"-a", action,
		}

		if extras != "" {
			for _, extra := range strings.Split(extras, ";") {
				extra = strings.TrimSpace(extra)
				if extra == "" {
					continue
				}
				parts := strings.SplitN(extra, "=", 2)
				if len(parts) == 2 {
					args = append(args, "--es", parts[0], parts[1])
				}
			}
		}

		out, err := ExecuteCommand(serial, args...)
		if err != nil {
			return IntentErrorMsg{Error: fmt.Errorf("%w: %s", err, string(out))}
		}
		return IntentResultMsg{Output: strings.TrimSpace(string(out))}
	}
}
