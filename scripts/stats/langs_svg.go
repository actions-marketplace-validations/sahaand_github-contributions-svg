package main

import (
	"fmt"
	"sort"
)

func topLangsCard(langBytes map[string]int) string {
	if len(langBytes) == 0 {
		return svgWrap(300, 60, fmt.Sprintf(
			"  <text x=\"20\" y=\"38\" fill=\"%s\" font-family=\"monospace\">No language data</text>\n", colorSub,
		))
	}

	W, H := 300, 195

	type langEntry struct {
		name  string
		bytes int
	}
	var sorted []langEntry
	for name, b := range langBytes {
		sorted = append(sorted, langEntry{name, b})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].bytes > sorted[j].bytes })
	if len(sorted) > 6 {
		sorted = sorted[:6]
	}

	total := 0
	for _, e := range sorted {
		total += e.bytes
	}

	content := titleLine(W, "Top Languages")

	for i, e := range sorted {
		col := i % 2
		row := i / 2
		lx := 18 + col*144
		ly := 65 + row*28
		pct := float64(e.bytes) / float64(total) * 100
		clr := langColors[e.name]
		if clr == "" {
			clr = langColors["Other"]
		}
		content += fmt.Sprintf(
			"  <circle cx=\"%d\" cy=\"%d\" r=\"5\" fill=\"%s\"/>\n"+
				"  <text x=\"%d\" y=\"%d\" fill=\"%s\" font-size=\"11\" font-family=\"monospace\">%s %.1f%%</text>\n",
			lx+6, ly-5, clr,
			lx+16, ly, colorText, esc(e.name), pct,
		)
	}

	barX, barY, barH := 18, 170, 10
	barW := W - 36
	content += fmt.Sprintf(
		"  <rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" rx=\"5\" fill=\"%s\"/>\n",
		barX, barY, barW, barH, colorBorder,
	)
	off := float64(barX)
	for _, e := range sorted {
		sw := float64(barW) * float64(e.bytes) / float64(total)
		clr := langColors[e.name]
		if clr == "" {
			clr = langColors["Other"]
		}
		content += fmt.Sprintf(
			"  <rect x=\"%.1f\" y=\"%d\" width=\"%.1f\" height=\"%d\" fill=\"%s\"/>\n",
			off, barY, sw, barH, clr,
		)
		off += sw
	}

	return svgWrap(W, H, content)
}
