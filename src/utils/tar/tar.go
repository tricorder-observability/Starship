// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"path"

	"github.com/tricorder/src/utils/file"
)

// GZExtract extracts *.tar.gz file to destDir
func GZExtract(tarFilePath string, destDir string) error {
	reader, closer, err := file.Reader(tarFilePath)
	if err != nil {
		return fmt.Errorf("could not reader %s error: %v", tarFilePath, err)
	}

	defer closer.Close()

	uncompressedStream, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("could not craete gzip reader %s error: %v", tarFilePath, err)
	}

	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("iterate tar %s error: %v", tarFilePath, err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			expectedFile := path.Join(destDir, header.Name)
			err := file.Create(expectedFile)
			if err != nil {
				return fmt.Errorf("could not create file %s error: %v", expectedFile, err)
			}
			expectedWriter, expectedCloser, err := file.Writer(expectedFile)
			if err != nil {
				expectedCloser.Close()
				return fmt.Errorf("could not writer %s error: %v", expectedFile, err)
			}
			if _, err := io.Copy(expectedWriter, tarReader); err != nil {
				expectedCloser.Close()
				return fmt.Errorf("could not copy content from %s to %s error: %v", header.Name, expectedFile, err)
			}
			expectedCloser.Close()

		default:
			continue
		}
	}

	return nil
}
