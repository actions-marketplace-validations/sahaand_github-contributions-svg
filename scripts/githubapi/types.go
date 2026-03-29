package githubapi

// ContributionDay holds data for a single day in a GitHub contribution calendar.
// Level maps the contributionLevel string to an integer: 0=NONE, 1–4=FIRST–FOURTH_QUARTILE.
type ContributionDay struct {
	Date    string
	Count   int
	Level   int // 0–4
	Weekday int // 0=Sunday, 6=Saturday (as returned by GitHub GraphQL)
}

// LangInfo holds aggregated commit contribution data for a single language.
type LangInfo struct {
	Language      string
	Color         string
	Contributions int
}

// RepoNode holds star and fork counts for a single repository.
type RepoNode struct {
	ForkCount      int
	StargazerCount int
}

// ContribStats aggregates the numeric contribution totals returned by GitHub's
// contributionsCollection for a given period.
type ContribStats struct {
	TotalContributions                  int
	TotalCommitContributions            int
	TotalIssueContributions             int
	TotalPullRequestContributions       int
	TotalPullRequestReviewContributions int
	TotalRepositoryContributions        int
	IsHalloween                         bool
}

// levelToInt converts a GitHub contributionLevel string to an integer 0–4.
func LevelToInt(s string) int {
	switch s {
	case "FIRST_QUARTILE":
		return 1
	case "SECOND_QUARTILE":
		return 2
	case "THIRD_QUARTILE":
		return 3
	case "FOURTH_QUARTILE":
		return 4
	default:
		return 0
	}
}
