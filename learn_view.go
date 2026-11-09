package main

import (
	"fmt"
	"github.com/egoist/mygo/ui"
)

func (a *app) learnView(c *ui.Context) {
	stage := stageForLesson(a.lesson)
	ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
		for i, s := range courseStages {
			label := s.title
			done := 0
			for _, id := range s.lessons {
				if a.completed[id] {
					done++
				}
			}
			label += fmt.Sprintf("  %d/%d", done, len(s.lessons))
			c.Key(fmt.Sprintf("stage-%d", i))
			pill(c, label, stage == i, func() { a.selectLesson(s.lessons[0]) })
		}
	})
	ui.Row(c).FillWidth().Gap(24).AlignItems(ui.Start).Children(func() {
		ui.Column(c).Width(220).Gap(6).Children(func() {
			for _, i := range courseStages[stage].lessons {
				l := lessons[i]
				label := l.title
				if a.completed[i] {
					label += "  ✓"
				}
				// A list reads down its left edge, so lesson titles start there.
				c.Key(fmt.Sprintf("lesson-%d", i))
				pill(c, label, i == a.lesson, func() { a.selectLesson(i) }).AlignSelf(ui.Stretch).Justify(ui.Start)
			}
			secondary(c, a.progress())
		})
		panel(c).Grow(1).MinWidth(0).Children(func() {
			l := lessons[a.lesson]
			heading(c, l.title)
			ui.Row(c).Gap(8).Children(func() {
				for i, tab := range []string{"Overview", "Examples", "Checks"} {
					c.Key(fmt.Sprintf("lesson-tab-%d", i))
					pill(c, tab, a.lessonTab == i, func() { a.lessonTab = i; a.feedback = "" })
				}
			})
			switch a.lessonTab {
			case 0:
				body(c, l.body)
				a.tileRow(c, mustTiles(l.example))
				secondary(c, l.source)
				primaryAction(c, "See worked examples").OnClick(func() { a.lessonTab = 1 })
			case 1:
				examples := workedHands[a.lesson]
				ui.Row(c).Gap(8).Children(func() {
					for i := range examples {
						actionButton(c, fmt.Sprintf("Example %d", i+1)).OnClick(func() { a.lessonExample = i })
					}
				})
				e := examples[a.lessonExample]
				body(c, e.title)
				a.tileRow(c, mustTiles(e.hand))
				body(c, e.why)
				primaryAction(c, "Check your understanding").OnClick(func() { a.lessonTab = 2; a.feedback = "" })
			case 2:
				questions := lessonChecks(a.lesson)
				q := questions[a.lessonCheck]
				passed := a.lessonAnswers[lessonKey(a.lesson, a.lessonCheck)]
				secondary(c, fmt.Sprintf("Check %d of %d", a.lessonCheck+1, len(questions)))
				body(c, q.prompt)
				for i, choice := range q.choices {
					c.Key(fmt.Sprintf("answer-%d", i))
					actionButton(c, choice).MaxWidth(560).Disabled(passed).OnClick(func() { a.answerLesson(i) })
				}
				if a.feedback != "" {
					body(c, a.feedback)
				} else if passed {
					body(c, "Previously passed. "+q.why)
				}
				if a.lessonCheck+1 < len(questions) {
					primaryAction(c, "Next check").Disabled(!passed).OnClick(func() { a.lessonCheck++; a.answered = false; a.feedback = "" })
				} else if a.completed[a.lesson] {
					body(c, "Lesson complete. Both checks are saved.")
				}
			}
			order := courseOrder()
			ui.Row(c).Gap(10).Children(func() {
				actionButton(c, "Previous").Disabled(a.lesson == order[0]).OnClick(func() { a.advanceLesson(-1) })
				primaryAction(c, "Next lesson").Disabled(a.lesson == order[len(order)-1]).OnClick(func() { a.advanceLesson(1) })
				actionButton(c, "Practice this skill").OnClick(func() {
					mode := 0
					switch stageForLesson(a.lesson) {
					case 1:
						mode = 2
					case 2:
						mode = 3
					case 3:
						mode = 5
					}
					if a.lesson == 3 || a.lesson == 16 {
						mode = 1
					}
					if a.lesson == 7 || a.lesson == 22 {
						mode = 4
					}
					a.navigate("Exercises")
					a.exerciseMode = mode
					a.exercise = 0
				})
			})
		})
	})
	secondary(c, "Original exercises following Riichi Book I by Daina Chiba (CC BY-NC 3.0), with an added beginner rules introduction. Lessons are open for review at any time.")
}
