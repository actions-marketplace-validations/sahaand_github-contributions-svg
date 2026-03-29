package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"scripts/githubapi"
)

func fetchData() (StatsData, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return StatsData{}, fmt.Errorf("GITHUB_TOKEN is not set")
	}
	username := os.Getenv("GITHUB_USERNAME")
	if username == "" {
		username = "saeedata"
	}

	fmt.Println("  → user info")
	userResp, err := githubapi.RestGet(token, "/users/"+username)
	if err != nil {
		return StatsData{}, fmt.Errorf("user info: %w", err)
	}

	fmt.Println("  → repositories")
	repos, err := githubapi.RestAll(token, "/users/"+username+"/repos?type=owner&sort=updated")
	if err != nil {
		return StatsData{}, fmt.Errorf("repos: %w", err)
	}

	totalStars := 0
	langBytes := map[string]int{}
	for _, r := range repos {
		totalStars += githubapi.AsInt(r, "stargazers_count")
		fork, _ := r["fork"].(bool)
		lang, _ := r["language"].(string)
		if !fork && lang != "" {
			size := githubapi.AsInt(r, "size")
			if size == 0 {
				size = 1
			}
			langBytes[lang] += size
		}
	}

	fmt.Println("  → contribution data")
	gqlResult, err := githubapi.GqlQuery(token, `
		query($login: String!) {
		  user(login: $login) {
		    contributionsCollection {
		      totalCommitContributions
		      contributionCalendar {
		        totalContributions
		        weeks {
		          contributionDays {
		            date
		            contributionCount
		          }
		        }
		      }
		    }
		    pullRequests(states: MERGED) { totalCount }
		    issues { totalCount }
		  }
		}
	`, map[string]any{"login": username})
	if err != nil {
		return StatsData{}, fmt.Errorf("graphql: %w", err)
	}

	gqlData, _ := gqlResult["data"].(map[string]any)
	gu, _ := gqlData["user"].(map[string]any)
	cc, _ := gu["contributionsCollection"].(map[string]any)
	cal, _ := cc["contributionCalendar"].(map[string]any)
	prs, _ := gu["pullRequests"].(map[string]any)
	issues, _ := gu["issues"].(map[string]any)

	// Collect contribution days up to today, then sort descending.
	today := time.Now().Format("2006-01-02")
	type contribDay struct {
		date  string
		count int
	}
	var days []contribDay
	if weeks, ok := cal["weeks"].([]any); ok {
		for _, w := range weeks {
			week, _ := w.(map[string]any)
			if ds, ok := week["contributionDays"].([]any); ok {
				for _, d := range ds {
					day, _ := d.(map[string]any)
					dt, _ := day["date"].(string)
					cnt := githubapi.AsInt(day, "contributionCount")
					if dt != "" && dt <= today {
						days = append(days, contribDay{dt, cnt})
					}
				}
			}
		}
	}
	sort.Slice(days, func(i, j int) bool { return days[i].date > days[j].date })

	curStreak, longest, curRun := 0, 0, 0
	counting := true
	for _, d := range days {
		if d.count > 0 {
			if counting {
				curStreak++
			}
			curRun++
			if curRun > longest {
				longest = curRun
			}
		} else {
			counting = false
			curRun = 0
		}
	}

	name, _ := userResp["name"].(string)
	if name == "" {
		name = username
	}

	return StatsData{
		Name:         name,
		Followers:    githubapi.AsInt(userResp, "followers"),
		PublicRepos:  len(repos),
		Stars:        totalStars,
		Commits:      githubapi.AsInt(cc, "totalCommitContributions"),
		PRs:          githubapi.AsInt(prs, "totalCount"),
		Issues:       githubapi.AsInt(issues, "totalCount"),
		TotalContrib: githubapi.AsInt(cal, "totalContributions"),
		CurStreak:    curStreak,
		Longest:      longest,
		LangBytes:    langBytes,
	}, nil
}
