package main

import (
	"fmt"

	"scripts/githubapi"
)

// fetchGrid fetches the GitHub contribution grid for the given user via the
// GitHub GraphQL API (the exact query used by Platane/snk).
func fetchGrid(username, token string) ([]cell, error) {
	weeks, err := githubapi.FetchContributionCalendar(token, username)
	if err != nil {
		return nil, fmt.Errorf("fetch contribution calendar: %w", err)
	}

	var cells []cell
	for col, week := range weeks {
		for _, day := range week {
			cells = append(cells, cell{
				col:   col,
				row:   day.Weekday, // 0=Sunday as returned by GitHub GraphQL
				level: day.Level,
				date:  day.Date,
			})
		}
	}
	return cells, nil
}
