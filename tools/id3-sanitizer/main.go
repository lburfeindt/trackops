package main

import (
	"fmt"
	"io"
	"os"

	"github.com/lburfeindt/trackops/tools/id3-sanitizer/internal/sanitizer"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: id3-sanitizer <mp3>")
		return 1
	}
	changed, err := sanitizer.File(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "id3-sanitizer: %v\n", err)
		return 1
	}
	status := "unchanged"
	if changed {
		status = "updated"
	}
	fmt.Fprintf(stdout, "%s: %s\n", status, args[0])
	return 0
}
