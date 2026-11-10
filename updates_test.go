package main

import (
	"context"
	"errors"
	"testing"

	"github.com/egoist/mygo/ui"
)

// updateTester is the app with a fake updater whose workers publish into
// a queue the test drains, so each step happens when the test says.
type updateTester struct {
	a         *app
	tt        *ui.Tester
	queue     chan func()
	finds     int
	found     *release
	findErr   error
	install   chan error
	relaunchs int
}

func newUpdateTester(t *testing.T) *updateTester {
	ut := &updateTester{a: newApp(), queue: make(chan func(), 64), install: make(chan error, 1)}
	a := ut.a
	a.update = func(fn func()) { ut.queue <- fn }
	a.updates.current = "0.1.1"
	a.updates.enabled = func() bool { return true }
	a.updates.relaunch = func() { ut.relaunchs++ }
	a.updates.find = func(context.Context) (*release, error) { ut.finds++; return ut.found, ut.findErr }
	ut.tt = ui.NewTester(a.view, 1280, 900)
	return ut
}

func (ut *updateTester) release(version string) *release {
	return &release{version: version, notes: "## What's new\n\n- Faster tables", install: func(ctx context.Context, progress func(int64, int64)) error {
		progress(5<<20, 10<<20)
		select {
		case err := <-ut.install:
			if err == nil {
				progress(10<<20, 10<<20)
			}
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}}
}

// next runs the next update the workers published.
func (ut *updateTester) next() {
	(<-ut.queue)()
	ut.tt.Frame()
}

func (ut *updateTester) click(t *testing.T, s string) {
	t.Helper()
	if err := ut.tt.Click(s); err != nil {
		t.Fatal(err)
	}
}

func (ut *updateTester) has(t *testing.T, s string) {
	t.Helper()
	if !ut.tt.HasText(s) {
		t.Fatalf("missing %q", s)
	}
}

func TestUpdateButtonNeedsUpdater(t *testing.T) {
	a := newApp()
	tt := ui.NewTester(a.view, 1280, 900)
	if _, ok := tt.Find("Check for updates"); ok {
		t.Fatal("update button shown without the updater")
	}
}

func TestBackgroundCheckOffersUpdate(t *testing.T) {
	ut := newUpdateTester(t)
	ut.found = ut.release("0.2.0")
	ut.a.checkForUpdates(false)
	ut.tt.Frame()
	if ut.tt.HasText("Checking for updates…") {
		t.Fatal("a background check showed the dialog")
	}
	ut.next()
	ut.has(t, "Riichi Studio 0.2.0 is available")
	ut.has(t, "Faster tables")

	ut.click(t, "Download and install")
	ut.next() // the first progress
	ut.has(t, "Downloading 0.2.0…")
	ut.has(t, "5.0 MB of 10.0 MB")

	ut.click(t, "Hide")
	ut.has(t, "Downloading update…")
	ut.install <- nil
	ut.next() // the last progress
	ut.next() // installed
	ut.has(t, "Restart to finish updating")
	ut.click(t, "Restart now")
	if ut.relaunchs != 1 {
		t.Fatalf("relaunched %d times, want 1", ut.relaunchs)
	}
}

func TestLaterKeepsVersionInNavigation(t *testing.T) {
	ut := newUpdateTester(t)
	ut.found = ut.release("0.2.0")
	ut.a.checkForUpdates(false)
	ut.next()
	ut.click(t, "Later")
	if ut.tt.HasText("Riichi Studio 0.2.0 is available") {
		t.Fatal("Later left the dialog open")
	}
	// The next background check does not offer the version again.
	ut.a.checkForUpdates(false)
	ut.next()
	if ut.tt.HasText("Riichi Studio 0.2.0 is available") {
		t.Fatal("a background check offered a version put off")
	}
	ut.click(t, "Update to 0.2.0")
	ut.has(t, "Riichi Studio 0.2.0 is available")
}

func TestUserCheckReportsOutcome(t *testing.T) {
	ut := newUpdateTester(t)
	ut.click(t, "Check for updates")
	ut.has(t, "Checking for updates…")
	ut.next()
	ut.has(t, "You're up to date")
	ut.has(t, "Riichi Studio 0.1.1 is the newest version.")
	ut.click(t, "OK")

	ut.findErr = errors.New("offline")
	ut.click(t, "Check for updates")
	ut.next()
	ut.has(t, "Couldn't check for updates")
	ut.has(t, "offline")
	ut.findErr, ut.found = nil, ut.release("0.2.0")
	ut.click(t, "Try again")
	ut.next()
	ut.has(t, "Riichi Studio 0.2.0 is available")
	if ut.finds != 3 {
		t.Fatalf("checked %d times, want 3", ut.finds)
	}
}

func TestBackgroundCheckFailsQuietly(t *testing.T) {
	ut := newUpdateTester(t)
	ut.findErr = errors.New("offline")
	ut.a.checkForUpdates(false)
	ut.next()
	if ut.a.updates.open {
		t.Fatal("a failed background check showed the dialog")
	}
	ut.has(t, "Check for updates")
}

func TestUnsupportedBuildSaysSo(t *testing.T) {
	ut := newUpdateTester(t)
	ut.a.updates.enabled = func() bool { return false }
	ut.click(t, "Check for updates")
	ut.has(t, "Updates aren't available")
	if ut.finds != 0 {
		t.Fatal("an unsupported build checked")
	}
}

func TestCancelDownloadReturnsToOffer(t *testing.T) {
	ut := newUpdateTester(t)
	ut.found = ut.release("0.2.0")
	ut.a.checkForUpdates(false)
	ut.next()
	ut.click(t, "Download and install")
	ut.next()
	ut.click(t, "Cancel")
	ut.has(t, "Download and install")
	ut.next() // the canceled install's result is ignored
	ut.has(t, "Download and install")

	ut.click(t, "Download and install")
	ut.next()
	ut.install <- errors.New("bad signature")
	ut.next()
	ut.has(t, "The update couldn't be installed. Your current version is unchanged.")
	ut.has(t, "bad signature")
}
