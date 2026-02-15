package components

import (
	"time"

	"github.com/charmbracelet/lipgloss"

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

	toastStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(0, 2)

	if t.IsError {
		toastStyle = toastStyle.BorderForeground(Error)
		return toastStyle.Render(ErrorStyle.Render("✗ " + t.Message))
	}

	toastStyle = toastStyle.BorderForeground(Success)
	return toastStyle.Render(StatusConnected.Render("✓ " + t.Message))
}
