package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

const ANGLE = 30.0

// atan converts an arctan result to degrees.
func atan(val float64) float64 {
	return (math.Atan(val) * 360) / 2 / math.Pi
}

// toEpochDays returns the number of days since the Unix epoch for a given time.
func toEpochDays(t time.Time) int {
	return int(math.Floor(float64(t.Unix()) / (24 * 60 * 60)))
}

// create3DContrib renders the isometric 3D contribution calendar into b.
func create3DContrib(b *strings.Builder, userInfo *UserInfo, x, y, width, height float64) {
	if len(userInfo.Calendar) == 0 {
		return
	}

	firstDate := userInfo.Calendar[0].Time
	sundayOfFirstWeek := toEpochDays(firstDate) - int(firstDate.Weekday())
	weekcount := int(math.Ceil((float64(len(userInfo.Calendar)) + float64(firstDate.Weekday())) / 7.0))

	dx := width / 64
	dy := dx * math.Tan(ANGLE*((2*math.Pi)/360))
	dxx := dx * 0.9
	dyy := dy * 0.9

	offsetX := dx * 7
	offsetY := height - float64(weekcount+7)*dy

	b.WriteString(`<g>` + "\n")

	type sortedDay struct {
		cal  ContributionDay
		week int
		dow  int
		zIdx int
	}
	var sorted []sortedDay
	for _, cal := range userInfo.Calendar {
		week := (toEpochDays(cal.Time) - sundayOfFirstWeek) / 7
		dow := int(cal.Time.Weekday())
		sorted = append(sorted, sortedDay{cal, week, dow, week + dow})
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].zIdx != sorted[j].zIdx {
			return sorted[i].zIdx < sorted[j].zIdx
		}
		return sorted[i].week < sorted[j].week
	})

	for _, d := range sorted {
		baseX := offsetX + float64(d.week-d.dow)*dx
		baseY := offsetY + float64(d.week+d.dow)*dy

		calHeight := 3.0
		if d.cal.Count > 0 {
			calHeight = math.Log10(float64(d.cal.Count)/20+1)*144 + 3
		}

		contribLevel := d.cal.Level

		fmt.Fprintf(b, `<g transform="translate(%.2f %.2f)">`+"\n", baseX, baseY-calHeight)
		fmt.Fprintf(b, `<animateTransform attributeName="transform" type="translate" values="%.2f %.2f;%.2f %.2f" dur="3s" repeatCount="1" fill="freeze"/>`+"\n",
			baseX, baseY-3, baseX, baseY-calHeight)

		widthTop := dxx
		fmt.Fprintf(b, `<rect stroke="none" x="0" y="0" width="%.2f" height="%.2f" transform="skewY(%.2f) skewX(%.2f) scale(%.2f %.2f)" class="cont-top-%d"></rect>`+"\n",
			widthTop, widthTop, -ANGLE, atan(dxx/2/dyy), dxx/widthTop, (2*dyy)/widthTop, contribLevel)

		widthLeft := dxx
		scaleLeft := math.Sqrt(dxx*dxx+dyy*dyy) / widthLeft
		heightLeft := calHeight / scaleLeft
		fmt.Fprintf(b, `<rect stroke="none" x="0" y="0" width="%.2f" height="%.2f" transform="skewY(%.2f) scale(%.2f %.2f)" class="cont-left-%d"><animate attributeName="height" values="%.2f;%.2f" dur="3s" repeatCount="1" fill="freeze"/></rect>`+"\n",
			widthLeft, heightLeft, ANGLE, dxx/widthLeft, scaleLeft, contribLevel, 3.0/scaleLeft, heightLeft)

		widthRight := dxx
		scaleRight := math.Sqrt(dxx*dxx+dyy*dyy) / widthRight
		heightRight := calHeight / scaleRight
		fmt.Fprintf(b, `<rect stroke="none" x="0" y="0" width="%.2f" height="%.2f" transform="translate(%.2f %.2f) skewY(%.2f) scale(%.2f %.2f)" class="cont-right-%d"><animate attributeName="height" values="%.2f;%.2f" dur="3s" repeatCount="1" fill="freeze"/></rect>`+"\n",
			widthRight, heightRight, dxx, dyy, -ANGLE, dxx/widthRight, scaleRight, contribLevel, 3.0/scaleRight, heightRight)

		b.WriteString(`</g>` + "\n")
	}

	b.WriteString(`</g>` + "\n")
}
