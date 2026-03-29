package main

import "fmt"

func trophiesCard(data StatsData) string {
	type entry struct {
		cat string
		val int
	}
	entries := []entry{
		{"Stars", data.Stars},
		{"Followers", data.Followers},
		{"Commits", data.Commits},
		{"Repos", data.PublicRepos},
		{"Pull Requests", data.PRs},
		{"Issues", data.Issues},
	}

	const (
		TW   = 148
		TH   = 108
		COLS = 3
		PAD  = 12
	)
	nrows := (len(entries) + COLS - 1) / COLS
	W := PAD + COLS*(TW+PAD)
	H := 55 + PAD + nrows*(TH+PAD)

	content := fmt.Sprintf(
		"  <text x=\"%d\" y=\"32\" fill=\"%s\" font-size=\"15\" font-weight=\"bold\" font-family=\"monospace\" text-anchor=\"middle\">GitHub Trophies</text>\n"+
			"  <line x1=\"20\" y1=\"42\" x2=\"%d\" y2=\"42\" stroke=\"%s\" stroke-width=\"1\"/>\n",
		W/2, colorTitle, W-20, colorBorder,
	)

	for i, e := range entries {
		col := i % COLS
		row := i / COLS
		cx := PAD + col*(TW+PAD)
		cy := 50 + PAD + row*(TH+PAD)
		rank := getRank(e.cat, e.val)
		clr := rankColors[rank]
		if clr == "" {
			clr = colorSub
		}
		content += fmt.Sprintf(
			"  <g transform=\"translate(%d,%d)\">\n"+
				"    <rect width=\"%d\" height=\"%d\" rx=\"6\" fill=\"%s\" stroke=\"%s\" stroke-width=\"1.5\"/>\n"+
				"    <text x=\"%d\" y=\"32\" fill=\"%s\" font-size=\"18\" font-family=\"monospace\" text-anchor=\"middle\">&#9670;</text>\n"+
				"    <text x=\"%d\" y=\"52\" fill=\"%s\" font-size=\"10\" font-family=\"monospace\" text-anchor=\"middle\">%s</text>\n"+
				"    <rect x=\"40\" y=\"58\" width=\"%d\" height=\"18\" rx=\"3\" fill=\"%s30\"/>\n"+
				"    <text x=\"%d\" y=\"71\" fill=\"%s\" font-size=\"12\" font-weight=\"bold\" font-family=\"monospace\" text-anchor=\"middle\">%s</text>\n"+
				"    <text x=\"%d\" y=\"92\" fill=\"%s\" font-size=\"10\" font-family=\"monospace\" text-anchor=\"middle\">%s</text>\n"+
				"  </g>\n",
			cx, cy,
			TW, TH, colorBG, clr,
			TW/2, clr,
			TW/2, colorText, esc(e.cat),
			TW-80, clr,
			TW/2, clr, rank,
			TW/2, colorSub, commaf(e.val),
		)
	}
	return svgWrap(W, H, content)
}
