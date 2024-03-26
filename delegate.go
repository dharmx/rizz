package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

var (
	style             = lip.NewStyle()
	titleStyle        = style.MarginLeft(2)
	itemStyle         = style.PaddingLeft(4)
	selectedItemStyle = style.PaddingLeft(2).Foreground(lip.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = style.Margin(1, 0, 2, 4)
)

type itemDelegate struct{}

func (delegate itemDelegate) Height() int { return 1 }

func (delegate itemDelegate) Spacing() int { return 0 }

func (delegate itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (delegate itemDelegate) Render(writer io.Writer, model list.Model, index int, listItem list.Item) {
	file, e := listItem.(item)
	if !e {
		return
	}

	callback := itemStyle.Render
	if index == model.Index() {
		callback = func(contents ...string) string {
			return selectedItemStyle.Render("┃ " + strings.Join(contents, " "))
		}
	}
	fmt.Fprint(writer, callback(file.title))
}
