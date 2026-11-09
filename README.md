# Riichi Studio

Riichi Studio is a desktop app for learning Japanese riichi mahjong. It takes you from recognizing the tiles to reading hand value and making table decisions, with a solo practice table and a scoring lab along the way. It is built with Go and [MyGo](https://github.com/egoist/mygo) and runs offline on Windows, macOS and Linux.

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

Build the Windows app and its installer:

```powershell
go tool mygo build
```

This writes `build/windows-amd64/Riichi Studio.exe` and `riichi-studio-setup-<version>-amd64.exe`, an installer made with NSIS (mygo downloads NSIS when it is not installed). The installer needs no admin rights: it installs into `%LOCALAPPDATA%\Programs\Riichi Studio`, adds a Start menu shortcut, and registers an uninstaller under Settings > Apps. Add `-skip-installer` to build only the executable.

Before packing the installer, mygo runs `scripts/sign-windows.ps1` on each executable (`windows.signCommand` in `mygo.json`). It applies the app icon and, when `WINDOWS_CERTIFICATE_FILE` (and `WINDOWS_CERTIFICATE_PASSWORD`, if the certificate has one) is set, signs with `signtool`. The release workflow sets these from the `WINDOWS_CERTIFICATE` (base64 `.pfx`) and `WINDOWS_CERTIFICATE_PASSWORD` repository secrets. Unsigned builds work, but SmartScreen warns on first run.

### Updates

The app updates itself from this repository's GitHub releases (`updates` in `mygo.json`). It checks 10 seconds after launch and daily, and **Check for updates** in the navigation bar checks on demand; a newer version shows there as **Update to x.y.z**. The update window shows the release notes and downloads, verifies and installs the update, then relaunches.

Updates are signed with an Ed25519 key from `go tool mygo keygen`. `mygo.json` holds the public key; the release workflow reads the secret key from the `MYGO_UPDATER_PRIVATE_KEY` repository secret and refuses to release without it, since installed apps only accept updates signed with that key. Keep a backup of the secret key: losing it strands existing installs. To build signed updates locally, set `MYGO_UPDATER_PRIVATE_KEY` to the key's contents before `go tool mygo build`.

The per-user installer and the portable zip can update themselves; development builds (`go run`) cannot.

### Scoring engine

The Agari scoring engine ships with the app as a resource, so the app needs no network access or extra installs. Each platform's build lives in its own directory of `resources/`, which mygo merges into that platform's app: `resources/windows-amd64/agari/agari.exe` is in the repository, and `scripts/fetch-agari.sh` downloads the macOS and Linux builds from Agari's release, checks them against pinned SHA-256 sums and puts them in `resources/darwin-*/agari` and `resources/linux-*/agari` (not committed). Run it once before `go run .`, `go test` or `go tool mygo build` on those systems; with no arguments it fetches the build for your machine. An `agari` on your PATH is used when none is bundled.

### macOS and Linux

```sh
scripts/fetch-agari.sh darwin-arm64 darwin-amd64
go tool mygo build -platform darwin/universal        # on a Mac: Riichi Studio.app and a .dmg

scripts/fetch-agari.sh linux-amd64 linux-arm64
go tool mygo build -platform linux/amd64,linux/arm64 # .tar.gz with install.sh, and .deb
```

The Linux packages need GTK 3 only (`linux.webview` is `false` in `mygo.json`). Without a Developer ID certificate the macOS app is signed ad hoc: Gatekeeper blocks a downloaded copy until you allow it under System Settings > Privacy & Security, while the release's `install.sh` (`curl -fsSL https://github.com/Ahdeyyy/riichi-studio/releases/latest/download/install.sh | sh`) installs it without that prompt, and is also the self-updating install on Linux.

## Workflows

- **CI** (`.github/workflows/ci.yml`) runs on pushes to `master` and on pull requests. It runs the tests and `mygo vet` on Windows, macOS and Linux, then builds the macOS and Linux packages and checks that each holds the engine for its platform.
- **Release** (`.github/workflows/release.yml`) is started by hand with a version. It builds, signs and publishes the Windows installer and portable zip, the universal macOS disk image, and the Linux archives and Debian packages for x86-64 and ARM64, with the signed update files for all of them, then commits the version and tags it.

Optional secrets sign the macOS app with a Developer ID and notarize it: `MACOS_CERTIFICATE` (a base64 `.p12` of a Developer ID Application certificate), `MACOS_CERTIFICATE_PASSWORD`, and for notarization `APPLE_ID`, `APPLE_TEAM_ID` and `APPLE_APP_PASSWORD` (an app-specific password).

## Test

```powershell
$env:CGO_ENABLED='0'
go test ./...
go tool mygo vet .
```

`go test -run '^$' -bench . -benchmem` times a frame of the busiest pages and the hand analysis; `go test -run TestStartupCost -v` reports launch time and memory.

`go run -tags dev .` adds a toggle to the running app. It switches between demo, worst-case, empty and single-item data, which is useful when checking layout changes.

## Attribution

The lessons and strategy content are adapted, in original wording, from **Riichi Book I: A Mahjong Strategy Primer for European Players** by **Daina Chiba**, available at [github.com/dainachiba/RiichiBooks](https://github.com/dainachiba/RiichiBooks). The book is licensed under [Creative Commons Attribution-NonCommercial 3.0](https://creativecommons.org/licenses/by-nc/3.0/). The adapted material carries the same terms: keep this attribution and do not use it commercially. The beginner rules introduction and the exercises are original to this app. The book itself is not included.

Riichi Studio also uses:

- [Agari](https://github.com/agari-industries/agari) v0.27.0 for scoring, under the MIT license (`resources/agari/LICENSE`).
- Tile artwork from [FluffyStuff/riichi-mahjong-tiles](https://github.com/FluffyStuff/riichi-mahjong-tiles), released under CC0 1.0 (`resources/tiles/LICENSE.md`).

[SOURCES.md](SOURCES.md) has the book sections each lesson draws on, the pinned versions and checksums, and the rule scope.

## License

The source code is licensed under the [MIT License](LICENSE). The lesson content adapted from Riichi Book I stays under CC BY-NC 3.0, and the bundled engine and artwork keep their own licenses, as listed above.
