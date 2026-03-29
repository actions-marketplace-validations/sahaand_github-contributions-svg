package main

import (
	"fmt"
	"math"
	"strings"
)

var rangeLabels = []string{"1", "10", "100", "1K", "10K"}

func toLevel(value int) float64 {
	if value < 1 {
		return 0.8
	}
	result := math.Log10(float64(value))
	if result > 5 {
		result = 5
	}
	return result + 1
}

func createRadarContrib(b *strings.Builder, userInfo *UserInfo, x, y, width, height float64) {
	radius := (height / 2) * 0.8
	cx := width / 2
	cy := (height / 2) * 1.1

	levels := float64(len(rangeLabels))

	data := []struct {
		Name  string
		Value int
	}{
		{"Commit", userInfo.TotalCommitContributions},
		{"Issue", userInfo.TotalIssueContributions},
		{"PullReq", userInfo.TotalPullRequestContributions},
		{"Review", userInfo.TotalPullRequestReviewContributions},
		{"Repo", userInfo.TotalRepositoryContributions},
	}
	total := float64(len(data))

	posX := func(level, num float64) float64 {
		return radius * (level / levels) * math.Sin((num/total)*2*math.Pi)
	}
	posY := func(level, num float64) float64 {
		return radius * (level / levels) * -math.Cos((num/total)*2*math.Pi)
	}

	fmt.Fprintf(b, `<g transform="translate(%.2f, %.2f)">`+"\n", x+cx, y+cy)

	for j := 0; j < len(rangeLabels); j++ {
		lvl := float64(j + 1)
		for i := 0; i < len(data); i++ {
			x1, y1 := posX(lvl, float64(i)), posY(lvl, float64(i))
			x2, y2 := posX(lvl, float64(i+1)), posY(lvl, float64(i+1))
			fmt.Fprintf(b, `<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" class="stroke-weak" style="stroke-dasharray: 4 4; stroke-width: 1px;" />`+"\n",
				x1, y1, x2, y2)
		}
	}

	for i, d := range rangeLabels {
		yPos := -radius * ((float64(i) + 1) / levels)
		fontSize := radius / 12
		xPos := radius / 50
		fmt.Fprintf(b, `<text x="%.2f" y="%.2f" dominant-baseline="auto" text-anchor="start" style="font-size: %.2fpx" class="fill-weak">%s</text>`+"\n",
			xPos, yPos, fontSize, d)
	}

	for i, d := range data {
		fmt.Fprintf(b, `<g class="axis">`)
		x1, y1 := posX(1, float64(i)), posY(1, float64(i))
		x2, y2 := posX(levels, float64(i)), posY(levels, float64(i))
		fmt.Fprintf(b, `<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" class="stroke-weak" style="stroke-dasharray: 4 4; stroke-width: 1px;" />`+"\n",
			x1, y1, x2, y2)

		lx, ly := posX(1.25*levels, float64(i)), posY(1.17*levels, float64(i))
		fmt.Fprintf(b, `<text x="%.2f" y="%.2f" text-anchor="middle" dominant-baseline="middle" style="font-size: %.2fpx" class="fill-fg">%s<title>%d</title></text>`+"\n",
			lx, ly, radius/7.5, d.Name, d.Value)
		fmt.Fprintf(b, `</g>`+"\n")
	}

	var pts []string
	for i, d := range data {
		lvl := toLevel(d.Value)
		px, py := posX(lvl, float64(i)), posY(lvl, float64(i))
		pts = append(pts, fmt.Sprintf("%.2f,%.2f", px, py))
	}
	pointsStr := strings.Join(pts, " ")

	// The polygon's parent <g> is already translated to the radar center (0,0),
	// so scaling from 0→1 around the origin makes it grow outward from the center.
	fmt.Fprintf(b, `<polygon class="radar" points="%s">`+"\n", pointsStr)
	// Grow from center with ease-out spline (matches the 3D bar rise animation style).
	b.WriteString(`<animateTransform attributeName="transform" type="scale" from="0" to="1" dur="1.5s" begin="0.2s" fill="freeze" calcMode="spline" keyTimes="0;1" keySplines="0.2 0.8 0.4 1"/>` + "\n")
	// Fade in simultaneously so the polygon doesn't pop in abruptly at scale=0.
	b.WriteString(`<animate attributeName="opacity" from="0" to="1" dur="0.8s" begin="0.2s" fill="freeze"/>` + "\n")
	b.WriteString(`</polygon>` + "\n")

	b.WriteString("</g>\n")
}
