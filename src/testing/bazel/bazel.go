package testing

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/tricorder/src/utils/common"
	"github.com/tricorder/src/utils/file"
	"github.com/tricorder/src/utils/log"
)

// Bazel-specific testing APIs.

// TestFilePath returns the absolute path to the path relative to the root of the repo.
func TestFilePath(repoRootRelPath string) string {
	// Note that we need to append tricorder as the repo's root directory.
	repoRelPath := fmt.Sprintf("tricorder/%s", repoRootRelPath)

	runFilePath, err := runfiles.Rlocation(repoRelPath)
	if err != nil {
		log.Fatalf("Could not find runfile file '%s', error: %v", repoRootRelPath, err)
	}
	if !file.Exists(runFilePath) {
		log.Fatalf("Runfile '%s' obtained from '%s' does not exist", runFilePath, repoRootRelPath)
	}
	return runFilePath
}

// TestBinaryPath returns the absolute path to the path relative to the root of the repo.
func TestBinaryPath(repoRootRelPath string) string {
	pathForMyTool, ok := bazel.FindBinary(filepath.Dir(repoRootRelPath), filepath.Base(repoRootRelPath))
	if !ok {
		log.Fatalf("Could not find binary, %s", repoRootRelPath)
	}
	return pathForMyTool
}

// CreateTmpDir returns a path to a newly created temporary directory.
func CreateTmpDir() string {
	prefix := "tricorder-"
	dir, err := bazel.NewTmpDir(prefix)
	if err != nil {
		log.Fatalf(
			"While creating tmp dir with prefix '%s', bazel failed to create the directory, error: %v",
			prefix,
			err,
		)
	}
	return dir
}

// GetTmpFile returns a random file path under temp directory.
// The file is not created.
func GetTmpFile() string {
	return path.Join(CreateTmpDir(), common.RandStr(10))
}

// CreateTmpFile returns a path to a file under the temporary directory.
// Also returns a function that delete the file, so you can use defer to automate the cleanup:
// f, cleaner := CreateTmpFile()
// defer cleaner()
func CreateTmpFile() string {
	f := GetTmpFile()
	openedFile, err := os.Create(f)
	if err != nil {
		log.Fatalf("While creating temp file at '%s', failed to create the file, error: %v", f, err)
	}
	err = openedFile.Close()
	if err != nil {
		log.Fatalf("While creating temp file at '%s', failed to close the file after creation, error: %v", f, err)
	}
	return f
}

func CreateTmpFileWithContent(content string) string {
	f := GetTmpFile()
	const defaultPerf = 0o644
	err := os.WriteFile(f, []byte(content), defaultPerf)
	if err != nil {
		log.Fatalf("Failed to write to file '%s', error: %v", f, err)
	}
	return f
}

// ReadTestFile reads a test file pointed to by a relative path.
func ReadTestFile(relPath string) (string, error) {
	return file.Read(TestFilePath(relPath))
}

// ReadTestBinFile reads a test file pointed to by a relative path and returns byte slice.
func ReadTestBinFile(relPath string) ([]byte, error) {
	return file.ReadBin(TestFilePath(relPath))
}
