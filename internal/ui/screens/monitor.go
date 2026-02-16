package screens

import (
	"fmt"
	"strings"
	"time"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PerfMonitor struct {
	state *state.AppState

	// Stats
	lastStats    adb.SystemStats
	currentStats adb.SystemStats
	hasHistory   bool

	// Calculated values
	cpuPercent float64
	rxRate     uint64 // bytes per second
	txRate     uint64 // bytes per second

	// Ticker
	sub chan struct{}
}

type TickMsg time.Time

func NewPerfMonitor(state *state.AppState) *PerfMonitor {
	return &PerfMonitor{
		state: state,
		sub:   make(chan struct{}),
	}
}

func (m *PerfMonitor) Init() tea.Cmd {
	return tea.Batch(
		adb.GetSystemStatsCmd(m.state.DeviceSerial()),
		m.tickCmd(),
	)
}

func (m *PerfMonitor) tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *PerfMonitor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, nil // handled by parent or just stop? Parent handles navigation.
		}

	case TickMsg:
		if !m.state.HasDevice() {
			return m, m.tickCmd()
		}
		return m, tea.Batch(
			adb.GetSystemStatsCmd(m.state.DeviceSerial()),
			m.tickCmd(),
		)

	case adb.SystemStatsMsg:
		if msg.Error != nil {
			// Handle error silently or show toast? For now silent catch up next tick
			return m, nil
		}

		newStats := msg.Stats

		if m.hasHistory {
			// Calculate CPU
			deltaTotal := newStats.CPUTotal - m.currentStats.CPUTotal
			deltaIdle := newStats.CPUIdle - m.currentStats.CPUIdle

			if deltaTotal > 0 {
				used := deltaTotal - deltaIdle
				m.cpuPercent = float64(used) / float64(deltaTotal) * 100
			}

			// Calculate Net Rate (assuming 1 second tick)
			// Handle wrap around roughly
			if newStats.NetRxBytes >= m.currentStats.NetRxBytes {
				m.rxRate = newStats.NetRxBytes - m.currentStats.NetRxBytes
			}
			if newStats.NetTxBytes >= m.currentStats.NetTxBytes {
				m.txRate = newStats.NetTxBytes - m.currentStats.NetTxBytes
			}
		}

		m.lastStats = m.currentStats
		m.currentStats = newStats
		m.hasHistory = true
	}

	return m, nil
}

func (m *PerfMonitor) View() string {
	if !m.state.HasDevice() {
		return components.RenderNoDevice(m.state, "Performance Monitor")
	}

	// Styles
	// Styles
	labelStyle := components.StatusMuted.Copy().Width(12)
	barStyle := lipgloss.NewStyle().Background(components.Primary)
	barEmptyStyle := lipgloss.NewStyle().Background(lipgloss.Color("#1f2937")) // dark gray (Tailwind gray-800)

	// 1. CPU
	cpuBar := renderProgressBar(m.cpuPercent, 40, barStyle, barEmptyStyle)
	cpuRow := lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render("CPU Use"),
		cpuBar,
		fmt.Sprintf(" %.1f%%", m.cpuPercent),
	)

	// 2. Memory
	memPercent := 0.0
	memLabel := "0 / 0 MB"
	if m.currentStats.MemTotal > 0 {
		memPercent = float64(m.currentStats.MemUsed) / float64(m.currentStats.MemTotal) * 100
		memLabel = fmt.Sprintf("%s / %s",
			adb.FormatFileSize(fmt.Sprintf("%d", m.currentStats.MemUsed*1024)),
			adb.FormatFileSize(fmt.Sprintf("%d", m.currentStats.MemTotal*1024)),
		)
	}
	memBar := renderProgressBar(memPercent, 40, barStyle, barEmptyStyle)
	memRow := lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render("Memory"),
		memBar,
		fmt.Sprintf(" %.1f%% (%s)", memPercent, memLabel),
	)

	// 3. Network
	rxStr := adb.FormatFileSize(fmt.Sprintf("%d", m.rxRate)) + "/s"
	txStr := adb.FormatFileSize(fmt.Sprintf("%d", m.txRate)) + "/s"
	netRow := lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render("Network"),
		fmt.Sprintf("↓ %s   ↑ %s", rxStr, txStr),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		"",
		cpuRow,
		"",
		memRow,
		"",
		netRow,
	)

	return components.RenderLayout(
		m.state,
		"Performance Monitor",
		content,
		components.Help("esc", "back"),
	)
}

func renderProgressBar(percent float64, width int, filled, empty lipgloss.Style) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	fillCount := int(float64(width) * percent / 100.0)
	emptyCount := width - fillCount

	bar := ""
	if fillCount > 0 {
		bar += filled.Render(strings.Repeat(" ", fillCount))
	}
	if emptyCount > 0 {
		bar += empty.Render(strings.Repeat(" ", emptyCount))
	}
	return bar
}
