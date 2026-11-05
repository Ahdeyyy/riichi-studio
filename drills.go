package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
)

type waitDrill struct {
	shape   string
	answers []int
	why     string
}

var waitDrills = []waitDrill{
	{"45s", []int{20, 23}, "3 or 6 bamboo finishes the run. A two-sided wait has up to eight available copies."},
	{"46p", []int{13}, "Only 5 dots fills the gap: 4–5–6. This is a closed wait (kanchan)."},
	{"12m", []int{2}, "Only 3 characters completes this edge wait (penchan). 9 does not wrap around to 1."},
	{"88s", []int{25}, "A third 8 bamboo completes the set. A pair has only two remaining copies."},
	{"78p", []int{14, 17}, "6 or 9 dots completes the run. Two-sided waits can include a terminal."},
	{"79m", []int{7}, "Only 8 characters completes the run. Keep track of gaps between tiles."},
}

func (a *app) waitTrainer(c *ui.Context) {
	ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
		for i, name := range []string{"Waits", "Hand waits", "Discards", "Yaku", "Scoring", "Decisions", "Mixed review"} {
			c.Key(fmt.Sprintf("exercise-mode-%d", i))
			pill(c, name, a.exerciseMode == i, func() {
				a.exerciseMode = i
				a.exercise = 0
				a.answered = false
				a.feedback = ""
				if i == 6 {
					a.startReview()
				}
			})
		}
	})
	if a.exerciseMode == 2 {
		a.discardExerciseView(c)
		return
	}
	if a.exerciseMode >= 1 {
		if a.exerciseMode == 6 {
			a.mixedReviewView(c)
			return
		}
		topic := map[int]string{1: "Hand waits", 3: "Yaku", 4: "Scoring", 5: "Decisions"}[a.exerciseMode]
		pool := questionsFor(topic)
		a.questionExerciseView(c, pool[a.exercise%len(pool)], fmt.Sprintf("%d of %d", a.exercise%len(pool)+1, len(pool)), func() { a.exercise = (a.exercise + 1) % len(pool); a.answered = false; a.feedback = "" })
		return
	}
	d := waitDrills[a.exercise%len(waitDrills)]
	panel(c).FillWidth().Children(func() {
		heading(c, "Find the finishing tile")
		body(c, "Choose a tile that turns this two-tile shape into a run or a set. Other groups in the hand are already complete.")
		a.tileRow(c, mustTiles(d.shape))
		secondary(c, fmt.Sprintf("Exercise %d of %d", a.exercise%len(waitDrills)+1, len(waitDrills)))
		ui.Row(c).Wrap().Gap(10).Children(func() {
			suit := mustTiles(d.shape)[0] / 9
			count := 9
			if suit == 3 {
				count = 7
			}
			for n := 0; n < count; n++ {
				t := suit*9 + n
				c.Key(fmt.Sprintf("choice-%d", t))
				a.tile(c, t, false, "Choose "+tileName(t), func() {
					if a.answered {
						return
					}
					correct := false
					for _, v := range d.answers {
						if v == t {
							correct = true
						}
					}
					if correct {
						a.answered = true
						a.feedback = "Correct. " + d.why
					} else {
						a.feedback = "That tile does not complete this shape. Try another tile."
					}
					a.recordExercise("shape:"+d.shape, correct)
				})
			}
		})
		if a.feedback != "" {
			body(c, a.feedback)
		}
		primaryAction(c, "Next exercise").OnClick(func() { a.exercise = (a.exercise + 1) % len(waitDrills); a.answered = false; a.feedback = "" })
	})
	panel(c).FillWidth().Children(func() {
		heading(c, "Look for both sides")
		body(c, "Adjacent middle tiles usually make a stronger wait than a gap or an edge. In a real hand, subtract copies already in your hand or visible on the table.")
		actionButton(c, "Read about waits").OnClick(func() { a.navigate("Learn"); a.selectLesson(3) })
	})
}
