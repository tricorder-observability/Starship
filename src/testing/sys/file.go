package sys

import (
	"os"

	"github.com/tricorder/src/utils/log"
)

// MustRemoveAll crashes if failed to remove the path.
func MustRemoveAll(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatalf("Failed to remove '%s', error: %v", path, err)
	}
}
