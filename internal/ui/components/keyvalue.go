package components

import "github.com/charmbracelet/lipgloss"

type KeyValueRow struct {
	Key   string
	Value string
}

func KeyValueList(rows []KeyValueRow) string {
	label := StatusMuted.
		Align(lipgloss.Right)

	value := StatusMuted.
		Align(lipgloss.Left)

	out := ""
	for _, r := range rows {
		out += label.Render(r.Key) + " " + value.Render(r.Value) + "\n"
	}
	return out
}
