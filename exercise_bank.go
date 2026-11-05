package main

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
)

type reviewQuestion struct {
	key, topic, hand string
	q                checkQuestion
}

var reviewBank []reviewQuestion

var exerciseBankReady bool

func ensureExerciseBank() {
	if exerciseBankReady {
		return
	}
	exerciseBankReady = true
	for i := range lessons {
		topic := "Decisions"
		if stageForLesson(i) == 1 {
			topic = "Shape"
		}
		if stageForLesson(i) == 2 {
			topic = "Yaku"
		}
		if i == 7 || i == 22 {
			topic = "Scoring"
		}
		for q, question := range lessonChecks(i) {
			reviewBank = append(reviewBank, reviewQuestion{fmt.Sprintf("check:%d:%d", i, q), topic, "", question})
		}
	}
	for han := 1; han <= 5; han++ {
		for _, fu := range []int{30, 40} {
			for _, dealer := range []bool{false, true} {
				value := ronPoints(han, fu, dealer)
				seat := "non-dealer"
				if dealer {
					seat = "dealer"
				}
				why := fmt.Sprintf("%d han and %d fu: %s ron is %s points, rounded up to 100. No honba or deposits are included.", han, fu, seat, points(value))
				if han == 5 || han == 4 && fu == 40 {
					why = "This reaches mangan: 8,000 non-dealer ron or 12,000 dealer ron. Fu no longer raises this payment."
				}
				reviewBank = append(reviewBank, reviewQuestion{fmt.Sprintf("score:ron:%d:%d:%t", han, fu, dealer), "Scoring", "", checkQuestion{fmt.Sprintf("%d han, %d fu, %s ron. With no honba or deposits, what is the payment?", han, fu, seat), []string{points(value), points(value + 100), points(value + 1000)}, 0, why}})
			}
		}
	}
	for han := 1; han <= 4; han++ {
		for _, dealer := range []bool{false, true} {
			seat := "non-dealer"
			if dealer {
				seat = "dealer"
			}
			correct := tsumoPayment(han, 30, dealer)
			wrong1 := tsumoPayment(han+1, 30, dealer)
			wrong2 := tsumoPayment(han, 30, !dealer)
			reviewBank = append(reviewBank, reviewQuestion{fmt.Sprintf("score:tsumo:%d:%t", han, dealer), "Scoring", "", checkQuestion{fmt.Sprintf("%d han, 30 fu, %s tsumo. No honba or deposits: how is it paid?", han, seat), []string{correct, wrong1, wrong2}, 0, "Round each payer's share up to 100 separately. The dealer pays twice each non-dealer's share when a non-dealer wins."}})
		}
	}
	for _, h := range []string{"123m456p23467s55p", "123m456p234s678s5p", "123m456p789s55p77z", "1122m3344p5566s7z", "19m19p19s1234567z", "2345m456p789s777z", "34567s123m456p55z"} {
		w := waits(mustTiles(h))
		correct := waitNames(w)
		if len(w) == 13 {
			correct = "All 13 terminal and honor types"
		}
		reviewBank = append(reviewBank, reviewQuestion{"wait:" + h, "Hand waits", h, checkQuestion{"Which tiles complete this 13-tile shape? Ignore yaku and furiten for this shape check.", []string{correct, "Only East", "No tile completes this shape"}, 0, "Check every decomposition, including seven pairs and thirteen orphans where relevant. Complete-shape waits: " + correct + "."}})
	}
	for _, d := range waitDrills {
		correct := waitNames(d.answers)
		wrong := tileCode((d.answers[0]/9)*9 + (d.answers[0]%9+1)%9)
		if d.answers[0] >= 27 {
			wrong = tileCode(27 + (d.answers[0]-27+1)%7)
		}
		if len(d.answers) > 1 {
			wrong = tileCode(d.answers[0])
		}
		reviewBank = append(reviewBank, reviewQuestion{"shape:" + d.shape, "Shape", d.shape, checkQuestion{"Name all tile types that complete this block as a run or set.", []string{correct, wrong, "No tile can complete it"}, 0, d.why}})
	}
	for _, d := range discardDrills {
		analysis := analyzeDiscards(d.hand)
		best := analysis[0]
		var bestTiles []int
		wrong := "Discard any tile"
		for _, option := range analysis {
			if option.shanten == best.shanten && option.ukeire == best.ukeire {
				bestTiles = append(bestTiles, option.tile)
			} else {
				wrong = tileCode(option.tile)
			}
		}
		reviewBank = append(reviewBank, reviewQuestion{"discard:" + d.hand, "Shape", d.hand, checkQuestion{"With no opponent threat, which discards minimize shanten and then maximize improving copies?", []string{waitNames(bestTiles), wrong, "Break any complete run"}, 0, fmt.Sprintf("Best choices: %s. They leave shanten %d with %d improving copies. %s", waitNames(bestTiles), best.shanten, best.ukeire, d.why)}})
	}
	// Shuffle choice order once, retaining the correct index. Saved exercise
	// results use stable keys rather than screen position or option order.
	for i := range reviewBank {
		q := &reviewBank[i].q
		q.choices = append([]string(nil), q.choices...)
		correct := q.choices[q.correct]
		rand.Shuffle(len(q.choices), func(a, b int) { q.choices[a], q.choices[b] = q.choices[b], q.choices[a] })
		for j, c := range q.choices {
			if c == correct {
				q.correct = j
			}
		}
	}
}
func basePoints(han, fu int) int {
	if han >= 5 {
		return 2000
	}
	return min(2000, fu*(1<<(han+2)))
}
func round100(v int) int { return (v + 99) / 100 * 100 }
func ronPoints(han, fu int, dealer bool) int {
	factor := 4
	if dealer {
		factor = 6
	}
	return round100(basePoints(han, fu) * factor)
}
func tsumoPayment(han, fu int, dealer bool) string {
	b := basePoints(han, fu)
	if dealer {
		return fmt.Sprintf("%s from each opponent", points(round100(b*2)))
	}
	return fmt.Sprintf("%s from East, %s from each other", points(round100(b*2)), points(round100(b)))
}
func waitNames(w []int) string {
	names := []string{}
	for _, t := range w {
		names = append(names, tileCode(t))
	}
	return strings.Join(names, ", ")
}
func questionsFor(topic string) []reviewQuestion {
	ensureExerciseBank()
	var out []reviewQuestion
	for _, q := range reviewBank {
		if q.topic == topic {
			out = append(out, q)
		}
	}
	return out
}
func (a *app) startReview() {
	ensureExerciseBank()
	pool := append([]reviewQuestion(nil), reviewBank...)
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	sort.SliceStable(pool, func(i, j int) bool {
		aR, bR := a.exerciseRecords[pool[i].key], a.exerciseRecords[pool[j].key]
		priority := func(r exerciseRecord) float64 {
			if r.Attempts == 0 {
				return .5
			}
			v := 1 - float64(r.Correct)/float64(r.Attempts)
			if !r.LastCorrect {
				v += 1
			}
			return v
		}
		return priority(aR) > priority(bR)
	})
	a.reviewSession = pool[:min(10, len(pool))]
	a.reviewIndex = 0
	a.reviewCorrect = 0
	a.reviewMissed = nil
	a.answered = false
	a.feedback = ""
}

