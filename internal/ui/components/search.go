package components

import (
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
)

type SearchState struct {
	Active bool
	Query  string
}

func (s *SearchState) Start() {
	s.Active = true
	s.Query = ""
}

func (s *SearchState) Clear() {
	s.Active = false
	s.Query = ""
}

func (s *SearchState) HandleKey(msg tea.KeyMsg) bool {
	if !s.Active {
		return false
	}

	switch msg.String() {
	case "esc":
		s.Clear()
	case "enter":
		s.Active = false
	case "backspace":
		if len(s.Query) > 0 {
			s.Query = s.Query[:len(s.Query)-1]
		}
	default:
		if len(msg.Runes) == 1 && unicode.IsPrint(msg.Runes[0]) {
			s.Query += string(msg.Runes)
		}
	}

	return true
}
