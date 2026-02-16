package components

import (
	"strings"

	"github.com/SakshhamTheCoder/adbt/internal/state"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LayoutWithScrollProps struct {
	Title             string
	StaticContent     string
	ScrollableContent string
	Footer            string
}

var viewports = make(map[string]*viewport.Model)
var viewportsReady = make(map[string]bool)

func RenderLayoutWithScrollableSection(state *state.AppState, props LayoutWithScrollProps) string {
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

	screenKey := props.Title
	vp, exists := viewports[screenKey]

	if !exists || !viewportsReady[screenKey] {
		newVp := viewport.New(contentWidth, scrollableHeight)
		viewports[screenKey] = &newVp
		vp = &newVp
		viewportsReady[screenKey] = true
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

func UpdateViewport(screenTitle string, msg tea.Msg) tea.Cmd {
	vp, exists := viewports[screenTitle]
	if !exists {
		return nil
	}

	var cmd tea.Cmd
	*vp, cmd = vp.Update(msg)
	return cmd
}

func ViewportGotoTop(screenTitle string) {
	if vp, exists := viewports[screenTitle]; exists {
		vp.GotoTop()
	}
}

func ViewportGotoBottom(screenTitle string) {
	if vp, exists := viewports[screenTitle]; exists {
		vp.GotoBottom()
	}
}

func RenderNoDevice(state *state.AppState, title string) string {
	return RenderLayoutWithScrollableSection(state, LayoutWithScrollProps{
		Title:             title,
		ScrollableContent: StatusDisconnected.Render("No device selected"),
		Footer:            Help("esc", "back"),
	})
}

func ViewportEnsureVisible(screenTitle string, line int) {
	vp, exists := viewports[screenTitle]
	if !exists {
		return
	}

	top := vp.YOffset
	bottom := vp.YOffset + vp.Height - 1

	if line < top {
		vp.YOffset = line
	} else if line > bottom {
		vp.YOffset = line - vp.Height + 1
	}
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
