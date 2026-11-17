package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	AccessKey   string   `toml:"access_key"`
	Queries     []string `toml:"queries"`
	Accounts    []string `toml:"accounts"`
	Orientation string   `toml:"orientation"`
	HistorySize int      `toml:"history_size"`
}

var defaultQueries = []string{"woodblock print", "japanese print"}

var defaultAccounts = []string{
	"artchicago", "austriannationallibrary", "aucklandmuseum", "birminghammuseumstrust",
	"bostonpubliclibrary", "britishlibrary", "clevelandart", "danishlibrary", "denverartmuseum",
	"galtmuseum", "gofyn", "libraryofcongress", "mcgilllibrary", "metropolitanmuseum",
	"nationalgallery", "nationallibraryofaustralia", "nationallibrarynorway",
	"natlibraryofscotland", "natmus", "nhm", "nypl", "seamuseum", "tepapa", "thewalters",
}

func loadConfig() (config, error) {
	cfg := config{
		Queries:     defaultQueries,
		Accounts:    defaultAccounts,
		Orientation: "landscape",
		HistorySize: 100,
	}
	path := filepath.Join(xdgDir("XDG_CONFIG_HOME", ".config"), "woodblock", "config.toml")
	if _, err := toml.DecodeFile(path, &cfg); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return cfg, err
	}
	if key := os.Getenv("UNSPLASH_ACCESS_KEY"); key != "" {
		cfg.AccessKey = key
	}
	if cfg.AccessKey == "" {
		return cfg, errors.New("set access_key in " + path + " or UNSPLASH_ACCESS_KEY")
	}
	if len(cfg.Queries) == 0 {
		return cfg, errors.New("queries must not be empty in " + path)
	}
	return cfg, nil
}

func xdgDir(env, fallback string) string {
	if dir := os.Getenv(env); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, fallback)
}
