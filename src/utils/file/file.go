// Package file provides utility APIs for interacting with file system
package file

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

const defaultPerm = 0o644

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func List(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	res := []string{}
	for _, f := range files {
		res = append(res, f.Name())
	}
	return res
}

// Create a file with the specified path.
// All parent dirs will be created, and the filename is used to create a file.
func Create(p string) error {
	err := os.MkdirAll(path.Dir(p), 0o777)
	if err != nil {
		return fmt.Errorf("while creating file at '%s', couldn't make dir, error: %v", p, err)
	}
	_, err = os.Create(p)
	if err != nil {
		return fmt.Errorf("while creating file at '%s', couldn't create the file under base, error: %v", p, err)
	}
	return nil
}

// CreateDir with the specified path.
// All parent dirs will be created
func CreateDir(p string) error {
	if Exists(p) {
		return nil
	}
	err := os.MkdirAll(p, 0o777)
	if err != nil {
		return fmt.Errorf("while creating dir '%s', couldn't make dir, error: %v", p, err)
	}
	return nil
}

// Append writes content to the specified file, return error if failed.
func Append(p string, content string) error {
	if !Exists(p) {
		return fmt.Errorf("while appending to file '%s', file does not exit", p)
	}
	// https://pkg.go.dev/os#pkg-constants
	f, err := os.OpenFile(p,
		// Append only
		os.O_APPEND|
			// Write only
			os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("while appending to file '%s', filed to open it, error: %v", p, err)
	}
	defer f.Close()

	s, err := f.WriteString(content)
	if err != nil {
		return fmt.Errorf("while appending to file '%s', failed to write to it, error: %v", p, err)
	}
	if s != len(content) {
		return fmt.Errorf("while appending to file '%s', only wrote %d bytes out of %d", p, s, len(content))
	}
	return nil
}

// ReadBin wraps os.ReadFile().
func ReadBin(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// Read returns content as string.
func Read(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	return string(content), err
}

// Write writes content to the file at filePath.
func Write(filePath string, content string) error {
	if !Exists(filePath) {
		if err := Create(filePath); err != nil {
			return fmt.Errorf("while writing file, create %s file failed error: %v", filePath, err)
		}
	}
	return os.WriteFile(filePath, []byte(content), defaultPerm)
}

// ReadLines reads content of the input file into string slice. Each line is one element of the result.
func ReadLines(filePath string) ([]string, error) {
	result, err := Read(filePath)
	if err != nil {
		return nil, fmt.Errorf("while reading file into lines, failed to read file, error: %v", err)
	}
	return strings.Split(result, "\n"), nil
}

// Copy copy srcPath file to dstPath
func Copy(srcPath string, dstPath string) error {
	if !Exists(srcPath) {
		return fmt.Errorf("while copying file , %s file does not exit", srcPath)
	}

	err := Create(dstPath)
	if err != nil {
		return fmt.Errorf("while copying file, create %s file failed error: %v", dstPath, err)
	}

	buf, err := Read(srcPath)
	if err != nil {
		return fmt.Errorf("while copying file, reading %s file failed error: %v", srcPath, err)
	}

	err = Write(dstPath, buf)
	if err != nil {
		return fmt.Errorf("while copying file, writing %s file failed error: %v", dstPath, err)
	}
	return nil
}

// Reader return reader and closer object
func Reader(filePath string) (io.Reader, io.Closer, error) {
	filePathDir := path.Dir(filePath)
	if !Exists(filePathDir) {
		if err := os.MkdirAll(filePathDir, 0o777); err != nil {
			return nil, nil, fmt.Errorf("while creating dir '%s', couldn't make dir, error: %v", filePathDir, err)
		}
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("while read file to reader, failed to open file, error: %v", err)
	}

	return file, file, nil
}

// Writer return writer and closer object
func Writer(filePath string) (io.Writer, io.Closer, error) {
	filePathDir := path.Dir(filePath)
	if !Exists(filePathDir) {
		err := os.MkdirAll(filePathDir, 0o777)
		if err != nil {
			return nil, nil, fmt.Errorf("while creating dir '%s', couldn't make dir, error: %v", filePathDir, err)
		}
	}

	file, err := os.OpenFile(filePath,
		// Write only
		os.O_WRONLY, 0o600)
	if err != nil {
		return nil, nil, fmt.Errorf("while write file to writer, failed to open file, error: %v", err)
	}

	return file, file, nil
}

// ReadSymLink read the symbolic link
func ReadSymLink(linkPath string) (string, error) {
	return os.Readlink(linkPath)
}

// CreateSymLink create dstPath symlink to srcPath
func CreateSymLink(srcPath, dstPath string) error {
	if !Exists(srcPath) {
		return fmt.Errorf("while create symbol link, %s file does not exit", srcPath)
	}

	tmpDir := path.Dir(dstPath)
	if !Exists(tmpDir) {
		err := CreateDir(tmpDir)
		if err != nil {
			return err
		}
	}

	return os.Symlink(srcPath, dstPath)
}

// Contains checks if the file contains the specified string.
func Contains(filePath, content string) bool {
	contents, err := Read(filePath)
	if err != nil {
		return false
	}
	return strings.Contains(contents, content)
}
