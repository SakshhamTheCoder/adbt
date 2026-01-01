package adb

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

/* ---------- models ---------- */

type FileEntry struct {
	Name  string
	Path  string
	IsDir bool
}

/* ---------- messages ---------- */

type FilesLoadedMsg struct {
	Path  string
	Files []FileEntry
	Error error
}

type FileActionResultMsg struct {
	Action string
	Error  error
}

/* ---------- list directory ---------- */

func ListFilesCmd(serial, path string) tea.Cmd {
	return func() tea.Msg {
		out, err := ExecuteCommand(
			serial,
			"shell",
			"ls",
			"-p",
			path,
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
			if line == "" {
				continue
			}

			isDir := strings.HasSuffix(line, "/")
			name := strings.TrimSuffix(line, "/")

			files = append(files, FileEntry{
				Name:  name,
				Path:  path + "/" + name,
				IsDir: isDir,
			})
		}

		return FilesLoadedMsg{
			Path:  path,
			Files: files,
		}
	}
}

/* ---------- delete ---------- */

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

/* ---------- pull ---------- */

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
