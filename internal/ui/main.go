package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type view string

const (
	signInView   view = "signInView"
	channelsView view = "channelsView"
	chatView     view = "chatView"
)

type MainModel struct {
	currentView   view
	models        map[view]tea.Model
	width, height int
}

func NewMainModel() MainModel {
	m := make(map[view]tea.Model, 0)
	m[signInView] = newSignInModel()
	m[channelsView] = newChannelsModel()
	m[chatView] = newChatModel()

	return MainModel{
		currentView: signInView,
		models:      m,
	}
}

func (m MainModel) getCurrentModel() tea.Model {
	return m.models[m.currentView]
}

func (m MainModel) Init() tea.Cmd {
	return m.getCurrentModel().Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return m, tea.Quit
		}
	}
	m.models[m.currentView], cmd = m.getCurrentModel().Update(msg)
	return m, cmd
}

func (m MainModel) View() string {
	return m.getCurrentModel().View()
}
