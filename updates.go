package main

import (
	"context"
	"fmt"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
)

// The app checks this long after launch, then every updateInterval.
const (
	firstUpdateCheck = 10 * time.Second
	updateInterval   = 24 * time.Hour
)

// release is a newer version that a check found.
type release struct {
	version, notes string
	install        func(ctx context.Context, progress func(downloaded, total int64)) error
}

// updatePhase is what the update dialog shows.
type updatePhase int

const (
	updateIdle updatePhase = iota
	updateChecking
	updateUpToDate
	updateUnsupported
	updateCheckFailed
	updateOffer
	updateDownloading
	updateInstallFailed
	updateReady
)

// updateActions reach mygo.Updater, which only the running app uses: tests
// build the app without it, and leave find nil. Every field is read and
// written on the UI thread; workers publish through app.update.
type updateActions struct {
	find     func(context.Context) (*release, error)
	enabled  func() bool
	relaunch func()
	current  string // the running version

	open              bool // the dialog shows
	phase             updatePhase
	release           *release
	available         string // a newer version found by the last check, or ""
	later             string // the version the user put off, which background checks do not offer again
	downloaded, total int64
	err               string
	cancel            context.CancelFunc
	notes             *markdown.Options // a pointer: the app is copied when fixtures swap
}

// watchUpdates hooks up mygo.Updater (updates in mygo.json) and checks
// shortly after launch and daily; a newer version opens the update dialog.
func (a *app) watchUpdates() {
	a.updates.find = func(ctx context.Context) (*release, error) {
		up, err := mygo.Updater.Check(ctx)
		if err != nil || up == nil {
			return nil, err
		}
		return &release{version: up.Version, notes: up.Notes, install: up.Install}, nil
	}
	a.updates.enabled = mygo.Updater.Enabled
	a.updates.relaunch = mygo.App.Relaunch
	a.updates.current = mygo.App.Version()
	publish := a.update
	go func() {
		if !mygo.Updater.Enabled() {
			return // development and package-manager builds
		}
		wait := firstUpdateCheck
		for {
			time.Sleep(wait)
			publish(func() { a.checkForUpdates(false) })
			wait = updateInterval
		}
	}()
}

// checkForUpdates starts a check. One the user asked for shows the dialog
// from the start; a background check shows it only for a newer version.
func (a *app) checkForUpdates(user bool) {
	u := &a.updates
	if u.find == nil {
		return
	}
	switch u.phase {
	case updateChecking, updateDownloading, updateReady:
		// Under way, or installed: show where it is.
		if user {
			u.open = true
		}
		return
	}
	if user {
		u.later = ""
		u.open = true
		if !u.enabled() {
			u.phase = updateUnsupported
			return
		}
	}
	u.phase, u.err = updateChecking, ""
	ctx, cancel := context.WithCancel(context.Background())
	u.cancel = cancel
	publish := a.update
	go func() {
		r, err := u.find(ctx)
		publish(func() {
			if ctx.Err() != nil {
				return // canceled
			}
			u.cancel = nil
			cancel()
			a.checked(r, err)
		})
	}()
}

// checked shows the outcome of a check.
func (a *app) checked(r *release, err error) {
	u := &a.updates
	switch {
	case err != nil:
		u.phase, u.err = updateCheckFailed, err.Error()
		if !u.open {
			a.settleUpdate() // a background check fails quietly
		}
	case r == nil:
		u.phase, u.release, u.available = updateUpToDate, nil, ""
		if !u.open {
			u.phase = updateIdle
		}
	default:
		u.phase, u.release, u.available = updateOffer, r, r.version
		u.open = u.open || r.version != u.later
	}
}

// installUpdate downloads, verifies and installs the offered release. The
// running app is untouched until it relaunches.
func (a *app) installUpdate() {
	u := &a.updates
	r := u.release
	if r == nil || u.phase != updateOffer && u.phase != updateInstallFailed {
		return
	}
	u.phase, u.err, u.downloaded, u.total = updateDownloading, "", 0, 0
	ctx, cancel := context.WithCancel(context.Background())
	u.cancel = cancel
	publish := a.update
	go func() {
		var last time.Time
		err := r.install(ctx, func(downloaded, total int64) {
			// A redraw for every chunk would flood the window.
			if now := time.Now(); downloaded >= total || now.Sub(last) >= 100*time.Millisecond {
				last = now
				publish(func() {
					if ctx.Err() == nil {
						u.downloaded, u.total = downloaded, total
					}
				})
			}
		})
		publish(func() {
			if ctx.Err() != nil {
				return // canceled
			}
			u.cancel = nil
			cancel()
			if err != nil {
				u.phase, u.err, u.open = updateInstallFailed, err.Error(), true
				return
			}
			u.phase, u.open = updateReady, true
		})
	}()
}

// cancelUpdate stops a check or a download.
func (a *app) cancelUpdate() {
	u := &a.updates
	if u.cancel != nil {
		u.cancel()
		u.cancel = nil
	}
	switch u.phase {
	case updateChecking:
		u.open = false
		a.settleUpdate()
	case updateDownloading:
		u.phase = updateOffer
	}
}

// settleUpdate leaves a finished check: the dialog offers the version an
// earlier check found, if any.
func (a *app) settleUpdate() {
	if a.updates.release != nil {
		a.updates.phase = updateOffer
	} else {
		a.updates.phase = updateIdle
	}
}

// remindLater closes the dialog; background checks skip this version until
// the user checks again.
func (a *app) remindLater() {
	if a.updates.release != nil {
		a.updates.later = a.updates.release.version
	}
	a.updates.open = false
}

