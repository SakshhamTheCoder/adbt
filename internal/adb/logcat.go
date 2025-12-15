package adb

import (
	"bufio"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type LogcatSession struct {
	cmd     *exec.Cmd
	scanner *bufio.Scanner
}

type LogcatStartedMsg struct {
	Session *LogcatSession
}

type LogcatLineMsg struct {
	Line string
}

type LogcatErrorMsg struct {
	Error error
}

type LogcatStoppedMsg struct{}

func StartLogcatCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("adb", "-s", serial, "logcat")

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return LogcatErrorMsg{Error: err}
		}

		if err := cmd.Start(); err != nil {
			return LogcatErrorMsg{Error: err}
		}

		scanner := bufio.NewScanner(stdout)

		return LogcatStartedMsg{
			Session: &LogcatSession{
				cmd:     cmd,
				scanner: scanner,
			},
		}
	}
}

func NextLogcatLineCmd(s *LogcatSession) tea.Cmd {
	return func() tea.Msg {
		if s.scanner.Scan() {
			return LogcatLineMsg{Line: s.scanner.Text()}
		}

		if err := s.scanner.Err(); err != nil {
			return LogcatErrorMsg{Error: err}
		}

		_ = s.cmd.Wait()
		return LogcatStoppedMsg{}
	}
}

func StopLogcat(s *LogcatSession) {
	if s != nil && s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
}
