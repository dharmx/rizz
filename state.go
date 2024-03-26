package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list       list.Model
	choice     list.Item
	quitting   bool
	fullscreen bool
	root       string
}

func (model model) Init() tea.Cmd {
	return nil
}

func (model model) View() string {
	return "\n" + model.list.View()
}

func (model model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "q", tea.KeyEscape.String(), tea.KeyCtrlC.String():
			return model, tea.Quit
		case "j":
			index := model.list.Index() + 1
			if index == len(model.list.Items()) {
				index = 0
			}
			model.list.Select(index)
		case "k":
			index := model.list.Index() - 1
			if index == -1 {
				index = len(model.list.Items()) - 1
			}
			model.list.Select(index)
		case "enter":
			if selectedItem, e := model.list.SelectedItem().(item); e {
				model.choice = selectedItem
			}
			return model, tea.Quit
		case tea.KeyCtrlF.String():
			if model.fullscreen {
				model.fullscreen = false
				return model, tea.ExitAltScreen
			}
			model.fullscreen = true
			return model, tea.EnterAltScreen
		}
	}
	return model, nil
}
