package screens

import (
	"path/filepath"
	"time"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type Files struct {
	state *state.AppState

	path   string
	files  []adb.FileEntry
	cursor int

	confirm components.ConfirmPrompt
	toast   components.Toast
}

func NewFiles(state *state.AppState) *Files {
	return &Files{
		state: state,
		path:  "/sdcard",
	}
}

func (f *Files) Init() tea.Cmd {
	if !f.state.HasDevice() {
		return nil
	}
	return adb.ListFilesCmd(f.state.DeviceSerial(), f.path)
}

func (f *Files) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	f.toast.Update(msg)

	if f.confirm.Visible {
		switch msg.(type) {

		case components.ConfirmYesMsg:
			entry := f.files[f.cursor]
			f.confirm.Hide()
			return f, adb.DeleteFileCmd(
				f.state.DeviceSerial(),
				entry.Path,
			)

		case components.ConfirmNoMsg:
			f.confirm.Hide()
			return f, nil
		}

		return f, f.confirm.Update(msg)
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if f.cursor > 0 {
				f.cursor--
				components.ViewportEnsureVisible("Files", f.cursor)
			}

		case "down", "j":
			if f.cursor < len(f.files)-1 {
				f.cursor++
				components.ViewportEnsureVisible("Files", f.cursor)
			}

		case "enter":
			if len(f.files) == 0 {
				return f, nil
			}

			entry := f.files[f.cursor]
			if entry.IsDir {
				f.path = entry.Path
				f.cursor = 0
				components.ViewportGotoTop("Files")
				return f, adb.ListFilesCmd(
					f.state.DeviceSerial(),
					f.path,
				)
			}

		case "backspace":
			parent := filepath.Dir(f.path)
			if parent != f.path {
				f.path = parent
				f.cursor = 0
				components.ViewportGotoTop("Files")
				return f, adb.ListFilesCmd(
					f.state.DeviceSerial(),
					f.path,
				)
			}

		case "d":
			if len(f.files) == 0 {
				return f, nil
			}
			entry := f.files[f.cursor]
			f.confirm.Show("Delete file:\n" + entry.Name)

		case "r":
			f.cursor = 0
			components.ViewportGotoTop("Files")
			return f, adb.ListFilesCmd(
				f.state.DeviceSerial(),
				f.path,
			)
		}

	case adb.FilesLoadedMsg:
		if msg.Error != nil {
			var cmd tea.Cmd
			f.toast, cmd = components.ShowToast(
				"Failed to load files",
				true,
				3*time.Second,
			)
			return f, cmd
		}

		f.files = msg.Files
		if f.cursor >= len(f.files) {
			f.cursor = 0
		}
		components.ViewportGotoTop("Files")

	case adb.FileActionResultMsg:
		if msg.Error != nil {
			var cmd tea.Cmd
			f.toast, cmd = components.ShowToast(
				msg.Action+" failed",
				true,
				3*time.Second,
			)
			return f, cmd
		}

		var cmd tea.Cmd
		f.toast, cmd = components.ShowToast(
			msg.Action+" successful",
			false,
			2*time.Second,
		)
		return f, tea.Batch(
			cmd,
			adb.ListFilesCmd(f.state.DeviceSerial(), f.path),
		)
	}

	return f, nil
}

func (f *Files) View() string {
	if !f.state.HasDevice() {
		return components.RenderNoDevice(f.state, "Files")
	}

	body := components.FileList(f.files, f.cursor)

	if f.confirm.Visible {
		body += "\n\n" + f.confirm.View()
	}

	if f.toast.Visible {
		body += "\n\n" + f.toast.View()
	}

	return components.RenderLayoutWithScrollableSection(
		f.state,
		components.LayoutWithScrollProps{
			Title:             "Files",
			StaticContent:     components.StatusMuted.Render("Path: "+f.path) + "\n",
			ScrollableContent: body,
			Footer: components.Help("enter", "open") + "  " +
				components.Help("backspace", "up") + "  " +
				components.Help("d", "delete") + "  " +
				components.Help("r", "refresh") + "  " +
				components.Help("esc", "back"),
		},
	)
}
