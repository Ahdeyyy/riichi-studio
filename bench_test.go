package main

import (
	"math/rand/v2"
	"runtime"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// Steady frames of the busiest pages, rendered on the CPU by ui.Tester.
// Run with: go test -run '^$' -bench . -benchmem
func benchFrame(b *testing.B, prepare func(*app)) {
	a := newApp()
	prepare(a)
	tt := ui.NewTester(a.view, 1280, 860)
	for range 5 {
		tt.Frame()
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		tt.Frame()
	}
}

func BenchmarkPracticeFrame(b *testing.B) { benchFrame(b, func(a *app) {}) }
func BenchmarkPracticeHintFrame(b *testing.B) {
	benchFrame(b, func(a *app) { a.resetTable(true); a.selected = 3; a.hint = true })
}
func BenchmarkWorstCaseFrame(b *testing.B) {
	benchFrame(b, func(a *app) { a.applyFixture("Worst case") })
}
func BenchmarkLearnFrame(b *testing.B)     { benchFrame(b, func(a *app) { a.navigate("Learn") }) }
func BenchmarkExercisesFrame(b *testing.B) { benchFrame(b, func(a *app) { a.navigate("Exercises") }) }

func BenchmarkShanten(b *testing.B) {
	hand := mustTiles("1239m4567p2466s1z5z")
	b.ReportAllocs()
	for b.Loop() {
		shanten(hand)
	}
}

// TestStartupCost reports what launching costs before the first frame.
func TestStartupCost(t *testing.T) {
	start := time.Now()
	a := newApp()
	t.Logf("newApp: %v", time.Since(start))
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	t.Logf("retained heap after newApp: %.1f MB", float64(m.HeapAlloc)/(1<<20))
	runtime.KeepAlive(a)
}

// The suit tables must agree with the whole-hand walk on every hand.
func TestShantenMatchesWalk(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	var wall []int
	for tile := 0; tile < 34; tile++ {
		wall = append(wall, tile, tile, tile, tile)
	}
	for i := 0; i < 20000; i++ {
		r.Shuffle(len(wall), func(a, b int) { wall[a], wall[b] = wall[b], wall[a] })
		hand := wall[:1+r.IntN(14)]
		if i%3 == 0 {
			// Random hands are rarely close to complete: draw from two suits.
			hand = hand[:0:0]
			var n [34]int
			for len(hand) < 13+i%2 {
				if tile := r.IntN(18); n[tile] < 4 {
					n[tile]++
					hand = append(hand, tile)
				}
			}
		}
		var counts [34]int
		for _, tile := range hand {
			counts[tile]++
		}
		if got, want := shanten(hand), shantenWalk(counts); got != want {
			t.Fatalf("%s: shanten %d, walk %d", notation(hand), got, want)
		}
	}
}
