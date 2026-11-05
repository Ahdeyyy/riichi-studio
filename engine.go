package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//go:embed resources/agari/agari.exe resources/agari/LICENSE
var engineFiles embed.FS
var enginePath string

func prepareEngine() error {
	if runtime.GOOS != "windows" {
		p, e := exec.LookPath("agari")
		enginePath = p
		return e
	}
	dir, e := os.UserCacheDir()
	if e != nil {
		return e
	}
	dir = filepath.Join(dir, "riichi-tutor", "agari-0.27.0")
	if e = os.MkdirAll(dir, 0755); e != nil {
		return e
	}
	enginePath = filepath.Join(dir, "agari.exe")
	data, e := engineFiles.ReadFile("resources/agari/agari.exe")
	if e != nil {
		return e
	}
	// Always restore the known embedded engine instead of trusting cache content.
	return os.WriteFile(enginePath, data, 0700)
}

type scoreInput struct {
	hand, win, dora       string
	tsumo, riichi, dealer bool
}

func scoreHand(in scoreInput) (string, error) {
	tiles, e := parseTiles(in.hand)
	if e != nil {
		return "", e
	}
	if len(tiles) != 14 {
		return "", fmt.Errorf("Enter exactly 14 tiles, including the winning tile. This lab scores closed hands without kans.")
	}
	w, e := parseTiles(in.win)
	if e != nil || len(w) != 1 {
		return "", fmt.Errorf("Enter one winning tile, such as 8s.")
	}
	found := false
	for _, t := range tiles {
		if t == w[0] {
			found = true
		}
	}
	if !found {
		return "", fmt.Errorf("The winning tile must be part of your hand.")
	}
	if shanten(tiles) != -1 {
		return "", fmt.Errorf("This hand is not complete. Build four groups and a pair, seven pairs, or thirteen orphans.")
	}
	if in.dora != "" {
		d, e := parseTiles(in.dora)
		if e != nil {
			return "", e
		}
		if len(d) > 5 {
			return "", fmt.Errorf("Use at most five dora indicators.")
		}
		var n [34]int
		for _, t := range append(tiles, d...) {
			n[t]++
			if n[t] > 4 {
				return "", fmt.Errorf("Your hand and indicators exceed four copies of %s.", tileName(t))
			}
		}
	}
	if enginePath == "" {
		return "", fmt.Errorf("Agari is unavailable. On non-Windows systems, install the agari CLI and restart.")
	}
	args := []string{in.hand, "--win", in.win, "--round", "e", "--seat", "s", "--json", "--no-color"}
	if in.dealer {
		args[6] = "e"
	}
	if in.tsumo {
		args = append(args, "--tsumo")
	}
	if in.riichi {
		args = append(args, "--riichi")
	}
	if in.dora != "" {
		args = append(args, "--dora", in.dora)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, enginePath, args...)
	hideEngineWindow(command)
	out, e := command.CombinedOutput()
	if e != nil {
		return "", fmt.Errorf("Agari could not score this hand: %s", strings.TrimSpace(string(out)))
	}
	var result struct {
		Interpretations []struct {
			Yaku []struct {
				Name string
				Han  int
			}
			Han, Fu    int
			Dora       struct{ Total int }
			ScoreLevel string `json:"score_level"`
			Payment    struct {
				Total         int
				FromDealer    int `json:"from_dealer"`
				FromNonDealer int `json:"from_non_dealer"`
				FromDiscarder int `json:"from_discarder"`
			}
		}
	}
	if e = json.Unmarshal(out, &result); e != nil {
		return "", e
	}
	if len(result.Interpretations) == 0 {
		return "", fmt.Errorf("No scoring interpretation found.")
	}
	best := 0
	for i, r := range result.Interpretations {
		if r.Payment.Total > result.Interpretations[best].Payment.Total {
			best = i
		}
	}
	r := result.Interpretations[best]
	if len(r.Yaku) == 0 {
		return "", fmt.Errorf("The shape is complete but has no yaku. Dora alone cannot qualify a win.")
	}
	text := fmt.Sprintf("%d han / %d fu\n\n", r.Han, r.Fu)
	for _, y := range r.Yaku {
		text += fmt.Sprintf("%s   %d han\n", y.Name, y.Han)
	}
	if r.Dora.Total > 0 {
		text += fmt.Sprintf("Dora bonuses   %d han\n", r.Dora.Total)
	}
	text += fmt.Sprintf("\n%d points total", r.Payment.Total)
	if in.tsumo {
		if in.dealer {
			text += fmt.Sprintf("\n%d from each opponent", r.Payment.Total/3)
		} else {
			text += fmt.Sprintf("\n%d from East; %d from each other player", r.Payment.FromDealer, r.Payment.FromNonDealer)
		}
	} else {
		text += "\nPaid by the discarder"
	}
	return text, nil
}
func (a *app) calculate(in scoreInput) {
	if a.busy {
		return
	}
	a.busy = true
	a.scoreResult = ""
	a.scorePage = a.page
	publish := a.update
	go func() {
		result, e := scoreHand(in)
		if e != nil {
			result = e.Error()
		}
		if publish != nil {
			publish(func() { a.scoreResult = result; a.busy = false })
		}
	}()
}
