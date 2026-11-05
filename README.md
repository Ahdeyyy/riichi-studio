# Riichi Studio

A native Windows riichi mahjong tutor built with Go and MyGo. A full-width felt table with a wood rail, ivory tiles with colored backs, detailed suit artwork, and a compact navigation bar that shares the title bar's band. Every surface is opaque; the window uses no backdrop blur. Tiles lift on hover and selection, new hands deal in sequence, and physical tiles animate from the hand into the discard area. Motion follows the operating system's reduced-motion preference. Includes 28 lessons in four stages, 56 worked examples, 56 lesson checks, 180 exercise items with saved results and targeted review, a tile glossary and offline Agari scoring.

Run from this directory:

```powershell
go run .
```

Build a standalone Windows executable:

```powershell
go tool mygo build -skip-installer
./scripts/apply-icon.ps1 -Executable 'build/windows-amd64/Riichi Studio.exe'
```

The scoring engine is embedded; no Rust installation or network access is needed at runtime. Requires the Go version in go.mod to build. Windows x86-64 is the bundled scoring target; other systems need an `agari` CLI on PATH.

The original app icon was chosen from 20 drafts through two Claude Opus 5.5 judging rounds. SVG, PNG and ICO sources, review notes and size tests are in [design/icons](design/icons/README.md).

Start with **Learn** if you have never played. In **Practice**, select East and declare riichi for the guided first win. Use **New random hand** for solo training. Hints compare hand shape; the table does not simulate opponents or full matches. **Scoring lab** accepts 14-tile closed hands in standard m/p/s/z notation.

Progress is stored in the operating system's user configuration folder under `riichi-tutor/progress.json`.

The learning path progresses through **First hands**, **Hand shape**, **Yaku and value**, and **Table decisions**. Each lesson has an overview, two worked examples and two checks. Examples cover calls, furiten, shanten, ukeire, complex waits, pinfu, value sets, dora, alternative yaku, exact scoring cases, riichi judgment, suji, push/fold and placement goals. Lessons remain open for review at any time.

**Exercises** offers 79 two-tile block waits (including honors), seven full-hand wait problems, ten discard problems, yaku quizzes, 28 generated payment questions and decision scenarios. These and the lesson checks form a 180-item mixed-review bank. Ten-question review sessions prioritize past mistakes; correct/attempt counts are saved. Discard problems compare shanten and then improving copies, counting the discarded tile as visible; they assume no other visible tiles and do not assess danger or hand value.

Detailed check progress and exercise records are stored in `riichi-tutor/learning.json`. Previously passed single-check lessons migrate as passed first checks; each added second check remains available to finish. Completing this material builds a foundation for casual play, but the solo practice table still does not simulate opponents or a full match.

Tile artwork uses bundled 600 × 800 PNG images from FluffyStuff's CC0 riichi tile set. The **Tile labels** checkbox applies everywhere: 1–9 for suited tiles and E/S/W/N/Wh/G/R for honors. Labels overlay the top-left corner inside each tile without changing the artwork's position or size. Labels are on by default; your choice is saved in `riichi-tutor/preferences.json` in the same configuration folder. All artwork works offline.

The jade title bar supports dragging, double-click maximize, minimize, restore and close. Windows retains its resizing edges and snap-layout support. Action buttons fit their text rather than stretching across panels.

Checks:

```powershell
$env:CGO_ENABLED='0'
go test ./...
go tool mygo vet .
```

See [SOURCES.md](SOURCES.md) for source attribution, licenses, and the rule scope.

Layout fixtures (Demo, Worst case, Empty, One) live in `fixtures.go`. To flip between them in the running app, start a dev build; release builds leave the toggle out:

```powershell
go run -tags dev .
```
