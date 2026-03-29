package main

import (
	"fmt"
	"math"
	"strings"
)

// lerp linearly interpolates between lo and hi at parameter u ∈ [0,1].
func lerp(u, lo, hi float64) float64 { return lo + u*(hi-lo) }

// min4 returns the smaller of two ints.
func min4(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// segmentSize returns the (size, borderRadius) for snake body segment i.
// i=0 is the head (largest). Matches snk's quadratic tapering formula.
func segmentSize(i int) (float64, float64) {
	dMin := float64(sizeDot) * 0.8
	dMax := float64(sizeCell) * 0.9
	iMax := float64(min4(bodyLen, 4))
	u := math.Pow(1-math.Min(float64(i), iMax)/iMax, 2)
	s := lerp(u, dMin, dMax)
	r := math.Min(4.5, 4*s/float64(sizeDot))
	return s, r
}

// buildSVG produces the complete animated SVG string from the contribution grid.
func buildSVG(cells []cell) string {
	numCols := 0
	for _, c := range cells {
		if c.col >= numCols {
			numCols = c.col + 1
		}
	}
	numRows := 7

	headPath := buildHeadPath(cells, numCols, numRows)
	eatEvents := buildEatEvents(headPath, cells)

	N := len(headPath)
	if N == 0 {
		return ""
	}
	totalMs := N * stepMs

	svgW := (numCols + 2) * sizeCell
	svgH := (numRows + 5) * sizeCell
	vx := -sizeCell
	vy := -2 * sizeCell

	cellMap := make(map[point]cell, len(cells))
	for _, c := range cells {
		cellMap[point{c.col, c.row}] = c
	}

	eatMap := make(map[point]eatEvent, len(eatEvents))
	for _, e := range eatEvents {
		eatMap[e.pos] = e
	}

	cellID := make(map[point]int, len(eatEvents))
	for i, e := range eatEvents {
		cellID[e.pos] = i
	}

	var b strings.Builder

	// ── SVG header ────────────────────────────────────────────────────────────
	fmt.Fprintf(&b,
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="%d %d %d %d" width="%d" height="%d">`+"\n",
		vx, vy, svgW, svgH, svgW, svgH)
	b.WriteString("<desc>Generated with snk-go</desc>\n")

	// ── CSS ───────────────────────────────────────────────────────────────────
	b.WriteString("<style>\n")

	fmt.Fprintf(&b, ":root{--ce:%s;--cs:%s", colorEmpty, colorSnake)
	for i, c := range levelColors {
		fmt.Fprintf(&b, ";--c%d:%s", i, c)
	}
	b.WriteString("}\n")

	fmt.Fprintf(&b, ".c{fill:var(--ce);animation:c_ linear %dms infinite}\n", totalMs)
	fmt.Fprintf(&b, "@keyframes c_{0%%,100%%{fill:var(--ce)}}\n")

	for _, e := range eatEvents {
		id := cellID[e.pos]
		t := float64(e.stepIdx) / float64(N) * 100.0
		lvl := e.level
		fmt.Fprintf(&b,
			".c.c%d{animation-name:c%d}"+
				"@keyframes c%d{0%%{fill:var(--c%d)}%.4f%%{fill:var(--c%d)}%.4f%%{fill:var(--ce)}100%%{fill:var(--ce)}}\n",
			id, id, id, lvl, t-0.0001, lvl, t+0.0001)
	}

	fmt.Fprintf(&b, ".s{fill:var(--cs);animation:s_ linear %dms infinite}\n", totalMs)
	fmt.Fprintf(&b, "@keyframes s_{}\n")

	for i := 0; i < bodyLen; i++ {
		fmt.Fprintf(&b, ".s.s%d{animation-name:s%d}\n", i, i)
		fmt.Fprintf(&b, "@keyframes s%d{", i)
		prev := point{-1, -1}
		for k := 0; k < N; k++ {
			bi := k - i
			if bi < 0 {
				bi = 0
			}
			pos := headPath[bi]
			if pos == prev && k != 0 && k != N-1 {
				continue
			}
			prev = pos
			t := float64(k) / float64(N) * 100.0
			fmt.Fprintf(&b, "%.4f%%{transform:translate(%dpx,%dpx)}",
				t, pos.col*sizeCell, pos.row*sizeCell)
		}
		pos0 := headPath[0]
		fmt.Fprintf(&b, "100%%{transform:translate(%dpx,%dpx)}}\n", pos0.col*sizeCell, pos0.row*sizeCell)
	}

	totalNonEmpty := len(eatEvents)
	if totalNonEmpty > 0 {
		gridW := numCols * sizeCell
		cellW := float64(gridW) / float64(totalNonEmpty)
		fmt.Fprintf(&b, ".st{fill:var(--ce);animation:st_ linear %dms infinite}\n", totalMs)
		fmt.Fprintf(&b, "@keyframes st_{0%%,100%%{fill:var(--ce)}}\n")
		for i, e := range eatEvents {
			t := float64(e.stepIdx) / float64(N) * 100.0
			fmt.Fprintf(&b,
				".st.st%d{animation-name:st%d}"+
					"@keyframes st%d{0%%{fill:var(--ce)}%.4f%%{fill:var(--ce)}%.4f%%{fill:var(--c%d)}100%%{fill:var(--c%d)}}\n",
				i, i, i, t-0.0001, t+0.0001, e.level, e.level)
			_ = cellW
		}
	}

	b.WriteString("</style>\n")

	// ── Grid cells ────────────────────────────────────────────────────────────
	for col := 0; col < numCols; col++ {
		for row := 0; row < numRows; row++ {
			p := point{col, row}
			cx := col*sizeCell + margin
			cy := row*sizeCell + margin

			class := "c"
			if _, ok := eatMap[p]; ok {
				class = fmt.Sprintf("c c%d", cellID[p])
			}

			fmt.Fprintf(&b,
				`<rect class="%s" x="%d" y="%d" width="%d" height="%d" rx="2" ry="2"/>`+"\n",
				class, cx, cy, sizeDot, sizeDot)
		}
	}

	// ── Progress bar ──────────────────────────────────────────────────────────
	if totalNonEmpty > 0 {
		gridW := numCols * sizeCell
		cellW := float64(gridW) / float64(totalNonEmpty)
		stackY := numRows*sizeCell + sizeDot
		for i := range eatEvents {
			x := float64(i) * cellW
			fmt.Fprintf(&b,
				`<rect class="st st%d" x="%.2f" y="%d" width="%.2f" height="%d"/>`+"\n",
				i, x, stackY, cellW, sizeDot)
		}
	}

	// ── Snake body segments (tail first, head last = drawn on top) ────────────
	startPos := headPath[0]
	startX := startPos.col*sizeCell + margin
	startY := startPos.row*sizeCell + margin

	for i := bodyLen - 1; i >= 0; i-- {
		s, r := segmentSize(i)
		off := int((float64(sizeDot) - s) / 2)
		fmt.Fprintf(&b,
			`<rect class="s s%d" x="%d" y="%d" width="%.1f" height="%.1f" rx="%.1f" ry="%.1f"/>`+"\n",
			i, startX+off, startY+off, s, s, r, r)
	}

	b.WriteString("</svg>\n")
	return b.String()
}
