package components

import (
	"strings"

	"adbt/internal/state"
)

type LayoutProps struct {
	Title  string
	Body   string
	Footer string
}

func RenderLayout(state *state.AppState, props LayoutProps) string {
	var b strings.Builder

	b.WriteString(RenderHeader(state, props.Title) + "\n")

	contentWidth := max(state.Width-4, 20)

	b.WriteString(
		ContentStyle.
			Width(contentWidth).
			Render(props.Body) + "\n",
	)

	if props.Footer != "" {
		b.WriteString(FooterStyle.Render(props.Footer))
	}

	return b.String()
}
