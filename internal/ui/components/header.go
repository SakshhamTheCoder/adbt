package components

import (
	"fmt"

	"github.com/SakshhamTheCoder/adbt/internal/state"
)

func RenderHeader(appState *state.AppState, screenName string) string {
	title := fmt.Sprintf("ADBT  |  %s", screenName)

	if appState.HasDevice() {
		device := appState.SelectedDevice
		title += StatusConnected.Render("  ● ")
		title += StatusMuted.Render(device.DisplayName())
	} else {
		title += StatusDisconnected.Render("  ● ")
		title += StatusMuted.Render("No device")
	}

	width := appState.Width - 4
	if width < 20 {
		width = 20
	}

	return HeaderStyle.Width(width).Render(title)
}
