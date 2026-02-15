package screens

import (
	"strings"
	"unicode"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

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
	search      string
	searching   bool
}

func NewLogcat(state *state.AppState) *Logcat {
	return &Logcat{state: state}
}

func (l *Logcat) Init() tea.Cmd {
	if !l.state.HasDevice() {
		return nil
	}
	l.running = true
	return adb.StartLogcatCmd(l.state.DeviceSerial())
}

func (l *Logcat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case adb.LogcatStartedMsg:
		l.session = msg.Session
		return l, adb.NextLogcatLineCmd(l.session)

	case adb.LogcatLineMsg:
		l.lines = append(l.lines, msg.Line)
		if len(l.lines) > 1000 {
			l.lines = l.lines[len(l.lines)-1000:]
		}
		if l.running {
			if !l.searching {
				components.ViewportGotoBottom("Logcat")
			}
			return l, adb.NextLogcatLineCmd(l.session)
		}

	case adb.LogcatStoppedMsg:
		l.running = false

	case tea.KeyMsg:
		if l.searching {
			switch msg.String() {
			case "esc":
				l.searching = false
				l.search = ""
				return l, tea.Batch()
			case "enter":
				l.searching = false
			case "backspace":
				if len(l.search) > 0 {
					l.search = l.search[:len(l.search)-1]
				}
			default:
				if len(msg.Runes) == 1 && unicode.IsPrint(msg.Runes[0]) {
					l.search += string(msg.Runes)
				}
			}
			return l, nil
		}

		switch msg.String() {
		case "c":
			l.lines = nil
			components.ViewportGotoTop("Logcat")
		case "s":
			l.running = !l.running
			if l.running {
				return l, adb.NextLogcatLineCmd(l.session)
			}
		case "f":
			l.filterLevel = (l.filterLevel + 1) % len(logLevels)
		case "/":
			l.searching = true
			l.search = ""
		default:
			return l, components.UpdateViewport("Logcat", msg)
		}
	}

	return l, nil
}

func (l *Logcat) View() string {
	if !l.state.HasDevice() {
		return components.RenderNoDevice(l.state, "Logcat")
	}

	filtered := l.filteredLines()

	maxWidth := l.state.Width - 8
	if maxWidth < 20 {
		maxWidth = 20
	}
	truncStyle := lipgloss.NewStyle().MaxWidth(maxWidth)

	var body strings.Builder
	for _, line := range filtered {
		styled := colorLogLine(line)
		if l.search != "" {
			styled = highlightSearch(styled, l.search)
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
			statusLine.WriteString(components.HelpKeyStyle.Render(name))
		} else {
			statusLine.WriteString(components.StatusMuted.Render(name))
		}

		if i < len(logLevels)-1 {
			statusLine.WriteString(components.StatusMuted.Render(" / "))
		}
	}

	if l.searching {
		statusLine.WriteString("  ")
		statusLine.WriteString(components.HelpKeyStyle.Render("search: ") + l.search + "▌")
	} else if l.search != "" {
		statusLine.WriteString("  ")
		statusLine.WriteString(components.StatusMuted.Render("search: \"" + l.search + "\""))
	}

	statusLine.WriteString("\n")

	return components.RenderLayoutWithScrollableSection(l.state, components.LayoutWithScrollProps{
		Title:             "Logcat",
		StaticContent:     statusLine.String(),
		ScrollableContent: body.String(),
		Footer: components.Help("c", "clear") + "  " +
			components.Help("s", "start/stop") + "  " +
			components.Help("f", "filter") + "  " +
			components.Help("/", "search") + "  " +
			components.Help("esc", "back"),
	})
}

/* ---------- helpers ---------- */

func (l *Logcat) filteredLines() []string {
	minLevel := logLevels[l.filterLevel]
	if minLevel == "" && l.search == "" {
		return l.lines
	}

	result := make([]string, 0, len(l.lines))
	for _, line := range l.lines {
		if minLevel != "" && !lineMatchesLevel(line, minLevel) {
			continue
		}
		if l.search != "" && !strings.Contains(
			strings.ToLower(line),
			strings.ToLower(l.search),
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
