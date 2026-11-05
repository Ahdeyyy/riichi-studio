# Riichi Studio icon

The selected mark is **Table declaration**: an ivory riichi deposit stick with a single vermilion dot and a narrow tan lower edge on jade. It depicts an object specific to riichi, with no letters, tile glyphs, learning badges, or decorative texture.

## Selection

1. Created 20 distinct rough motifs, preserved in `drafts/` and `round-1.png`.
2. Claude provider, model `claude-opus-5-5`, independently judged all 20. Shortlist: 08 Riichi deposit, 05 The missing tile, 18 Table compass, 01 Signal tile. See `round-1-review.md` for all scores, objections, and suggested changes.
3. Chose the strongest direction, 08, and made four refinements of proportion, angle, edge and dot size. See `round-2.png` and `refinements/`.
4. A fresh Claude Opus 5.5 judging task selected A, Table declaration, at 8.5/10. See `round-2-review.md`. Its remaining findings concerned pixel alignment, which were applied to the final small exports. No additional judge was used for those export fixes.

The judge rejected the generic book, sprout, finance-like bars and dashboard grid. The original capsule was lengthened and given flatter ends to remove its medical association. The horizontal version beat the diagonal version because its small exports stay cleaner and avoid a slash reading. Final small exports were manually snapped to the pixel grid, including the 20/24/40/48 px Windows DPI sizes.

## Deliverables

- `../../resources/icon.svg`: editable 1024px SVG master (512-unit viewBox).
- `../../resources/icon.png`: 1024 x 1024 RGBA PNG, transparent beyond the jade shape.
- `../../resources/icon.ico`: nine PNG frames at 16, 20, 24, 32, 40, 48, 64, 128 and 256 px.
- `../../resources/icon-window.png`: 128px runtime image using the pixel-aligned small geometry.
- `final/`: individual PNGs at 16 through 1024 px and the small SVG masters.
- `final-proof.png`: final mark and true-size tests on light and dark backgrounds.

The icon is integrated into the native window/title bar/taskbar, executable resources and in-app header. MyGo initially creates Windows resources from the 1024px PNG. `scripts/apply-icon.ps1` replaces those icon frames with the optically adjusted ICO while preserving the manifest and version information. The final executable was checked by extracting its associated Windows icon.

## Reproduce

The app build needs no Node dependencies. To regenerate the design artifacts:

```powershell
npm install --prefix .cache/icon-tools @resvg/resvg-js@2.6.2
node design/icons/explore.cjs
node design/icons/refine.cjs
node design/icons/finalize.cjs
go tool mygo build -skip-installer
./scripts/apply-icon.ps1 -Executable 'build/windows-amd64/Riichi Studio.exe'
```

All marks were drawn as original SVG geometry and rasterized with resvg. No stock icon or AI-generated raster asset was used. The app's separate mahjong tile images retain their original attribution in SOURCES.md.