type discardDrill struct{ hand, why string }

var discardDrills = []discardDrill{
	{"123m456p23467s55p1z", "Preserve the three runs, a pair and a two-sided block before keeping an isolated honor."},
	{"234m678p34567s88p4z", "Look for five useful blocks and compare the honor with the connected number tiles."},
	{"5578m2455p133s777z", "An overlapping connection can be redundant. Compare useful draws, not just how connected a tile looks."},
	{"3346m2245667p777z", "Keep enough support for the weakest block, and preserve a head candidate."},
	{"123m345p567s2256p4z", "Avoid breaking a completed run when another tile is outside the useful blocks."},
	{"234567m3459p67s55z", "An isolated terminal can be weaker than a two-sided shape even when it shares a suit with a complete run."},
	{"234m456p345s22p78s1z", "Compare shanten first. If tied, count how many improving copies remain unseen."},
	{"1122m3344p5566s7z9m", "Compare the seven-pairs route as well as standard groups. Six pairs are close to seven-pairs tenpai."},
	{"19m19p19s1234567z5m", "A special hand can be much faster than forcing four groups. Thirteen orphans uses terminals and honors."},
	{"345m456p567s22z13m9p", "Keep the gap block when it is needed for the fourth group; compare the floating terminal's contribution."},
}

type discardOption struct{ tile, shanten, ukeire int }

var discardCache = map[string][]discardOption{}

func analyzeDiscards(h string) []discardOption {
	if cached, ok := discardCache[h]; ok {
		return cached
	}
	hand := mustTiles(h)
	var counts [34]int
	for _, t := range hand {
		counts[t]++
	}
	var out []discardOption
	seen := map[int]bool{}
	for i, t := range hand {
		if seen[t] {
			continue
		}
		seen[t] = true
		rest := append([]int(nil), hand[:i]...)
		rest = append(rest, hand[i+1:]...)
		s := shanten(rest)
		n := 0
		for draw := 0; draw < 34; draw++ {
			if counts[draw] < 4 && shanten(append(append([]int(nil), rest...), draw)) < s {
				n += 4 - counts[draw]
			}
		}
		out = append(out, discardOption{t, s, n})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].shanten != out[j].shanten {
			return out[i].shanten < out[j].shanten
		}
		return out[i].ukeire > out[j].ukeire
	})
	discardCache[h] = out
	return out
}
