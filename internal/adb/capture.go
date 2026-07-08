package adb

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type CaptureResultMsg struct {
	Action string // "screenshot" or "recording"
	Path   string
	Error  error
}

type ScreenRecordStartedMsg struct {
	Session *ScreenRecordSession
	Error   error
}

type ScreenRecordSession struct {
	cmd        *exec.Cmd
	serial     string
	remotePath string
}

// ScreenshotCmd grabs a screenshot via `exec-out screencap -p`, which streams the
// raw PNG to stdout — avoiding the CRLF corruption shell redirection causes on Windows.
func ScreenshotCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("adb", "-s", serial, "exec-out", "screencap", "-p").Output()
		if err != nil {
			return CaptureResultMsg{Action: "screenshot", Error: err}
		}

		local := filepath.Join(DefaultSaveDir(), TimestampedName("screenshot", ".png"))
		if err := os.WriteFile(local, out, 0o644); err != nil {
			return CaptureResultMsg{Action: "screenshot", Error: err}
		}
		return CaptureResultMsg{Action: "screenshot", Path: local}
	}
}

// StartScreenRecordCmd launches screenrecord on the device, writing to /sdcard.
// The process is long-running and held in the returned session until stopped.
func StartScreenRecordCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		remote := "/sdcard/" + TimestampedName("adbt-rec", ".mp4")
		cmd := exec.Command("adb", "-s", serial, "shell", "screenrecord", remote)
		if err := cmd.Start(); err != nil {
			return ScreenRecordStartedMsg{Error: err}
		}
		return ScreenRecordStartedMsg{
			Session: &ScreenRecordSession{
				cmd:        cmd,
				serial:     serial,
				remotePath: remote,
			},
		}
	}
}

// StopScreenRecordCmd interrupts the recording so the file is finalized, pulls it
// to the host, then removes the on-device copy.
func StopScreenRecordCmd(s *ScreenRecordSession) tea.Cmd {
	return func() tea.Msg {
		if s == nil {
			return CaptureResultMsg{Action: "recording", Error: fmt.Errorf("no active recording")}
		}

		// SIGINT lets screenrecord write the final moov atom; killing adb alone would not.
		_, _ = ExecuteCommand(s.serial, "shell", "pkill", "-INT", "screenrecord")
		_ = s.cmd.Wait()
		time.Sleep(time.Second) // give the device a moment to flush the file

		local := filepath.Join(DefaultSaveDir(), path.Base(s.remotePath))
		if _, err := ExecuteCommand(s.serial, "pull", s.remotePath, local); err != nil {
			return CaptureResultMsg{Action: "recording", Error: err}
		}
		_, _ = ExecuteCommand(s.serial, "shell", "rm", "-f", s.remotePath)

		return CaptureResultMsg{Action: "recording", Path: local}
	}
}
