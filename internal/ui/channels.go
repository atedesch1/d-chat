package ui

import tea "github.com/charmbracelet/bubbletea"

type channelsModel struct {
}

func newChannelsModel() *channelsModel {
	return &channelsModel{}
}

func (m *channelsModel) Init() tea.Cmd {
	return nil
}

func (m *channelsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	return m, cmd
}

func (m *channelsModel) View() string {
	return string(channelsView)
}
