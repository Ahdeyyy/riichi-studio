package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"image"
	"image/color"
	"strings"
	"sync"
	"time"
)

func body(c *ui.Context, s string)      { ui.Text(c, s).FontSize(16).TextColor(ink).FillWidth() }
func secondary(c *ui.Context, s string) { ui.Text(c, s).FontSize(14).TextColor(muted).FillWidth() }
func heading(c *ui.Context, s string)   { ui.Text(c, s).FontSize(24).Bold().TextColor(ink) }

func (a *app) view(c *ui.Context) {
	c.SetTheme(theme)
	ui.Column(c).Fill().Gradient(pageTop, pageBottom, 180).Children(func() {
		a.titlebar(c)
		ui.Row(c).FillWidth().Padding(4, 28, 10).Gap(6).AlignItems(ui.Center).Background(chrome).BorderWidth(0, 0, 1, 0).BorderColor(hairline).Children(func() {
			for _, page := range []string{"Practice", "Learn", "Exercises", "Scoring lab", "Reference"} {
				c.Key(page)
				pill(c, page, a.page == page, func() { a.navigate(page) })
			}
			ui.Spacer(c).Grow(1)
			a.updateButton(c)
			ui.Checkbox(c, &a.showTileLabels, "Tile labels").OnChange(a.savePreferences)
		})
		ui.Scroll(c).FillWidth().Grow(1).Children(func() {
			c.Key(a.page)
			ui.Column(c).FillWidth().Padding(22, 28, 28).Gap(18).Transition(ui.ElementTransition{Duration: 260 * time.Millisecond, Enter: &ui.Motion{Y: 12}}).Children(func() {
				switch a.page {
				case "Learn":
					a.learnView(c)
				case "Wait trainer", "Exercises":
					a.waitTrainer(c)
				case "Scoring lab":
					a.scoringView(c)
				case "Reference":
					a.referenceView(c)
				default:
					a.practiceView(c)
				}
			})
		})
		a.devToolbar(c)
	})
}

