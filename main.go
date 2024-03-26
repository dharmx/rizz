package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func main() {
	// TODO: Use CLI args for root and fullscreen.
	if _, e := runProgram(".", false); e != nil {
		fmt.Println("ERROR RUNNING PROGRAM ", e)
		os.Exit(1)
	}
}

func runProgram(root string, fullscreen bool) (tea.Model, error) {
	state := model{root: root, fullscreen: fullscreen}
	if absolutePath, e := filepath.Abs(state.root); e != nil {
		return nil, e
	} else {
		state.list = addFileItems(absolutePath, state)
		var options []tea.ProgramOption
		if state.fullscreen {
			options = append(options, tea.WithAltScreen())
		}

		return tea.NewProgram(state, options...).Run()
	}
}

func addFileItems(absolutePath string, state model) list.Model {
	delegate := itemDelegate{}
	width, height := style.GetMaxWidth(), 10
	if newWidth, newHeight, e := term.GetSize(0); e != nil || state.fullscreen {
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
	return files
}
