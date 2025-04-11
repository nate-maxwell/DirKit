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

// Helper function for determining if a path exists on disk or not.
// Args:
//
//	path(string): The path to check
//
// Returns:
//
//	bool: True if the path exists on disk else false.
//	error: The fs.ErrNotExist error if the path does not exist else nil.
func pathExists(path string) (bool, error) {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return false, err
	}
	return true, nil
}

// A helper function to determine if a path is a directory not not.
// Args:
//
//	path(string): The path to check.
//
// Returns:
//
//	bool: True if the path is a directory else false.
//	error: A os.IsNotExist error if the path does not exists, else nil.
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

// Gets the content names, or full path for contents, of a directory.
// Args:
//
//	path(string): Directory path to list the contents of.
//	fullPath(bool): To return string names or full paths of directory contents.
//
// Returns:
//
//	[]string: String names or full paths of directory contents.
//	error: Any error created from attempting to read the directory, else nil.
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

// Creates a directory from teh given path.
// Args:
//
//	path(string): The directory path to create.
//
// Returns:
//
//	error: Any error created while attempting to create the directory, else nil.
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

// Creates a directory with today's date as the name.
// Args:
//
//	path(styring): The path to create the new folder in.
//
// Returns:
//
//	error: Any error created while attempting to create the directory, else nil.
func CreateDatedDirectory(path string) error {
	datePath := filepath.Join(path, GetDate())
	err := CreateDirectory(datePath)
	if err != nil {
		return err
	}
	return nil
}

// Deletes a directory and its contents as long as they are within the safety path.
// Args:
//
//	folderPath(string): The folder path to delete.
//
// Returns:
//
//	error: the *PathError created from os.RemoveAll if one was created, else nil.
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

// Removes specified file as long as it is within the safety path.
// Args:
//
//	filepath(string): The path to the file you wish to delete.
//
// Returns:
//
//	error: A custom error if the filepath was not within the safety path or a *PathError err from
//	os.Remove, else Nil.
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

// Delete all files in a directory as long as they are within the safety path.
// Args:
//
//	directory_path(string): The path to the directory.
//
// Returns:
//
//	any *PathError crated from DeleteSafeFile or errors from GetDirContents, else nil.
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

// Copy file into a separate destination folder.
// Args:
//
//	source(string): File path of the file to copy.
//	dest(string): File path to copy the file too, optionally can have different name.
//
// Returns:
//
//	error: *PathError crated from os module or possible other error from io module else nil.
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

// Copy contents of a folder to the given destination.
// Args:
//
//	sourcePath(string): Folder path to the folder that is to be copied.
//	destination(string): Folder path to copy the folder + contents to.
//
// Returns:
//
//	error: Any relevant errors created durring process, usually os *PathErrors else nil.
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

// Returns string: 'yyyymmdd'.
func GetDate() string {
	return time.Now().Format("20060102")
}

// Returns string: 'HH:MM:SS:XX', X is microsecond.
func GetTime() string {
	return time.Now().Format("15:04:05:00")
}

// Exports a string map to json file path.
// Args:
//
//	fielpath(string): The file path to place the .json file.
//	data(map[string]interface{}): Any map with string keys and values that can be converted to strings.
//	overWrite(bool): To overwrite json file if it already exists in path.
//
// Returns:
//
//	error: Any relevant error from the json handling or file writing process.
func ExportMapToJson(filePath string, data map[string]interface{}, overWrite bool) error {
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
