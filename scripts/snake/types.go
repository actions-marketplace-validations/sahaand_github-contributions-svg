package main

// ── Layout / sizing (matches Platane/snk defaults) ────────────────────────────

const (
	sizeCell = 16                       // px per grid cell
	sizeDot  = 12                       // px per dot (contribution square)
	margin   = (sizeCell - sizeDot) / 2 // = 2
	stepMs   = 100                      // ms per animation step
	bodyLen  = 4                        // fixed snake body length (segments visible)
)

// ── Color palette ─────────────────────────────────────────────────────────────
// Transparent background, dark-navy empty cells, light→dark green for levels 1–4.

const (
	colorEmpty = "#161b22" // --ce  eaten / empty cell
	colorSnake = "purple"  // --cs  snake body
)

var levelColors = [5]string{
	"#161b22", // --c0  level 0 (same as empty)
	"#9be9a8", // --c1  light green
	"#40c463", // --c2  medium green
	"#30a14e", // --c3  medium-dark green
	"#216e39", // --c4  dark green
}

// ── Data types ────────────────────────────────────────────────────────────────

type cell struct {
	col, row, level int
	date            string
}

type point struct{ col, row int }

type eatEvent struct {
	stepIdx int
	pos     point
	level   int
}
