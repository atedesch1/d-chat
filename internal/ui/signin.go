package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type signInModel struct {
	textInput textinput.Model

	username     string
	hasConfirmed bool
	isTaken      bool
	loading      bool
}

func newSignInModel() *signInModel {
	ti := textinput.New()
	ti.Placeholder = "Username"
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = 20

	return &signInModel{
		textInput: ti,

		username:     "",
		isTaken:      false,
		loading:      false,
		hasConfirmed: false,
	}
}

type EnteredValidUsernameMsg string

func (m *signInModel) enterUsernameCmd() tea.Msg {
	// validate username
	if m.username == "taken" {
		m.isTaken = true
		return nil
	}
	return EnteredValidUsernameMsg(m.username)
}

func (m *signInModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *signInModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.hasConfirmed = true
			m.username = m.textInput.Value()
			cmds = append(cmds, m.enterUsernameCmd)
		}

		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *signInModel) View() string {
	input := m.textInput.Value()

	hint := ""

	if m.hasConfirmed && m.isTaken && input == m.username {
		hint = fmt.Sprintf("Username %s is taken!", m.username)
	} else if input != "" {
		hint = fmt.Sprintf("Press ENTER to join as %s!", input)
	}

	return fmt.Sprintf(
		"Enter a unique username\n\n%s\n\n%s",
		m.textInput.View(),
		hint,
	) + "\n"
}
