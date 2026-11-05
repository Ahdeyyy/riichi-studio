package main

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
)

var honorNames = []string{"East", "South", "West", "North", "White dragon", "Green dragon", "Red dragon"}
var honorFaces = []string{"東", "南", "西", "北", "白", "發", "中"}

func tileName(t int) string {
	if t >= 27 {
		return honorNames[t-27]
	}
	return fmt.Sprintf("%d %s", t%9+1, []string{"characters", "dots", "bamboo"}[t/9])
}
func tileCode(t int) string {
	if t >= 27 {
		return fmt.Sprintf("%dz", t-26)
	}
	return fmt.Sprintf("%d%s", t%9+1, []string{"m", "p", "s"}[t/9])
}
func parseTiles(s string) ([]int, error) {
	var tiles, digits []int
	for _, r := range strings.ReplaceAll(s, " ", "") {
		if r >= '1' && r <= '9' {
			digits = append(digits, int(r-'0'))
			continue
		}
		suit := strings.IndexRune("mpsz", r)
		if suit < 0 || len(digits) == 0 {
			return nil, fmt.Errorf("Use numbers followed by m, p, s or z, such as 123m456p.")
		}
		for _, n := range digits {
			if suit == 3 && n > 7 {
				return nil, fmt.Errorf("Honors use 1z through 7z.")
			}
			tiles = append(tiles, suit*9+n-1)
		}
		digits = nil
	}
	if len(digits) > 0 {
		return nil, fmt.Errorf("Add a suit letter after the numbers.")
	}
	var counts [34]int
	for _, t := range tiles {
		counts[t]++
		if counts[t] > 4 {
			return nil, fmt.Errorf("There are only four copies of %s.", tileName(t))
		}
	}
	return tiles, nil
}
func mustTiles(s string) []int {
	t, e := parseTiles(s)
	if e != nil {
		panic(e)
	}
	return t
}
func notation(tiles []int) string {
	v := append([]int(nil), tiles...)
	sort.Ints(v)
	s := ""
	for suit := 0; suit < 4; suit++ {
		d := ""
		for _, t := range v {
			if t/9 == suit {
				d += fmt.Sprint(t%9 + 1)
			}
		}
		if d != "" {
			s += d + string("mpsz"[suit])
		}
	}
	return s
}

