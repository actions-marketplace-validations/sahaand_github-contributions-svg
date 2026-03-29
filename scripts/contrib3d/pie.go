package main

import (
	"fmt"
	"math"
	"strings"
)

func createPieLanguage(b *strings.Builder, userInfo *UserInfo, x, y, width, height float64) {
	if userInfo.TotalContributions == 0 {
		return
	}

	var languages []LangInfo
	if len(userInfo.Languages) > 5 {
		languages = append([]LangInfo{}, userInfo.Languages[:5]...)
	} else {
		languages = append([]LangInfo{}, userInfo.Languages...)
	}

	sumContrib := 0
	for _, l := range languages {
		sumContrib += l.Contributions
	}
	otherContrib := userInfo.TotalCommitContributions - sumContrib
	if otherContrib > 0 {
		languages = append(languages, LangInfo{
			Language:      "other",
			Color:         "#444444",
			Contributions: otherContrib,
		})
	}

	// pie function
	radius := height / 2
	margin := radius / 10
	row := 8.0
	offset := (row-float64(len(languages)))/2.0 + 0.5
	fontSize := height / row / 1.5

	fmt.Fprintf(b, `<g transform="translate(%.2f, %.2f)">`+"\n", x, y)

	groupLabelX := radius * 2.1
	fmt.Fprintf(b, `<g transform="translate(%.2f, 0)">`+"\n", groupLabelX)

	for i, d := range languages {
		py := (float64(i)+offset)*(height/row) - fontSize/2
		fmt.Fprintf(b, `<rect x="0" y="%.2f" width="%.2f" height="%.2f" fill="%s" class="stroke-bg" stroke-width="1px"/>`+"\n",
			py, fontSize, fontSize, d.Color)
	}

	for i, d := range languages {
		py := (float64(i) + offset) * (height / row)
		fmt.Fprintf(b, `<text dominant-baseline="middle" x="%.2f" y="%.2f" class="fill-fg" font-size="%.2fpx">%s</text>`+"\n",
			fontSize*1.2, py, fontSize, d.Language)
	}
	b.WriteString("</g>\n")

	b.WriteString(fmt.Sprintf(`<g transform="translate(%.2f, %.2f)">`+"\n", radius, radius))

	total := 0
	for _, d := range languages {
		total += d.Contributions
	}

	innerRadius := radius / 2.0
	outerRadius := radius - margin
	startAngle := 0.0

	for _, d := range languages {
		delta := (float64(d.Contributions) / float64(total)) * 2 * math.Pi
		if delta == 0 {
			continue
		}
		endAngle := startAngle + delta

		path := arcPath(innerRadius, outerRadius, startAngle, endAngle)
		b.WriteString(fmt.Sprintf(`<path d="%s" style="fill:%s" class="stroke-bg" stroke-width="2px"><title>%s %d</title></path>`+"\n",
			path, d.Color, d.Language, d.Contributions))

		startAngle = endAngle
	}
	b.WriteString("</g>\n")
	b.WriteString("</g>\n")
}

func arcPath(ir, or, start, end float64) string {
	// SVG angles: 0 is up. x = r * sin(a), y = -r * cos(a)
	x0_out := or * math.Sin(start)
	y0_out := -or * math.Cos(start)
	x1_out := or * math.Sin(end)
	y1_out := -or * math.Cos(end)

	x0_in := ir * math.Sin(start)
	y0_in := -ir * math.Cos(start)
	x1_in := ir * math.Sin(end)
	y1_in := -ir * math.Cos(end)

	largeArc := 0
	if end-start > math.Pi {
		largeArc = 1
	}

	// For a full circle, we need two arcs because SVG arcs can't draw a 360 circle in one sweep easily.
	// But in pie charts, a single language rarely has 100%, if it does, delta is 2*Pi.
	if end-start >= 2*math.Pi-0.001 {
		return fmt.Sprintf(`M0,%.2f A%.2f,%.2f 0 1,1 0,%.2f A%.2f,%.2f 0 1,1 0,%.2f M0,%.2f A%.2f,%.2f 0 1,0 0,%.2f A%.2f,%.2f 0 1,0 0,%.2f Z`,
			-or, or, or, or, or, or, -or,
			-ir, ir, ir, ir, ir, ir, -ir)
	}

	return fmt.Sprintf(`M%.2f,%.2f A%.2f,%.2f 0 %d,1 %.2f,%.2f L%.2f,%.2f A%.2f,%.2f 0 %d,0 %.2f,%.2f Z`,
		x0_out, y0_out, or, or, largeArc, x1_out, y1_out,
		x1_in, y1_in, ir, ir, largeArc, x0_in, y0_in)
}
