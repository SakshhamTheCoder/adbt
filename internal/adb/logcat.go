package adb

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

type LogcatSession struct {
	cmd     *exec.Cmd
	scanner *bufio.Scanner
	mu      sync.Mutex
	stopped bool
}

type LogcatStartedMsg struct {
	Session *LogcatSession
}

// Line/Stopped/Error messages carry the session they came from so the UI can
// ignore stale output after the stream is restarted with a new filter.
type LogcatLineMsg struct {
	Session *LogcatSession
	Line    string
}

type LogcatErrorMsg struct {
	Session *LogcatSession
	Error   error
}

type LogcatStoppedMsg struct {
	Session *LogcatSession
}

type PidResolvedMsg struct {
	Pkg   string
	Pid   string
	Error error
}

type LogcatSavedMsg struct {
	Path  string
	Error error
}

// StartLogcatCmd starts a logcat stream. When pid is non-empty the stream is
// limited to that process via --pid.
func StartLogcatCmd(serial, pid string) tea.Cmd {
	return func() tea.Msg {
		args := []string{"-s", serial, "logcat"}
		if pid != "" {
			args = append(args, "--pid="+pid)
		}
		cmd := exec.Command("adb", args...)

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
			return LogcatLineMsg{Session: s, Line: s.scanner.Text()}
		}

		if err := s.scanner.Err(); err != nil {
			return LogcatErrorMsg{Session: s, Error: err}
		}

		_ = s.Stop()
		return LogcatStoppedMsg{Session: s}
	}
}

// ResolvePidCmd resolves a package name to its (single) running PID.
func ResolvePidCmd(serial, pkg string) tea.Cmd {
	return func() tea.Msg {
		out, err := ExecuteCommand(serial, "shell", "pidof", "-s", pkg)
		if err != nil {
			return PidResolvedMsg{Pkg: pkg, Error: err}
		}

		pid := strings.TrimSpace(string(out))
		if fields := strings.Fields(pid); len(fields) > 0 {
			pid = fields[0]
		}
		return PidResolvedMsg{Pkg: pkg, Pid: pid}
	}
}

// SaveLogcatCmd writes the given lines to a timestamped file in the default save dir.
func SaveLogcatCmd(lines []string) tea.Cmd {
	return func() tea.Msg {
		path := filepath.Join(DefaultSaveDir(), TimestampedName("logcat", ".txt"))
		data := strings.Join(lines, "\n")
		if len(lines) > 0 {
			data += "\n"
		}
		if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
			return LogcatSavedMsg{Error: err}
		}
		return LogcatSavedMsg{Path: path}
	}
}

func (s *LogcatSession) Stop() error {
	if s == nil || s.cmd == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return nil
	}
	s.stopped = true

	if s.cmd.Process == nil {
		return nil
	}

	if s.cmd.ProcessState == nil || !s.cmd.ProcessState.Exited() {
		_ = s.cmd.Process.Kill()
	}

	return s.cmd.Wait()
}
