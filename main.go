package main

// Imports and Globals {{{
import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	titleStyle        = lip.NewStyle().MarginLeft(2)
	itemStyle         = lip.NewStyle().PaddingLeft(4)
	selectedItemStyle = lip.NewStyle().PaddingLeft(2).Foreground(lip.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lip.NewStyle().Margin(1, 0, 2, 4)
)

// }}}

type Model struct {
	list       list.Model
	choice     list.Item
	quitting   bool
	fullscreen bool
	root       string
}

func main() {
	model := Model{root: "/home/dharmx", fullscreen: true}

	if absolutePath, e := filepath.Abs(model.root); e != nil {
		fmt.Println("ERROR RUNNING PROGRAM ", e)
		os.Exit(1)
	} else {
		delegate := itemDelegate{}
		style := lip.NewStyle()
		width, height := style.GetMaxWidth(), 10
		if newWidth, newHeight, e := term.GetSize(0); e != nil || model.fullscreen {
			width, height = newWidth, newHeight
		}

		files := list.New([]list.Item{}, delegate, width, height)
		files.Title = absolutePath

		if directory, e := os.Open(absolutePath); e != nil {
			fmt.Println("ERROR RUNNING PROGRAM ", e)
			os.Exit(1)
		} else {
			if names, e := directory.Readdirnames(0); e != nil {
				fmt.Println("ERROR RUNNING PROGRAM ", e)
				os.Exit(1)
			} else {
				for index, name := range names {
					filetype := false
					if info, e := os.Stat(fmt.Sprintf("%s/%s", absolutePath, name)); e == nil {
						filetype = info.IsDir()
					}
					files.InsertItem(index, item{title: name, filetype: filetype})
				}
			}
		}
		model.list = files

		var options []tea.ProgramOption
		if model.fullscreen {
			options = append(options, tea.WithAltScreen())
		}

		if _, e := tea.NewProgram(model, options...).Run(); e != nil {
			fmt.Println("ERROR RUNNING PROGRAM ", e)
			os.Exit(1)
		}
	}
}

// List Item {{{
type item struct {
	title    string
	filetype bool
}

func (file item) FilterValue() string {
	return file.title
}

// }}}

// Delegate {{{
type itemDelegate struct{}

func (delegate itemDelegate) Height() int { return 1 }

func (delegate itemDelegate) Spacing() int { return 0 }

func (delegate itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (delegate itemDelegate) Render(writer io.Writer, model list.Model, index int, listItem list.Item) {
	file, e := listItem.(item)
	if !e {
		return
	}

	contents := fmt.Sprintf("%s", file.title)
	callback := itemStyle.Render
	if index == model.Index() {
		callback = func(symbol ...string) string {
			return selectedItemStyle.Render("┃ " + strings.Join(symbol, " "))
		}
	}
	fmt.Fprint(writer, callback(contents))
}

// }}}

// Startup {{{
func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) View() string {
	return "\n" + model.list.View()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
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

// }}}
