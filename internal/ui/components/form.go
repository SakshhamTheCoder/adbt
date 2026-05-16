package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type FormSubmitMsg struct {
	Values []string
}

type FormCancelMsg struct{}

type FormField struct {
	Label       string
	Value       string
	Type        FormFieldType
	Placeholder string
	Options     []string
	Suggestions []string
}

type FormFieldType int

const (
	FormFieldText FormFieldType = iota
	FormFieldSelect
	FormFieldAutocomplete
)

type FormModal struct {
	Visible         bool
	Title           string
	Fields          []FormField
	Cursor          int
	scrollOffset    int
	inputs          []textinput.Model
	optionCursors   []int
	suggestionOpen  bool
	suggestionIndex int
}

const (
	formContentWidth  = 68
	formPickerWidth   = 56
	formInputWidth    = 34
	formVisibleFields = 6
)

func (f *FormModal) Show(title string, fields []FormField) {
	f.Visible = true
	f.Title = title
	f.Fields = fields
	f.Cursor = 0
	f.scrollOffset = 0
	f.inputs = make([]textinput.Model, len(fields))
	f.optionCursors = make([]int, len(fields))
	f.suggestionOpen = false
	f.suggestionIndex = 0

	for idx, field := range fields {
		if field.Type == FormFieldSelect && len(field.Options) > 0 {
			f.optionCursors[idx] = selectedOptionIndex(field.Options, field.Value)
			if f.optionCursors[idx] < 0 {
				f.optionCursors[idx] = 0
			}
			continue
		}

		input := textinput.New()
		input.SetValue(field.Value)
		input.Placeholder = field.Placeholder
		input.Prompt = ""
		input.TextStyle = HelpKeyStyle
		input.PlaceholderStyle = StatusMuted
		input.Cursor.Style = HelpKeyStyle
		input.CharLimit = 0
		input.Width = formInputWidth
		if idx == 0 {
			input.Focus()
		}
		f.inputs[idx] = input
	}
	f.focusCursor()
	f.refreshSuggestions()
}

func (f *FormModal) Hide() {
	*f = FormModal{}
}

