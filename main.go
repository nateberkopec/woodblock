package main

import (
	"fmt"
	"io"
	"os"
)

var version = "dev"

func main() {
	command, args := "show", os.Args[1:]
	if len(args) > 0 {
		command, args = args[0], args[1:]
	}

	var err error
	switch command {
	case "show":
		err = show()
	case "open":
		err = openLink(args)
	case "next":
		err = next()
	case "back":
		err = back()
	case "version", "--version":
		fmt.Println(version)
	case "help", "-h", "--help":
		usage(os.Stdout)
	default:
		usage(os.Stderr)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "woodblock:", err)
		os.Exit(1)
	}
}

func usage(w io.Writer) {
	fmt.Fprint(w, `usage: woodblock [command]

commands:
  show                    show the current wallpaper (default)
  open [source|unsplash]  open the source page (e.g. a museum), or the Unsplash page
  next                    set a new wallpaper
  back                    go back to the previous wallpaper
  version                 print the version
`)
}
