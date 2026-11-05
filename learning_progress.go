package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type exerciseRecord struct {
	Attempts, Correct int
	LastCorrect       bool
}
type learningSave struct {
	Answers   map[string]bool           `json:"answers"`
	Exercises map[string]exerciseRecord `json:"exercises"`
}

func (a *app) learningPath() string {
	if a.savePath == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(a.savePath), "learning.json")
}
func (a *app) loadLearning() {
	data, err := os.ReadFile(a.learningPath())
	var saved learningSave
	if err == nil && json.Unmarshal(data, &saved) == nil {
		if saved.Answers != nil {
			a.lessonAnswers = saved.Answers
		}
		if saved.Exercises != nil {
			a.exerciseRecords = saved.Exercises
		}
	} else {
		// Old versions had one check per lesson. Preserve those passed checks.
		for i, done := range a.completed {
			if done && i >= 0 && i < len(lessons) {
				a.lessonAnswers[lessonKey(i, 0)] = true
			}
		}
	}
	a.recomputeCompletion()
}
func (a *app) recomputeCompletion() {
	a.completed = map[int]bool{}
	for i := range lessons {
		all := true
		for q := range lessonChecks(i) {
			if !a.lessonAnswers[lessonKey(i, q)] {
				all = false
			}
		}
		if all {
			a.completed[i] = true
		}
	}
}
func (a *app) saveLearning() {
	if a.learningPath() == "" {
		return
	}
	data, _ := json.Marshal(learningSave{a.lessonAnswers, a.exerciseRecords})
	if err := os.MkdirAll(filepath.Dir(a.learningPath()), 0755); err != nil {
		a.feedback += " Learning progress could not be saved."
		return
	}
	if err := os.WriteFile(a.learningPath(), data, 0600); err != nil {
		a.feedback += " Learning progress could not be saved."
	}
}
func (a *app) recordExercise(key string, correct bool) {
	r := a.exerciseRecords[key]
	r.Attempts++
	if correct {
		r.Correct++
	}
	r.LastCorrect = correct
	a.exerciseRecords[key] = r
	a.saveLearning()
}
func (a *app) selectLesson(i int) {
	a.lesson = i
	a.lessonTab = 0
	a.lessonExample = 0
	a.lessonCheck = 0
	a.feedback = ""
	a.answered = false
	for q := range lessonChecks(i) {
		if !a.lessonAnswers[lessonKey(i, q)] {
			a.lessonCheck = q
			break
		}
	}
}
func (a *app) answerLesson(answer int) {
	q := lessonChecks(a.lesson)[a.lessonCheck]
	if a.lessonAnswers[lessonKey(a.lesson, a.lessonCheck)] {
		return
	}
	if answer == q.correct {
		a.feedback = "Correct. " + q.why
		a.answered = true
		a.lessonAnswers[lessonKey(a.lesson, a.lessonCheck)] = true
		a.recomputeCompletion()
		a.saveProgress()
		a.saveLearning()
	} else {
		a.feedback = "Not yet. " + q.why
	}
}
func (a *app) advanceLesson(delta int) {
	order := courseOrder()
	for at, i := range order {
		if i == a.lesson && at+delta >= 0 && at+delta < len(order) {
			a.selectLesson(order[at+delta])
			return
		}
	}
}
