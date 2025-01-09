package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/charmbracelet/bubbles/list"
)

type ApiResponse struct {
	Info    Info        `json:"info"`
	Results []Character `json:"results"`
}

type Info struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}

type Character struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Type    string `json:"type"`
	Gender  string `json:"gender"`
}

func (c Character) Title() string       { return c.Name }
func (c Character) FilterValue() string { return c.Name }
func (c Character) Description() string {
	return fmt.Sprintf("%s • %s • %s", c.Species, c.Status, c.Gender)
}

func FetchCharacters(term string) ([]list.Item, error) {
	results := make([]list.Item, 0)
	if term == "" {
		// Return empty list when no search term
		return results, nil
	}

	// URL encode the search term
	escapedTerm := url.QueryEscape(term)
	url := fmt.Sprintf("https://rickandmortyapi.com/api/character/?name=%s", escapedTerm)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle 404 case specially since it's a valid response when no results are found
	if resp.StatusCode == 404 {
		return results, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	for _, character := range apiResp.Results {
		results = append(results, character)
	}

	return results, nil
}
