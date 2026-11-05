package main

import (
	"github.com/egoist/mygo/ui"
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
