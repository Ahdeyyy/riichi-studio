package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
	"strings"
)

func init() {
	seen := map[string]bool{}
	for _, d := range waitDrills {
		seen[d.shape] = true
	}
	add := func(h string, answers []int, why string) {
		if !seen[h] {
			waitDrills = append(waitDrills, waitDrill{h, answers, why})
			seen[h] = true
		}
	}
	for suit := 0; suit < 3; suit++ {
		for low := 0; low < 9; low++ {
			for high := low; high <= min(8, low+2); high++ {
				shape := []int{suit*9 + low, suit*9 + high}
				var answers []int
				for n := 0; n < 9; n++ {
					ranks := []int{low, high, n}
					valid := low == high && n == low
					if low != high {
						a, b, c := ranks[0], ranks[1], ranks[2]
						if c < a {
							a, c = c, a
						}
						if c < b {
							b, c = c, b
						}
						if b < a {
							a, b = b, a
						}
						valid = b == a+1 && c == b+1
					}
					if valid {
						answers = append(answers, suit*9+n)
					}
				}
				h := notation(shape)
				why := "Completion tiles: " + waitNames(answers) + ". A run stays in one suit and does not wrap; a set uses three identical tiles."
				add(h, answers, why)
			}
		}
	}
	for t := 27; t < 34; t++ {
		add(notation([]int{t, t}), []int{t}, "A third identical honor completes the set. Honors cannot form runs; only two other copies remain before visible tiles are counted.")
	}
}

func (a *app) questionExerciseView(c *ui.Context, r reviewQuestion, position string, next func()) {
	panel(c).FillWidth().Children(func() {
		heading(c, r.topic)
		secondary(c, position)
		if r.hand != "" {
			a.tileRow(c, mustTiles(r.hand))
		}
		body(c, r.q.prompt)
		for i, choice := range r.q.choices {
			actionButton(c, choice).Key(fmt.Sprintf("drill-answer-%d", i)).MaxWidth(580).Disabled(a.answered).OnClick(func() {
				correct := i == r.q.correct
				a.answered = true
				if a.exerciseMode == 6 {
					if correct {
						a.reviewCorrect++
					} else {
						a.reviewMissed = append(a.reviewMissed, r)
					}
				}
				prefix := "Correct. "
				if !correct {
					prefix = "Review this one. The answer is " + r.q.choices[r.q.correct] + ". "
				}
				a.feedback = prefix + r.q.why
				a.recordExercise(r.key, correct)
			})
		}
		if a.feedback != "" {
			body(c, a.feedback)
		}
		primaryAction(c, "Next exercise").Disabled(!a.answered).OnClick(next)
		record := a.exerciseRecords[r.key]
		if record.Attempts > 0 {
			secondary(c, fmt.Sprintf("This exercise: %d correct in %d attempts", record.Correct, record.Attempts))
		}
	})
}
func (a *app) mixedReviewView(c *ui.Context) {
	if len(a.reviewSession) == 0 {
		panel(c).FillWidth().Children(func() {
			heading(c, "A ten-question review")
			body(c, "A mix of shapes, yaku, payments and decisions. Questions you have missed are shown first; every answer includes an explanation.")
			primaryAction(c, "Start review").OnClick(a.startReview)
		})
		return
	}
	if a.reviewIndex >= len(a.reviewSession) {
		panel(c).FillWidth().Children(func() {
			heading(c, "Review complete")
			body(c, fmt.Sprintf("%d of %d correct in this session. Your results are saved.", a.reviewCorrect, len(a.reviewSession)))
			for _, r := range a.reviewMissed {
				body(c, r.q.prompt)
				secondary(c, "Answer: "+r.q.choices[r.q.correct]+". "+r.q.why)
			}
			if len(a.reviewMissed) > 0 {
				actionButton(c, "Retry mistakes").OnClick(func() {
					a.reviewSession = append([]reviewQuestion(nil), a.reviewMissed...)
					a.reviewMissed = nil
					a.reviewCorrect = 0
					a.reviewIndex = 0
					a.answered = false
					a.feedback = ""
				})
			}
			primaryAction(c, "Another review").OnClick(a.startReview)
		})
		return
	}
	a.questionExerciseView(c, a.reviewSession[a.reviewIndex], fmt.Sprintf("Review %d of %d", a.reviewIndex+1, len(a.reviewSession)), func() { a.reviewIndex++; a.answered = false; a.feedback = "" })
}
func (a *app) discardExerciseView(c *ui.Context) {
	d := discardDrills[a.exercise%len(discardDrills)]
	hand := mustTiles(d.hand)
	analysis := analyzeDiscards(d.hand)
	best := analysis[0]
	panel(c).FillWidth().Children(func() {
		heading(c, "Choose a discard")
		secondary(c, fmt.Sprintf("Hand %d of %d", a.exercise%len(discardDrills)+1, len(discardDrills)))
		body(c, "No opponent threat and no visible tiles besides this hand. Minimize shanten, then maximize improving draws. Click a tile to discard it.")
		ui.Row(c).Wrap().Gap(6).Children(func() {
			for i, t := range hand {
				c.Key(fmt.Sprintf("discard-choice-%d", i))
				a.tile(c, t, false, "Discard "+tileName(t), func() {
					if a.answered {
						return
					}
					a.answered = true
					var picked discardOption
					var accepted []string
					for _, r := range analysis {
						if r.tile == t {
							picked = r
						}
						if r.shanten == best.shanten && r.ukeire == best.ukeire {
							accepted = append(accepted, tileName(r.tile))
						}
					}
					correct := picked.shanten == best.shanten && picked.ukeire == best.ukeire
					prefix := "Correct. "
					if !correct {
						prefix = "Compare the alternatives. "
					}
					a.feedback = fmt.Sprintf("%sBest choices: %s. They leave shanten %d with %d improving copies. Your choice leaves shanten %d with %d improving copies. %s", prefix, strings.Join(accepted, ", "), best.shanten, best.ukeire, picked.shanten, picked.ukeire, d.why)
					a.recordExercise("discard:"+d.hand, correct)
				})
			}
		})
		if a.feedback != "" {
			body(c, a.feedback)
		}
		primaryAction(c, "Next hand").Disabled(!a.answered).OnClick(func() { a.exercise = (a.exercise + 1) % len(discardDrills); a.answered = false; a.feedback = "" })
		secondary(c, "This drill ranks shape and theoretical availability only. It does not replace yaku, value or defense judgment at a real table.")
	})
}
