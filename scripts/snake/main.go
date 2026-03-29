// snake generates an animated snake SVG from a GitHub contribution grid.
//
// The snake uses BFS pathfinding to seek and eat every non-empty cell.
// A progress bar below the grid fills as cells are eaten.
// After all cells are cleared the snake returns to (0,0) and the loop repeats.
//
// Usage:
//
//	go run ./scripts/snake --username=<user> --token=<token> [--output=snake.svg]
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	username := flag.String("username", "", "GitHub username (required)")
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub token (required for GraphQL API)")
	output := flag.String("output", "snake.svg", "output SVG file")
	flag.Parse()

	if *username == "" {
		*username = os.Getenv("GITHUB_USERNAME")
	}
	if *username == "" {
		fmt.Fprintln(os.Stderr, "usage: snake --username=<user> --token=<token> [--output=<file>]")
		os.Exit(1)
	}
	if *token == "" {
		fmt.Fprintln(os.Stderr, "error: GitHub token required (set --token or GITHUB_TOKEN)")
		os.Exit(1)
	}

	cells, err := fetchGrid(*username, *token)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fetch:", err)
		os.Exit(1)
	}
	if len(cells) == 0 {
		fmt.Fprintln(os.Stderr, "no contribution data found")
		os.Exit(1)
	}

	svg := buildSVG(cells)

	if err := os.WriteFile(*output, []byte(svg), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "write:", err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s  (%d bytes)\n", *output, len(svg))
}