// installing reports whether the download finished and the update is being
// put in place, which can no longer be canceled.
func (u *updateActions) installing() bool {
	return u.phase == updateDownloading && u.total > 0 && u.downloaded >= u.total
}

// updateButton checks for updates on demand, or shows how far an update
// got and reopens its dialog.
func (a *app) updateButton(c *ui.Context) {
	u := &a.updates
	if u.find == nil {
		return
	}
	c.Key("updates")
	show := func() { a.checkForUpdates(true) }
	switch {
	case u.phase == updateReady:
		accentPill(c, "Restart to update", show)
	case u.phase == updateDownloading:
		pill(c, "Downloading update…", false, show)
	case u.available != "":
		accentPill(c, "Update to "+u.available, func() { a.settleUpdate(); u.later, u.open = "", true })
	default:
		pill(c, "Check for updates", false, show)
	}
}

// accentPill is a primary button sized as the pills beside it.
func accentPill(c *ui.Context, label string, click func()) {
	ui.PrimaryButton(c, label).AlignSelf(ui.Start).Padding(9, 16).FontSize(14).OnClick(click)
}

// updateDialog shows the check, the offer, the download and the restart
// over the app. Closing it leaves a download running.
func (a *app) updateDialog(c *ui.Context) {
	u := &a.updates
	ui.DialogBase(c, &u.open, func(backdrop, dialog ui.Element) {
		backdrop.Background(ui.RGBA(0, 0, 0, 0.5))
		dialog.Width(540).MaxHeight(640).Padding(26).Gap(16).Radius(16).Background(panelFill).Border(1, hairline).Shadow(0, 12, 36, 0, ui.RGBA(0, 0, 0, 0.4))
		a.updateContent(c)
	})
}

func (a *app) updateContent(c *ui.Context) {
	u := &a.updates
	title := func(s string) { ui.Text(c, s).FontSize(20).Bold().TextColor(ink).FillWidth() }
	buttons := func(fn func()) { ui.Row(c).FillWidth().Gap(10).Justify(ui.End).Margin(6, 0, 0, 0).Children(fn) }
	done := func() { u.open = false; a.settleUpdate() }
	ok := func() { buttons(func() { ui.PrimaryButton(c, "OK").OnClick(done) }) }
	version := ""
	if u.release != nil {
		version = u.release.version
	}
	switch u.phase {
	case updateChecking:
		title("Checking for updates…")
		ui.Progress(c, -1).Label("Checking for updates")
		buttons(func() { ui.Button(c, "Cancel").OnClick(a.cancelUpdate) })
	case updateUpToDate:
		title("You're up to date")
		body(c, fmt.Sprintf("Riichi Studio %s is the newest version.", u.current))
		ok()
	case updateUnsupported:
		title("Updates aren't available")
		body(c, "This copy of Riichi Studio can't update itself. Development builds, and copies installed by a package manager, are updated another way.")
		ok()
	case updateCheckFailed:
		title("Couldn't check for updates")
		body(c, "Check your connection and try again.")
		secondary(c, u.err)
		buttons(func() {
			ui.Button(c, "Close").OnClick(done)
			ui.PrimaryButton(c, "Try again").OnClick(func() { a.settleUpdate(); a.checkForUpdates(true) })
		})
	case updateOffer, updateInstallFailed:
		title(fmt.Sprintf("Riichi Studio %s is available", version))
		if u.phase == updateInstallFailed {
			body(c, "The update couldn't be installed. Your current version is unchanged.")
			secondary(c, u.err)
		} else {
			secondary(c, fmt.Sprintf("You have %s. The update downloads in the background and takes effect when the app restarts.", u.current))
		}
		if u.release.notes != "" {
			ui.Scroll(c).FillWidth().MaxHeight(300).Padding(14, 16).Radius(10).Background(pageBottom).Border(1, hairline).Children(func() {
				if u.notes == nil {
					u.notes = &markdown.Options{}
				}
				markdown.View(c, u.release.notes, u.notes)
			})
		}
		buttons(func() {
			ui.Button(c, "Later").OnClick(a.remindLater)
			label := "Download and install"
			if u.phase == updateInstallFailed {
				label = "Try again"
			}
			ui.PrimaryButton(c, label).OnClick(a.installUpdate)
		})
	case updateDownloading:
		if u.installing() {
			title(fmt.Sprintf("Installing %s…", version))
			ui.Progress(c, -1).Label("Installing the update")
			secondary(c, "Verifying the download and putting it in place.")
		} else {
			title(fmt.Sprintf("Downloading %s…", version))
			value := -1.0
			detail := "Starting the download."
			if u.total > 0 {
				value = float64(u.downloaded) / float64(u.total)
				detail = fmt.Sprintf("%s of %s", megabytes(u.downloaded), megabytes(u.total))
			}
			ui.Progress(c, value).Label("Downloading the update")
			secondary(c, detail)
		}
		buttons(func() {
			ui.Button(c, "Hide").OnClick(func() { u.open = false })
			ui.Button(c, "Cancel").Disabled(u.installing()).OnClick(a.cancelUpdate)
		})
	case updateReady:
		title("Restart to finish updating")
		body(c, fmt.Sprintf("Riichi Studio %s is installed. Restart now, or it starts the next time you open the app. Your progress is saved.", version))
		buttons(func() {
			ui.Button(c, "Later").OnClick(func() { u.open = false })
			ui.PrimaryButton(c, "Restart now").OnClick(func() {
				if u.relaunch != nil {
					u.relaunch()
				}
			})
		})
	}
}

func megabytes(n int64) string { return fmt.Sprintf("%.1f MB", float64(n)/(1<<20)) }
