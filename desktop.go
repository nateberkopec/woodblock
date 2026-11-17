package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

func imageDir() (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cache, "woodblock")
	return dir, os.MkdirAll(dir, 0o755)
}

func apply(w wallpaper) error {
	dir, err := imageDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, w.ID+".jpg")
	if _, err := os.Stat(path); err != nil {
		if err := download(w.ImageURL, path); err != nil {
			return err
		}
	}
	if err := setEverySpace(path); err == nil {
		return nil
	}
	script := `tell application "System Events" to tell every desktop to set picture to ` + strconv.Quote(path)
	if output, err := exec.Command("osascript", "-e", script).CombinedOutput(); err != nil {
		return fmt.Errorf("osascript: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func onScreenID() string {
	output, err := exec.Command("osascript", "-e", `tell application "System Events" to get picture of desktop 1`).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(filepath.Base(strings.TrimSpace(string(output))), ".jpg")
}

func pruneImages(keep []string) error {
	dir, err := imageDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !slices.Contains(keep, strings.TrimSuffix(entry.Name(), ".jpg")) {
			os.Remove(filepath.Join(dir, entry.Name()))
		}
	}
	return nil
}

func openURL(target string) error {
	return exec.Command("open", target).Run()
}
