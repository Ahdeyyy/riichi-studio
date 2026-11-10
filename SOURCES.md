# Learning and scoring sources

The strategy curriculum is adapted in original wording from **Riichi Book I: A Mahjong Strategy Primer for European Players**, Daina Chiba, version 13 (the supplied PDF). The book is licensed under [Creative Commons Attribution–NonCommercial 3.0](https://creativecommons.org/licenses/by-nc/3.0/). The added novice introduction is original app content. Adapted lesson text should retain attribution and the noncommercial condition. The PDF itself is not redistributed.

Book source: https://github.com/dainachiba/RiichiBooks

Lessons reference §3.2 (tiles and groups), §3.4 (waits), Chapter 4 (five-block method), Chapter 5 (yaku), Chapter 6 (scoring), Chapter 7 (riichi), Chapter 8 (defense), and Chapter 9 (calls). Exercises are newly composed. The book explicitly assumes prior knowledge of the rules; the first lessons bridge that gap. This app is an introduction, not a reproduction of the entire book.

The expanded course also follows Chapter 10's focus on match placement. It contains 28 lessons, 56 worked examples and two checks per lesson. The 180-item exercise bank combines those checks, generated block-completion problems, full-hand waits, discard comparisons and payment arithmetic. Generated shape answers use the app's hand analyzer; concrete payment examples are regression-checked against Agari. Strategy scenarios state assumptions and teach useful defaults rather than claiming a universal best move. The beginner rules material is original, with explicit limits where the practice table does not implement a rule.

Scoring uses [Agari](https://github.com/agari-industries/agari), v0.27.0, under the MIT license, included at `resources/agari/LICENSE`. The Windows x86-64 executable comes from the upstream release. SHA-256 of its release archive: `552bb5cf5b02feef317777ef0ad0face2a85d1a56af25037d9029dfc00e8d4df`, verified against the upstream checksum. No modifications were made to the executable. It ships with the Windows app as a resource (`resources/windows-amd64/agari`). The macOS and Linux apps ship the upstream builds for their platforms, which `scripts/fetch-agari.sh` downloads from the same release and checks against the upstream SHA-256 sums pinned in that script.

## Tile artwork

Tile markings, fronts and backs use the original **600 × 800 PNG exports** from [FluffyStuff/riichi-mahjong-tiles](https://github.com/FluffyStuff/riichi-mahjong-tiles), pinned to commit `26e127ba2117f45cdce5ea0225748cc0cfad3169`. These assets are released under CC0 1.0 (public domain); the original declaration is included at `resources/tiles/LICENSE.md`. The source PNGs are bundled unchanged; at launch the app scales them to 192 × 256 pixels, the largest size it shows them at, to save memory. The app composites the front, back and markings at runtime, overlays optional labels in the tiles' top-left corners, and retains its interaction animations. No AI-generated tile artwork is used.

## Rule scope

The practice table is solo, closed-hand draw-and-discard training, not a four-player match. It has no opponent simulation, ron, open calls, kans, red fives, dead wall, honba, pot or match standings. The practice wall includes every tile outside your starting hand, and is not a regulation live wall. A first guided hand deliberately supplies a winning draw; random hands shuffle the remaining tiles. Riichi has a teaching deposit without a persistent points account. Hints minimize shanten; they do not assess safety or maximize ukeire. Winning shapes include standard hands, seven pairs, and thirteen orphans.

The lab scores closed 14-tile hands without kans or red fives. It uses Agari's default four-player rules, East round and South seat (East seat when Dealer is checked), with riichi, self-draw and dora indicator controls. It does not validate a historical game log, furiten, or how a declaration was reached. Payments exclude honba and riichi deposits. Advanced rule variants from the book are not all represented.
