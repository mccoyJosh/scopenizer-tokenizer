package util

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
)

var PathSeparator = func() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n\n", fmt.Sprintf("error: %s", err))
	debug.PrintStack()

	os.Exit(1)
}

// CheckIfErrorWithMessage should be used to naively panics if an error is not nil with a message.
func CheckIfErrorWithMessage(err error, message string) {
	if err == nil {
		return
	}

	fmt.Println(message)
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n\n", fmt.Sprintf("error: %s", err))
	debug.PrintStack()

	os.Exit(1)
}

const OwnerPermRw = 0777

func MakeDir(path string) {
	if err := os.MkdirAll(path, OwnerPermRw); err != nil {
		if !os.IsExist(err) {
			CheckIfErrorWithMessage(err, "Unable to create directory")
		}
		err = os.Chmod(path+PathSeparator()+"submission", OwnerPermRw)
		if err != nil {
			CheckIfErrorWithMessage(err, "Unable to update directory mode")
		}
	}
}

func CheckIfDirExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func RemoveDir(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		Error(fmt.Sprintf("Error in deleting directory %s", path), err)
	} else {
		Debug(fmt.Sprintf("Successfully deleted directory %s", path))
	}
}

// CreateFileWithInfo
// Creates a file (with the full path provided) and adds the
// string to that file
func CreateFileWithInfo(fullPath string, data string) error {
	// CREATES THE FILE
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	// WRITES TO THE FILE
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func MoveDirectory(originalPath string, moveToPath string) error {
	err := os.Rename(originalPath, moveToPath)
	return err
}

func RemoveFile(path string) {
	err := os.Remove(path)
	if err != nil {
		Error(fmt.Sprintf("Failed to remove file %s", path), err)
	}
}

func RemoveDirectoriesInThisDirectory(path string) {
	err := filepath.WalkDir(path, func(pathTemp string, d fs.DirEntry, err error) error {
		if path != pathTemp {
			if d.IsDir() {
				RemoveDir(pathTemp)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func RemoveEverythingInThisDirectory(path string, excludeExtensionRegex string) {
	err := filepath.WalkDir(path, func(pathTemp string, d fs.DirEntry, err error) error {
		if path != pathTemp {
			if d.IsDir() {
				RemoveDir(pathTemp)
			} else {
				removeItem := true
				if excludeExtensionRegex != "" {
					b, err := regexp.MatchString(excludeExtensionRegex, d.Name())
					if err != nil {
						Error("Regex In Removal Of Everything In This Directory", err)
					} else if b {
						removeItem = !b
					}
				}
				if removeItem {
					RemoveFile(pathTemp)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func ZipSource(source, target string) error {
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func GetContentOfFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		Error("Failed to find file: " + path, err)
	}
	return string(content)
}

// ReplaceSnippetFromFilesInDirectory
// This will go through all the files in a given directory
// (given the fileRegex can be found in the file name) and
// will find and replace the provided snippet from the file
// with the given replaceSnippet
func ReplaceSnippetFromFilesInDirectory(dirPath string, fileRegex string, removeSnippet string, replaceSnippet string) error {
	err := filepath.WalkDir(dirPath, func(pathTemp string, d fs.DirEntry, err error) error {
		if strings.Contains(d.Name(), fileRegex) {
			originalContent := GetContentOfFile(pathTemp)
			newContent := strings.Replace(originalContent, removeSnippet, replaceSnippet, -1)

			err = CreateFileWithInfo(pathTemp, newContent)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err

}

func ReplaceStringBetweenStringsInclusiveFromFilesInDirectory(dirPath string, fileRegex string, startString string, endString string, replaceSnippet string) error {
	err := filepath.WalkDir(dirPath, func(pathTemp string, d fs.DirEntry, err error) error {
		if strings.Contains(d.Name(), fileRegex) {
			content := GetContentOfFile(pathTemp)
			continueReplacing := true
			for continueReplacing {
				positionOfStartString := strings.Index(content, startString)
				if positionOfStartString == -1 {
					continueReplacing = false
					continue
				}

				positionOfEndString := -1
				for i := positionOfStartString; i+len(endString) < len(content) && positionOfEndString == -1; i++ {
					if content[i:i+len(endString)] == endString {
						positionOfEndString = i
					}
				}
				locatedString := ""

				if positionOfEndString >= 0 && positionOfStartString >= 0 {
					locatedString = content[positionOfStartString : positionOfEndString+len(endString)]
					content = strings.Replace(content, locatedString, replaceSnippet, -1)
				} else {
					continueReplacing = false
				}
			}

			err = CreateFileWithInfo(pathTemp, content)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err

}

func FindValueBetweenStringsFromFilesInDirectory(dirPath string, fileRegex string, startString string, endString string) (string, bool, error) {
	stringFound := ""
	itemFound := false

	err := filepath.WalkDir(dirPath, func(pathTemp string, d fs.DirEntry, err error) error {
		if strings.Contains(d.Name(), fileRegex) {
			content := GetContentOfFile(pathTemp)

			positionOfStartString := strings.Index(content, startString)
			positionOfEndString := strings.Index(content, endString)

			if positionOfEndString >= 0 && positionOfStartString >= 0 {
				itemFound = true
				stringFound = content[positionOfStartString+len(startString) : positionOfEndString]
			}
		}
		return nil
	})
	return stringFound, itemFound, err

}

// GetTextOfFile
// Reads and returns the text found in a file
func GetTextOfFile(filePath string) (string, error) {
	ret := ""
	var err error = nil
	content, thisError := ioutil.ReadFile(filePath)
	if thisError != nil {
		err = thisError
	} else {
		ret = string(content)
	}
	return ret, err
}

func RunCommandWithSingleArgGetOutputIgnoreErrors(command string, arg string, runLocation string) string {
	cmd := exec.Command(command, arg)
	cmd.Dir = runLocation
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func GetNamesOfFilesAndDirs(dirPath string, depth int) ([]string, []string, error) {
	dirs := make([]string, 0)
	files := make([]string, 0)
	initialDepth := strings.Count(dirPath, "/") + 1

	Debug("Getting Name of Files and Dirs in: " + dirPath)
	err := filepath.WalkDir(dirPath, func(pathTemp string, d fs.DirEntry, err error) error {
		// if initial + 1 then smallest depth
		tempDepthCount := strings.Count(pathTemp, "/")
		if tempDepthCount >= initialDepth && tempDepthCount <= initialDepth+depth {
			if d.IsDir() {
				dirs = append(dirs, d.Name())
			} else {
				files = append(files, d.Name())
			}
		}
		return nil
	})
	return dirs, files, err

}
