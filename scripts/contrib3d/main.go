package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	username := flag.String("username", os.Getenv("USERNAME"), "GitHub username")
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub token (required)")
	output := flag.String("output", "3d-contrib-green.svg", "output SVG file")
	// year=0 means "no filter" → same rolling ~52-week window GitHub's UI uses.
	// Pass a specific year (e.g. --year=2024) to show a full calendar year instead.
	year := flag.Int("year", 0, "contribution year (default: last 52 weeks, matching GitHub's contribution graph)")
	flag.Parse()

	if *username == "" || *token == "" {
		fmt.Fprintln(os.Stderr, "usage: contrib3d --username=<user> --token=<token> [--output=<file>] [--year=<year>]")
		os.Exit(1)
	}

	var yearPtr *int
	if *year > 0 {
		yearPtr = year
	}
	userInfo, err := fetchGraphQL(*token, *username, yearPtr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fetch error:", err)
		os.Exit(1)
	}

	svgStr := createSvg(userInfo)

	if err := os.WriteFile(*output, []byte(svgStr), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "write error:", err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s\n", *output)
}
