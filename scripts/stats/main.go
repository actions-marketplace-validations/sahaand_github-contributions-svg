package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	outputDir := flag.String("output-dir", "assets", "directory to write generated SVG files")
	flag.Parse()

	username := os.Getenv("GITHUB_USERNAME")
	if username == "" {
		username = "saeedata"
	}
	fmt.Printf("Generating README assets for @%s...\n", username)

	data, err := fetchData()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  stars=%d  commits=%d  prs=%d  streak=%d\n",
		data.Stars, data.Commits, data.PRs, data.CurStreak)

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}

	files := map[string]string{
		filepath.Join(*outputDir, "stats.svg"):     statsCard(data),
		filepath.Join(*outputDir, "streak.svg"):    streakCard(data),
		filepath.Join(*outputDir, "trophies.svg"):  trophiesCard(data),
		filepath.Join(*outputDir, "top-langs.svg"): topLangsCard(data.LangBytes),
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("  wrote %s\n", path)
	}
	fmt.Println("Done.")
}
