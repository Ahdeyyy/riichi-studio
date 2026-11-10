package main

import (
	"github.com/egoist/mygo/plugins/updater"
	"github.com/egoist/mygo/ui"
)

// updateActions reach the updater plugin, which only the running app uses:
// tests build the app without it, and leave check nil.
type updateActions struct {
	check     func()
	available string // a newer version found by the last check, or ""
}

// watchUpdates follows the plugin's background checks, so the navigation
// offers a version as soon as one is found.
func (a *app) watchUpdates() {
	a.updates.check = updater.CheckForUpdates
	updater.OnChange(func() {
		v := updater.AvailableVersion()
		a.update(func() { a.updates.available = v })
	})
}

// updateButton checks for updates on demand, or names the version that a
// check found. Either way the plugin's window takes it from there.
func (a *app) updateButton(c *ui.Context) {
	if a.updates.check == nil {
		return
	}
	c.Key("updates")
	if a.updates.available != "" {
		// Sized as the pills beside it.
		ui.PrimaryButton(c, "Update to "+a.updates.available).AlignSelf(ui.Start).Padding(9, 16).FontSize(14).OnClick(a.updates.check)
		return
	}
	pill(c, "Check for updates", false, a.updates.check)
}
