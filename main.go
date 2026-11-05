package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"log"
	"os"
	"path/filepath"
)

//go:embed resources/icon.png
var appIconPNG []byte

//go:embed resources/icon-window.png
var appWindowIconPNG []byte
var appIconBitmap *ui.Bitmap

type app struct {
	page                                  string
	lesson, exercise                      int
	lessonTab, lessonCheck, lessonExample int
	exerciseMode, reviewIndex             int
	lessonAnswers                         map[string]bool
	exerciseRecords                       map[string]exerciseRecord
	reviewSession                         []reviewQuestion
	reviewCorrect                         int
	reviewMissed                          []reviewQuestion
	feedback                              string
	answered                              bool
	completed                             map[int]bool
	savePath                              string
	hand, discards, wall                  []int
	handIDs, discardIDs                   []int
	dealID, nextTileID                    int
	selected, turn                        int
	riichi, won, hint                     bool
	tableMessage                          string
	notation, winning, indicator          string
	tsumo, declared, dealer               bool
	scoreResult                           string
	scorePage                             string
	busy                                  bool
	showTileLabels                        bool
	update                                func(func())
	chrome                                windowActions
	devData                               int
}

func newApp() *app {
	loadTileAssets()
	if appIconBitmap == nil {
		var err error
		appIconBitmap, err = ui.DecodeBitmap(appWindowIconPNG)
		if err != nil {
			panic(err)
		}
	}
	a := &app{page: "Practice", completed: map[int]bool{}, selected: -1, notation: "123m456p234678s55p", winning: "8s", tsumo: true, declared: true}
	a.resetTable(false)
	a.showTileLabels = true
	a.lessonAnswers = map[string]bool{}
	a.exerciseRecords = map[string]exerciseRecord{}
	ensureExerciseBank()
	for _, d := range discardDrills {
		analyzeDiscards(d.hand)
	}
	return a
}
func (a *app) loadProgress() {
	dir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	a.savePath = filepath.Join(dir, "riichi-tutor", "progress.json")
	data, err := os.ReadFile(a.savePath)
	if err == nil {
		_ = json.Unmarshal(data, &a.completed)
	}
	if a.completed == nil {
		a.completed = map[int]bool{}
	}
	a.loadPreferences()
	a.loadLearning()
}

func (a *app) preferencesPath() string {
	if a.savePath == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(a.savePath), "preferences.json")
}
func (a *app) loadPreferences() {
	data, err := os.ReadFile(a.preferencesPath())
	if err != nil {
		return
	}
	var preferences struct {
		TileLabels *bool `json:"tileLabels"`
	}
	if json.Unmarshal(data, &preferences) == nil && preferences.TileLabels != nil {
		a.showTileLabels = *preferences.TileLabels
	}
}
func (a *app) savePreferences() {
	path := a.preferencesPath()
	if path == "" {
		return
	}
	data, _ := json.Marshal(struct {
		TileLabels bool `json:"tileLabels"`
	}{a.showTileLabels})
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Printf("Preferences: %v", err)
		return
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		log.Printf("Preferences: %v", err)
	}
}
func (a *app) saveProgress() {
	if a.savePath == "" {
		return
	}
	data, _ := json.Marshal(a.completed)
	if err := os.MkdirAll(filepath.Dir(a.savePath), 0755); err != nil {
		a.feedback += " Progress could not be saved."
		return
	}
	if err := os.WriteFile(a.savePath, data, 0600); err != nil {
		a.feedback += " Progress could not be saved."
	}
}
func main() {
	a := newApp()
	a.loadProgress()
	if err := prepareEngine(); err != nil {
		log.Printf("Scoring engine: %v", err)
	}
	mygo.App.WhenReady(func() {
		w := mygo.NewWindow(mygo.WindowOptions{Title: "Riichi Studio", Width: 1280, Height: 860, MinWidth: 980, MinHeight: 700, StateKey: "main", TitleBarStyle: mygo.TitleBarHidden, TitleBarHeight: 40, HideWindowButtons: true, BackgroundColor: "#13261F", Content: ui.View(a.view)})
		a.chrome = windowActions{minimize: w.Minimize, toggleMaximize: w.ToggleMaximize, close: w.Close, maximized: w.IsMaximized}
		if err := w.SetIcon(appWindowIconPNG); err != nil {
			log.Printf("Window icon: %v", err)
		}
		a.update = func(fn func()) { w.Update(fn) }
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
func (a *app) navigate(page string) {
	a.page = page
	a.feedback = ""
	a.answered = false
	if page == "Learn" {
		a.selectLesson(a.lesson)
	}
}
func (a *app) progress() string {
	// Count true entries: progress.json is parsed as-is and may hold false.
	done := 0
	for _, ok := range a.completed {
		if ok {
			done++
		}
	}
	return fmt.Sprintf("%d of %d lessons complete", done, len(lessons))
}