// Shape analysis enumerates melds, incomplete groups and a head, then also
// checks seven pairs and thirteen orphans. It does not imply a legal yaku.
func shanten(tiles []int) int {
	var counts [34]int
	for _, t := range tiles {
		counts[t]++
	}
	best := 8
	pairs, unique, orphans, orphanPair := 0, 0, 0, 0
	for i, n := range counts {
		if n > 0 {
			unique++
			if i >= 27 || i%9 == 0 || i%9 == 8 {
				orphans++
				if n >= 2 {
					orphanPair = 1
				}
			}
		}
		if n >= 2 {
			pairs++
		}
	}
	best = min(best, 6-pairs+max(0, 7-unique), 13-orphans-orphanPair)
	var walk func(int, int, int, int)
	walk = func(start, melds, taatsu, head int) {
		for start < 34 && counts[start] == 0 {
			start++
		}
		if start == 34 {
			best = min(best, 8-2*melds-min(taatsu, 4-melds)-head)
			return
		}
		i := start
		if counts[i] >= 3 {
			counts[i] -= 3
			walk(i, melds+1, taatsu, head)
			counts[i] += 3
		}
		if i < 27 && i%9 <= 6 && counts[i+1] > 0 && counts[i+2] > 0 {
			counts[i]--
			counts[i+1]--
			counts[i+2]--
			walk(i, melds+1, taatsu, head)
			counts[i]++
			counts[i+1]++
			counts[i+2]++
		}
		if counts[i] >= 2 {
			counts[i] -= 2
			if head == 0 {
				walk(i, melds, taatsu, 1)
			}
			walk(i, melds, taatsu+1, head)
			counts[i] += 2
		}
		if i < 27 {
			for d := 1; d <= 2; d++ {
				if i%9+d < 9 && counts[i+d] > 0 {
					counts[i]--
					counts[i+d]--
					walk(i, melds, taatsu+1, head)
					counts[i]++
					counts[i+d]++
				}
			}
		}
		counts[i]--
		walk(i, melds, taatsu, head)
		counts[i]++
	}
	walk(0, 0, 0, 0)
	return best
}
func waits(tiles []int) []int {
	var out []int
	var n [34]int
	for _, t := range tiles {
		n[t]++
	}
	for t := 0; t < 34; t++ {
		if n[t] < 4 && shanten(append(append([]int(nil), tiles...), t)) == -1 {
			out = append(out, t)
		}
	}
	return out
}
func bestDiscards(hand []int) ([]int, int) {
	best := 9
	var result []int
	seen := map[int]bool{}
	for i, t := range hand {
		if seen[t] {
			continue
		}
		seen[t] = true
		rest := append([]int(nil), hand[:i]...)
		rest = append(rest, hand[i+1:]...)
		s := shanten(rest)
		if s < best {
			best = s
			result = nil
		}
		if s == best {
			result = append(result, t)
		}
	}
	return result, best
}
func (a *app) resetTable(random bool) {
	a.hand = mustTiles("123m456p23467s55p1z")
	if random {
		a.hand = nil
	}
	counts := [34]int{}
	for _, t := range a.hand {
		counts[t]++
	}
	a.wall = nil
	for t := 0; t < 34; t++ {
		for n := counts[t]; n < 4; n++ {
			a.wall = append(a.wall, t)
		}
	}
	rand.Shuffle(len(a.wall), func(i, j int) { a.wall[i], a.wall[j] = a.wall[j], a.wall[i] })
	if random {
		a.hand = append(a.hand, a.wall[:14]...)
		a.wall = a.wall[14:]
	} else {
		for i, t := range a.wall {
			if t == 25 {
				a.wall[0], a.wall[i] = a.wall[i], a.wall[0]
				break
			}
		}
	}
	sort.Ints(a.hand)
	a.dealID++
	a.nextTileID = a.dealID * 1000
	a.handIDs = nil
	a.discardIDs = nil
	for range a.hand {
		a.handIDs = append(a.handIDs, a.nextTileID)
		a.nextTileID++
	}
	a.discards = nil
	a.selected = -1
	a.turn = 1
	a.riichi = false
	a.won = false
	a.hint = false
	a.scoreResult = ""
	a.tableMessage = "Select a tile, then discard. Keep four groups and a pair in mind."
}
func (a *app) discard() {
	if a.selected < 0 || a.selected >= len(a.hand) || a.won || len(a.wall) == 0 {
		return
	}
	t := a.hand[a.selected]
	a.discardIDs = append(a.discardIDs, a.handIDs[a.selected])
	a.handIDs = append(a.handIDs[:a.selected], a.handIDs[a.selected+1:]...)
	a.handIDs = append(a.handIDs, a.nextTileID)
	a.nextTileID++
	a.hand = append(a.hand[:a.selected], a.hand[a.selected+1:]...)
	a.discards = append(a.discards, t)
	a.hand = append(a.hand, a.wall[0])
	a.wall = a.wall[1:]
	a.turn++
	a.selected = -1
	a.hint = false
	if a.riichi {
		a.selected = len(a.hand) - 1
	}
	a.tableMessage = "You drew " + tileName(a.hand[len(a.hand)-1]) + ". Choose your next discard."
	if shanten(a.hand) == -1 {
		a.tableMessage = "Your hand is complete. Claim tsumo to see the yaku and payment."
	}
	if len(a.wall) == 0 {
		a.tableMessage = "The practice wall is empty. Start a new hand."
	}
}
