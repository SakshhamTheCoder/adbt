package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Toast struct {
	Message string
	IsError bool
	Visible bool
}

type clearToastMsg struct{}

func ShowToast(msg string, isError bool, d time.Duration) (Toast, tea.Cmd) {
	return Toast{
			Message: msg,
			IsError: isError,
			Visible: true,
		},
		tea.Tick(d, func(time.Time) tea.Msg {
			return clearToastMsg{}
		})
}

func (t *Toast) Update(msg tea.Msg) {
	if _, ok := msg.(clearToastMsg); ok {
		t.Visible = false
		t.Message = ""
		t.IsError = false
	}
}

func (t *Toast) View() string {
	if !t.Visible {
		return ""
	}

	if t.IsError {
		return ErrorStyle.Render(t.Message)
	}

	return StatusConnected.Render(t.Message)
}
