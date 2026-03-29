package main

import "strings"

const (
	colorBG     = "#1a1b27"
	colorBorder = "#414868"
	colorTitle  = "#7aa2f7"
	colorText   = "#c0caf5"
	colorSub    = "#a9b1d6"
	colorOrange = "#e0af68"
	colorGreen  = "#9ece6a"
	colorPurple = "#bb9af7"
	colorRed    = "#f7768e"
	colorCyan   = "#7dcfff"
)

func esc(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func getRank(category string, value int) string {
	th, ok := thresholds[category]
	if !ok {
		return "-"
	}
	rank := "-"
	for i, t := range th {
		if value >= t {
			rank = ranks[i]
		}
	}
	return rank
}
