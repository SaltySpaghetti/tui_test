package main

import (
	"fmt"
	"os"
	"time"
	"tui/test/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (app app) Init() tea.Cmd {
	return nil
}

var appStyle = lipgloss.NewStyle().Padding(1, 2)

type searchResultMsg []list.Item

func (app app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return app, tea.Quit
		default:
			if msg.Type == tea.KeyRunes || msg.Type == tea.KeyBackspace || msg.Type == tea.KeySpace {
				app.searchInput, cmd = app.searchInput.Update(msg)
				cmds = append(cmds, cmd)

				cmd = app.handleSearch()
				cmds = append(cmds, cmd)

				return app, tea.Batch(cmds...)
			}
		}
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		app.resultsList.SetSize(msg.Width-h, msg.Height-v)
		return app, nil

	case searchResultMsg:
		cmd := app.resultsList.SetItems(msg)
		return app, cmd
	}

	// Update the components
	app.resultsList, cmd = app.resultsList.Update(msg)
	cmds = append(cmds, cmd)

	return app, tea.Batch(cmds...)
}

func (app app) handleSearch() tea.Cmd {
	return tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
		app.state.ApiState.Status = models.Loading
		characters, err := models.FetchCharacters(app.searchInput.Value())
		if err != nil {
			app.state.ApiState.Status = models.Error
			app.state.ApiState.Error = err
			return err
		}

		return searchResultMsg(characters)
	})
}

func (app app) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		app.searchInput.View(),
		app.resultsList.View(),
	)
}

// var searchBar, content string
// switch app.state.Step {
// case models.Searching:
// 	searchBar = app.searchInput.View()

// 	switch app.state.ApiState.Status {
// 	case models.Idle:
// 		content = ""
// 	case models.Loading:
// 		content = "Loading..."
// 	case models.Success:
// 		content = app.resultsList.View()
// 	case models.Error:
// 		content = fmt.Sprintf("Error: %v", app.state.ApiState.Error)
// 	}

// 	return fmt.Sprintf(
// 		"%s\n%s",
// 		searchBar,
// 		content,
// 	)
// case models.Downloading:
// 	content = "Downloading..."
// 	return fmt.Sprintf(
// 		"%s",
// 		content,
// 	)
// default:
// 	return fmt.Sprintf("Something bad happened")
// }

type app struct {
	state       models.State
	searchInput textinput.Model
	resultsList list.Model
}

func main() {
	searchInput := textinput.New()
	searchInput.Prompt = "Search: "
	searchInput.Placeholder = "Type to search..."
	searchInput.Focus()

	resultsList := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	resultsList.Title = "Results"
	resultsList.SetShowPagination(false)

	app := app{
		state: models.State{
			Step: models.Searching,
		},
		searchInput: searchInput,
		resultsList: resultsList,
	}

	if _, err := tea.NewProgram(app, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
