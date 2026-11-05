# Icon round 1 — independent design review

Method: viewed `round-1.png`, a fresh 16/32/64 px downsample of all 20 drafts, and drafts 01, 05, 08, 10, 12, 15, 16 and 18 at about 250 px, plus their SVG geometry.

## Candid summary

The set is clean and the palette works: jade `#173F36`, ivory `#F5EDDA` and vermilion `#CB493B` look good together and feel like a real table. But most motifs are **generic icon vocabulary in mahjong colors**: a book, sprout, chat bubble, hashtag, bar chart, compass star, letter R, origami crane. Swap the palette and they belong to any app. **Only one draft (08) has an idea that riichi players would recognize and that no other app owns.** Three more are good enough to refine. None says "learning" without falling back on clichés (book, steps, lamp). My view: the icon should own *riichi*, and the app name and context can carry "learning". The best tie to learning is a concept like "ready / tenpai / the missing piece", not a classroom object.

## Scores (out of 10, holistic: small-size legibility, distinctiveness, riichi fit, learning fit, balance, simplicity, unwanted associations)

| # | Name | Score | One-line verdict |
|---|------|------:|------------------|
| 08 | Riichi deposit | **7.5** | The only truly ownable idea. Reads as a bandage/pill at its current proportions; fixable. |
| 05 | The missing tile | **6.5** | Good "wait / what's missing" learning idea. Tall bars plus a raised one read as a bar chart. |
| 18 | Table compass | **6.0** | Strong silhouette at 16 px, a top-down table. Risks reading as a road sign / warning diamond. |
| 01 | Signal tile | **6.0** | A clearly ownable tile shape, but it reads as phone / ID card, and the corner red dot is a notification badge. |
| 16 | Single circle | 5.5 | The best 16 px legibility in the set, but it *is* the Target logo / a record button. Disqualifying. |
| 10 | Enso wait | 5.0 | Reads as ©, "C", or a power/refresh glyph. Enso is generic Zen branding. |
| 09 | Tile fan | 5.0 | Legible, but it's playing cards: poker/solitaire. |
| 07 | Wind rose | 4.5 | Navigation app / Safari / map pin. East wind is too subtle. |
| 03 | Dragon ribbon | 4.0 | Reads as the numeral "3" or a "Z". Not 中. |
| 04 | Bamboo lesson | 4.0 | Sprout = wellness/eco/habit app. |
| 06 | Five blocks | 4.0 | App launcher / dashboard grid. |
| 19 | Riichi ligature | 4.0 | A letter monogram; "R" belongs to anything. |
| 02 | Folded lesson | 3.5 | A generic open book (edtech cliché). |
| 17 | Paper crane | 3.5 | Generic Japanese-ness, nothing to do with mahjong. Muddy at 16 px. |
| 11 | Sparrow | 3.0 | Nice 1-sou lore, but an unreadable blob below 64 px. |
| 12 | Pinfu weave | 3.0 | A hashtag. |
| 13 | Tile conversation | 3.0 | A chat app. |
| 15 | Open tile | 3.0 | Reads as "n", "i", or trousers. |
| 14 | Study light | 2.5 | Unclear triangle/lamp; weak at every size. |
| 20 | Ready steps | 2.5 | A bar chart: finance/analytics. Exactly what the user dislikes. |

## Shortlist (4)

### 08 Riichi deposit — PROVISIONAL WINNER
**Why it wins.** The 1000-point stick with its single red dot is *the* physical act of riichi: you declare by putting it on the table. No other icon in the set (or on the app stores) uses it. It means "I am ready" (tenpai), which is also the learner's goal, so riichi fit and learning fit come from one object. It is one shape with one accent, uses the exact palette (ivory stick, vermilion dot, jade felt), and survives 16 px as a clear diagonal. It is minimal without being generic.

**Flaws.** At a 4:1 aspect ratio with fully round ends and a −35° tilt, it reads as a **band-aid, pill, or capsule** (a medical association) and slightly as a "disabled / strikethrough" slash. At 16 px the dot (r = 23/512) disappears. The stick also leaves the corners empty, so it looks small.

**Refinements.**
- Proportions: go from 316×78 to about **380×60** (roughly 6.3:1). That is thinner and longer, like a real tenbou, which removes the band-aid read. At 16 px, 60/512 is still about 1.9 px.
- Ends: change `rx` from 39 (full pill) to about **12–14**. Squared-off sticks read as a tool or token, not a medicine.
- Angle: try **−30° and 0°** (horizontal, as it lies on the table in front of the player). Horizontal kills the "slash" association completely. If you keep the diagonal, raise its length so it fills the 464 squircle diagonally with about 70 px of clearance.
- Dot: r ≈ 20 inside a 60-tall stick (one-third margin top and bottom). Make a **hand-tuned 16 px variant** where the dot becomes a full-height vermilion square segment (about 2×2 px) so the accent survives.
- Optional depth: a 6–8/512 lower edge in a warm shadow ivory (about `#D9C9A3`, the tan already used in 09/14) on the bottom side gives the tactile bone/plastic feel of the table without detail.
- Do NOT add a second object, text, or a book. Its strength is that it's a single object.

### 05 The missing tile
**Strength:** the clearest *learning* concept: "which tile completes the hand?" is the core skill of riichi.
**Flaws:** 91×217 tiles (1:2.4) are bars, not tiles, and with the raised red one the whole thing reads as a bar chart or "pause plus record".
**Refinements:** use real tile proportions (about **3:4, e.g. 104×140**), all on **one baseline**, in a row of three centered optically. Make the third tile an **empty slot**: a vermilion outline of 26–30/512 stroke with `rx` matching the tiles, or a solid vermilion tile shifted about 24 px right with a bigger gap. Don't raise it vertically; that's what causes the chart read. Leave a gap of about 20 px between the ivory tiles so they stay separate at 32 px.

### 18 Table compass
**Strength:** a bold, geometric silhouette that reads as a top-down table with four seats. Very stable at 16 px.
**Flaws:** a solid diamond frame = road/warning sign or "decision" flowchart node. The center tile (54×78) shrinks to a dot at 32 px.
**Refinements:** break the frame into **four separate ivory bars** (the four walls/hands) with notches about 28/512 wide at the corners. That reads as a table of walls, not a sign. Enlarge the center tile to about **72×100**. Optionally, replace the center tile with a short horizontal riichi stick, which links it to 08 if you want a family.

### 01 Signal tile
**Strength:** the most direct "tile" silhouette, while staying abstract (no hand-drawn tile face).
**Flaws:** the bottom bar looks like a phone's home indicator, and a **top-right red dot = unread-notification badge**. It reads as a phone, card, or document.
**Refinements:** remove the bar. Move the dot to the **optical center** and enlarge it (r ≈ 44) so it reads as a minimal 1-pin. Add a thickness slab at the base of the tile (about 26/512, tan `#D9C9A3`) so it reads as a physical mahjong tile, not a screen. Rotate it about −8° to make it feel tactile rather than UI.

## Next-round suggestions
- Push 08 hardest: produce 3–4 variants (horizontal, −30°, with/without depth edge, dot vs center band), each at 1024/64/32/16 px with hand-hinted small sizes.
- Test every finalist on light *and* dark OS surfaces (taskbar, dock). The jade squircle helps there, but check that ivory art doesn't bloom at 16 px.
- Keep the vermilion to **one accent per icon**. The best drafts already follow this rule.
