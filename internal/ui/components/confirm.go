package components

import tea "github.com/charmbracelet/bubbletea"

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

func (c *ConfirmPrompt) View() string {
	if !c.Visible {
		return ""
	}

	out := ""
	out += ErrorStyle.Render(
		"Are you sure?\n" +
			c.Message +
			"\n\n",
	)
	out += HelpKeyStyle.Render("[y]") + " Yes    "
	out += HelpKeyStyle.Render("[n]") + " No"

	return out
}
