package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"strings"
	"time"
)

func body(c *ui.Context, s string)      { ui.Text(c, s).FontSize(16).TextColor(ink).FillWidth() }
func secondary(c *ui.Context, s string) { ui.Text(c, s).FontSize(14).TextColor(muted).FillWidth() }
func heading(c *ui.Context, s string)   { ui.Text(c, s).FontSize(24).Bold().TextColor(ink) }

func (a *app) view(c *ui.Context) {
	c.SetTheme(appTheme())
	ui.Column(c).Fill().Gradient(pageTop, pageBottom, 180).Children(func() {
		a.titlebar(c)
		ui.Row(c).FillWidth().Padding(4, 28, 10).Gap(6).AlignItems(ui.Center).Background(chrome).BorderWidth(0, 0, 1, 0).BorderColor(hairline).Children(func() {
			for _, page := range []string{"Practice", "Learn", "Exercises", "Scoring lab", "Reference"} {
				c.Key(page)
				pill(c, page, a.page == page, func() { a.navigate(page) })
			}
			ui.Spacer(c).Grow(1)
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
		ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() { heading(c, "Your practice table"); ui.Text(c, "Solo hand").FontSize(14).TextColor(muted) })
		ui.Row(c).Gap(12).Children(func() {
			ui.Button(c, "Restart guide").Disabled(a.busy).OnClick(func() { a.resetTable(false) })
			ui.Button(c, "New random hand").Disabled(a.busy).OnClick(func() { a.resetTable(true) })
		})
	})
	ui.Column(c).FillWidth().Padding(9).Radius(24).Gradient(ui.Hex("#685244"), ui.Hex("#2C2723"), 100).Border(.5, ui.Hex("#8A7561")).Shadow(0, 12, 22, 0, ui.RGBA(0, 0, 0, .23)).Children(func() {
		ui.Column(c).FillWidth().Padding(22, 24, 24).Gap(16).Radius(16).Gradient(ui.Hex("#245D50"), ui.Hex("#184236"), 100).InsetShadow(0, 3, 9, 0, ui.RGBA(0, 0, 0, .3)).DrawOver(feltTexture).Children(func() {
			ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
				ui.Text(c, fmt.Sprintf("Draw %d", a.turn)).FontSize(14).TextColor(feltText)
				ui.Text(c, fmt.Sprintf("%d tiles remaining", len(a.wall))).FontSize(14).TextColor(feltText)
			})
			ui.Box(c).FillWidth().Height(48).Draw(drawWall)
			ui.Column(c).FillWidth().MinHeight(148).Center().Gap(12).Children(func() {
				if len(a.discards) == 0 {
					ui.Box(c).Size(100, 100).Draw(drawCompass)
				} else {
					ui.Text(c, "Discards").FontSize(14).TextColor(feltLabel)
					ui.Row(c).Wrap().Gap(4).Justify(ui.Center).MaxWidth(580).Children(func() {
						for i, t := range a.discards {
							c.Key(a.discardIDs[i])
							a.tableTile(c, t, false, tileName(t), func() {}, 40, a.discardIDs[i], 0)
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

func feltTexture(p *ui.Painter, r ui.Rect) {
	// Fixed grain has no per-frame randomness or moving noise.
	for y := float32(4); y < r.H; y += 6 {
		for x := float32(4); x < r.W; x += 7 {
			v := int(x*17+y*31) % 11
			if v < 4 {
				p.Fill(ui.Rect{X: r.X + x, Y: r.Y + y, W: 1, H: 1}, ui.RGBA(228, 239, 209, .035), 0)
			}
		}
	}
}
func drawWall(p *ui.Painter, r ui.Rect) {
	width := float32(27)
	start := r.X + (r.W-18*(width+2))/2
	for i := 0; i < 18; i++ {
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
