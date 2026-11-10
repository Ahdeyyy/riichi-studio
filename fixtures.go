package main

// UI fixtures enter through app state, the same boundary saved progress
// and table play use. They are the regression data for layout changes;
// the dev toggle (build tag dev) switches between them.
var fixtureNames = []string{"Demo", "Worst case", "Empty", "One"}

func (a *app) applyFixture(name string) {
	page := a.page
	*a = *fixtureBase(a)
	a.page = page
	switch name {
	case "Worst case":
		// A solo table draws the whole wall: 104 discards after 104 draws.
		a.dealDiscards(104)
		a.riichi = true
		a.tableMessage = "Riichi declared. Your hand is locked; discard each new draw or claim tsumo."
		a.hint = true
		for i := range lessons {
			if i%3 != 2 {
				a.completed[i] = true
			}
		}
		ensureExerciseBank()
		for _, r := range reviewBank {
			a.exerciseRecords[r.key] = exerciseRecord{Attempts: 312, Correct: 297}
		}
		// A counted yakuman, the longest breakdown Agari reports here.
		a.notation, a.winning, a.indicator, a.dealer = "234m234p234s22456s", "6s", "1m1p1s4z", true
		a.scorePage = "Scoring lab"
		a.scoreResult = "13 han / 20 fu\n\nRiichi   1 han\nMenzen Tsumo   1 han\nPinfu   1 han\nTanyao   1 han\nSanshoku Doujun   2 han\nIipeikou   1 han\nDora bonuses   6 han\n\n48,000 points total\n16,000 from each opponent"
		a.feedback = "Compare the alternatives. Best choices: 1 characters, 9 characters, East wind, White dragon. They leave shanten 2 with 46 improving copies. Your choice leaves shanten 3 with 112 improving copies."
	case "Empty":
		// An exhaustive draw: nothing left to draw, nothing scored.
		a.dealDiscards(122)
		a.feedback = ""
	case "One":
		a.dealDiscards(121)
		a.discards, a.discardIDs = a.discards[:1], a.discardIDs[:1]
		a.completed[0] = true
		ensureExerciseBank()
		for _, r := range reviewBank {
			a.exerciseRecords[r.key] = exerciseRecord{Attempts: 1, Correct: 1, LastCorrect: true}
		}
	}
}

// fixtureBase is a fresh app keeping the running window's hooks. It has
// no save path, so fixture state never overwrites saved progress.
func fixtureBase(old *app) *app {
	a := newApp()
	a.update, a.chrome, a.updates, a.devData = old.update, old.chrome, old.updates, old.devData
	return a
}

// dealDiscards draws and discards n wall tiles in table order.
func (a *app) dealDiscards(n int) {
	for i := 0; i < n && len(a.wall) > 0; i++ {
		a.discards = append(a.discards, a.wall[0])
		a.discardIDs = append(a.discardIDs, a.nextTileID)
		a.nextTileID++
		a.wall = a.wall[1:]
		a.turn++
	}
}
