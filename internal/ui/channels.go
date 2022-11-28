package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type channelsModel struct {
	channels []string
	cursor   int
	loading  bool
}

func newChannelsModel() *channelsModel {
	return &channelsModel{
		channels: []string{},
		cursor:   0,
		loading:  true,
	}
}

type channelNamesMsg []string

func (m *channelsModel) getChannelsCmd() tea.Msg {
	return channelNamesMsg{
		"ch1", "ch2", "ch3",
	}
}

type ChannelChosenMsg string

func (m *channelsModel) chooseChannelCmd() tea.Msg {
	return ChannelChosenMsg(m.channels[m.cursor])
}

func (m *channelsModel) Init() tea.Cmd {
	return m.getChannelsCmd
}

func (m *channelsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case channelNamesMsg:
		m.loading = false
		m.channels = msg

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}

		case tea.KeyDown:
			if m.cursor+1 < len(m.channels) {
				m.cursor++
			}

		case tea.KeyEnter:
			cmds = append(cmds, m.chooseChannelCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *channelsModel) View() string {
	s := "Choose a channel from below:\n\n"

	if m.loading {
		s += "  Loading..."
	} else {
		for i := 0; i < len(m.channels); i++ {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			channel := m.channels[i]
			s += fmt.Sprintf("%s %s\n", cursor, channel)
		}
	}
	return s
}
