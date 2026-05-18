package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	Primary = lipgloss.Color("#706bff") // ANSI 69 (Indigo/Purplish Blue)
	Success = lipgloss.Color("#22c55e") // Green
	Warning = lipgloss.Color("#eab308") // Yellow
	Error   = lipgloss.Color("#ef4444") // Red

	FgHigh      = lipgloss.Color("#f3f4f6") // Light Gray
	FgMuted     = lipgloss.Color("#9ca3af") // Mid Gray
	BgDim       = lipgloss.Color("#1f2937") // Dark Gray
	BorderColor = lipgloss.Color("#374151") // Border Gray

	HeaderStyle = lipgloss.NewStyle().
			Foreground(FgHigh).
			Border(lipgloss.NormalBorder(), false, false, true, false). // Bottom border
			BorderForeground(BorderColor).
			Padding(0, 1)

	ContentStyle = lipgloss.NewStyle().
			Foreground(FgHigh).
			Padding(0, 2)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary)

	ListItemStyle = lipgloss.NewStyle().
			Foreground(FgHigh).
			PaddingLeft(2)

	ListItemSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#111827")). // Dark Gray
				Background(Primary).
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1)

	TabActiveStyle = lipgloss.NewStyle().
			Foreground(FgHigh).
			Background(BorderColor).
			Bold(true).
			Padding(0, 1)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(FgMuted).
				Padding(0, 1)

	LogoStyle = lipgloss.NewStyle().
			Background(Primary).
			Foreground(lipgloss.Color("#111827")).
			Bold(true).
			Padding(0, 1)

	ScreenTitleStyle = lipgloss.NewStyle().
				Foreground(FgHigh).
				Bold(true)

	FormBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(0, 2)

	FormPickerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)

	FormSubmitFocusedStyle = lipgloss.NewStyle().
				Background(Primary).
				Foreground(lipgloss.Color("#111827")).
				Bold(true)

	FormSubmitBlurredStyle = lipgloss.NewStyle().
				Background(BgDim).
				Foreground(FgMuted)

	PerfBarStyle      = lipgloss.NewStyle().Foreground(Primary)
	PerfBarEmptyStyle = lipgloss.NewStyle().Foreground(BorderColor)

	StatusConnected = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	StatusDisconnected = lipgloss.NewStyle().
				Foreground(Error).
				Bold(true)

	StatusMuted = lipgloss.NewStyle().
			Foreground(FgMuted)

	DimensionGuardStyle = lipgloss.NewStyle().
				Align(lipgloss.Center, lipgloss.Center).
				Foreground(FgMuted)

	FooterStyle = lipgloss.NewStyle().
			Foreground(FgMuted).
			Border(lipgloss.NormalBorder(), true, false, false, false). // Top border
			BorderForeground(BorderColor).
			Padding(0, 1)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(FgMuted)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	LogVerbose = lipgloss.NewStyle().Foreground(FgMuted)
	LogDebug   = lipgloss.NewStyle().Foreground(Primary)
	LogInfo    = lipgloss.NewStyle().Foreground(Success)
	LogWarn    = lipgloss.NewStyle().Foreground(Warning)
	LogError   = lipgloss.NewStyle().Foreground(Error)
	LogFatal   = lipgloss.NewStyle().Foreground(Error).Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)
)

func Help(key, desc string) string {
	return HelpKeyStyle.Render(key) + " " + HelpDescStyle.Render(desc)
}

func JoinHelp(items ...[2]string) string {
	var parts []string
	for _, item := range items {
		parts = append(parts, Help(item[0], item[1]))
	}
	return strings.Join(parts, "  ")
}

func SectionTitle(title string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(Primary).Render("■ " + title)
}
