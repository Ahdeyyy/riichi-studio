# Icon round 2: final-selection review

Method: viewed `round-2.png` (16/32/64 px on light and dark), each refinement at 256 px, and A/B/C/D at 16 and 32 px enlarged with nearest-neighbour scaling so I could see the actual pixels. I also checked the SVG geometry of the 1024, 32 and 16 px masters. I compared against the flaws listed in `round-1-review.md`.

## Decision

**Winner: A, "Table declaration". 8.5 / 10. Ready to ship once the 16 and 32 px exports are snapped to the pixel grid (an export fix, not a redesign).**

Runner-up: **C, "Placed stick"** (8.0).

| Variant | Score | Verdict |
|---|---:|---|
| **A** Table declaration | **8.5** | Clearly a tenbou, clearly physical, and the cleanest at 16 px. No slash, no pill. |
| C Placed stick | 8.0 | The best-looking at 256/1024: lively and tactile. Loses at 16 px, where the diagonal staircases and the red turns pink. |
| D Quiet horizontal | 7.5 | Crisp, but without the edge it reads as a flat UI glyph ("minus" / slider track). The shorter, thicker bar with a bigger dot is the most Hinomaru-adjacent of the four. |
| B Clean diagonal | 7.0 | At −30° with no depth it is still the most "slash / strikethrough / disabled" of the set. Weakest. |

## Why A

- **Round-1 flaws are fixed.** At 380×64 (about 5.9:1) with rx 9, it is a stick, not a band-aid or capsule. The flat ends read as a token or tool. Being horizontal removes the slash and "prohibited" readings entirely.
- **The tan edge is what makes it work.** A 7-unit `#D9C9A3` lip under the stick turns it from a glyph into an object lying on felt. That is the exact difference between A and D, and it removes D's "minus button" reading. It is also the only detail, and it costs nothing at small sizes.
- **Riichi meaning.** Horizontal is how the stick actually sits once declared: placed flat in front of the player, parallel to the table edge. Any riichi player recognizes it; outsiders see a calm, ownable mark. "Ready" (tenpai) is the learner's goal, so it carries the learning meaning without a book or lamp cliché.
- **Unwanted associations.** Medical: gone. Slash: gone. Japanese flag: low risk. The bar is about 6:1, not 3:2, and the dot (r 22 in a 64-high stick) leaves clear ivory above and below, so it reads as an inset marker, not a disc on a field. **Keep the dot at this size. Do not enlarge it at 256/1024**, because a bigger dot pushes the design toward D and the flag.
- **Small sizes.** Horizontal maps onto the pixel grid. At 16 px the red stays true vermilion. On the diagonals the red blends with ivory and goes salmon. On a Windows taskbar, where a desktop Go app lives at 16–32 px most of the time, this decides it.
- **Balance.** The stick spans about 82% of the squircle width. It fills the field without crowding, and it holds on both light and dark surfaces in the contact sheet. The empty jade above and below reads as intentional calm, not as "too small" (the round-1 complaint applied to the shorter diagonal pill).

### Why not C

C is the more characterful illustration, and at 1024 it would arguably win a beauty contest. But a tilt only adds energy when the icon is large. At 16 px the −24° stick becomes a jagged staircase, the red block turns into an anti-aliased pink smudge, and a faint slash reading comes back. A gives up a little charm for a big gain in legibility and zero ambiguity. If a future large-format use (splash, store art, README hero) wants motion, C's angle can serve as a secondary illustration of the same mark.

## Remaining refinements (A only; all are export hinting, not design changes)

Seen at the pixel level, A-16 currently has **soft half-rows**. The ivory spans y 212–300 in 512-space, which is 6.6–9.4 px at 16 px, so you get a greyish top row and a muddy tan/grey bottom row around only two crisp ivory rows. A-32 has the same problem (ivory 218–294 → 13.6–18.4 px), and its dot is a ~3 px blur with a pink halo.

1. **16 px master (1 px = 32 units):** snap every edge to multiples of 32.
   - Ivory: `x 64–448, y 224–288` (2 full rows, 12 px wide).
   - Tan edge: `x 64–448, y 288–320` (1 full row). The physical lip doubles as the extra weight the optical-thickening was trying to add, and it does so without blur.
   - Red: `x 224–288, y 224–288`, rx 0 (a crisp 2×2, full ivory height, as intended).
   - Ends: rx 0 at this size. Rounding only makes grey corner pixels.
2. **32 px master (1 px = 16 units):** ivory `x 64–448, y 224–288` (4 rows), tan `y 288–304` (1 row), red `x 224–288, y 224–288` with rx ≈ 16 (a 4×4 block with corner pixels knocked back, so it still reads as a dot). Look at it at 100% next to a 2×2 hard square centred at `240–272`, and keep whichever reads as "dot" rather than "segment". I expect the 4×4 rounded one.
3. **64 px and up:** no changes. Optionally snap the 64 px ivory to `y 224–288` and tan to `288–296` for crisp edges. The 256/1024 geometry is final.
4. **Windows .ico completeness:** Windows also requests **24, 48 and 20/40** px (DPI scaling). Make sure the ICO contains snapped 24 and 48 frames (24: 1 px = 21.33 units, so hand-place ivory at 2 rows, tan 1 row, red 2×2; 48: derive from the 32 rules ×1.5) rather than letting Windows resample the 32 or 256. Unhinted resampling would bring back exactly the blur fixed above.
5. Keep D's file around only as the monochrome/flat fallback (for example, a single-colour tray or notification glyph, where the tan edge would be lost anyway).

Nothing else. Do not add texture, tenbou end markings, extra dots, text, or a second object. The icon now has one object, one accent and one physical cue, which is what it should be.

## Ship recommendation

**Ship A.** Apply the pixel-snapping in items 1–2 (and the ICO frames in item 4) at export time. No further design rounds are needed.
