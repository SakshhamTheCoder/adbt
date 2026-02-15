package components

import (
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type ConfirmYesMsg struct{}
type ConfirmNoMsg struct{}

type ConfirmPrompt struct {
	Visible bool
	Message string
}

func (c *ConfirmPrompt) Show(message string) {
	c.Visible = true
	c.Message = message
}

func (c *ConfirmPrompt) Hide() {
	c.Visible = false
	c.Message = ""
}

func (c *ConfirmPrompt) Update(msg tea.Msg) tea.Cmd {
	if !c.Visible {
		return nil
	}

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "y", "enter":
			return func() tea.Msg { return ConfirmYesMsg{} }
		case "n", "esc":
			return func() tea.Msg { return ConfirmNoMsg{} }
		}
	}

	return nil
}

var confirmBoxStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(Error).
	Padding(0, 2)

func (c *ConfirmPrompt) View() string {
	if !c.Visible {
		return ""
	}

	content := ErrorStyle.Render("⚠ "+c.Message) + "\n" +
		HelpKeyStyle.Render("[y]") + " " + HelpDescStyle.Render("Yes") + "  " +
		HelpKeyStyle.Render("[n]") + " " + HelpDescStyle.Render("No")

	return confirmBoxStyle.Render(content)
}
