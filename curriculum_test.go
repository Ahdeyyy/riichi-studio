package main

import (
	"encoding/json"
	"fmt"
	"github.com/egoist/mygo/ui"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCurriculumIntegrity(t *testing.T) {
	a := newApp()
	_ = a
	if len(lessons) != 28 || len(courseOrder()) != 28 {
		t.Fatalf("unexpected course size %d", len(lessons))
	}
	seen := map[int]bool{}
	for _, i := range courseOrder() {
		if seen[i] {
			t.Fatalf("lesson %d repeated", i)
		}
		seen[i] = true
		qs := lessonChecks(i)
		if len(qs) != 2 {
			t.Fatal("every lesson needs two checks")
		}
		for _, q := range qs {
			if q.prompt == "" || q.why == "" || q.correct < 0 || q.correct >= len(q.choices) {
				t.Fatalf("invalid lesson check %d", i)
			}
		}
		tiles, err := parseTiles(lessons[i].example)
		if err != nil || len(tiles) > 14 && i != 0 {
			t.Fatalf("lesson %d sample: %v", i, err)
		}
		ex, ok := workedHands[i]
		if !ok {
			t.Fatalf("lesson %d lacks worked examples", i)
		}
		for _, e := range ex {
			ts, err := parseTiles(e.hand)
			if err != nil || len(ts) > 14 || e.why == "" || e.title == "" {
				t.Fatalf("lesson %d worked example %q invalid: %v", i, e.title, err)
			}
		}
	}
	keys := map[string]bool{}
	for _, r := range reviewBank {
		if keys[r.key] {
			t.Fatalf("duplicate exercise key %s", r.key)
		}
		keys[r.key] = true
		choices := map[string]bool{}
		for _, c := range r.q.choices {
			if c == "" || choices[c] {
				t.Fatalf("duplicate/empty choice in %s", r.key)
			}
			choices[c] = true
		}
		if r.q.correct < 0 || r.q.correct >= len(r.q.choices) {
			t.Fatal("invalid correct index")
		}
	}
	if len(waitDrills) != 79 || len(reviewBank) != 180 {
		t.Fatalf("bank counts: %d block waits, %d total", len(waitDrills), len(reviewBank))
	}
	for _, d := range discardDrills {
		h, err := parseTiles(d.hand)
		if err != nil || len(h) != 14 {
			t.Fatalf("discard hand %s: %d tiles, %v", d.hand, len(h), err)
		}
	}
}

func TestWaitExercisesAgainstCompletion(t *testing.T) {
	newApp()
	for _, d := range waitDrills {
		shape := mustTiles(d.shape)
		for draw := 0; draw < 34; draw++ {
			expected := false
			for _, v := range d.answers {
				if v == draw {
					expected = true
				}
			}
			counts := [34]int{}
			for _, v := range append(append([]int(nil), shape...), draw) {
				counts[v]++
			}
			valid := counts[draw] == 3
			if !valid {
				for i := 0; i < 27; i++ {
					if i%9 <= 6 && counts[i] == 1 && counts[i+1] == 1 && counts[i+2] == 1 {
						valid = true
					}
				}
			}
			if expected != valid {
				t.Fatalf("%s draw %s: %t != %t", d.shape, tileCode(draw), expected, valid)
			}
		}
	}
	for _, q := range questionsFor("Hand waits") {
		h := mustTiles(q.hand)
		if len(h) != 13 || len(waits(h)) == 0 {
			t.Fatalf("bad full-hand wait %s", q.hand)
		}
	}
}

func TestLearningProgressMigrationAndReview(t *testing.T) {
	a := newApp()
	a.savePath = filepath.Join(t.TempDir(), "progress.json")
	a.completed[0] = true
	a.saveProgress()
	a.loadLearning()
	if !a.lessonAnswers[lessonKey(0, 0)] || a.completed[0] {
		t.Fatal("old first check was not preserved with the new check pending")
	}
	a.selectLesson(0)
	if a.lessonCheck != 1 {
		t.Fatal("did not select pending check")
	}
	a.answerLesson(extraChecks[0].correct)
	if !a.completed[0] {
		t.Fatal("lesson did not complete")
	}
	a.recordExercise(reviewBank[0].key, false)
	a.recordExercise(reviewBank[1].key, true)
	b := newApp()
	b.savePath = a.savePath
	b.loadLearning()
	if !b.completed[0] || b.exerciseRecords[reviewBank[0].key].Attempts != 1 {
		t.Fatal("mastery did not survive reload")
	}
	b.startReview()
	if len(b.reviewSession) != 10 || b.reviewSession[0].key != reviewBank[0].key {
		t.Fatal("review did not prioritize the missed item")
	}
	data, err := os.ReadFile(b.learningPath())
	var s learningSave
	if err != nil || json.Unmarshal(data, &s) != nil {
		t.Fatal("invalid saved learning state")
	}
}

func TestScoringExerciseAnswers(t *testing.T) {
	for _, c := range []struct {
		h, fu  int
		dealer bool
		want   int
	}{{1, 30, false, 1000}, {2, 30, false, 2000}, {3, 30, false, 3900}, {4, 30, false, 7700}, {4, 40, false, 8000}, {5, 30, true, 12000}, {2, 25, false, 1600}} {
		if got := ronPoints(c.h, c.fu, c.dealer); got != c.want {
			t.Fatalf("%+v gives %d", c, got)
		}
	}
	if tsumoPayment(3, 20, false) != "1300 from East, 700 from each other" {
		t.Fatal("tsumo rounding wrong")
	}
	if err := prepareEngine(); err != nil {
		t.Fatal(err)
	}
	s, err := scoreHand(scoreInput{"234m456p234678s55p", "8s", "", false, false, false})
	if err != nil || !strings.Contains(s, "2000 points") {
		t.Fatalf("placement example scoring %s %v", s, err)
	}
	if shanten(mustTiles("123m456p24s67s55p1z")) != 1 {
		t.Fatal("worked one-shanten example wrong")
	}
	if shanten(mustTiles("1239m4567p2466s1z")) != 2 {
		t.Fatal("push/fold example is not two-shanten")
	}
	if len(waits(mustTiles("11m22p33s445566z7z"))) != 1 {
		t.Fatal("seven pairs sample is not tenpai")
	}
}

func TestExpandedExerciseFlows(t *testing.T) {
	a := newApp()
	a.page = "Exercises"
	tt := ui.NewTester(a.view, 1280, 1000)
	for _, mode := range []int{1, 3, 4, 5} {
		a.exerciseMode = mode
		a.exercise = 0
		a.answered = false
		tt.Frame()
		topic := map[int]string{1: "Hand waits", 3: "Yaku", 4: "Scoring", 5: "Decisions"}[mode]
		q := questionsFor(topic)[0]
		if err := tt.Click(q.q.choices[q.q.correct]); err != nil {
			t.Fatal(err)
		}
		if !a.answered || !a.exerciseRecords[q.key].LastCorrect {
			t.Fatalf("%s answer not recorded", topic)
		}
		if err := tt.Click("Next exercise"); err != nil {
			t.Fatal(err)
		}
	}
	a.exerciseMode = 2
	a.exercise = 0
	a.answered = false
	tt.Frame()
	best := analyzeDiscards(discardDrills[0].hand)[0]
	if err := tt.Click("Discard " + tileName(best.tile)); err != nil {
		t.Fatal(err)
	}
	if !a.exerciseRecords["discard:"+discardDrills[0].hand].LastCorrect {
		t.Fatal("discard answer not accepted")
	}
	a.exerciseMode = 0
	a.exercise = len(waitDrills) - 1
	a.answered = false
	tt.Frame()
	if err := tt.Click("Choose Red dragon"); err != nil {
		t.Fatal(err)
	}
	if !a.answered {
		t.Fatal("honor wait not accepted")
	}
	t.Log(fmt.Sprintf("%d lessons, %d worked examples, %d exercises", len(lessons), 2*len(workedHands), len(reviewBank)))
}

func TestMixedReviewAndRetry(t *testing.T) {
	a := newApp()
	a.page = "Exercises"
	a.exerciseMode = 6
	a.reviewSession = append([]reviewQuestion(nil), reviewBank[:2]...)
	tt := ui.NewTester(a.view, 1280, 1000)
	for i := 0; i < 2; i++ {
		q := a.reviewSession[i]
		choice := q.q.correct
		if i == 0 {
			choice = (choice + 1) % len(q.q.choices)
		}
		if err := tt.Click(q.q.choices[choice]); err != nil {
			t.Fatal(err)
		}
		if err := tt.Click("Next exercise"); err != nil {
			t.Fatal(err)
		}
	}
	if a.reviewCorrect != 1 || len(a.reviewMissed) != 1 || !tt.HasText("1 of 2 correct") {
		t.Fatal("session score or missed-question list wrong")
	}
	if err := tt.Click("Retry mistakes"); err != nil {
		t.Fatal(err)
	}
	if len(a.reviewSession) != 1 || a.reviewIndex != 0 || a.answered {
		t.Fatal("mistake retry was not reset")
	}
	q := a.reviewSession[0]
	if err := tt.Click(q.q.choices[q.q.correct]); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Next exercise"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("1 of 1 correct") {
		t.Fatal("retry score wrong")
	}
}
