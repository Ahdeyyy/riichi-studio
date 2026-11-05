package main

import (
	"github.com/egoist/mygo/ui"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestRenderExpandedCourse(t *testing.T) {
	if os.Getenv("RIICHI_RENDER") == "" {
		t.Skip("set RIICHI_RENDER=1 for layout review")
	}
	_ = os.MkdirAll(".cache", 0755)
	render := func(a *app, name string, width, height int) {
		t.Helper()
		f, err := os.Create(filepath.Join(".cache", name+".png"))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err = png.Encode(f, ui.Render(a.view, width, height, 1)); err != nil {
			t.Fatal(err)
		}
	}
	a := newApp()
	a.page = "Learn"
	a.selectLesson(27)
	a.lessonTab = 1
	render(a, "course-examples", 1280, 900)
	a.lessonTab = 2
	render(a, "course-checks-small", 980, 800)
	a.page = "Exercises"
	a.exerciseMode = 2
	render(a, "discard-exercises", 1280, 900)
	a.exerciseMode = 4
	render(a, "scoring-exercises", 980, 800)
	a.exerciseMode = 0
	a.exercise = len(waitDrills) - 1
	render(a, "honor-exercises", 980, 800)
	a.exerciseMode = 6
	a.startReview()
	render(a, "mixed-review", 1280, 900)
}
