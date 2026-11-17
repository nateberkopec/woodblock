package main

import (
	"errors"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"howett.net/plist"
)

func wallpaperIndexPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "com.apple.wallpaper", "Store", "Index.plist")
}

func setEverySpace(image string) error {
	path := wallpaperIndexPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var index map[string]any
	if _, err := plist.Unmarshal(data, &index); err != nil {
		return err
	}
	configuration, err := plist.Marshal(map[string]any{
		"type": "imageFile",
		"url":  map[string]any{"relative": (&url.URL{Scheme: "file", Path: image}).String()},
	}, plist.BinaryFormat)
	if err != nil {
		return err
	}
	choice := map[string]any{
		"Provider":      "com.apple.wallpaper.choice.image",
		"Files":         []any{},
		"Configuration": configuration,
	}
	now := time.Now()
	updated := 0
	setDesktop := func(entry any) {
		desktop, ok := dig(entry, "Desktop")
		if !ok {
			return
		}
		content, ok := dig(desktop, "Content")
		if !ok {
			return
		}
		content["Choices"] = []any{choice}
		content["Shuffle"] = "$null"
		desktop["LastSet"] = now
		updated++
	}
	setDisplays := func(entry any) {
		displays, _ := dig(entry, "Displays")
		for _, display := range displays {
			setDesktop(display)
		}
	}

	setDesktop(index["SystemDefault"])
	setDesktop(index["AllSpacesAndDisplays"])
	setDisplays(index)
	spaces, _ := dig(index, "Spaces")
	for _, space := range spaces {
		setDesktop(space.(map[string]any)["Default"])
		setDisplays(space)
	}
	if updated == 0 {
		return errors.New("no desktops found in " + path)
	}

	out, err := plist.Marshal(index, plist.BinaryFormat)
	if err != nil {
		return err
	}
	tmp := path + ".woodblock"
	if err := os.WriteFile(tmp, out, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		return err
	}
	return exec.Command("killall", "WallpaperAgent").Run()
}

func dig(entry any, key string) (map[string]any, bool) {
	parent, ok := entry.(map[string]any)
	if !ok {
		return nil, false
	}
	child, ok := parent[key].(map[string]any)
	return child, ok
}
