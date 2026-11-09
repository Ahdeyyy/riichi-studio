package main

import (
	"github.com/egoist/mygo/ui"
	"image/color"
	"testing"
	"time"
)

func TestTileGuideDoesNotMoveArtwork(t *testing.T) {
	a := newApp()
	tt := ui.NewTester(a.view, 1280, 900)
	tile, ok := tt.Find("East tile 14")
	if !ok {
		t.Fatal("tile missing")
	}
	guide, ok := tt.Find("E")
	if !ok || guide.X < tile.X || guide.Y < tile.Y || guide.X+guide.W > tile.X+tile.W/2 || guide.Y+guide.H > tile.Y+tile.H/3 {
		t.Fatalf("guide is not inside the tile's top-left corner: %v, %v", guide, tile)
	}
	var pixels []color.RGBA
	beforeImage := tt.Image()
	for y := int(tile.Y) + 1; y < int(tile.Y+tile.H)-1; y++ {
		for x := int(tile.X) + 1; x < int(tile.X+tile.W)-1; x++ {
			pixels = append(pixels, beforeImage.RGBAAt(x, y))
		}
	}
	if err := tt.Click("Tile labels"); err != nil {
		t.Fatal(err)
	}
	after, ok := tt.Find("East tile 14")
	if !ok || after != tile {
		t.Fatalf("tile bounds moved: %v -> %v", tile, after)
	}
	at := 0
	afterImage := tt.Image()
	for y := int(tile.Y) + 1; y < int(tile.Y+tile.H)-1; y++ {
		for x := int(tile.X) + 1; x < int(tile.X+tile.W)-1; x++ {
			insideOverlay := float32(x) >= guide.X-1 && float32(x) <= guide.X+guide.W+1 && float32(y) >= guide.Y-1 && float32(y) <= guide.Y+guide.H+1
			if !insideOverlay && pixels[at] != afterImage.RGBAAt(x, y) {
				t.Fatalf("tile artwork changed at %d,%d when guide toggled", x, y)
			}
			at++
		}
	}
	if err := tt.Click("Tile labels"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("East tile 14"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(600 * time.Millisecond)
	tt.Frame()
	selectedGuide, _ := tt.Find("E")
	selectedTile, _ := tt.Find("East tile 14")
	if selectedGuide.X-selectedTile.X != guide.X-tile.X || selectedGuide.Y-selectedTile.Y != guide.Y-tile.Y {
		t.Fatalf("guide did not stay anchored to the selected tile: %v, %v", selectedGuide, selectedTile)
	}
}

func TestCustomWindowControls(t *testing.T) {
	a := newApp()
	minimized, closed, maximized := false, false, false
	a.chrome = windowActions{minimize: func() { minimized = true }, close: func() { closed = true }, toggleMaximize: func() { maximized = !maximized }, maximized: func() bool { return maximized }}
	tt := ui.NewTester(a.view, 1280, 900)
	if err := tt.Click("Minimize"); err != nil {
		t.Fatal(err)
	}
	if !minimized {
		t.Fatal("minimize action missing")
	}
	if err := tt.Click("Maximize"); err != nil {
		t.Fatal(err)
	}
	if !maximized {
		t.Fatal("maximize action missing")
	}
	if _, ok := tt.Find("Restore"); !ok {
		t.Fatal("restore control missing after maximizing")
	}
	if err := tt.Click("Restore"); err != nil {
		t.Fatal(err)
	}
	if maximized {
		t.Fatal("restore action missing")
	}
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if !closed {
		t.Fatal("close action missing")
	}
}

func TestActionsUseContentWidth(t *testing.T) {
	a := newApp()
	tt := ui.NewTester(a.view, 1280, 1000)
	for _, tc := range []struct {
		page, button string
		max          float32
	}{{"Scoring lab", "Calculate score", 200}, {"Wait trainer", "Next exercise", 200}, {"Wait trainer", "Read about waits", 220}, {"Learn", "Only 1 and 9 in each suit", 300}} {
		a.navigate(tc.page)
		if tc.page == "Learn" {
			a.lessonTab = 2
		}
		tt.Frame()
		r, ok := tt.Find(tc.button)
		if !ok {
			t.Fatalf("missing %s", tc.button)
		}
		if r.W > tc.max {
			t.Errorf("%s stretches to %.0fpx", tc.button, r.W)
		}
	}
}

func TestUpdateButton(t *testing.T) {
	a := newApp()
	tt := ui.NewTester(a.view, 1280, 900)
	if _, ok := tt.Find("Check for updates"); ok {
		t.Fatal("update button shown without the updater")
	}
	checks := 0
	a.updates.check = func() { checks++ }
	tt.Frame()
	if err := tt.Click("Check for updates"); err != nil {
		t.Fatal(err)
	}
	a.updates.available = "0.2.0"
	tt.Frame()
	if err := tt.Click("Update to 0.2.0"); err != nil {
		t.Fatal(err)
	}
	if checks != 2 {
		t.Fatalf("checks = %d, want 2", checks)
	}
}
