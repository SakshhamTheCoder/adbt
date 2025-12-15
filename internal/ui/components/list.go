package components

import (
	"fmt"
	"strings"

	"adbt/internal/adb"
)

func DeviceList(devices []adb.Device, cursor int) string {
	if len(devices) == 0 {
		return ErrorStyle.Render("No devices found. Please connect a device with USB debugging enabled.")
	}

	var b strings.Builder
	b.WriteString(TitleStyle.Render("Connected Devices") + "\n\n")

	for i, device := range devices {
		cursorChar := " "
		style := ListItemStyle

		if i == cursor {
			cursorChar = "›"
			style = ListItemSelectedStyle
		}

		status := StatusDisconnected.Render("●")
		if device.IsConnected() {
			status = StatusConnected.Render("●")
		}

		line := fmt.Sprintf("%s %s %s", cursorChar, status, device.DisplayName())

		if device.IsConnected() && device.Android != "" {
			line += StatusMuted.Render(fmt.Sprintf(" - Android %s", device.Android))
		} else {
			line += StatusMuted.Render(fmt.Sprintf(" - %s", device.State))
		}

		b.WriteString(style.Render(line) + "\n")
	}

	return b.String()
}
