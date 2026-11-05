//go:build dev

package main

import "github.com/egoist/mygo/ui"

// devToolbar swaps fixtures at the app-state boundary. Build or run with
// -tags dev; release builds compile the no-op in devtoggle.go.
func (a *app) devToolbar(c *ui.Context) {
	ui.Row(c).Absolute().Left(0).Right(0).Bottom(16).Justify(ui.Center).Children(func() {
		ui.Row(c).Padding(4).Radius(10).Background(ui.Hex("#2B2B2E")).Border(1, ui.Hex("#45454B")).Children(func() {
			ui.Segmented(c, &a.devData, fixtureNames...).Label("Fixture data").OnChange(func() { a.applyFixture(fixtureNames[a.devData]) })
		})
	})
}
