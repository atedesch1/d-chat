package ui

import tea "github.com/charmbracelet/bubbletea"

type chatModel struct {
}

func newChatModel() *chatModel {
	return &chatModel{}
}

func (m *chatModel) Init() tea.Cmd {
	return nil
}

func (m *chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	return m, cmd
}

func (m *chatModel) View() string {
	return string(chatView)
}
