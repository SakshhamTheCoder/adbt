package components

import (
	"strings"

	"github.com/SakshhamTheCoder/adbt/internal/state"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type ScrollableLayoutProps struct {
	Title             string
	StaticContent     string
	ScrollableContent string
	Footer            string
	Viewport          *viewport.Model
}

type LayoutWithScrollProps = ScrollableLayoutProps

func RenderLayoutWithScrollableSection(state *state.AppState, props ScrollableLayoutProps) string {
	headerStr := RenderHeader(state, props.Title)
	headerHeight := lipgloss.Height(headerStr)

	footerStr := ""
	footerHeight := 0
	if props.Footer != "" {
		// Render footer without scroll info first to get its wrapped height
		tempFooter := FooterStyle.Width(state.Width).Render(props.Footer)
		footerHeight = lipgloss.Height(tempFooter)
	}

	contentHeight := max(state.Height-headerHeight-footerHeight, 5)

	staticLines := strings.Count(props.StaticContent, "\n")
	scrollableHeight := max(contentHeight-staticLines, 3)

	vp := props.Viewport
	if vp == nil {
		temp := viewport.New(state.Width-4, scrollableHeight)
		vp = &temp
	} else if vp.Width == 0 && vp.Height == 0 {
		*vp = viewport.New(state.Width-4, scrollableHeight)
	}

	vp.Width = state.Width - 4
	vp.Height = scrollableHeight
	vp.SetContent(props.ScrollableContent)

	// Re-render footer with scroll info if needed
	if props.Footer != "" {
		footerText := props.Footer
		if vp.TotalLineCount() > vp.Height {
			percentage := int(vp.ScrollPercent() * 100)
			scrollInfo := StatusMuted.Render(" │ ") +
				StatusMuted.Render(string(rune('0'+percentage/10))) +
				StatusMuted.Render(string(rune('0'+percentage%10))) +
				StatusMuted.Render("%")
			footerText += scrollInfo
		}
		footerStr = FooterStyle.Width(state.Width).Render(footerText)
	}

	contentStr := ContentStyle.
		Width(state.Width).
		Height(contentHeight).
		Render(props.StaticContent + vp.View())

	if footerStr == "" {
		return lipgloss.JoinVertical(lipgloss.Left, headerStr, contentStr)
	}
	return lipgloss.JoinVertical(lipgloss.Left, headerStr, contentStr, footerStr)
}

func RenderLayout(state *state.AppState, title, content, footer string) string {
	return RenderLayoutWithScrollableSection(state, LayoutWithScrollProps{
		Title:             title,
		ScrollableContent: content,
		Footer:            footer,
	})
}

func RenderNoDevice(state *state.AppState, title string) string {
	return RenderLayoutWithScrollableSection(state, LayoutWithScrollProps{
		Title:             title,
		ScrollableContent: StatusDisconnected.Render("No device selected"),
		Footer:            Help("esc", "back"),
	})
}

func RenderOverlay(base, overlay string, s *state.AppState) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")

	overlayW := lipgloss.Width(overlay)
	overlayH := len(overlayLines)

	x := (s.Width - overlayW) / 2
	y := (s.Height - overlayH) / 2

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	leftPad := strings.Repeat(" ", x)

	for i, oLine := range overlayLines {
		row := y + i
		if row < 0 || row >= len(baseLines) {
			continue
		}

		baseLines[row] = leftPad + oLine
	}

	return strings.Join(baseLines, "\n")
}

func RenderFormOverlay(base string, form FormModal, s *state.AppState) string {
	formView := form.View()
	rendered := RenderOverlay(base, formView, s)
	if !form.PickerVisible() {
		return rendered
	}

	picker := form.PickerView()
	return RenderOverlay(rendered, picker, s)
}

func RenderOverlayAt(base, overlay string, x, y int) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	for i, oLine := range overlayLines {
		row := y + i
		if row < 0 || row >= len(baseLines) {
			continue
		}

		prefix := strings.Repeat(" ", x)
		baseLines[row] = prefix + oLine
	}

	return strings.Join(baseLines, "\n")
}
