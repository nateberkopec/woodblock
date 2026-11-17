package main

import (
	"errors"
	"fmt"
)

var errNoWallpaper = errors.New("no wallpaper set yet; run woodblock next")

func show() error {
	h, err := loadHistory()
	if err != nil {
		return err
	}
	w, ok := h.current()
	if !ok {
		return errNoWallpaper
	}
	fmt.Println(w.Description)
	fmt.Println()
	fmt.Println("Credit:   " + w.Credit)
	fmt.Println("Unsplash: " + w.UnsplashURL)
	if w.SourceURL != "" {
		fmt.Println("Source:   " + w.SourceURL)
	}
	fmt.Println("Set:      " + w.SetAt.Local().Format("2006-01-02 15:04"))
	return nil
}

func openLink(args []string) error {
	h, err := loadHistory()
	if err != nil {
		return err
	}
	w, ok := h.current()
	if !ok {
		return errNoWallpaper
	}
	target := w.SourceURL
	switch {
	case len(args) == 1 && args[0] == "unsplash", target == "":
		target = w.UnsplashURL
	case len(args) == 1 && args[0] == "source":
	case len(args) > 0:
		return errors.New("usage: woodblock open [source|unsplash]")
	}
	return openURL(target)
}

func next() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	h, err := loadHistory()
	if err != nil {
		return err
	}
	u := newUnsplash(cfg.AccessKey)
	p, err := pick(u, cfg, h.ids())
	if err != nil {
		return err
	}
	w := newWallpaper(p)
	if err := apply(w); err != nil {
		return err
	}
	if err := u.trackDownload(p); err != nil {
		return err
	}
	h.push(w, cfg.HistorySize)
	if err := h.save(); err != nil {
		return err
	}
	if err := pruneImages(h.ids()); err != nil {
		return err
	}
	return show()
}

func back() error {
	h, err := loadHistory()
	if err != nil {
		return err
	}
	if h.Current < 1 || h.Current >= len(h.Wallpapers) {
		return errors.New("no earlier wallpaper in history")
	}
	h.Current--
	if err := apply(h.Wallpapers[h.Current]); err != nil {
		return err
	}
	if err := h.save(); err != nil {
		return err
	}
	return show()
}
