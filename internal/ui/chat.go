package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chatModel struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style

	loading       bool
	width, height int
}

func newChatModel() *chatModel {
	return &chatModel{
		messages:    []string{},
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		loading:     true,
	}
}

func (m *chatModel) setInterface() {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(m.width)
	ta.SetHeight(m.height / 5)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(m.width, m.height/5*4-5)
	vp.SetContent(`Welcome to the chat room!
	Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	m.textarea = ta
	m.viewport = vp

}

type RequestViewportDimensionsMsg string

func (m *chatModel) getViewportDimensionsCmd() tea.Msg {
	return RequestViewportDimensionsMsg("request")
}

func (m *chatModel) Init() tea.Cmd {
	// return textarea.Blink
	return m.getViewportDimensionsCmd
}

func (m *chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if !m.loading {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		if m.loading {
			m.loading = false
			m.setInterface()
		}

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEnter:
			if !m.loading {
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.textarea.Reset()
				m.viewport.GotoBottom()
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *chatModel) View() string {
	if m.loading {
		return " Loading..."
	}

	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
