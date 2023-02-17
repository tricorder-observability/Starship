package sqlite

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	// Import sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	SqliteDBFileName = "tricorder.db"
)

// PrepareSqliteDbFile will prepare a sqlite db file according the specified dir path.
// sqlite db file absolute path will be returned.
func PrepareSqliteDbFile(dirPath string) (string, error) {
	// check is the dir is existed.
	if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
		zap.S().Infof("directory: %s does not exist, create it now.", dirPath)
		// If it is a multi tier folder, recursively create all folders
		// otherwise an error will be reported if the folder does not exist
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			zap.S().Panicf("directory: %s does not exist and could not create it either.", dirPath)
			return "", err
		}
	}

	// check the db file is existed, if not, create it.
	var sqliteDbFilePath string
	if strings.HasSuffix(dirPath, "/") {
		sqliteDbFilePath = fmt.Sprintf("%s%s", dirPath, SqliteDBFileName)
	} else {
		sqliteDbFilePath = fmt.Sprintf("%s/%s", dirPath, SqliteDBFileName)
	}

	_, err := os.Open(sqliteDbFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			zap.S().Infof("Creating %s...", sqliteDbFilePath)
			// Create SQLite file
			file, err := os.Create(SqliteDBFileName)
			if err != nil {
				zap.S().Panicf("can not create sqlite db file: %s.", err.Error())
			}
			_ = file.Close()
			zap.S().Infof("%s created.", sqliteDbFilePath)
			return sqliteDbFilePath, nil
		}
		zap.S().Panic("do not have the create permissions.")
		return "", err
	}

	// db file is existed, return directly.
	return sqliteDbFilePath, nil
}
