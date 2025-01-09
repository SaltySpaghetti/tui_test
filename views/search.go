package views

import (
	"fmt"
	"time"
	"tui/test/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = SearchModel{}

type searchResultMsg []list.Item

type apiErrorMsg error

type SearchModel struct {
	apiState    models.ApiState
	searchInput textinput.Model
	resultsList list.Model
}

var appStyle = lipgloss.NewStyle().Padding(1, 2)

func NewSearchModel() SearchModel {
	searchInput := textinput.New()
	searchInput.Prompt = "Search: "
	searchInput.Placeholder = "Type to search..."
	searchInput.Focus()

	resultsList := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	resultsList.Title = "Results"
	resultsList.SetShowPagination(false)

	return SearchModel{
		searchInput: searchInput,
		resultsList: resultsList,
	}
}

func (m SearchModel) Init() tea.Cmd {
	return nil
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyRunes || msg.Type == tea.KeyBackspace || msg.Type == tea.KeySpace {
			m.searchInput, cmd = m.searchInput.Update(msg)
			cmds = append(cmds, cmd)

			cmd = m.handleSearch()
			cmds = append(cmds, cmd)

			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.resultsList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case searchResultMsg:
		cmd := m.resultsList.SetItems(msg)
		return m, cmd

	case apiErrorMsg:
		m.apiState.Status = models.Error
		m.apiState.Error = msg
		return m, nil
	}

	// Update the components
	m.resultsList, cmd = m.resultsList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)

}

func (m SearchModel) handleSearch() tea.Cmd {
	m.apiState.Status = models.Loading
	return tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
		characters, err := models.FetchCharacters(m.searchInput.Value())
		if err != nil {
			return apiErrorMsg(err)
		}

		return searchResultMsg(characters)
	})
}

// View implements tea.Model.
func (m SearchModel) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		m.searchInput.View(),
		m.resultsList.View(),
	)
}
