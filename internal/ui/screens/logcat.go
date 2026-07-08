package screens

import (
	"strings"
	"time"

	"github.com/SakshhamTheCoder/adbt/internal/adb"
	"github.com/SakshhamTheCoder/adbt/internal/state"
	"github.com/SakshhamTheCoder/adbt/internal/ui/components"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var logLevels = []string{"", "V", "D", "I", "W", "E", "F"}

type Logcat struct {
	state   *state.AppState
	lines   []string
	session *adb.LogcatSession
	running bool

	filterLevel int
	pidFilter   string
	search      components.SearchState
	viewport    viewport.Model

	filterForm components.FormModal
	toast      components.Toast
}

func NewLogcat(state *state.AppState) *Logcat {
	return &Logcat{
		state:    state,
		viewport: viewport.New(0, 0),
	}
}

func (l *Logcat) Init() tea.Cmd {
	if !l.state.HasDevice() {
		return nil
	}
	l.running = true
	return tea.Batch(tea.SetWindowTitle(components.ShellTitle(l.state, "Logcat")), adb.StartLogcatCmd(l.state.DeviceSerial(), l.pidFilter))
}

// CapturingText keeps "q" out of the global quit handler while typing a filter or search.
func (l *Logcat) CapturingText() bool {
	return l.search.Active || l.filterForm.Visible
}

func (l *Logcat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l.toast.Update(msg)

	if l.filterForm.Visible {
		switch m := msg.(type) {
		case components.FormSubmitMsg:
			l.filterForm.Hide()
			return l, l.applyFilterInput(m.Values)
		case components.FormCancelMsg:
			l.filterForm.Hide()
			return l, nil
		}
		return l, l.filterForm.Update(msg)
	}

	switch msg := msg.(type) {

	case adb.LogcatStartedMsg:
		l.session = msg.Session
		// Some terminals briefly switch the title to the child process name.
		return l, tea.Batch(tea.SetWindowTitle(components.ShellTitle(l.state, "Logcat")), adb.NextLogcatLineCmd(l.session))

	case adb.LogcatLineMsg:
		if msg.Session != l.session {
			return l, nil // stale line from a restarted stream
		}
		l.lines = append(l.lines, msg.Line)
		if len(l.lines) > 1000 {
			l.lines = l.lines[len(l.lines)-1000:]
		}
		if l.running {
			if !l.search.Active {
				l.gotoBottom()
			}
			return l, adb.NextLogcatLineCmd(l.session)
		}

	case adb.LogcatStoppedMsg:
		if msg.Session == l.session {
			l.running = false
		}

	case adb.LogcatErrorMsg:
		if msg.Session == l.session || msg.Session == nil {
			l.running = false
			var cmd tea.Cmd
			l.toast, cmd = components.ShowToast("Logcat error: "+msg.Error.Error(), true, 3*time.Second)
			return l, cmd
		}

	case adb.PidResolvedMsg:
		if msg.Error != nil || msg.Pid == "" {
			var cmd tea.Cmd
			l.toast, cmd = components.ShowToast(msg.Pkg+" is not running", true, 3*time.Second)
			return l, cmd
		}
		l.pidFilter = msg.Pid
		return l, l.restart()

	case adb.LogcatSavedMsg:
		var cmd tea.Cmd
		if msg.Error != nil {
			l.toast, cmd = components.ShowToast("Save failed: "+msg.Error.Error(), true, 3*time.Second)
		} else {
			l.toast, cmd = components.ShowToast("Saved to "+msg.Path, false, 3*time.Second)
		}
		return l, cmd

	case tea.KeyMsg:
		if l.search.Active {
			l.search.HandleKey(msg)
			return l, consumeKeyCmd()
		}

		switch msg.String() {
		case "c":
			l.lines = nil
			l.gotoTop()
		case "s":
			l.running = !l.running
			if l.running && l.session != nil {
				return l, adb.NextLogcatLineCmd(l.session)
			}
		case "right":
			l.filterLevel = (l.filterLevel + 1) % len(logLevels)
		case "left":
			l.filterLevel = (l.filterLevel + len(logLevels) - 1) % len(logLevels)
		case "/":
			l.search.Start()
		case "p":
			l.filterForm.Show("Filter by package or PID", []components.FormField{
				{Label: "Package or PID", Value: l.pidFilter},
			})
		case "w":
			return l, adb.SaveLogcatCmd(l.filteredLines())
		case "esc":
			if l.search.Query != "" {
				l.search.Clear()
				return l, consumeKeyCmd()
			}
		default:
			return l, l.updateViewport(msg)
		}
	}

	return l, nil
}

// applyFilterInput sets the pid filter from a package name or numeric PID and restarts.
func (l *Logcat) applyFilterInput(values []string) tea.Cmd {
	input := ""
	if len(values) > 0 {
		input = strings.TrimSpace(values[0])
	}

	if input == "" {
		l.pidFilter = ""
		return l.restart()
	}

	if isAllDigits(input) {
		l.pidFilter = input
		return l.restart()
	}

	return adb.ResolvePidCmd(l.state.DeviceSerial(), input)
}

// restart stops the current stream and starts a fresh one with the current filter.
func (l *Logcat) restart() tea.Cmd {
	old := l.session
	l.session = nil
	l.lines = nil
	l.running = true
	l.gotoTop()

	var stopCmd tea.Cmd
	if old != nil {
		stopCmd = func() tea.Msg {
			_ = old.Stop()
			return nil
		}
	}
	return tea.Batch(stopCmd, adb.StartLogcatCmd(l.state.DeviceSerial(), l.pidFilter))
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func (l *Logcat) View() string {
	if !l.state.HasDevice() {
		return components.RenderNoDevice(l.state, "Logcat")
	}

	filtered := l.filteredLines()

	maxWidth := max(l.state.Width-8, 20)
	truncStyle := lipgloss.NewStyle().MaxWidth(maxWidth)

	var body strings.Builder
	for _, line := range filtered {
		styled := colorLogLine(line)
		if l.search.Query != "" {
			styled = highlightSearch(styled, l.search.Query)
		}
		body.WriteString(truncStyle.Render(styled) + "\n")
	}

	var statusLine strings.Builder
	if l.running {
		statusLine.WriteString(components.StatusConnected.Render("● streaming"))
	} else {
		statusLine.WriteString(components.StatusMuted.Render("● paused"))
	}

	statusLine.WriteString("  ")
	for i, level := range logLevels {
		name := level
		if name == "" {
			name = "All"
		}

		if i == l.filterLevel {
			statusLine.WriteString(components.TabActiveStyle.Render(name))
		} else {
			statusLine.WriteString(components.TabInactiveStyle.Render(name))
		}

		if i < len(logLevels)-1 {
			statusLine.WriteString(" ")
		}
	}

	if l.search.Active {
		statusLine.WriteString("  ")
		statusLine.WriteString(components.HelpKeyStyle.Render("search: ") + l.search.Query + "▌")
	} else if l.search.Query != "" {
		statusLine.WriteString("  ")
		statusLine.WriteString(components.StatusMuted.Render("search: \"" + l.search.Query + "\""))
	}

	if l.pidFilter != "" {
		statusLine.WriteString("  ")
		statusLine.WriteString(components.HelpKeyStyle.Render("pid:" + l.pidFilter))
	}

	statusLine.WriteString("\n")

	rendered := components.RenderLayoutWithScrollableSection(l.state, components.LayoutWithScrollProps{
		Title:             "Logcat",
		StaticContent:     statusLine.String(),
		ScrollableContent: body.String(),
		Footer: components.JoinHelp(
			[2]string{"c", "clear"},
			[2]string{"s", "start/stop"},
			[2]string{"←/→", "level"},
			[2]string{"p", "pid filter"},
			[2]string{"/", "search"},
			[2]string{"w", "save"},
			[2]string{"esc", "back"},
		),
		Viewport: &l.viewport,
	})

	if l.filterForm.Visible {
		rendered = components.RenderFormOverlay(rendered, l.filterForm, l.state)
	}

	if l.toast.Visible {
		rendered = components.RenderOverlay(rendered, l.toast.View(), l.state)
	}

	return rendered
}

/* ---------- helpers ---------- */

func (l *Logcat) filteredLines() []string {
	minLevel := logLevels[l.filterLevel]
	if minLevel == "" && l.search.Query == "" {
		return l.lines
	}

	result := make([]string, 0, len(l.lines))
	for _, line := range l.lines {
		if minLevel != "" && !lineMatchesLevel(line, minLevel) {
			continue
		}
		if l.search.Query != "" && !strings.Contains(
			strings.ToLower(line),
			strings.ToLower(l.search.Query),
		) {
			continue
		}
		result = append(result, line)
	}
	return result
}

func lineMatchesLevel(line, minLevel string) bool {
	priority := extractPriority(line)
	if priority == "" {
		return true
	}
	return priorityRank(priority) >= priorityRank(minLevel)
}

func extractPriority(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return ""
	}

	tag := fields[4]
	if len(tag) >= 2 && tag[1] == '/' {
		return string(tag[0])
	}
	return ""
}

func priorityRank(level string) int {
	switch level {
	case "V":
		return 0
	case "D":
		return 1
	case "I":
		return 2
	case "W":
		return 3
	case "E":
		return 4
	case "F":
		return 5
	}
	return -1
}

func colorLogLine(line string) string {
	p := extractPriority(line)
	switch p {
	case "V":
		return components.LogVerbose.Render(line)
	case "D":
		return components.LogDebug.Render(line)
	case "I":
		return components.LogInfo.Render(line)
	case "W":
		return components.LogWarn.Render(line)
	case "E":
		return components.LogError.Render(line)
	case "F":
		return components.LogFatal.Render(line)
	}
	return line
}

func (l *Logcat) Cleanup() tea.Cmd {
	l.running = false
	session := l.session
	return func() tea.Msg {
		_ = session.Stop()
		return nil
	}
}

func (l *Logcat) updateViewport(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	l.viewport, cmd = l.viewport.Update(msg)
	return cmd
}

func (l *Logcat) gotoTop() {
	l.viewport.GotoTop()
}

func (l *Logcat) gotoBottom() {
	l.viewport.GotoBottom()
}

func highlightSearch(line, term string) string {
	lower := strings.ToLower(line)
	lowerTerm := strings.ToLower(term)

	idx := strings.Index(lower, lowerTerm)
	if idx == -1 {
		return line
	}

	before := line[:idx]
	match := line[idx : idx+len(term)]
	after := line[idx+len(term):]

	return before + components.WarningStyle.Render(match) + after
}
