package filestores

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileSystem struct {
	systemPath string
}

const serverPrefix = "srv:/"

func NewFileSystem(systemPath string) *FileSystem {
	if err := os.MkdirAll(systemPath, 0755); err != nil {
		panic(err)
	}
	return &FileSystem{
		systemPath: systemPath,
	}
}

func (c *FileSystem) Save(input Storable) (string, error) {
	filename, nestedDirs := filenameAndNestedDirs(input)
	if err := c.createNestedDirs(nestedDirs); err != nil {
		return "", err
	}
	fileDst, err := os.Create(fmt.Sprintf("%s/%s/%s", c.systemPath, nestedDirs, filename))
	if err != nil {
		return "", err
	}
	defer fileDst.Close()

	if _, err := io.Copy(fileDst, input.GetFile()); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s/%s",
		serverPrefix,
		filepath.Join(c.systemPath, nestedDirs, filename),
	), nil
}

func (c *FileSystem) Delete(filepath string) error {
	return os.Remove(c.removeEnpointFrom(filepath))
}

func (c *FileSystem) removeEnpointFrom(file string) string {
	path := filepath.Join(strings.Split(file, "/")[3:]...)
	return fmt.Sprintf("%s/%s", c.systemPath, path)
}

func (c *FileSystem) createNestedDirs(nestedDirs string) error {
	if err := os.MkdirAll(filepath.Join("./", c.systemPath, nestedDirs), 0755); err != nil {
		return err
	}
	return nil
}

func filenameAndNestedDirs(input Storable) (string, string) {
	filenameSplited := strings.Split(input.Filename(), "/")
	filename := filenameSplited[len(filenameSplited)-1]
	nestedDirs := strings.Join(filenameSplited[:len(filenameSplited)-1], "/")
	return filename, nestedDirs
}
