package screens

import "github.com/charmbracelet/bubbles/viewport"

func ensureViewportLineVisible(vp *viewport.Model, line int) {
	top := vp.YOffset
	bottom := vp.YOffset + vp.Height - 1

	if line < top {
		vp.YOffset = line
	} else if line > bottom {
		vp.YOffset = line - vp.Height + 1
	}
}
