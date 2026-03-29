package main

import (
	"fmt"
	"time"

	"scripts/githubapi"
)

const maxReposOneQuery = 100

// ContributionDay holds per-day data including a parsed Time for rendering.
type ContributionDay struct {
	Count int
	Level int
	Date  string
	Time  time.Time
}

// LangInfo holds aggregated commit contribution data for a single language.
type LangInfo struct {
	Language      string
	Color         string
	Contributions int
}

// UserInfo holds all data needed to render the 3D contribution graph.
type UserInfo struct {
	TotalContributions                  int
	TotalCommitContributions            int
	TotalIssueContributions             int
	TotalPullRequestContributions       int
	TotalPullRequestReviewContributions int
	TotalRepositoryContributions        int
	TotalForkCount                      int
	TotalStargazerCount                 int
	Calendar                            []ContributionDay
	Languages                           []LangInfo
	IsHalloween                         bool
}

// fetchGraphQL fetches all data needed for the 3D contribution graph.
func fetchGraphQL(token, userName string, year *int) (*UserInfo, error) {
	weeks, stats, langs, err := githubapi.FetchContributionCalendarWithYear(token, userName, year)
	if err != nil {
		return nil, fmt.Errorf("contribution calendar: %w", err)
	}

	repos, err := githubapi.FetchRepositories(token, userName, maxReposOneQuery)
	if err != nil {
		return nil, fmt.Errorf("repositories: %w", err)
	}

	// Flatten weeks into a sorted calendar slice with parsed Time.
	var calendar []ContributionDay
	for _, week := range weeks {
		for _, d := range week {
			t, _ := time.Parse("2006-01-02", d.Date)
			calendar = append(calendar, ContributionDay{
				Count: d.Count,
				Level: d.Level,
				Date:  d.Date,
				Time:  t,
			})
		}
	}

	// Convert shared LangInfo to local LangInfo.
	languages := make([]LangInfo, len(langs))
	for i, l := range langs {
		languages[i] = LangInfo{
			Language:      l.Language,
			Color:         l.Color,
			Contributions: l.Contributions,
		}
	}

	// Aggregate fork and star counts.
	var totalForks, totalStars int
	for _, r := range repos {
		totalForks += r.ForkCount
		totalStars += r.StargazerCount
	}

	return &UserInfo{
		TotalContributions:                  stats.TotalContributions,
		TotalCommitContributions:            stats.TotalCommitContributions,
		TotalIssueContributions:             stats.TotalIssueContributions,
		TotalPullRequestContributions:       stats.TotalPullRequestContributions,
		TotalPullRequestReviewContributions: stats.TotalPullRequestReviewContributions,
		TotalRepositoryContributions:        stats.TotalRepositoryContributions,
		TotalForkCount:                      totalForks,
		TotalStargazerCount:                 totalStars,
		Calendar:                            calendar,
		Languages:                           languages,
		IsHalloween:                         stats.IsHalloween,
	}, nil
}
