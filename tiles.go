package main

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"sync"
)

var honorNames = []string{"East", "South", "West", "North", "White dragon", "Green dragon", "Red dragon"}
var honorFaces = []string{"東", "南", "西", "北", "白", "發", "中"}

// Names and codes are read for every tile of every frame, so they are
// built once.
var tileNames, tileCodes = func() (names, codes [34]string) {
	for t := range names {
		if t >= 27 {
			names[t], codes[t] = honorNames[t-27], fmt.Sprintf("%dz", t-26)
			continue
		}
		names[t] = fmt.Sprintf("%d %s", t%9+1, []string{"characters", "dots", "bamboo"}[t/9])
		codes[t] = fmt.Sprintf("%d%c", t%9+1, "mps"[t/9])
	}
	return
}()

func tileName(t int) string { return tileNames[t] }
func tileCode(t int) string { return tileCodes[t] }
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
//
// The suits cannot share a group, so each is enumerated on its own and the
// four results are combined, rather than walking their product. The view
// asks for the shanten of the hand on every frame, and the drills ask
// hundreds of times at launch.
func shanten(tiles []int) int {
	var counts [34]int
	for _, t := range tiles {
		counts[t]++
	}
	if len(tiles) > 14 {
		return shantenWalk(counts)
	}
	pairs, unique, orphans, orphanPair := 0, 0, 0, 0
	for i, n := range counts {
		if n > 4 {
			return shantenWalk(counts)
		}
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
	best := min(8, 6-pairs+max(0, 7-unique), 13-orphans-orphanPair)
	// all[head][melds] is the most incomplete groups the suits so far give.
	all := noShapes
	all[0][0] = 0
	for suit := 0; suit < 4; suit++ {
		shapes := suitShapes(counts[suit*9:min(34, suit*9+9)], suit < 3)
		next := noShapes
		for h := 0; h < 2; h++ {
			for m := 0; m < 5; m++ {
				if all[h][m] < 0 {
					continue
				}
				for sh := 0; h+sh < 2; sh++ {
					for sm := 0; m+sm < 5; sm++ {
						if t := shapes[sh][sm]; t >= 0 {
							next[h+sh][m+sm] = max(next[h+sh][m+sm], all[h][m]+t)
						}
					}
				}
			}
		}
		all = next
	}
	for h := 0; h < 2; h++ {
		for m := 0; m < 5; m++ {
			if t := int(all[h][m]); t >= 0 {
				best = min(best, 8-2*m-min(t, 4-m)-h)
			}
		}
	}
	return best
}

// suitTable[head][melds] is the most incomplete groups one suit's tiles
// make beside that many melds and heads, or -1 when they cannot make them.
type suitTable [2][5]int8

var noShapes = suitTable{{-1, -1, -1, -1, -1}, {-1, -1, -1, -1, -1}}

// Tables are kept by the suit's counts: hands repeat suits far more often
// than whole hands. Scoring runs off the UI thread, hence the lock.
var suitTables = struct {
	sync.Mutex
	m map[uint32]suitTable
}{m: map[uint32]suitTable{}}

func suitShapes(counts []int, runs bool) suitTable {
	var c [9]int
	key := uint32(0)
	for i, n := range counts {
		c[i] = n
		key = key*5 + uint32(n)
	}
	if runs {
		key |= 1 << 31
	}
	suitTables.Lock()
	table, ok := suitTables.m[key]
	suitTables.Unlock()
	if ok {
		return table
	}
	table = noShapes
	var walk func(int, int, int, int)
	walk = func(i, melds, taatsu, head int) {
		for i < 9 && c[i] == 0 {
			i++
		}
		if i == 9 {
			if melds < 5 {
				table[head][melds] = max(table[head][melds], int8(taatsu))
			}
			return
		}
		if c[i] >= 3 {
			c[i] -= 3
			walk(i, melds+1, taatsu, head)
			c[i] += 3
		}
		if runs && i <= 6 && c[i+1] > 0 && c[i+2] > 0 {
			c[i]--
			c[i+1]--
			c[i+2]--
			walk(i, melds+1, taatsu, head)
			c[i]++
			c[i+1]++
			c[i+2]++
		}
		if c[i] >= 2 {
			c[i] -= 2
			if head == 0 {
				walk(i, melds, taatsu, 1)
			}
			walk(i, melds, taatsu+1, head)
			c[i] += 2
		}
		if runs {
			for d := 1; d <= 2; d++ {
				if i+d < 9 && c[i+d] > 0 {
					c[i]--
					c[i+d]--
					walk(i, melds, taatsu+1, head)
					c[i]++
					c[i+d]++
				}
			}
		}
		c[i]--
		walk(i, melds, taatsu, head)
		c[i]++
	}
	walk(0, 0, 0, 0)
	suitTables.Lock()
	if len(suitTables.m) >= 1<<15 {
		clear(suitTables.m)
	}
	suitTables.m[key] = table
	suitTables.Unlock()
	return table
}

// shantenWalk is the same analysis by one walk over the whole hand. It
// takes any number of tiles, and is the reference the suit tables are
// tested against, but needs tens of microseconds for a 14-tile hand.
func shantenWalk(counts [34]int) int {
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
