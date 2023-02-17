package sqlite

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDBFile(t *testing.T) {
	dir, _ := os.Getwd()
	testCases := []struct {
		caseStr        string
		dirPath        string
		wantDBFilePath string
		err            error
	}{
		{
			caseStr:        "successful create db file with dir suffix",
			dirPath:        fmt.Sprintf("%s/", dir),
			wantDBFilePath: fmt.Sprintf("%s/%s", dir, SqliteDBFileName),
			err:            nil,
		},
		{
			caseStr:        "successful create db file without suffix",
			dirPath:        dir,
			wantDBFilePath: fmt.Sprintf("%s/%s", dir, SqliteDBFileName),
			err:            nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseStr, func(t *testing.T) {
			dbFilePath, err := PrepareSqliteDbFile(tc.dirPath)
			if err != nil {
				assert.Equal(t, true, strings.Contains(err.Error(), tc.err.Error()))
			}
			assert.Equal(t, tc.wantDBFilePath, dbFilePath)
			// clean up created file.
			_ = os.Remove(dbFilePath)
		})
	}
}
