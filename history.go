package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type wallpaper struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Credit      string    `json:"credit"`
	UnsplashURL string    `json:"unsplash_url"`
	SourceURL   string    `json:"source_url,omitempty"`
	ImageURL    string    `json:"image_url"`
	SetAt       time.Time `json:"set_at"`
}

type history struct {
	Wallpapers []wallpaper `json:"wallpapers"`
	Current    int         `json:"current"`
}

func historyPath() string {
	return filepath.Join(xdgDir("XDG_STATE_HOME", ".local/state"), "woodblock", "history.json")
}

func loadHistory() (history, error) {
	var h history
	data, err := os.ReadFile(historyPath())
	if errors.Is(err, fs.ErrNotExist) {
		return h, nil
	}
	if err != nil {
		return h, err
	}
	return h, json.Unmarshal(data, &h)
}

func (h history) save() error {
	path := historyPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func (h history) current() (wallpaper, bool) {
	if h.Current < 0 || h.Current >= len(h.Wallpapers) {
		return wallpaper{}, false
	}
	return h.Wallpapers[h.Current], true
}

func (h history) ids() []string {
	ids := make([]string, len(h.Wallpapers))
	for i, w := range h.Wallpapers {
		ids[i] = w.ID
	}
	return ids
}

func (h *history) push(w wallpaper, size int) {
	h.Wallpapers = append(h.Wallpapers, w)
	if len(h.Wallpapers) > size {
		h.Wallpapers = h.Wallpapers[len(h.Wallpapers)-size:]
	}
	h.Current = len(h.Wallpapers) - 1
}
