package main

import "fmt"

func statsCard(data StatsData) string {
	W, H := 460, 215
	type item struct {
		label string
		val   int
		color string
	}
	items := []item{
		{"Stars", data.Stars, colorOrange},
		{"Commits (year)", data.Commits, colorPurple},
		{"Pull Requests", data.PRs, colorCyan},
		{"Issues", data.Issues, colorRed},
		{"Repositories", data.PublicRepos, colorGreen},
		{"Followers", data.Followers, colorTitle},
	}

	rows := titleLine(W, esc(data.Name)+"'s GitHub Stats")
	for i, it := range items {
		col := i % 2
		row := i / 2
		lx := 30 + col*220
		vx := 215 + col*220
		y := 80 + row*40
		rows += fmt.Sprintf(
			"  <text x=\"%d\" y=\"%d\" fill=\"%s\" font-size=\"13\" font-family=\"monospace\">%s</text>\n"+
				"  <text x=\"%d\" y=\"%d\" fill=\"%s\" font-size=\"13\" font-family=\"monospace\" text-anchor=\"end\" font-weight=\"bold\">%s</text>\n",
			lx, y, colorSub, esc(it.label),
			vx, y, it.color, commaf(it.val),
		)
	}
	return svgWrap(W, H, rows)
}
