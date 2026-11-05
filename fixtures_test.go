package main

import (
	"github.com/egoist/mygo/ui"
	"strings"
	"testing"
)

func TestFixturesRenderEveryPage(t *testing.T) {
	for _, name := range fixtureNames {
		for _, page := range []string{"Practice", "Learn", "Exercises", "Scoring lab", "Reference"} {
			a := newApp()
			tt := ui.NewTester(a.view, 980, 700)
			a.page = page
			a.applyFixture(name)
			a.navigate(page)
			tt.Frame()
			if _, ok := tt.Find(page); !ok {
				t.Fatalf("%s / %s: navigation missing", name, page)
			}
		}
	}
}

func TestWorstCaseText(t *testing.T) {
	for _, tc := range [][2]string{
		{plural(1, "tile", "tiles"), "1 tile"},
		{plural(0, "tile", "tiles"), "0 tiles"},
		{points(48000), "48,000"},
		{points(700), "700"},
		{points(1000000), "1,000,000"},
		{engineMessage([]byte("  \n")), "no details were reported."},
		{engineMessage([]byte("bad hand\npanic: stack trace")), "bad hand"},
	} {
		if tc[0] != tc[1] {
			t.Errorf("got %q, want %q", tc[0], tc[1])
		}
	}
	if long := engineMessage([]byte(strings.Repeat("x", 400))); len([]rune(long)) != 160 {
		t.Errorf("engine message not capped: %d runes", len([]rune(long)))
	}
	a := newApp()
	a.completed = map[int]bool{0: true, 1: false, 2: false}
	if got := a.progress(); got != "1 of 28 lessons complete" {
		t.Errorf("progress counts false entries: %q", got)
	}
}

func TestPondKeepsTableHeight(t *testing.T) {
	for _, felt := range []float32{858, 1158} {
		for n := 1; n <= 122; n++ {
			size, row := pondLayout(n, felt)
			perRow := int((row + 4) / (size + 4))
			if rows := (n + perRow - 1) / perRow; float32(rows)*(size*1.48+4) > 160 {
				t.Fatalf("%d discards at %.0f wide overflow: %d rows of %.0f", n, felt, rows, size)
			}
		}
	}
}
