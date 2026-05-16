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
	var b strings.Builder

	b.WriteString(RenderHeader(state, props.Title) + "\n")

	contentWidth := max(state.Width-4, 20)
	contentHeight := state.Height - 7
	if contentHeight < 10 {
		contentHeight = 10
	}

	staticLines := strings.Count(props.StaticContent, "\n")
	scrollableHeight := contentHeight - staticLines - 1
	if scrollableHeight < 5 {
		scrollableHeight = 5
	}

	vp := props.Viewport
	if vp == nil {
		temp := viewport.New(contentWidth, scrollableHeight)
		vp = &temp
	} else if vp.Width == 0 && vp.Height == 0 {
		*vp = viewport.New(contentWidth, scrollableHeight)
	}

	vp.Width = contentWidth
	vp.Height = scrollableHeight
	vp.SetContent(props.ScrollableContent)

	var combinedContent strings.Builder
	combinedContent.WriteString(props.StaticContent)
	combinedContent.WriteString(vp.View())

	b.WriteString(
		ContentStyle.
			Width(contentWidth).
			Height(contentHeight).
			Render(combinedContent.String()) + "\n",
	)

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

		b.WriteString(FooterStyle.Render(footerText))
	}

	return b.String()
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
