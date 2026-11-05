# Riichi Studio

Riichi Studio is a desktop app for learning Japanese riichi mahjong. It takes you from recognizing the tiles to reading hand value and making table decisions, with a solo practice table and a scoring lab along the way. It is built with Go and [MyGo](https://github.com/egoist/mygo) and runs offline on Windows.

## Features

- **Learn**: 28 lessons in four stages (First hands, Hand shape, Yaku and value, Table decisions). Each lesson has an overview, two worked examples and two checks.
- **Practice**: a solo draw-and-discard table. Select a tile, discard it, declare riichi and claim tsumo. Hints suggest discards by hand shape. The first hand is a guided win; **New random hand** deals a fresh one.
- **Exercises**: waits, full-hand waits, discard choices, yaku, payments and decisions. A ten-question mixed review shows your past mistakes first.
- **Scoring lab**: enter a closed 14-tile hand in m/p/s/z notation, such as `123m456p234678s55p`, and see its han, fu, yaku and payment.
- **Reference**: the full tile set and a glossary of terms.
- **Tile labels**: optional corner labels (1–9, E/S/W/N, Wh/G/R) for players who don't read the characters yet.
- **Saved progress**: lessons, checks, exercise results and preferences are saved in your user configuration folder under `riichi-tutor/`.

The practice table is single-player training. It does not simulate opponents, calls or full matches. [SOURCES.md](SOURCES.md) lists the full rule scope.

## Run locally

You need the Go version listed in `go.mod`.

```powershell
go run .
```

## Build

Build a standalone Windows executable and apply the app icon:

```powershell
go tool mygo build -skip-installer
./scripts/apply-icon.ps1 -Executable 'build/windows-amd64/Riichi Studio.exe'
```

The Agari scoring engine is embedded for Windows x86-64, so the app needs no network access or extra installs. On other systems, install the `agari` CLI and put it on your PATH.

## Test

```powershell
$env:CGO_ENABLED='0'
go test ./...
go tool mygo vet .
```

`go run -tags dev .` adds a toggle to the running app. It switches between demo, worst-case, empty and single-item data, which is useful when checking layout changes.

## Attribution

The lessons and strategy content are adapted, in original wording, from **Riichi Book I: A Mahjong Strategy Primer for European Players** by **Daina Chiba**, available at [github.com/dainachiba/RiichiBooks](https://github.com/dainachiba/RiichiBooks). The book is licensed under [Creative Commons Attribution-NonCommercial 3.0](https://creativecommons.org/licenses/by-nc/3.0/). The adapted material carries the same terms: keep this attribution and do not use it commercially. The beginner rules introduction and the exercises are original to this app. The book itself is not included.

Riichi Studio also uses:

- [Agari](https://github.com/agari-industries/agari) v0.27.0 for scoring, under the MIT license (`resources/agari/LICENSE`).
- Tile artwork from [FluffyStuff/riichi-mahjong-tiles](https://github.com/FluffyStuff/riichi-mahjong-tiles), released under CC0 1.0 (`resources/tiles/LICENSE.md`).

[SOURCES.md](SOURCES.md) has the book sections each lesson draws on, the pinned versions and checksums, and the rule scope.

## License

The source code is licensed under the [MIT License](LICENSE). The lesson content adapted from Riichi Book I stays under CC BY-NC 3.0, and the bundled engine and artwork keep their own licenses, as listed above.
