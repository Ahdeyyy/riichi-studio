package main

import (
	"embed"
	"fmt"
	"github.com/egoist/mygo/ui"
	"sync"
	"time"
)

// The original 600 x 800 PNG exports are bundled unchanged. See SOURCES.md.
//
//go:embed resources/tiles/*.png
var tileAssets embed.FS
var tileImages [34]*ui.Bitmap
var tileFront, tileBack *ui.Bitmap
var tileAssetsOnce sync.Once

func tileAssetName(t int) string {
	if t < 27 {
		return fmt.Sprintf("%s%d", []string{"Man", "Pin", "Sou"}[t/9], t%9+1)
	}
	return []string{"Ton", "Nan", "Shaa", "Pei", "Haku", "Hatsu", "Chun"}[t-27]
}
func loadTileAssets() {
	tileAssetsOnce.Do(func() {
		decode := func(name string) *ui.Bitmap {
			data, err := tileAssets.ReadFile("resources/tiles/" + name + ".png")
			if err != nil {
				panic(err)
			}
			b, err := ui.DecodeBitmap(data)
			if err != nil {
				panic(err)
			}
			return b
		}
		for t := range tileImages {
			tileImages[t] = decode(tileAssetName(t))
		}
		tileFront = decode("Front")
		tileBack = decode("Back")
	})
}
func tileCornerLabel(t int) string {
	if t < 27 {
		return fmt.Sprint(t%9 + 1)
	}
	return []string{"E", "S", "W", "N", "Wh", "G", "R"}[t-27]
}

type physicalTile struct{ ID int }

func (a *app) tile(c *ui.Context, t int, selected bool, label string, click func()) {
	a.tableTile(c, t, selected, label, click, 60, -1, 0)
}
func (a *app) tableTile(c *ui.Context, t int, selected bool, label string, click func(), width float32, id int, delay time.Duration) {
	b := ui.ButtonBase(c).Label(label).Tooltip(tileName(t)).Size(width, width*1.48).OnClick(click)
	lift, rotation := float32(0), float32(0)
	if b.Hovered() {
		lift = -7
		rotation = -1
	}
	if selected {
		lift = -16
		rotation = 0
	}
	b.AnimatePose(ui.Pose{Y: lift, Rotate: rotation}, ui.Animation{Spring: ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: .12}})
	if id >= 0 {
		entry := &ui.Motion{Y: -35, Rotate: -5, Scale: .92}
		if delay < 0 {
			delay = 0
			entry = &ui.Motion{X: 50, Y: -220, Rotate: 8, Scale: .85}
		}
		anim := ui.Animation{Duration: 340 * time.Millisecond, Ease: ui.EaseOut, Delay: delay}
		b.Transition(ui.ElementTransition{Animation: &anim, Position: true, Size: true, Match: physicalTile{id}, Enter: entry, Exit: &ui.Motion{Y: -80, Scale: .8}})
	}
	showLabels := a.showTileLabels
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		face := ui.Rect{X: r.X, Y: r.Y, W: r.W, H: r.H - 5}
		p.Shadow(face, 7, 1, 5, 8, 0, ui.RGBA(0, 0, 0, .32))
		// Both the tile surface and the markings come from the supplied PNGs.
		p.Image(tileBack, ui.Rect{X: r.X, Y: r.Y + 5, W: r.W, H: r.H - 5}, ui.Contain)
		p.Image(tileFront, face, ui.Contain)
		art := ui.Rect{X: r.X + r.W*.12, Y: r.Y + face.H*.10, W: r.W * .76, H: face.H * .80}
		p.Image(tileImages[t], art, ui.Contain)
		if selected {
			p.Stroke(face, ui.Hex("#E8C87F"), 7, 1.5)
		}
	})
	if showLabels {
		b.Children(func() {
			// Absolute children paint over the tile without taking layout space.
			label := tileCornerLabel(t)
			// 14 DIP from a 56 DIP tile up; smaller pond tiles scale it down
			// so the label stays in its corner.
			size := min(14, max(8, width/4))
			inset := size * 4 / 14
			ui.Text(c, label).Absolute().Left(inset).Top(1).FontSize(size).Bold().TextColor(ui.Hex("#5B5145")).Draw(func(p *ui.Painter, r ui.Rect) {
				// Paint a 1 DIP ivory halo before the text's own foreground.
				// One accessible text element is retained; no duplicate labels.
				span := ui.Span{Text: label, Font: "Segoe UI", Size: size, Weight: 700, Color: ui.Hex("#FFFFF8")}
				for _, offset := range [][2]float32{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-.7, -.7}, {.7, -.7}, {-.7, .7}, {.7, .7}} {
					p.RichText(r.X+offset[0], r.Y+offset[1], 0, span)
				}
			})
		})
	}
}
func (a *app) tileRow(c *ui.Context, tiles []int) {
	ui.Row(c).Wrap().Gap(6).Children(func() {
		for i, t := range tiles {
			c.Key(fmt.Sprintf("tile-%d", i))
			a.tile(c, t, false, tileName(t), func() {})
		}
	})
}
