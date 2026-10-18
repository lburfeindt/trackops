package main

import (
	"fmt"
	"os"

	"github.com/lburfeindt/trackops/tools/id3-sanitizer/internal/gui"
	"github.com/ncruces/zenity"
)

func main() {
	if err := gui.Run(gui.NativeDialogs{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		_ = zenity.Error(err.Error(), zenity.Title("ID3 Sanitizer"))
		os.Exit(1)
	}
}
