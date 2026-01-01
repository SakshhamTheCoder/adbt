package components

import (
	"fmt"

	"adbt/internal/adb"
)

func FileList(files []adb.FileEntry, cursor int) string {
	var out string

	for i, f := range files {
		prefix := "  "
		if i == cursor {
			prefix = "› "
		}

		icon := "📄"
		if f.IsDir {
			icon = "📁"
		}

		line := fmt.Sprintf(
			"%s%s %s",
			prefix,
			icon,
			f.Name,
		)

		if i == cursor {
			out += ListItemSelectedStyle.Render(line)
		} else {
			out += ListItemStyle.Render(line)
		}

		out += "\n"
	}

	return out
}
