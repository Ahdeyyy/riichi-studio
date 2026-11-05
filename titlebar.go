package main

import "github.com/egoist/mygo/ui"

type windowActions struct {
	minimize, toggleMaximize, close func()
	maximized                       func() bool
}

func (a *app) titlebar(c *ui.Context) {
	ui.Row(c).FillWidth().Height(40).Shrink(0).Background(chrome).Children(func() {
		ui.Row(c).Grow(1).FillHeight().Padding(0, 16).Gap(10).AlignItems(ui.Center).DragWindow().Children(func() {
			ui.Image(c, appIconBitmap).Size(22, 22)
			ui.Text(c, "Riichi Studio").FontSize(14).TextColor(ink)
		})
		windowButton(c, "Minimize", false, func() {
			if a.chrome.minimize != nil {
				a.chrome.minimize()
			}
		}, func(p *ui.Painter, r ui.Rect) { x, y := r.X+r.W/2, r.Y+r.H/2; p.Line(x-5, y+3, x+5, y+3, 1, ink) })
		maximized := a.chrome.maximized != nil && a.chrome.maximized()
		label := "Maximize"
		if maximized {
			label = "Restore"
		}
		windowButton(c, label, true, func() {
			if a.chrome.toggleMaximize != nil {
				a.chrome.toggleMaximize()
			}
		}, func(p *ui.Painter, r ui.Rect) {
			x, y := r.X+r.W/2, r.Y+r.H/2
			if maximized {
				p.Stroke(ui.Rect{X: x - 3, Y: y - 6, W: 8, H: 8}, ink, 0, 1)
				p.Fill(ui.Rect{X: x - 5, Y: y - 3, W: 8, H: 8}, chrome, 0)
				p.Stroke(ui.Rect{X: x - 5, Y: y - 3, W: 8, H: 8}, ink, 0, 1)
			} else {
				p.Stroke(ui.Rect{X: x - 5, Y: y - 5, W: 10, H: 10}, ink, 0, 1)
			}
		})
		windowButton(c, "Close", false, func() {
			if a.chrome.close != nil {
				a.chrome.close()
			}
		}, func(p *ui.Painter, r ui.Rect) {
			x, y := r.X+r.W/2, r.Y+r.H/2
			p.StrokePolyline([][2]float32{{x - 5, y - 5}, {x + 5, y + 5}}, 1, ink)
			p.StrokePolyline([][2]float32{{x + 5, y - 5}, {x - 5, y + 5}}, 1, ink)
		})
	})
}

func windowButton(c *ui.Context, label string, maximize bool, click func(), draw func(*ui.Painter, ui.Rect)) {
	c.Key("window-" + label)
	b := ui.ButtonBase(c).Label(label).Tooltip(label).Size(46, 40).OnClick(click)
	if maximize {
		b.MaximizeButton()
	}
	if b.Hovered() {
		if label == "Close" {
			b.Background(ui.Hex("#B6473C"))
		} else {
			b.Background(pillHover)
		}
	}
	b.Draw(draw)
}

// Controls in a column otherwise stretch to its full width.
func actionButton(c *ui.Context, label string) ui.Element {
	return ui.Button(c, label).AlignSelf(ui.Start).Height(38).Padding(8, 16).FontSize(14)
}
func primaryAction(c *ui.Context, label string) ui.Element {
	return ui.PrimaryButton(c, label).AlignSelf(ui.Start).Height(38).Padding(8, 16).FontSize(14)
}
