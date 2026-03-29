package main

import (
	"fmt"
	"strings"
)

const (
	svgWidth    = 1280.0
	svgHeight   = 850.0
	pieHeight   = 200 * 1.3
	pieWidth    = pieHeight * 2
	radarWidth  = 400 * 1.3
	radarHeight = (radarWidth * 3) / 4
	radarX      = svgWidth - radarWidth - 40
)

const cssTemplate = `
* { font-family: "Ubuntu", "Helvetica", "Arial", sans-serif; }
.fill-bg { fill: transparent; }
.fill-fg { fill: #eeeeff; }
.fill-weak { fill: #aaaaaa; }
.fill-strong { fill: rgb(255,200,55); }
.stroke-bg { stroke: #00000f; }
.stroke-fg { stroke: #eeeeff; }
.stroke-weak { stroke: #aaaaaa; }
.radar { fill: #47a042; fill-opacity: 0.5; stroke: #47a042; stroke-width: 2px; }
.cont-top-0 { fill: #444444; }
.cont-left-0 { fill: #3f3f3f; }
.cont-right-0 { fill: #393939; }
.cont-top-1 { fill: #1b7d28; }
.cont-left-1 { fill: #197425; }
.cont-right-1 { fill: #156920; }
.cont-top-2 { fill: #24a736; }
.cont-left-2 { fill: #219b32; }
.cont-right-2 { fill: #1c8c2b; }
.cont-top-3 { fill: #2dd143; }
.cont-left-3 { fill: #29c23e; }
.cont-right-3 { fill: #23af36; }
.cont-top-4 { fill: #57da69; }
.cont-left-4 { fill: #51ca61; }
.cont-right-4 { fill: #45b754; }
`

// createSvg composes the full 3D contribution SVG for the given user info.
func createSvg(userInfo *UserInfo) string {
	var b strings.Builder

	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`+"\n",
		svgWidth, svgHeight, svgWidth, svgHeight)

	b.WriteString(`<style>`)
	b.WriteString(cssTemplate)
	b.WriteString(`</style>` + "\n")

	fmt.Fprintf(&b, `<rect x="0" y="0" width="%.0f" height="%.0f" class="fill-bg"></rect>`+"\n", svgWidth, svgHeight)

	create3DContrib(&b, userInfo, 0, 0, svgWidth, svgHeight)
	createRadarContrib(&b, userInfo, radarX, 70, radarWidth, radarHeight)
	createPieLanguage(&b, userInfo, 40, svgHeight-pieHeight-70, pieWidth, pieHeight)

	positionXContrib := (svgWidth * 3) / 10
	positionYContrib := svgHeight - 20

	b.WriteString(fmt.Sprintf(`<g><text style="font-size: 32px; font-weight: bold;" x="%.2f" y="%.2f" text-anchor="end" class="fill-strong">%d</text>`+"\n",
		positionXContrib, positionYContrib, userInfo.TotalContributions))

	b.WriteString(fmt.Sprintf(`<text style="font-size: 24px;" x="%.2f" y="%.2f" text-anchor="start" class="fill-fg">contributions</text>`+"\n",
		positionXContrib+10, positionYContrib))

	positionXStar := (svgWidth * 5) / 10
	fmt.Fprintf(&b, `<g transform="translate(%.2f, %.2f) scale(2)"><path fill-rule="evenodd" d="M8 .25a.75.75 0 01.673.418l1.882 3.815 4.21.612a.75.75 0 01.416 1.279l-3.046 2.97.719 4.192a.75.75 0 01-1.088.791L8 12.347l-3.766 1.98a.75.75 0 01-1.088-.79l.72-4.194L.818 6.374a.75.75 0 01.416-1.28l4.21-.611L7.327.668A.75.75 0 018 .25zm0 2.445L6.615 5.5a.75.75 0 01-.564.41l-3.097.45 2.24 2.184a.75.75 0 01.216.664l-.528 3.084 2.769-1.456a.75.75 0 01.698 0l2.77 1.456-.53-3.084a.75.75 0 01.216-.664l2.24-2.183-3.096-.45a.75.75 0 01-.564-.41L8 2.694v.001z" class="fill-fg"></path></g>`+"\n",
		positionXStar-32, positionYContrib-28)
	fmt.Fprintf(&b, `<text style="font-size: 32px; font-weight: bold;" x="%.2f" y="%.2f" text-anchor="start" class="fill-fg">%s<title>%d</title></text>`+"\n",
		positionXStar+10, positionYContrib, toScale(userInfo.TotalStargazerCount), userInfo.TotalStargazerCount)

	positionXFork := (svgWidth * 6) / 10
	fmt.Fprintf(&b, `<g transform="translate(%.2f, %.2f) scale(2)"><path fill-rule="evenodd" d="M5 3.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm0 2.122a2.25 2.25 0 10-1.5 0v.878A2.25 2.25 0 005.75 8.5h1.5v2.128a2.251 2.251 0 101.5 0V8.5h1.5a2.25 2.25 0 002.25-2.25v-.878a2.25 2.25 0 10-1.5 0v.878a.75.75 0 01-.75.75h-4.5A.75.75 0 015 6.25v-.878zm3.75 7.378a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm3-8.75a.75.75 0 100-1.5.75.75 0 000 1.5z" class="fill-fg"></path></g>`+"\n",
		positionXFork-32, positionYContrib-28)
	fmt.Fprintf(&b, `<text style="font-size: 32px; font-weight: bold;" x="%.2f" y="%.2f" text-anchor="start" class="fill-fg">%s<title>%d</title></text>`+"\n",
		positionXFork+4, positionYContrib, toScale(userInfo.TotalForkCount), userInfo.TotalForkCount)

	if len(userInfo.Calendar) > 0 {
		startDate := userInfo.Calendar[0].Time.Format("2006-01-02")
		endDate := userInfo.Calendar[len(userInfo.Calendar)-1].Time.Format("2006-01-02")
		period := startDate + " / " + endDate
		fmt.Fprintf(&b, `<text style="font-size: 16px;" x="%.2f" y="20" dominant-baseline="hanging" text-anchor="end" class="fill-weak">%s</text>`+"\n",
			svgWidth-20, period)
	}
	b.WriteString("</g>\n")
	b.WriteString("</svg>\n")

	return b.String()
}

// toScale formats large numbers with K suffix (e.g. 1500 → "1.5K").
func toScale(v int) string {
	if v < 1000 {
		return fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("%.1fK", float64(v)/1000.0)
}
