package main

// StatsData holds all fetched GitHub stats for a user.
type StatsData struct {
	Name         string
	Followers    int
	PublicRepos  int
	Stars        int
	Commits      int
	PRs          int
	Issues       int
	TotalContrib int
	CurStreak    int
	Longest      int
	LangBytes    map[string]int
}

var ranks = []string{"B", "A", "AA", "AAA", "S", "SS", "SSS", "SECRET"}

var thresholds = map[string][]int{
	"Stars":         {1, 10, 100, 500, 1000, 2000, 5000, 10000},
	"Followers":     {1, 10, 25, 50, 100, 250, 500, 1000},
	"Commits":       {1, 10, 100, 500, 1000, 2000, 5000, 10000},
	"Repos":         {1, 5, 10, 20, 50, 100, 200, 500},
	"Pull Requests": {1, 10, 50, 200, 500, 1000, 2000, 5000},
	"Issues":        {1, 10, 50, 200, 500, 1000, 2000, 5000},
}

var rankColors = map[string]string{
	"SECRET": "#ffa500",
	"SSS":    "#ffa500",
	"SS":     "#e8e8e8",
	"S":      "#ffd700",
	"AAA":    "#7aa2f7",
	"AA":     "#7aa2f7",
	"A":      "#9ece6a",
	"B":      "#a9b1d6",
}

var langColors = map[string]string{
	"Kotlin":     "#A97BFF",
	"Java":       "#b07219",
	"Go":         "#00ADD8",
	"Python":     "#3572A5",
	"Swift":      "#F05138",
	"Dart":       "#00B4AB",
	"TypeScript": "#3178c6",
	"JavaScript": "#f1e05a",
	"Rust":       "#dea584",
	"C++":        "#f34b7d",
	"C":          "#555555",
	"Shell":      "#89e051",
	"HTML":       "#e34c26",
	"CSS":        "#563d7c",
	"Ruby":       "#701516",
	"Other":      "#8b949e",
}
