package components

import "github.com/charmbracelet/lipgloss"

var (
	Primary = lipgloss.Color("#00D7FF")
	Success = lipgloss.Color("#10B981")
	Warning = lipgloss.Color("#F59E0B")
	Error   = lipgloss.Color("#EF4444")
	Muted   = lipgloss.Color("#6B7280")
	Border  = lipgloss.Color("#374151")

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(0, 2).
			Align(lipgloss.Center)

	ContentStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary)

	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	ListItemSelectedStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true).
				PaddingLeft(1)

	StatusConnected = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	StatusDisconnected = lipgloss.NewStyle().
				Foreground(Error).
				Bold(true)

	StatusMuted = lipgloss.NewStyle().
			Foreground(Muted)

	FooterStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Padding(0, 2)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(Muted)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	LogVerbose = lipgloss.NewStyle().Foreground(Muted)
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
