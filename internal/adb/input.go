package adb

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type InputResultMsg struct {
	Error error
}

// escapeInputText prepares a string for `adb shell input text`, which treats
// spaces as argument separators and uses %s to represent a literal space.
func escapeInputText(text string) string {
	return strings.ReplaceAll(text, " ", "%s")
}

func SendTextCmd(serial, text string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(
			serial,
			"shell",
			"input",
			"text",
			escapeInputText(text),
		)
		return InputResultMsg{Error: err}
	}
}

func SendKeyEventCmd(serial string, keycode int) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(
			serial,
			"shell",
			"input",
			"keyevent",
			strconv.Itoa(keycode),
		)
		return InputResultMsg{Error: err}
	}
}
