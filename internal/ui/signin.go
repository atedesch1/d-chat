package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type signInModel struct {
}

func newSignInModel() *signInModel {
	return &signInModel{}
}

func (m *signInModel) Init() tea.Cmd {
	return nil
}

func (m *signInModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	return m, cmd
}

func (m *signInModel) View() string {
	return string(signInView)
}
