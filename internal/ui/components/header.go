package components

import (
	"fmt"

	"github.com/SakshhamTheCoder/adbt/internal/state"
	"github.com/charmbracelet/lipgloss"
)

func RenderHeader(appState *state.AppState, screenName string) string {
	width := max(appState.Width, 40)

	appName := LogoStyle.Render("ADBT")

	statusStr := ""
	if device := appState.SelectedDevice(); device != nil {
		statusStr = StatusConnected.Render("●") + " " + StatusMuted.Render(device.DisplayName())
	} else {
		statusStr = StatusDisconnected.Render("●") + " " + StatusMuted.Render("No device")
	}

	screenTitle := ScreenTitleStyle.Render(screenName)
	separator := StatusMuted.Render(" │ ")

	content := appName + "  " + screenTitle + separator + statusStr

	return HeaderStyle.Width(width).Align(lipgloss.Left).Render(content)
}

func ShellTitle(appState *state.AppState, screenName string) string {
	title := fmt.Sprintf("ADBT  |  %s", screenName)
	if device := appState.SelectedDevice(); device != nil {
		return title + "  |  " + device.DisplayName()
	}
	return title + "  |  No device"
}