func (a *app) practiceView(c *ui.Context) {
	ui.Row(c).FillWidth().Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
		ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
			heading(c, "Your practice table")
			ui.Text(c, "Solo hand").FontSize(14).TextColor(muted).SingleLine().Shrink(0)
		})
		ui.Row(c).Gap(12).Children(func() {
			ui.Button(c, "Restart guide").Disabled(a.busy).OnClick(func() { a.resetTable(false) })
			ui.Button(c, "New random hand").Disabled(a.busy).OnClick(func() { a.resetTable(true) })
		})
	})
	ui.Column(c).FillWidth().Padding(9).Radius(24).Gradient(ui.Hex("#685244"), ui.Hex("#2C2723"), 100).Border(.5, ui.Hex("#8A7561")).Shadow(0, 12, 22, 0, ui.RGBA(0, 0, 0, .23)).Children(func() {
		ui.Column(c).FillWidth().Padding(22, 24, 24).Gap(16).Radius(16).Gradient(ui.Hex("#245D50"), ui.Hex("#184236"), 100).InsetShadow(0, 3, 9, 0, ui.RGBA(0, 0, 0, .3)).DrawOver(feltTexture).Children(func() {
			ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
				ui.Text(c, fmt.Sprintf("Draw %d", a.turn)).FontSize(14).TextColor(feltText)
				ui.Text(c, plural(len(a.wall), "tile", "tiles")+" remaining").FontSize(14).TextColor(feltText)
			})
			ui.Box(c).FillWidth().Height(48).Draw(wallPainter(len(a.wall)))
			ui.Column(c).FillWidth().MinHeight(148).Center().Gap(12).Children(func() {
				if len(a.discards) == 0 {
					ui.Box(c).Size(100, 100).Draw(drawCompass)
				} else {
					ui.Text(c, "Discards").FontSize(14).TextColor(feltLabel)
					// A long solo hand discards most of the wall. The pond keeps
					// the compass's height, so the hand and the Discard button
					// never move: tiles take the largest size whose rows fit.
					width, _ := c.Size()
					size, rowWidth := pondLayout(len(a.discards), width-122)
					pond := ui.Row(c).Wrap().Gap(4).Justify(ui.Center).MaxWidth(rowWidth)
					pond.Children(func() {
						for i, t := range a.discards {
							c.Key(a.discardIDs[i])
							a.tableTile(c, t, false, tileName(t), func() {}, size, a.discardIDs[i], 0)
						}
					})
				}
			})
			ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
				ui.Text(c, "Your hand").FontSize(14).TextColor(feltText)
				if a.riichi {
					ui.Text(c, "Riichi").FontSize(16).TextColor(jade)
				} else {
					ui.Text(c, "Select a tile to lift it").FontSize(14).TextColor(feltText)
				}
			})
			ui.Column(c).FillWidth().Padding(16, 0, 6).Children(func() {
				ui.Row(c).FillWidth().Justify(ui.Center).Gap(3).Children(func() {
					width, _ := c.Size()
					tileWidth := min(float32(64), (width-145)/14-3)
					for i, t := range a.hand {
						delay := time.Duration(i) * 22 * time.Millisecond
						if i == 13 && a.turn > 1 {
							delay = -1
						}
						c.Key(a.handIDs[i])
						if i == 13 && a.turn > 1 {
							ui.Box(c).Width(10)
						}
						a.tableTile(c, t, i == a.selected, fmt.Sprintf("%s tile %d", tileName(t), i+1), func() {
							if !a.won && (!a.riichi || i == len(a.hand)-1) {
								a.selected = i
							}
						}, tileWidth, a.handIDs[i], delay)
					}
				})
			})
		})
	})
	ui.Row(c).FillWidth().Gap(12).AlignItems(ui.Center).Children(func() {
		ui.PrimaryButton(c, "Discard tile").Height(44).Padding(0, 26).Disabled(a.selected < 0 || a.won || len(a.wall) == 0).OnClick(a.discard)
		ready := false
		if a.selected >= 0 && !a.riichi && !a.won && len(a.wall) >= 4 {
			r := append([]int(nil), a.hand[:a.selected]...)
			r = append(r, a.hand[a.selected+1:]...)
			ready = shanten(r) == 0
		}
		ui.Button(c, "Declare riichi").Height(44).Disabled(!ready).OnClick(func() {
			a.riichi = true
			a.discard()
			a.tableMessage = "Riichi declared. Your hand is locked; discard each new draw or claim tsumo."
		})
		ui.Button(c, "Tsumo").Height(44).Disabled(a.won || shanten(a.hand) != -1 || a.busy).OnClick(func() {
			a.won = true
			a.tableMessage = "Tsumo. A closed self-draw supplies a yaku."
			a.calculate(scoreInput{hand: notation(a.hand), win: tileCode(a.hand[len(a.hand)-1]), tsumo: true, riichi: a.riichi})
		})
		ui.Spacer(c).Grow(1)
		ui.Button(c, "Show hint").Height(44).OnClick(func() { a.hint = !a.hint })
	})
	ui.Column(c).FillWidth().Padding(12, 4).Gap(10).Children(func() {
		ui.Text(c, a.tableMessage).FontSize(16).FillWidth().ContentTransition(ui.ContentTransition{})
		if a.hint {
			ui.Column(c).FillWidth().Gap(8).Transition(ui.ElementTransition{Enter: &ui.Motion{Y: 8}, Exit: &ui.Motion{Y: -8}}).Children(func() {
				if shanten(a.hand) == -1 {
					body(c, "The shape is complete. Tsumo wins this closed hand.")
				} else {
					best, s := bestDiscards(a.hand)
					var names []string
					for _, t := range best {
						names = append(names, tileName(t))
					}
					body(c, "Suggested discards: "+strings.Join(names, ", "))
					if s == 0 {
						secondary(c, "This leaves tenpai. You may declare riichi with that discard.")
					} else {
						secondary(c, fmt.Sprintf("Shanten %d after discarding. Suggestions compare shape, not safety or hand value.", s))
					}
				}
			})
		}
		if a.busy && a.scorePage == "Practice" {
			body(c, "Calculating with Agari…")
		}
		if a.scoreResult != "" && a.scorePage == "Practice" {
			panel(c).FillWidth().Transition(ui.ElementTransition{Enter: &ui.Motion{Y: 18, Scale: .96}}).Children(func() {
				heading(c, "Tsumo")
				body(c, a.scoreResult)
				secondary(c, "Hand value excludes the separate riichi deposit.")
			})
		}
	})
	ui.Row(c).FillWidth().Gap(24).Children(func() {
		ui.Button(c, "Start with the tiles").OnClick(func() { a.navigate("Learn"); a.selectLesson(0) })
		ui.Button(c, "Explore five blocks").OnClick(func() { a.navigate("Learn"); a.selectLesson(4) })
		ui.Spacer(c).Grow(1)
		ui.Text(c, a.progress()).FontSize(14).TextColor(muted)
	})
}

