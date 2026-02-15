package adb

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type FileEntry struct {
	Name        string
	Path        string
	IsDir       bool
	Size        string
	Permissions string
}

type FilesLoadedMsg struct {
	Path  string
	Files []FileEntry
	Error error
}

type FileActionResultMsg struct {
	Action string
	Error  error
}

func ListFilesCmd(serial, path string) tea.Cmd {
	return func() tea.Msg {
		listPath := path
		if !strings.HasSuffix(listPath, "/") {
			listPath += "/"
		}

		out, err := ExecuteCommand(
			serial,
			"shell",
			"ls",
			"-la",
			listPath,
		)
		if err != nil {
			return FilesLoadedMsg{
				Path:  path,
				Error: err,
			}
		}

		lines := ParseLines(out)
		files := make([]FileEntry, 0, len(lines))

		for _, line := range lines {
			if line == "" || strings.HasPrefix(line, "total") {
				continue
			}

			entry := parseLsLine(line, path)
			if entry.Name == "" || entry.Name == "." || entry.Name == ".." {
				continue
			}

			files = append(files, entry)
		}

		return FilesLoadedMsg{
			Path:  path,
			Files: files,
		}
	}
}

func parseLsLine(line, parentPath string) FileEntry {
	fields := strings.Fields(line)
	if len(fields) < 7 {
		return FileEntry{}
	}

	perms := fields[0]
	isDir := len(perms) > 0 && perms[0] == 'd'

	var name string
	var size string

	if len(fields) >= 8 {
		size = fields[4]
		name = strings.Join(fields[7:], " ")
	} else {
		name = fields[len(fields)-1]
	}

	if isDir {
		size = ""
	}

	return FileEntry{
		Name:        name,
		Path:        parentPath + "/" + name,
		IsDir:       isDir,
		Size:        size,
		Permissions: perms,
	}
}

func FormatFileSize(sizeStr string) string {
	if sizeStr == "" {
		return ""
	}

	var n int64
	_, err := fmt.Sscanf(sizeStr, "%d", &n)
	if err != nil {
		return sizeStr
	}

	switch {
	case n >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(n)/float64(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/float64(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(n)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", n)
	}
}

func DeleteFileCmd(serial, path string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(
			serial,
			"shell",
			"rm",
			"-rf",
			path,
		)

		return FileActionResultMsg{
			Action: "delete",
			Error:  err,
		}
	}
}

func PullFileCmd(serial, remotePath, localPath string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(
			serial,
			"pull",
			remotePath,
			localPath,
		)

		return FileActionResultMsg{
			Action: "pull",
			Error:  err,
		}
	}
}

func PushFileCmd(serial, localPath, remotePath string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(
			serial,
			"push",
			localPath,
			remotePath,
		)

		return FileActionResultMsg{
			Action: "push",
			Error:  err,
		}
	}
}
