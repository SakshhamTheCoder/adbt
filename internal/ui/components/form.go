package components

import (
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
)

type FormSubmitMsg struct {
	Values []string
}

type FormCancelMsg struct{}

type FormField struct {
	Label string
	Value string
}

type FormModal struct {
	Visible bool
	Title   string
	Fields  []FormField
	Cursor  int
}

func (f *FormModal) Show(title string, fields []FormField) {
	f.Visible = true
	f.Title = title
	f.Fields = fields
	f.Cursor = 0
}

func (f *FormModal) Hide() {
	f.Visible = false
	f.Title = ""
	f.Fields = nil
	f.Cursor = 0
}

func (f *FormModal) Update(msg tea.Msg) tea.Cmd {
	if !f.Visible || len(f.Fields) == 0 {
		return nil
	}

	if f.Cursor < 0 {
		f.Cursor = 0
	}
	if f.Cursor >= len(f.Fields) {
		f.Cursor = len(f.Fields) - 1
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch key.String() {

	case "tab", "down", "j":
		if f.Cursor < len(f.Fields)-1 {
			f.Cursor++
		}

	case "shift+tab", "up", "k":
		if f.Cursor > 0 {
			f.Cursor--
		}

	case "enter":
		values := make([]string, len(f.Fields))
		for i := range f.Fields {
			values[i] = f.Fields[i].Value
		}
		return func() tea.Msg {
			return FormSubmitMsg{Values: values}
		}

	case "esc":
		return func() tea.Msg {
			return FormCancelMsg{}
		}

	case "backspace":
		field := &f.Fields[f.Cursor]
		if len(field.Value) > 0 {
			field.Value = field.Value[:len(field.Value)-1]
		}

	default:
		if len(key.Runes) == 1 && unicode.IsPrint(key.Runes[0]) {
			f.Fields[f.Cursor].Value += string(key.Runes)
		}
	}

	return nil
}

func (f *FormModal) View() string {
	if !f.Visible {
		return ""
	}

	var b strings.Builder

	b.WriteString(TitleStyle.Render(f.Title) + "\n\n")

	for i, field := range f.Fields {
		prefix := "  "
		if i == f.Cursor {
			prefix = "› "
		}

		value := field.Value
		if value == "" {
			value = StatusMuted.Render("<enter value>")
		}

		b.WriteString(
			prefix +
				StatusMuted.Render(field.Label+":") +
				" " +
				value +
				"\n",
		)
	}

	b.WriteString("\n")
	b.WriteString(
		HelpKeyStyle.Render("enter") + " submit    " +
			HelpKeyStyle.Render("esc") + " cancel    " +
			HelpKeyStyle.Render("tab") + " next",
	)

	return ContentStyle.Render(b.String())
}
