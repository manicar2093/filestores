package filestores

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileSystem struct {
	systemPath         string
	hostname           *url.URL
	isSimpleSystemPath bool
}

func NewFileSystem(systemPath, hostname string) *FileSystem {
	hostnameAsUrl, err := url.Parse(hostname)
	if err != nil {
		panic(err)
	}
	if hostnameAsUrl.Scheme == "" {
		panic(errors.New("hostname has no schema"))
	}
	if err := os.MkdirAll(systemPath, 0755); err != nil {
		panic(err)
	}
	instance := &FileSystem{
		systemPath: systemPath,
		hostname:   hostnameAsUrl,
	}
	if !path.IsAbs(systemPath) {
		instance.systemPath = fmt.Sprintf("/%s", path.Clean(systemPath))
	}
	if len(strings.Split(instance.systemPath, "/")) == 1 {
		instance.isSimpleSystemPath = true
	}

	return instance
}

func (c *FileSystem) Save(input Storable) (string, error) {
	info := input.GetStoreInfo()
	filename, nestedDirs, err := c.getFilenameAndCreateNestedDirs(input, info)
	if err != nil {
		return "", err
	}
	fileDst, err := os.Create(path.Join(".", c.systemPath, nestedDirs, filename))
	if err != nil {
		return "", err
	}
	defer fileDst.Close()

	if _, err := io.Copy(fileDst, info.Reader); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", c.hostname, filepath.Join(c.systemPath, nestedDirs, filename)), nil
}

func (c *FileSystem) Delete(filepath string) error {
	originalPath, err := c.removeEnpointFrom(filepath)
	if err != nil {
		return err
	}
	return os.Remove(originalPath)
}

func (c *FileSystem) Get(objectPath string) (ObjectInfo, error) {
	originalPath, err := c.removeEnpointFrom(objectPath)
	if err != nil {
		return ObjectInfo{}, err
	}
	file, err := os.Open(originalPath)
	if err != nil {
		return ObjectInfo{}, err
	}

	return FileToStoreInfo(file)
}

func (c *FileSystem) removeEnpointFrom(file string) (string, error) {
	u, err := url.Parse(file)
	if err != nil {
		return "", err
	}
	var filePath string
	if c.isSimpleSystemPath {
		filePath = filepath.Join(strings.Split(u.Path, "/")[2:]...)
	} else {
		filePath = strings.Replace(u.Path, c.systemPath, "", 1)
	}
	return path.Join(".", c.systemPath, filePath), nil
}

func (c *FileSystem) getFilenameAndCreateNestedDirs(input Storable, info ObjectInfo) (string, string, error) {
	filenameSplited := strings.Split(input.Filename(), "/")
	filename := fmt.Sprintf("%s%s", filenameSplited[len(filenameSplited)-1], info.Ext)
	nestedDirs := strings.Join(filenameSplited[:len(filenameSplited)-1], "/")
	if err := c.createNestedDirs(nestedDirs); err != nil {
		return "", "", err
	}
	return filename, nestedDirs, nil
}

func (c *FileSystem) createNestedDirs(nestedDirs string) error {
	if err := os.MkdirAll(filepath.Join(".", c.systemPath, nestedDirs), 0755); err != nil {
		return err
	}
	return nil
}
