package main

import (
	"fmt"
	"io"
	"os"

	"github.com/lburfeindt/trackops/tools/format-playlist/internal/playlist"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: format-playlist <playlist>")
		return 1
	}

	file, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "format-playlist: open playlist: %v\n", err)
		return 1
	}
	defer file.Close()

	titles, err := playlist.Titles(file)
	if err != nil {
		fmt.Fprintf(stderr, "format-playlist: read playlist: %v\n", err)
		return 1
	}
	for _, title := range titles {
		fmt.Fprintln(stdout, title)
	}
	return 0
}
