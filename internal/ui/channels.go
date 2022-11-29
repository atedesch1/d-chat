package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView      view = "listView"
	textInputView view = "textInputView"
)

type channelsModel struct {
	focusedView view

	channels []string
	cursor   int
	loading  bool

	textInput   textinput.Model
	channelName string
}

func newChannelsModel() *channelsModel {
	ti := textinput.New()
	ti.Placeholder = "Channel"
	ti.Prompt = ""
	ti.CharLimit = 16
	ti.Width = 20

	return &channelsModel{
		focusedView: listView,

		channels: []string{},
		cursor:   0,
		loading:  true,

		textInput:   ti,
		channelName: "",
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

type channelCreatedMsg string

func (m *channelsModel) createChannelCmd() tea.Msg {
	return channelCreatedMsg(m.channelName)
}

func (m *channelsModel) Init() tea.Cmd {
	return m.getChannelsCmd
}

func (m *channelsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case channelNamesMsg:
		m.loading = false
		m.channels = msg

	case channelCreatedMsg:
		cmds = append(cmds, m.createChannelCmd)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			switch m.focusedView {
			case listView:
				m.focusedView = textInputView
			case textInputView:
				m.focusedView = listView
			}

		case tea.KeyCtrlR:
			m.cursor = 0
			cmds = append(cmds, m.getChannelsCmd)

		case tea.KeyUp:
			if m.focusedView == listView && m.cursor > 0 {
				m.cursor--
			}

		case tea.KeyDown:
			if m.focusedView == listView && m.cursor+1 < len(m.channels) {
				m.cursor++
			}

		case tea.KeyEnter:
			switch m.focusedView {
			case listView:
				cmds = append(cmds, m.chooseChannelCmd)
			case textInputView:
				m.channelName = m.textInput.Value()
				m.textInput.Reset()
				cmds = append(cmds, m.createChannelCmd)
			}
		}
	}

	switch m.focusedView {
	case listView:

	case textInputView:
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
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
			if m.cursor == i && m.focusedView == listView {
				cursor = ">"
			}
			channel := m.channels[i]
			s += fmt.Sprintf("%s %s\n", cursor, channel)
		}

		cursor := " "
		if m.focusedView == textInputView {
			cursor = ">"
		}
		s += fmt.Sprintf("\n\nCreate Channel:\n%s %s\n", cursor, m.textInput.View())
	}
	return s
}
