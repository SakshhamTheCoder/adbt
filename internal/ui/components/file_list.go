package components

import (
	"fmt"
	"strings"

	"github.com/SakshhamTheCoder/adbt/internal/adb"
)

func FileList(files []adb.FileEntry, cursor int) string {
	var out strings.Builder

	for i, f := range files {
		prefix := "  "
		if i == cursor {
			prefix = "› "
		}

		icon := "📄"
		if f.IsDir {
			icon = "📁"
		}

		size := ""
		if !f.IsDir && f.Size != "" {
			size = " " + StatusMuted.Render(adb.FormatFileSize(f.Size))
		}

		line := fmt.Sprintf(
			"%s%s %s",
			prefix,
			icon,
			f.Name,
		)

		if i == cursor {
			out.WriteString(ListItemSelectedStyle.Render(line) + size)
		} else {
			out.WriteString(ListItemStyle.Render(line) + size)
		}

		out.WriteString("\n")
	}

	return out.String()
}
