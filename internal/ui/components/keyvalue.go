package components

import "strings"

type KeyValueRow struct {
	Key   string
	Value string
}

func KeyValueList(rows []KeyValueRow) string {
	var b strings.Builder
	for _, r := range rows {
		b.WriteString(StatusMuted.Render(r.Key))
		b.WriteString(" ")
		b.WriteString(r.Value)
		b.WriteString("\n")
	}
	return b.String()
}
