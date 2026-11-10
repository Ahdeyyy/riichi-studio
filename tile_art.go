package main

import (
	"embed"
	"fmt"
	"github.com/egoist/mygo/ui"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"sync"
	"sync/atomic"
	"time"
)

// The original 600 x 800 PNG exports are bundled unchanged. See SOURCES.md.
//
//go:embed resources/tiles/*.png
var tileAssets embed.FS
var tileImages [34]*ui.Bitmap
var tileFront, tileBack *ui.Bitmap
var tileAssetsOnce sync.Once

// Tiles show at most 64 DIPs wide, so their bitmaps are kept at three
// pixels a DIP: sharp on a 300% display, and a tenth of the 1.9 MB each
// full-size export holds in memory, and again on the GPU.
const tileBitmapWidth, tileBitmapHeight = 192, 256

func tileAssetName(t int) string {
	if t < 27 {
		return fmt.Sprintf("%s%d", []string{"Man", "Pin", "Sou"}[t/9], t%9+1)
	}
	return []string{"Ton", "Nan", "Shaa", "Pei", "Haku", "Hatsu", "Chun"}[t-27]
}
func loadTileAssets() {
	tileAssetsOnce.Do(func() {
		// Decoding and scaling 36 exports takes most of a second on one
		// core, before the window can show: use them all.
		var group sync.WaitGroup
		var failed atomic.Pointer[error]
		decode := func(name string, into **ui.Bitmap) {
			group.Go(func() {
				b, err := decodeTile(name)
				if err != nil {
					failed.CompareAndSwap(nil, &err)
					return
				}
				*into = b
			})
		}
		for t := range tileImages {
			decode(tileAssetName(t), &tileImages[t])
		}
		decode("Front", &tileFront)
		decode("Back", &tileBack)
		group.Wait()
		if err := failed.Load(); err != nil {
			panic(*err)
		}
	})
}
func decodeTile(name string) (*ui.Bitmap, error) {
	f, err := tileAssets.Open("resources/tiles/" + name + ".png")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	dst := image.NewRGBA(image.Rect(0, 0, tileBitmapWidth, tileBitmapHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)
	return ui.NewBitmap(dst), nil
}

var tileCornerLabels = func() (labels [34]string) {
	for t := range labels {
		if t < 27 {
			labels[t] = fmt.Sprint(t%9 + 1)
		} else {
			labels[t] = []string{"E", "S", "W", "N", "Wh", "G", "R"}[t-27]
		}
	}
	return
}()

func tileCornerLabel(t int) string { return tileCornerLabels[t] }

type physicalTile struct{ ID int }

var (
	labelInk     = ui.Hex("#5B5145")
	labelHalo    = ui.Hex("#FFFFF8")
	selectedEdge = ui.Hex("#E8C87F")
	tileShadow   = ui.RGBA(0, 0, 0, .32)
	tileSpring   = ui.Animation{Spring: ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: .12}}
)

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
	b.AnimatePose(ui.Pose{Y: lift, Rotate: rotation}, tileSpring)
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
		p.Shadow(face, 7, 1, 5, 8, 0, tileShadow)
		// Both the tile surface and the markings come from the supplied PNGs.
		p.Image(tileBack, ui.Rect{X: r.X, Y: r.Y + 5, W: r.W, H: r.H - 5}, ui.Contain)
		p.Image(tileFront, face, ui.Contain)
		art := ui.Rect{X: r.X + r.W*.12, Y: r.Y + face.H*.10, W: r.W * .76, H: face.H * .80}
		p.Image(tileImages[t], art, ui.Contain)
		if selected {
			p.Stroke(face, selectedEdge, 7, 1.5)
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
			// A 1 DIP ivory halo: half of the stroke shows around each glyph,
			// under its fill. One stroke replaces eight offset copies of the
			// label, which were a fifth of the table's painting.
			ui.Text(c, label).Absolute().Left(inset).Top(1).FontSize(size).Bold().TextColor(labelInk).TextStroke(2, labelHalo)
		})
	}
}
func (a *app) tileRow(c *ui.Context, tiles []int) {
	ui.Row(c).Wrap().Gap(6).Children(func() {
		for i, t := range tiles {
			c.Key(i)
			a.tile(c, t, false, tileName(t), func() {})
		}
	})
}