// The felt's grain is a fixed pattern of faint 1 DIP dots that repeats
// every 77 by 66 DIPs. Painting each dot was some five thousand fills a
// frame; a bitmap of 4 by 4 repeats covers the felt in a dozen images.
const grainWidth, grainHeight = 4 * 77, 4 * 66

var feltGrain = sync.OnceValue(func() *ui.Bitmap {
	img := image.NewNRGBA(image.Rect(0, 0, grainWidth, grainHeight))
	for y := 4; y < grainHeight; y += 6 {
		for x := 4; x < grainWidth; x += 7 {
			if (x*17+y*31)%11 < 4 {
				img.SetNRGBA(x, y, color.NRGBA{R: 228, G: 239, B: 209, A: 9})
			}
		}
	}
	// A dot stays a square at any display scale.
	return ui.NewBitmap(img).Pixelated()
})

func feltTexture(p *ui.Painter, r ui.Rect) {
	grain := feltGrain()
	p.Clip(r, 16, func() {
		for y := float32(0); y < r.H; y += grainHeight {
			for x := float32(0); x < r.W; x += grainWidth {
				p.Image(grain, ui.Rect{X: r.X + x, Y: r.Y + y, W: grainWidth, H: grainHeight}, ui.FillBox)
			}
		}
	})
}

// pondLayout picks the discard tile width and the row width for n discards
// across the felt's inner width. The 40 DIP size keeps its narrow, centered
// rows; smaller sizes use the felt's full width.
func pondLayout(n int, felt float32) (size, rowWidth float32) {
	const budget, gap = 160, 4
	for _, s := range []float32{40, 30, 26, 22} {
		w := felt
		if s == 40 {
			w = min(felt, 580)
		}
		perRow := max(1, int((w+gap)/(s+gap)))
		rows := (n + perRow - 1) / perRow
		if float32(rows)*(s*1.48+gap) <= budget {
			return s, w
		}
	}
	return 22, felt
}

// wallPainter draws one stack per seven tiles left, up to 18, so the wall
// shrinks as the hand goes on and is gone at an exhaustive draw.
func wallPainter(remaining int) func(*ui.Painter, ui.Rect) {
	return func(p *ui.Painter, r ui.Rect) { drawWall(p, r, min(18, (remaining+6)/7)) }
}
func drawWall(p *ui.Painter, r ui.Rect, stacks int) {
	width := float32(27)
	start := r.X + (r.W-float32(stacks)*(width+2))/2
	for i := 0; i < stacks; i++ {
		x := start + float32(i)*(width+2)
		y := r.Y + 4
		p.Shadow(ui.Rect{X: x, Y: y, W: width, H: 31}, 3, 0, 3, 3, 0, ui.RGBA(0, 0, 0, .2))
		p.Image(tileFront, ui.Rect{X: x, Y: y + 4, W: width, H: 36}, ui.Contain)
		p.Image(tileBack, ui.Rect{X: x, Y: y, W: width, H: 36}, ui.Contain)
	}
}
func drawCompass(p *ui.Painter, r ui.Rect) {
	p.Stroke(ui.Rect{X: r.X + 4, Y: r.Y + 4, W: r.W - 8, H: r.H - 8}, ui.RGBA(202, 184, 133, .3), 16, .5)
	p.Text(r.X+32, r.Y+25, "東", 32, ui.Hex("#D6C38F"))
	p.Text(r.X+19, r.Y+67, "East round", 14, ui.Hex("#ACC2AA"))
}
