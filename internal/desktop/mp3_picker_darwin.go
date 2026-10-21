//go:build darwin

package desktop

/*
#cgo LDFLAGS: -framework Cocoa -framework UniformTypeIdentifiers
#include <stdlib.h>

char *chooseMP3Inputs(void);
*/
import "C"

import (
	"context"
	"encoding/json"
	"fmt"
	"unsafe"
)

func selectMP3Inputs(_ context.Context) ([]string, error) {
	selection := C.chooseMP3Inputs()
	if selection == nil {
		return nil, fmt.Errorf("open MP3 picker failed")
	}
	defer C.free(unsafe.Pointer(selection))

	var paths []string
	if err := json.Unmarshal([]byte(C.GoString(selection)), &paths); err != nil {
		return nil, err
	}
	return paths, nil
}
