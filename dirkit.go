// # Dir Kit
//
// * A simple toolkit for folder and file handling that eliminates
// boilerplate or wraps commonly used functions in a consistent
// namespace for easy rememberance/importing.

package dirkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var safetyPath string = "D:/safety/" // Change on per-project needs

func pathExists(path string) (bool, error) {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return false, err
	}
	return true, nil
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func GetDirContents(path string, fullPath bool) ([]string, error) {
	var contents []string

	items, err := os.ReadDir(path)
	if err != nil {
		return make([]string, 0), err
	}
	for _, item := range items {
		var entry string
		if fullPath {
			entry = fmt.Sprintf("%s%s", path, item.Name())
		} else {
			entry = item.Name()
		}
		contents = append(contents, entry)
	}
	return contents, nil
}

func CreateDirectory(path string) error {
	exists, _ := pathExists(path)
	if !exists {
		err := os.Mkdir(path, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateDatedDirectory(path string) error {
	datePath := filepath.Join(path, GetDate())
	err := CreateDirectory(datePath)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSafeDirectory(folderPath string) error {
	if strings.HasPrefix(folderPath, safetyPath) {
		err := os.RemoveAll(folderPath)
		if err != nil {
			return err
		}
		return nil
	}
	errorMsg := fmt.Sprintf("folder path is not within %s", safetyPath)
	return errors.New(errorMsg)
}

func DeleteSafeFile(filepath string) error {
	if strings.HasPrefix(filepath, safetyPath) {
		err := os.Remove(filepath)
		if err != nil {
			return err
		}
		return nil
	}
	errorMsg := fmt.Sprintf("file path is not within %s", safetyPath)
	return errors.New(errorMsg)
}

func DeleteSafeFilesInDirectory(folderPath string) error {
	if strings.HasPrefix(folderPath, safetyPath) {
		files, err := GetDirContents(folderPath, true)
		if err != nil {
			return err
		}
		for _, file := range files {
			err := DeleteSafeFile(file)
			if err != nil {
				return err
			}
		}
		return nil
	}
	errorMsg := fmt.Sprintf("file path is not within %s", safetyPath)
	return errors.New(errorMsg)
}

func CopyFile(source string, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func CopyFolderContents(sourcePath string, destination string) error {
	sourcePath = filepath.Clean(sourcePath)
	destination = filepath.Clean(destination)

	err := CreateDirectory(destination)
	if err != nil {
		return err
	}

	curItems, err := GetDirContents(sourcePath, false)
	if err != nil {
		return err
	}

	for _, item := range curItems {
		curItemPath := filepath.Clean(filepath.Join(sourcePath, item))
		destPath := filepath.Clean(filepath.Join(destination, item))

		dir, err := isDir(curItemPath)
		if err != nil {
			return err
		}
		if dir {
			err := CopyFolderContents(curItemPath, destPath)
			if err != nil {
				return err
			}
		} else {
			err := CopyFile(curItemPath, destPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetDate() string {
	return time.Now().Format("20060102")
}

func GetTime() string {
	return time.Now().Format("15:04:05:00")
}

func ExportDataToJson(filePath string, data map[string]interface{}, overWrite bool) error {
	exists, err := pathExists(filePath)
	if err != nil {
		return err
	}

	if !exists || overWrite {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.Write(jsonData)
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}
