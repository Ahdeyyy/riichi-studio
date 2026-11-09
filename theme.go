package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"runtime"
	"strconv"
	"time"
)

// Every surface is opaque. Stacked translucent whites read differently on
// the page, inside a panel and over the gradient, so the same control
// changed color with its container.
var (
	ink   = ui.Hex("#F1E8D3")
	muted = ui.Hex("#B3C4B9")
	jade  = ui.Hex("#C9A568")
	cream = ui.Hex("#FFF8E6")

	chrome     = ui.Hex("#13261F") // title bar and navigation, one band
	pageTop    = ui.Hex("#1A332C")
	pageBottom = ui.Hex("#11221E")
	panelFill  = ui.Hex("#1D3830")
	hairline   = ui.Hex("#29473E")

	surface        = ui.Hex("#27483F")
	surfaceHover   = ui.Hex("#305449")
	surfacePressed = ui.Hex("#385E52")
	pillHover      = ui.Hex("#22413A")
	pillActive     = ui.Hex("#305449")

	feltText  = ui.Hex("#C6D5BE")
	feltLabel = ui.Hex("#ABC5B2")
)

// CSS ease: hover and selection fades are color changes, never movement.
var pillFade = ui.Animation{Duration: 150 * time.Millisecond, Ease: ui.CubicBezier(.25, .1, .25, 1)}

func appTheme() *ui.Theme {
	theme := ui.DarkTheme()
	// Elsewhere the theme keeps the system's own interface font.
	if runtime.GOOS == "windows" {
		theme.Font = "Segoe UI"
	}
	theme.FontSize = 16
	theme.Background = pageBottom
	theme.Text = ink
	theme.TextMuted = muted
	theme.Accent = jade
	theme.AccentText = ui.Hex("#1C2E28")
	theme.AccentHover = ui.Hex("#E0BC7E")
	theme.AccentPressed = ui.Hex("#B39156")
	theme.Radius = 8
	theme.Surface = surface
	theme.SurfaceHover = surfaceHover
	theme.SurfacePressed = surfacePressed
	// Matches the button face so buttons stay borderless; check boxes mix
	// it with the text color for a visible edge.
	theme.Border = surface
	return theme
}

func panel(c *ui.Context) ui.Element {
	return ui.Column(c).Padding(26).Gap(18).Radius(16).Background(panelFill).Border(1, hairline)
}

// pill is a tab-like choice: no face until hovered, filled while active.
// Call c.Key first; the pill builds its own label child.
func pill(c *ui.Context, label string, active bool, click func()) ui.Element {
	b := ui.ButtonBase(c).Label(label).AlignSelf(ui.Start).Padding(9, 16).Radius(8).AlignItems(ui.Center).OnClick(click)
	// The pill's own hue at zero alpha, so the fade does not pass through gray.
	bg := pillHover.Alpha(0)
	if active {
		bg = pillActive
	} else if b.Hovered() {
		bg = pillHover
	}
	b.Background(b.AnimateColor("pill", bg, pillFade)).Children(func() { ui.Text(c, label).FontSize(14) })
	return b
}

// plural pairs n with the singular or plural noun: "1 tile", "18 tiles".
func plural(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return fmt.Sprintf("%d %s", n, many)
}

// points groups thousands, as payments are read aloud: 48,000.
func points(n int) string {
	s := strconv.Itoa(max(n, -n))
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	if n < 0 {
		s = "-" + s
	}
	return s
}

// The theme never changes, so every frame sets the same one.
var theme = appTheme()