func (f *FormModal) Update(msg tea.Msg) tea.Cmd {
	if !f.Visible || len(f.Fields) == 0 {
		return nil
	}

	if f.Cursor < 0 {
		f.Cursor = 0
	}
	if f.Cursor > f.submitIndex() {
		f.Cursor = f.submitIndex()
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	switch key.String() {

	case "tab":
		f.suggestionOpen = false
		f.moveCursor(1)

	case "shift+tab":
		f.suggestionOpen = false
		f.moveCursor(-1)

	case "down":
		if f.suggestionOpen {
			f.moveChoice(1)
		} else {
			f.moveCursor(1)
		}

	case "up":
		if f.suggestionOpen {
			f.moveChoice(-1)
		} else {
			f.moveCursor(-1)
		}

	case "right":
		return f.updateCurrentInput(msg)

	case "left":
		return f.updateCurrentInput(msg)

	case "enter":
		if f.suggestionOpen {
			f.acceptActiveChoice()
			f.moveCursor(1)
			return nil
		}
		if f.isSubmitFocused() {
			return f.submitCmd()
		}
		if f.openChoices() {
			return nil
		}
		f.moveCursor(1)

	case "esc":
		if f.suggestionOpen {
			f.suggestionOpen = false
			return func() tea.Msg { return nil }
		}
		return func() tea.Msg {
			return FormCancelMsg{}
		}

	default:
		if f.currentField().Type != FormFieldSelect {
			return f.updateCurrentInput(msg)
		}
	}

	return nil
}

var formBoxStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(Primary).
	Padding(0, 2)

var formPickerStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(Primary).
	Padding(0, 1)

func (f *FormModal) View() string {
	if !f.Visible {
		return ""
	}

	lines := []string{f.fixedLine(TitleStyle.Render(f.Title)), f.fixedLine("")}

	start, end := f.visibleFieldRange()
	for i := start; i < end; i++ {
		field := f.Fields[i]
		prefix := "  "
		if i == f.Cursor {
			prefix = "› "
		}

		value := f.renderValue(i)

		lines = append(lines, f.fixedLine(prefix+StatusMuted.Render(field.Label+":")+" "+value))
	}

	if f.hasHiddenFieldsAbove() || f.hasHiddenFieldsBelow() {
		lines = append(lines, f.fixedLine(f.scrollHint()))
	}

	lines = append(lines, f.fixedLine(""))
	lines = append(lines, f.fixedLine(f.submitLine()))
	lines = append(lines, f.fixedLine(f.helpLine()))

	return formBoxStyle.
		Width(formContentWidth).
		MaxWidth(formContentWidth + 4).
		Render(strings.Join(lines, "\n"))
}

func (f *FormModal) PickerView() string {
	if !f.Visible || !f.suggestionOpen {
		return ""
	}

	field := f.currentField()
	var choices []string
	selected := 0
	switch field.Type {
	case FormFieldSelect:
		if len(field.Options) == 0 {
			return ""
		}
		choices = field.Options
		selected = f.optionCursors[f.Cursor]
	case FormFieldAutocomplete:
		choices = f.currentSuggestions()
		selected = f.suggestionIndex
	default:
		return ""
	}

	if len(choices) == 0 {
		return ""
	}
	choices, selected = visibleChoiceWindow(choices, selected, 5)
	return formPickerStyle.Render(strings.Join(f.pickerLines(choices, selected), "\n"))
}

func (f *FormModal) PickerVisible() bool {
	return f.Visible && f.suggestionOpen
}

func (f *FormModal) currentField() FormField {
	if len(f.Fields) == 0 || f.Cursor < 0 || f.Cursor >= len(f.Fields) {
		return FormField{}
	}
	return f.Fields[f.Cursor]
}

func (f *FormModal) moveCursor(delta int) {
	if len(f.Fields) == 0 {
		return
	}
	count := f.submitIndex() + 1
	f.Cursor = (f.Cursor + delta + count) % count
	f.suggestionOpen = false
	f.suggestionIndex = 0
	f.ensureCursorVisible()
	f.focusCursor()
	f.refreshSuggestions()
}

func (f *FormModal) focusCursor() {
	for i := range f.inputs {
		if i == f.Cursor && f.Fields[i].Type != FormFieldSelect {
			f.inputs[i].Focus()
		} else {
			f.inputs[i].Blur()
		}
	}
}

func (f *FormModal) moveChoice(delta int) bool {
	if f.isSubmitFocused() {
		return false
	}
	field := f.currentField()
	if field.Type == FormFieldSelect && len(field.Options) > 0 {
		f.optionCursors[f.Cursor] = (f.optionCursors[f.Cursor] + delta + len(field.Options)) % len(field.Options)
		return true
	}
	if f.suggestionOpen {
		matches := f.currentSuggestions()
		if len(matches) > 0 {
			f.suggestionIndex = (f.suggestionIndex + delta + len(matches)) % len(matches)
			return true
		}
	}
	return false
}

func (f *FormModal) updateCurrentInput(msg tea.Msg) tea.Cmd {
	if f.isSubmitFocused() || f.currentField().Type == FormFieldSelect {
		return nil
	}
	input, cmd := f.inputs[f.Cursor].Update(msg)
	f.inputs[f.Cursor] = input
	f.refreshSuggestions()
	return cmd
}

func (f *FormModal) openChoices() bool {
	if f.isSubmitFocused() {
		return false
	}
	field := f.currentField()
	switch field.Type {
	case FormFieldSelect:
		f.suggestionOpen = len(field.Options) > 0
		return f.suggestionOpen
	case FormFieldAutocomplete:
		matches := f.currentSuggestions()
		if len(matches) == 0 {
			return false
		}
		f.suggestionIndex = selectedOptionIndex(matches, f.inputs[f.Cursor].Value())
		if f.suggestionIndex < 0 {
			f.suggestionIndex = 0
		}
		f.suggestionOpen = true
		return true
	default:
		return false
	}
}

func (f *FormModal) acceptActiveChoice() bool {
	field := f.currentField()
	if field.Type == FormFieldSelect {
		f.suggestionOpen = false
		return true
	}
	if !f.suggestionOpen {
		return false
	}
	matches := f.currentSuggestions()
	if len(matches) == 0 {
		return false
	}
	f.inputs[f.Cursor].SetValue(matches[f.suggestionIndex])
	f.inputs[f.Cursor].CursorEnd()
	f.suggestionOpen = false
	return true
}

func (f *FormModal) refreshSuggestions() {
	if f.isSubmitFocused() {
		f.suggestionOpen = false
		return
	}
	field := f.currentField()
	if field.Type != FormFieldAutocomplete {
		f.suggestionOpen = false
		return
	}
	if !f.suggestionOpen {
		return
	}
	if exactIndex := selectedOptionIndex(f.currentSuggestions(), f.inputs[f.Cursor].Value()); exactIndex >= 0 {
		f.suggestionIndex = exactIndex
		return
	}
	if f.suggestionIndex >= len(f.currentSuggestions()) {
		f.suggestionIndex = 0
	}
}

func (f *FormModal) currentSuggestions() []string {
	field := f.currentField()
	if field.Type != FormFieldAutocomplete {
		return nil
	}
	query := strings.ToLower(f.inputs[f.Cursor].Value())
	if selectedOptionIndexFold(field.Suggestions, query) >= 0 {
		return field.Suggestions
	}

	var matches []string
	for _, suggestion := range field.Suggestions {
		if query == "" || strings.Contains(strings.ToLower(suggestion), query) {
			matches = append(matches, suggestion)
		}
	}
	if len(matches) > 6 {
		return matches[:6]
	}
	return matches
}

func (f *FormModal) valueAt(index int) string {
	field := f.Fields[index]
	if field.Type == FormFieldSelect && len(field.Options) > 0 {
		return field.Options[f.optionCursors[index]]
	}
	if index < len(f.inputs) {
		return f.inputs[index].Value()
	}
	return field.Value
}

func (f *FormModal) renderValue(index int) string {
	field := f.Fields[index]
	if field.Type == FormFieldSelect {
		if len(field.Options) == 0 {
			return StatusMuted.Render("<no options>")
		}
		return f.formValue(field.Options[f.optionCursors[index]], true) + " " + f.chooseIndicator()
	}
	if index < len(f.inputs) {
		f.inputs[index].Width = formInputWidth
		value := f.inputs[index].View()
		if f.inputs[index].Value() == "" && field.Placeholder == "" {
			value = f.formValue("<enter value>", false)
		}
		if field.Type == FormFieldAutocomplete {
			value += " " + f.chooseIndicator()
		}
		return value
	}
	return StatusMuted.Render("<enter value>")
}

func (f *FormModal) submitLine() string {
	label := "Submit"
	var rendered string
	if f.isSubmitFocused() {
		rendered = "› " + ListItemSelectedStyle.Render(label) + " ‹"
	} else {
		rendered = StatusMuted.Render(label)
	}
	return lipgloss.NewStyle().Width(formContentWidth).Align(lipgloss.Center).Render(rendered)
}

func (f *FormModal) pickerLines(choices []string, selected int) []string {
	lines := []string{f.fixedPickerLine(TitleStyle.Render("Choose"))}
	for idx, choice := range choices {
		prefix := "  "
		style := StatusMuted
		if idx == selected {
			prefix = "› "
			style = ListItemSelectedStyle
		}
		lines = append(lines, f.fixedPickerLine(prefix+style.Render(choice)))
	}
	return lines
}

func (f *FormModal) helpLine() string {
	if f.suggestionOpen {
		return Help("enter", "select") + "  " +
			Help("esc", "close") + "  " +
			Help("↑/↓", "choose")
	}
	if f.isSubmitFocused() {
		return Help("enter", "submit") + "  " +
			Help("esc", "cancel") + "  " +
			Help("↑/↓", "fields")
	}
	if f.currentField().Type == FormFieldSelect || f.currentField().Type == FormFieldAutocomplete {
		return Help("enter", "choose") + "  " +
			Help("esc", "cancel") + "  " +
			Help("tab", "next") + "  " +
			Help("↑/↓", "fields")
	}
	return Help("enter", "next") + "  " +
		Help("esc", "cancel") + "  " +
		Help("tab", "next") + "  " +
		Help("↑/↓", "fields")
}

func (f *FormModal) chooseIndicator() string {
	return StatusMuted.Render("[" + HelpKeyStyle.Render("enter") + StatusMuted.Render(" choose") + "]")
}

func (f *FormModal) formValue(value string, set bool) string {
	if set {
		return HelpKeyStyle.Render(value)
	}
	return StatusMuted.Render(value)
}

func (f *FormModal) fixedLine(line string) string {
	return lipgloss.NewStyle().
		Width(formContentWidth).
		MaxWidth(formContentWidth).
		Render(line)
}

func (f *FormModal) fixedPickerLine(line string) string {
	return lipgloss.NewStyle().
		Width(formPickerWidth - 4).
		MaxWidth(formPickerWidth - 4).
		Render(line)
}

func (f *FormModal) submitCmd() tea.Cmd {
	values := make([]string, len(f.Fields))
	for i := range f.Fields {
		values[i] = f.valueAt(i)
	}
	return func() tea.Msg {
		return FormSubmitMsg{Values: values}
	}
}

func (f *FormModal) submitIndex() int {
	return len(f.Fields)
}

func (f *FormModal) isSubmitFocused() bool {
	return f.Cursor == f.submitIndex()
}

func (f *FormModal) ensureCursorVisible() {
	if f.isSubmitFocused() {
		if len(f.Fields)-f.scrollOffset > formVisibleFields {
			f.scrollOffset = len(f.Fields) - formVisibleFields
		}
		if f.scrollOffset < 0 {
			f.scrollOffset = 0
		}
		return
	}
	if f.Cursor < f.scrollOffset {
		f.scrollOffset = f.Cursor
	}
	if f.Cursor >= f.scrollOffset+formVisibleFields {
		f.scrollOffset = f.Cursor - formVisibleFields + 1
	}
}

func (f *FormModal) visibleFieldRange() (int, int) {
	start := f.scrollOffset
	if start < 0 {
		start = 0
	}
	if start > len(f.Fields) {
		start = len(f.Fields)
	}
	end := start + formVisibleFields
	if end > len(f.Fields) {
		end = len(f.Fields)
	}
	return start, end
}

func (f *FormModal) hasHiddenFieldsAbove() bool {
	return f.scrollOffset > 0
}

func (f *FormModal) hasHiddenFieldsBelow() bool {
	_, end := f.visibleFieldRange()
	return end < len(f.Fields)
}

func (f *FormModal) scrollHint() string {
	switch {
	case f.hasHiddenFieldsAbove() && f.hasHiddenFieldsBelow():
		return StatusMuted.Render("  ↑ more fields  ↓ more fields")
	case f.hasHiddenFieldsAbove():
		return StatusMuted.Render("  ↑ more fields")
	case f.hasHiddenFieldsBelow():
		return StatusMuted.Render("  ↓ more fields")
	default:
		return ""
	}
}

func selectedOptionIndex(options []string, value string) int {
	for i, option := range options {
		if option == value {
			return i
		}
	}
	return -1
}

func selectedOptionIndexFold(options []string, lowerValue string) int {
	for i, option := range options {
		if strings.ToLower(option) == lowerValue {
			return i
		}
	}
	return -1
}

func visibleChoiceWindow(choices []string, selected, size int) ([]string, int) {
	if len(choices) <= size {
		return choices, selected
	}
	start := selected - size/2
	if start < 0 {
		start = 0
	}
	if start+size > len(choices) {
		start = len(choices) - size
	}
	return choices[start : start+size], selected - start
}
