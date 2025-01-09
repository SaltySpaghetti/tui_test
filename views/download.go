package views

import tea "github.com/charmbracelet/bubbletea"

var _ tea.Model = SearchModel{}

type DownloadModel struct {
}

func (m DownloadModel) Init() tea.Cmd {
	return nil
}

func (m DownloadModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View implements tea.Model.
func (m DownloadModel) View() string {
	return ""
}
