package screens

import (
	"strings"
	"time"
	"unicode"

	"github.com/SakshhamTheCoder/adbt/internal/adb"
	"github.com/SakshhamTheCoder/adbt/internal/state"
	"github.com/SakshhamTheCoder/adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

// Android keycodes used by the input screen.
const (
	keyBack    = 4
	keyHome    = 3
	keyRecents = 187
	keyEnter   = 66
	keyDel     = 67
	keyUp      = 19
	keyDown    = 20
	keyLeft    = 21
	keyRight   = 22
)

type Input struct {
	state  *state.AppState
	buffer string
	live   bool
	toast  components.Toast
}

func NewInput(state *state.AppState) *Input {
	return &Input{state: state}
}

func (m *Input) Init() tea.Cmd {
	return nil
}

// CapturingText tells the app this screen is consuming text input, so global
// single-key shortcuts (like "q") get forwarded here instead of quitting.
func (m *Input) CapturingText() bool {
	return m.state.HasDevice()
}

func (m *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.toast.Update(msg)

	switch msg := msg.(type) {
	case adb.InputResultMsg:
		if msg.Error != nil {
			var cmd tea.Cmd
			m.toast, cmd = components.ShowToast(
				"Input failed: "+msg.Error.Error(),
				true,
				3*time.Second,
			)
			return m, cmd
		}
		return m, nil

	case tea.KeyMsg:
		if !m.state.HasDevice() {
			return m, nil
		}

		serial := m.state.DeviceSerial()

		switch msg.String() {
		case "esc":
			return m, nil // let the app fall back to the dashboard
		case "ctrl+b":
			return m, adb.SendKeyEventCmd(serial, keyBack)
		case "ctrl+h":
			return m, adb.SendKeyEventCmd(serial, keyHome)
		case "ctrl+r":
			return m, adb.SendKeyEventCmd(serial, keyRecents)
		case "tab":
			m.live = !m.live
			return m, nil
		}

		if m.live {
			return m, m.handleLiveKey(msg, serial)
		}
		return m, m.handleTextKey(msg, serial)
	}

	return m, nil
}

func (m *Input) handleLiveKey(msg tea.KeyMsg, serial string) tea.Cmd {
	switch msg.String() {
	case "enter":
		return adb.SendKeyEventCmd(serial, keyEnter)
	case "backspace":
		return adb.SendKeyEventCmd(serial, keyDel)
	case "up":
		return adb.SendKeyEventCmd(serial, keyUp)
	case "down":
		return adb.SendKeyEventCmd(serial, keyDown)
	case "left":
		return adb.SendKeyEventCmd(serial, keyLeft)
	case "right":
		return adb.SendKeyEventCmd(serial, keyRight)
	case "space":
		return adb.SendTextCmd(serial, " ")
	}

	if len(msg.Runes) == 1 && unicode.IsPrint(msg.Runes[0]) {
		return adb.SendTextCmd(serial, string(msg.Runes))
	}
	return nil
}

func (m *Input) handleTextKey(msg tea.KeyMsg, serial string) tea.Cmd {
	switch msg.String() {
	case "enter":
		if strings.TrimSpace(m.buffer) == "" {
			return nil
		}
		text := m.buffer
		m.buffer = ""
		return adb.SendTextCmd(serial, text)
	case "backspace":
		if len(m.buffer) > 0 {
			r := []rune(m.buffer)
			m.buffer = string(r[:len(r)-1])
		}
		return nil
	case "space":
		m.buffer += " "
		return nil
	}

	if len(msg.Runes) == 1 && unicode.IsPrint(msg.Runes[0]) {
		m.buffer += string(msg.Runes)
	}
	return nil
}

func (m *Input) View() string {
	if !m.state.HasDevice() {
		return components.RenderNoDevice(m.state, "Input")
	}

	var body strings.Builder

	if m.live {
		body.WriteString(components.StatusConnected.Render("● live capture") + "\n\n")
		body.WriteString(components.StatusMuted.Render("Keys are forwarded to the device as you press them.") + "\n")
		body.WriteString(components.StatusMuted.Render("Each keystroke is a separate adb call, so it may lag.") + "\n")
	} else {
		body.WriteString(components.StatusMuted.Render("text mode") + "\n\n")
		body.WriteString(components.HelpKeyStyle.Render("send: ") + m.buffer + "▌\n")
		body.WriteString(components.StatusMuted.Render("Type a line and press enter to send it to the device.") + "\n")
	}

	body.WriteString("\n" + components.SectionTitle("Buttons") + "\n")
	body.WriteString(components.JoinHelp(
		[2]string{"ctrl+b", "back"},
		[2]string{"ctrl+h", "home"},
		[2]string{"ctrl+r", "recents"},
	) + "\n")

	footer := components.JoinHelp(
		[2]string{"enter", "send"},
		[2]string{"tab", "toggle live"},
		[2]string{"esc", "back"},
	)

	rendered := components.RenderLayout(m.state, "Input", body.String(), footer)

	if m.toast.Visible {
		rendered = components.RenderOverlay(rendered, m.toast.View(), m.state)
	}

	return rendered
}
