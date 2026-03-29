package main

import "fmt"

func streakCard(data StatsData) string {
	W, H := 460, 125
	type item struct {
		label string
		val   int
		color string
	}
	items := []item{
		{"Total Contributions", data.TotalContrib, colorOrange},
		{"Current Streak", data.CurStreak, colorGreen},
		{"Longest Streak", data.Longest, colorPurple},
	}

	content := titleLine(W, "Contribution Streak")
	colW := (W - 40) / 3
	for i, it := range items {
		cx := 20 + i*colW + colW/2
		content += fmt.Sprintf(
			"  <text x=\"%d\" y=\"78\" fill=\"%s\" font-size=\"28\" font-weight=\"bold\" font-family=\"monospace\" text-anchor=\"middle\">%s</text>\n"+
				"  <text x=\"%d\" y=\"100\" fill=\"%s\" font-size=\"11\" font-family=\"monospace\" text-anchor=\"middle\">%s</text>\n",
			cx, it.color, commaf(it.val),
			cx, colorSub, esc(it.label),
		)
		if i < 2 {
			lx := 20 + (i+1)*colW
			content += fmt.Sprintf(
				"  <line x1=\"%d\" y1=\"50\" x2=\"%d\" y2=\"105\" stroke=\"%s\" stroke-width=\"1\"/>\n",
				lx, lx, colorBorder,
			)
		}
	}
	return svgWrap(W, H, content)
}
