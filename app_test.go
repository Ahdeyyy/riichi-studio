package main

import (
	"github.com/egoist/mygo/ui"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestShapeAnalysis(t *testing.T) {
	for _, tc := range []struct {
		hand string
		want int
	}{{"123m456p234678s55p", -1}, {"123m456p23467s55p", 0}, {"1122m3344p5566s77z", -1}, {"19m19p19s1234567z", 0}, {"119m19p19s1234567z", -1}} {
		if got := shanten(mustTiles(tc.hand)); got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.hand, got, tc.want)
		}
	}
	w := waits(mustTiles("123m456p23467s55p"))
	if len(w) != 2 || w[0] != 22 || w[1] != 25 {
		t.Fatalf("waits: %v", w)
	}
	for _, s := range []string{"11111m", "8z", "12", "0m", "xyz"} {
		if _, e := parseTiles(s); e == nil {
			t.Errorf("accepted %q", s)
		}
	}
}
func TestGuidedTable(t *testing.T) {
	a := newApp()
	tt := ui.NewTester(a.view, 1280, 860)
	if e := tt.Click("East tile 14"); e != nil {
		t.Fatal(e)
	}
	if e := tt.Click("Declare riichi"); e != nil {
		t.Fatal(e)
	}
	if !a.riichi || len(a.discards) != 1 || a.discards[0] != 27 || shanten(a.hand) != -1 {
		t.Fatalf("riichi %v hand %v discards %v", a.riichi, a.hand, a.discards)
	}
	if len(a.hand) != 14 || len(a.wall)+len(a.hand)+len(a.discards) != 136 {
		t.Fatal("tile conservation failed")
	}
	// A locked hand must not let the player select an old tile.
	if e := tt.Click("1 characters tile 1"); e != nil {
		t.Fatal(e)
	}
	if a.selected != 13 {
		t.Fatalf("riichi selection changed: %d", a.selected)
	}
}
func TestLearningAndNavigation(t *testing.T) {
	a := newApp()
	a.savePath = filepath.Join(t.TempDir(), "progress.json")
	tt := ui.NewTester(a.view, 1280, 1000)
	if e := tt.Click("Learn"); e != nil {
		t.Fatal(e)
	}
	if e := tt.Click("Checks"); e != nil {
		t.Fatal(e)
	}
	if e := tt.Click("All winds and dragons"); e != nil {
		t.Fatal(e)
	}
	if a.completed[0] {
		t.Fatal("incorrect answer completed lesson")
	}
	if e := tt.Click("Only 1 and 9 in each suit"); e != nil {
		t.Fatal(e)
	}
	if !a.lessonAnswers[lessonKey(0, 0)] || a.completed[0] {
		t.Fatal("first check must be recorded without completing the whole lesson")
	}
	if e := tt.Click("Next check"); e != nil {
		t.Fatal(e)
	}
	if e := tt.Click("No, honors do not form runs"); e != nil {
		t.Fatal(e)
	}
	if !a.completed[0] {
		t.Fatal("lesson did not complete after both checks")
	}
	if _, e := os.ReadFile(a.savePath); e != nil {
		t.Fatal(e)
	}
	for _, page := range []string{"Exercises", "Scoring lab", "Reference", "Practice", "Learn"} {
		if e := tt.Click(page); e != nil {
			t.Fatal(e)
		}
		tt.Texts()
	}
}
func TestAgariScoring(t *testing.T) {
	if e := prepareEngine(); e != nil {
		t.Fatal(e)
	}
	for _, tc := range []struct {
		in   scoreInput
		want string
	}{
		{scoreInput{"123m456p234678s55p", "8s", "", true, true, false}, "2,700 points"},
		{scoreInput{"123m456p234678s55p", "8s", "", false, true, false}, "2,000 points"},
		{scoreInput{"123m456p234678s55p", "8s", "", true, true, true}, "3,900 points"},
	} {
		s, e := scoreHand(tc.in)
		if e != nil || !strings.Contains(s, tc.want) {
			t.Fatalf("score %q err %v; wanted %s", s, e, tc.want)
		}
	}
	if _, e := scoreHand(scoreInput{hand: "123m456p234678s55p", win: "9m"}); e == nil {
		t.Fatal("accepted absent winning tile")
	}
	if _, e := scoreHand(scoreInput{hand: "123m456p789s22255p", win: "2p"}); e == nil {
		t.Fatal("accepted a hand without yaku")
	}
	s, e := scoreHand(scoreInput{"123m456p234678s55p", "8s", "4p", true, true, false})
	if e != nil || !strings.Contains(s, "Dora bonuses   3 han") {
		t.Fatalf("dora score %s: %v", s, e)
	}
}
func TestWaitTrainer(t *testing.T) {
	a := newApp()
	a.page = "Wait trainer"
	tt := ui.NewTester(a.view, 1280, 900)
	if e := tt.Click("Choose 4 bamboo"); e != nil {
		t.Fatal(e)
	}
	if a.answered {
		t.Fatal("incorrect tile accepted")
	}
	if e := tt.Click("Choose 3 bamboo"); e != nil {
		t.Fatal(e)
	}
	if !a.answered {
		t.Fatal("correct tile rejected")
	}
	if e := tt.Click("Next exercise"); e != nil {
		t.Fatal(e)
	}
	if a.exercise != 1 || a.answered {
		t.Fatal("exercise did not advance")
	}
}
func TestTileMotionAndIdentity(t *testing.T) {
	a := newApp()
	tt := ui.NewTester(a.view, 1280, 900)
	start, _ := tt.Find("East tile 14")
	if err := tt.Click("East tile 14"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(600 * time.Millisecond)
	tt.Frame()
	lifted, _ := tt.Find("East tile 14")
	if lifted.Y > start.Y-12 {
		t.Fatalf("selected tile did not lift: %v -> %v", start, lifted)
	}
	id := a.handIDs[13]
	if err := tt.Click("Discard tile"); err != nil {
		t.Fatal(err)
	}
	if a.discardIDs[0] != id || len(a.handIDs) != 14 {
		t.Fatal("physical tile identity lost during discard")
	}
	tt.Advance(150 * time.Millisecond)
	tt.Frame()
	moving, ok := tt.Find("East")
	if !ok {
		t.Fatal("discard absent during motion")
	}
	tt.Advance(time.Second)
	tt.Frame()
	landed, _ := tt.Find("East")
	if moving.Y <= landed.Y {
		t.Fatalf("discard did not travel toward its destination: %v -> %v", moving, landed)
	}
	for _, page := range []string{"Learn", "Reference", "Practice"} {
		a.navigate(page)
		tt.Frame()
		tt.Advance(500 * time.Millisecond)
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if err := tt.Click("Restart guide"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if a.selected != -1 || len(a.discards) != 0 {
		t.Fatal("restart did not reset selection and discards")
	}
}
func TestRenderViews(t *testing.T) {
	if os.Getenv("RIICHI_RENDER") == "" {
		t.Skip("set RIICHI_RENDER=1 for layout snapshots")
	}
	_ = os.MkdirAll(".cache", 0755)
	for _, page := range []string{"Practice", "Learn", "Wait trainer", "Scoring lab", "Reference"} {
		a := newApp()
		a.page = page
		img := ui.Render(a.view, 1280, 900, 1)
		f, e := os.Create(filepath.Join(".cache", strings.ReplaceAll(page, " ", "-")+".png"))
		if e != nil {
			t.Fatal(e)
		}
		e = png.Encode(f, img)
		f.Close()
		if e != nil {
			t.Fatal(e)
		}
	}
	a := newApp()
	img := ui.Render(a.view, 980, 740, 1)
	f, err := os.Create(".cache/Practice-small.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	a.showTileLabels = false
	plain, err := os.Create(".cache/Practice-unlabelled.png")
	if err != nil {
		t.Fatal(err)
	}
	defer plain.Close()
	if err = png.Encode(plain, ui.Render(a.view, 1280, 900, 1)); err != nil {
		t.Fatal(err)
	}
}

func TestTileLabelsAndAssets(t *testing.T) {
	a := newApp()
	a.savePath = filepath.Join(t.TempDir(), "progress.json")
	for tile, b := range tileImages {
		if b == nil {
			t.Fatalf("missing image for %s", tileName(tile))
		}
		w, h := b.Size()
		if w != tileBitmapWidth || h != tileBitmapHeight {
			t.Fatalf("%s image size %dx%d", tileName(tile), w, h)
		}
	}
	tt := ui.NewTester(a.view, 1280, 900)
	if _, ok := tt.Find("E"); !ok {
		t.Fatal("East corner label missing")
	}
	if err := tt.Click("Tile labels"); err != nil {
		t.Fatal(err)
	}
	if a.showTileLabels {
		t.Fatal("tile labels did not turn off")
	}
	if _, ok := tt.Find("E"); ok {
		t.Fatal("East corner label remained visible")
	}
	b := newApp()
	b.savePath = a.savePath
	b.loadPreferences()
	if b.showTileLabels {
		t.Fatal("label preference was not restored")
	}
	if err := tt.Click("Tile labels"); err != nil {
		t.Fatal(err)
	}
	b.loadPreferences()
	if !b.showTileLabels {
		t.Fatal("enabled preference was not saved")
	}
	if tileCornerLabel(31) != "Wh" || tileCornerLabel(29) != "W" {
		t.Fatal("white dragon and west labels are ambiguous")
	}
}
