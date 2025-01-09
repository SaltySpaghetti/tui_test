package main

import (
	"fmt"
	"os"
	"tui/test/models"
	"tui/test/views"

	"github.com/Broderick-Westrope/charmutils"
	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	step          models.Step
	searchModel   views.SearchModel
	downloadModel views.DownloadModel
}

func (app app) Init() tea.Cmd {
	return nil
}

func (app app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return app, tea.Quit
		}
	}

	switch app.step {
	case models.Searching:
		cmd, err := charmutils.UpdateTypedModel(&app.searchModel, msg)
		if err != nil {
			panic(err)
		}
		return app, cmd

	case models.Downloading:
		cmd, err := charmutils.UpdateTypedModel(&app.downloadModel, msg)
		if err != nil {
			panic(err)
		}
		return app, cmd
	}
	return app, nil
}

func (app app) View() string {
	switch app.step {
	case models.Searching:
		return app.searchModel.View()

	case models.Downloading:
		return app.downloadModel.View()

	default:
		return "ERROR"
	}
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

func main() {
	app := app{
		step:        models.Searching,
		searchModel: views.NewSearchModel(),
	}

	if _, err := tea.NewProgram(app, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
