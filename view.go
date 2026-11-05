package main

import "github.com/egoist/mygo/ui"

func (a *app) scoringView(c *ui.Context) {
	secondary(c, "Closed-hand scoring with Agari. Four players, East round; South seat unless Dealer is selected. No honba or riichi pot.")
	panel(c).FillWidth().Children(func() {
		heading(c, "Explore a hand")
		body(c, "Hand notation")
		ui.TextInput(c, &a.notation).Label("Hand notation").Placeholder("123m456p234678s55p").FillWidth()
		secondary(c, "Include all 14 tiles. m = characters, p = dots, s = bamboo; 1–7z = East, South, West, North, White, Green, Red. No red fives or calls in this lab.")
		if ts, e := parseTiles(a.notation); e == nil && len(ts) > 0 && len(ts) <= 14 {
			a.tileRow(c, ts)
		}
		ui.Row(c).Gap(16).Children(func() {
			ui.Column(c).Gap(8).Children(func() { body(c, "Winning tile"); ui.TextInput(c, &a.winning).Label("Winning tile").Width(180) })
			ui.Column(c).Gap(8).Children(func() {
				body(c, "Dora indicators")
				ui.TextInput(c, &a.indicator).Label("Dora indicators").Placeholder("Optional, e.g. 4p").Width(240)
			})
		})
		ui.Row(c).Gap(24).Children(func() {
			ui.Checkbox(c, &a.tsumo, "Self-draw (tsumo)")
			ui.Checkbox(c, &a.declared, "Riichi")
			ui.Checkbox(c, &a.dealer, "Dealer")
		})
		primaryAction(c, "Calculate score").Disabled(a.busy).OnClick(func() { a.calculate(scoreInput{a.notation, a.winning, a.indicator, a.tsumo, a.declared, a.dealer}) })
	})
	panel(c).FillWidth().Children(func() {
		heading(c, "Score breakdown")
		if a.busy && a.scorePage == "Scoring lab" {
			body(c, "Calculating with Agari…")
		} else if a.scoreResult != "" && a.scorePage == "Scoring lab" {
			body(c, a.scoreResult)
		} else {
			body(c, "The example is a closed pinfu hand. Calculate it, then change the win type or dealer status to compare payments.")
		}
	})
}
func (a *app) referenceView(c *ui.Context) {
	panel(c).FillWidth().Children(func() {
		heading(c, "The tile set")
		for suit, name := range []string{"Characters", "Dots", "Bamboo", "Winds and dragons"} {
			body(c, name)
			n := 9
			if suit == 3 {
				n = 7
			}
			tiles := []int{}
			for i := 0; i < n; i++ {
				tiles = append(tiles, suit*9+i)
			}
			a.tileRow(c, tiles)
		}
	})
	for _, entry := range glossary {
		panel(c).FillWidth().Children(func() { heading(c, entry[0]); body(c, entry[1]) })
	}
	secondary(c, "Learning source: Riichi Book I by Daina Chiba, CC BY-NC 3.0. Scoring: Agari v0.27.0, MIT license. See SOURCES.md for links, attribution and rule scope.")
}
