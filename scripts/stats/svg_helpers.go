package main

import (
	"fmt"
	"strconv"
	"strings"
)

func svgWrap(W, H int, content string) string {
	return fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
  <rect x="0.5" y="0.5" width="%d" height="%d" rx="6" fill="%s" stroke="%s"/>
%s</svg>`,
		W, H, W, H, W-1, H-1, colorBG, colorBorder, content,
	)
}

func titleLine(W int, label string) string {
	return fmt.Sprintf(
		`  <text x="%d" y="32" fill="%s" font-size="15" font-weight="bold" font-family="monospace" text-anchor="middle">%s</text>
  <line x1="20" y1="42" x2="%d" y2="42" stroke="%s" stroke-width="1"/>
`,
		W/2, colorTitle, label, W-20, colorBorder,
	)
}

// commaf formats an integer with thousands separators (e.g. 1234567 → "1,234,567").
func commaf(n int) string {
	neg := n < 0
	if neg {
		n = -n
	}
	s := strconv.Itoa(n)
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}
