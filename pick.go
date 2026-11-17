package main

import (
	"errors"
	"math/rand/v2"
	"regexp"
	"slices"
	"strings"
	"time"
)

const pageAttempts = 4

var (
	photoID   = regexp.MustCompile(`\A[A-Za-z0-9_-]+\z`)
	sourceURL = regexp.MustCompile(`https?://\S+`)
)

func pick(u unsplash, cfg config, seen []string) (photo, error) {
	query := cfg.Queries[rand.IntN(len(cfg.Queries))]
	eligible := func(results []photo) []photo {
		return slices.DeleteFunc(slices.Clone(results), func(p photo) bool {
			if !photoID.MatchString(p.ID) || slices.Contains(seen, p.ID) {
				return true
			}
			return len(cfg.Accounts) > 0 && !slices.Contains(cfg.Accounts, p.User.Username)
		})
	}

	first, err := u.search(query, cfg.Orientation, 1)
	if err != nil {
		return photo{}, err
	}
	for range pageAttempts {
		results := first.Results
		if page := 1 + rand.IntN(max(first.TotalPages, 1)); page > 1 {
			next, err := u.search(query, cfg.Orientation, page)
			if err != nil {
				return photo{}, err
			}
			results = next.Results
		}
		if candidates := eligible(results); len(candidates) > 0 {
			return candidates[rand.IntN(len(candidates))], nil
		}
	}
	if candidates := eligible(first.Results); len(candidates) > 0 {
		return candidates[rand.IntN(len(candidates))], nil
	}
	return photo{}, errors.New("no unseen photos from allowed accounts for query " + query)
}

func newWallpaper(p photo) wallpaper {
	description := strings.TrimSpace(p.Description)
	if description == "" {
		description = p.AltDescription
	}
	return wallpaper{
		ID:          p.ID,
		Description: description,
		Credit:      p.User.Name,
		UnsplashURL: p.Links.HTML,
		SourceURL:   strings.TrimRight(sourceURL.FindString(p.Description), ".,;)"),
		ImageURL:    p.URLs.Full,
		SetAt:       time.Now(),
	}
}
