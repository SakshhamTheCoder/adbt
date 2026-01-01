package components

import (
	"strings"

	"adbt/internal/state"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type LayoutProps struct {
	Title          string
	Body           string
	Footer         string
	EnableViewport bool
}

type LayoutWithScrollProps struct {
	Title             string
	StaticContent     string
	ScrollableContent string
	Footer            string
}

var viewports = make(map[string]*viewport.Model)
var viewportsReady = make(map[string]bool)

func RenderLayout(state *state.AppState, props LayoutProps) string {
	var b strings.Builder

	b.WriteString(RenderHeader(state, props.Title) + "\n")

	contentWidth := max(state.Width-4, 20)
	contentHeight := state.Height - 7
	if contentHeight < 10 {
		contentHeight = 10
	}

	var bodyContent string

	if props.EnableViewport {

		screenKey := props.Title
		vp, exists := viewports[screenKey]

		if !exists || !viewportsReady[screenKey] {
			newVp := viewport.New(contentWidth, contentHeight)
			viewports[screenKey] = &newVp
			vp = &newVp
			viewportsReady[screenKey] = true
		}

		vp.Width = contentWidth
		vp.Height = contentHeight
		vp.SetContent(props.Body)
		bodyContent = vp.View()
	} else {
		bodyContent = props.Body
	}

	b.WriteString(
		ContentStyle.
			Width(contentWidth).
			Height(contentHeight).
			Render(bodyContent) + "\n",
	)

	if props.Footer != "" {
		footerText := props.Footer

		if props.EnableViewport {
			screenKey := props.Title
			if vp, exists := viewports[screenKey]; exists {
				if vp.TotalLineCount() > vp.Height {
					percentage := int(vp.ScrollPercent() * 100)
					scrollInfo := StatusMuted.Render(" │ ") +
						StatusMuted.Render(string(rune('0'+percentage/10))) +
						StatusMuted.Render(string(rune('0'+percentage%10))) +
						StatusMuted.Render("%")
					footerText += scrollInfo
				}
			}
		}

		b.WriteString(FooterStyle.Render(footerText))
	}

	return b.String()
}

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
	return RenderLayout(state, LayoutProps{
		Title:  title,
		Body:   StatusDisconnected.Render("No device selected"),
		Footer: Help("esc", "back"),
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
