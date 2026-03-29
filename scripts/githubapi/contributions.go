package githubapi

import "fmt"

// FetchContributionCalendar fetches the contribution calendar for the given user
// using the exact GraphQL query from Platane/snk.
// Returns weeks as [][]ContributionDay (outer = weeks, inner = days, day.Weekday 0=Sunday).
func FetchContributionCalendar(token, username string) ([][]ContributionDay, error) {
	const query = `
		query ($login: String!) {
			user(login: $login) {
				contributionsCollection {
					contributionCalendar {
						weeks {
							contributionDays {
								contributionCount
								contributionLevel
								weekday
								date
							}
						}
					}
				}
			}
		}
	`
	result, err := GqlQuery(token, query, map[string]any{"login": username})
	if err != nil {
		return nil, err
	}

	data, _ := result["data"].(map[string]any)
	user, _ := data["user"].(map[string]any)
	cc, _ := user["contributionsCollection"].(map[string]any)
	cal, _ := cc["contributionCalendar"].(map[string]any)

	return parseCalendarWeeks(cal), nil
}

// FetchContributionCalendarWithYear fetches the contribution calendar for a specific year
// using the yoshi389111/github-profile-3d-contrib style query.
// Also fetches all contribution totals and commitContributionsByRepository for language stats.
// Returns (weeks, stats, languages, error).
func FetchContributionCalendarWithYear(token, username string, year *int) ([][]ContributionDay, ContribStats, []LangInfo, error) {
	yearClause := ""
	if year != nil && *year > 0 {
		yearClause = fmt.Sprintf(`(from:"%d-01-01T00:00:00.000Z", to:"%d-12-31T23:59:59.000Z")`, *year, *year)
	}

	query := fmt.Sprintf(`
		query ($login: String!) {
			user(login: $login) {
				contributionsCollection%s {
					contributionCalendar {
						isHalloween
						totalContributions
						weeks {
							contributionDays {
								contributionCount
								contributionLevel
								weekday
								date
							}
						}
					}
					commitContributionsByRepository(maxRepositories: 100) {
						repository {
							primaryLanguage {
								name
								color
							}
						}
						contributions {
							totalCount
						}
					}
					totalCommitContributions
					totalIssueContributions
					totalPullRequestContributions
					totalPullRequestReviewContributions
					totalRepositoryContributions
				}
			}
		}
	`, yearClause)

	result, err := GqlQuery(token, query, map[string]any{"login": username})
	if err != nil {
		return nil, ContribStats{}, nil, err
	}

	data, _ := result["data"].(map[string]any)
	user, _ := data["user"].(map[string]any)
	cc, _ := user["contributionsCollection"].(map[string]any)
	cal, _ := cc["contributionCalendar"].(map[string]any)

	stats := ContribStats{
		TotalContributions:                  AsInt(cal, "totalContributions"),
		TotalCommitContributions:            AsInt(cc, "totalCommitContributions"),
		TotalIssueContributions:             AsInt(cc, "totalIssueContributions"),
		TotalPullRequestContributions:       AsInt(cc, "totalPullRequestContributions"),
		TotalPullRequestReviewContributions: AsInt(cc, "totalPullRequestReviewContributions"),
		TotalRepositoryContributions:        AsInt(cc, "totalRepositoryContributions"),
		IsHalloween:                         func() bool { v, _ := cal["isHalloween"].(bool); return v }(),
	}

	weeks := parseCalendarWeeks(cal)
	langs := parseLanguages(cc)

	return weeks, stats, langs, nil
}

// parseCalendarWeeks extracts [][]ContributionDay from a decoded contributionCalendar map.
func parseCalendarWeeks(cal map[string]any) [][]ContributionDay {
	rawWeeks, _ := cal["weeks"].([]any)
	weeks := make([][]ContributionDay, 0, len(rawWeeks))
	for _, w := range rawWeeks {
		week, _ := w.(map[string]any)
		rawDays, _ := week["contributionDays"].([]any)
		days := make([]ContributionDay, 0, len(rawDays))
		for _, d := range rawDays {
			day, _ := d.(map[string]any)
			date, _ := day["date"].(string)
			levelStr, _ := day["contributionLevel"].(string)
			weekday := AsInt(day, "weekday")
			days = append(days, ContributionDay{
				Date:    date,
				Count:   AsInt(day, "contributionCount"),
				Level:   LevelToInt(levelStr),
				Weekday: weekday,
			})
		}
		weeks = append(weeks, days)
	}
	return weeks
}

// parseLanguages extracts []LangInfo from a decoded contributionsCollection map.
func parseLanguages(cc map[string]any) []LangInfo {
	repoContribs, _ := cc["commitContributionsByRepository"].([]any)
	langMap := make(map[string]LangInfo)
	for _, entry := range repoContribs {
		e, _ := entry.(map[string]any)
		repo, _ := e["repository"].(map[string]any)
		primaryLang, _ := repo["primaryLanguage"].(map[string]any)
		if primaryLang == nil {
			continue
		}
		langName, _ := primaryLang["name"].(string)
		if langName == "" {
			continue
		}
		color, _ := primaryLang["color"].(string)
		if color == "" {
			color = "#444444"
		}
		contribs := AsInt(AsIntMap(e["contributions"]), "totalCount")
		if existing, ok := langMap[langName]; ok {
			existing.Contributions += contribs
			langMap[langName] = existing
		} else {
			langMap[langName] = LangInfo{Language: langName, Color: color, Contributions: contribs}
		}
	}
	langs := make([]LangInfo, 0, len(langMap))
	for _, l := range langMap {
		langs = append(langs, l)
	}
	// Sort by contributions descending.
	for i := 0; i < len(langs); i++ {
		for j := i + 1; j < len(langs); j++ {
			if langs[j].Contributions > langs[i].Contributions {
				langs[i], langs[j] = langs[j], langs[i]
			}
		}
	}
	return langs
}

// AsIntMap safely casts an any value to map[string]any.
func AsIntMap(v any) map[string]any {
	m, _ := v.(map[string]any)
	return m
}
