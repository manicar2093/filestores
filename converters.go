package filestores

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func FileToStoreInfo(file *os.File) (StoreInfo, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return StoreInfo{}, err
	}
	size := fileInfo.Size()
	fileBytes := make([]byte, size)
	if _, err := file.Read(fileBytes); err != nil {
		return StoreInfo{}, err
	}

	return StoreInfo{
		Ext:         filepath.Ext(file.Name()),
		ContentType: http.DetectContentType(fileBytes),
		Reader:      bytes.NewReader(fileBytes),
		Size:        size,
	}, nil
}

func FileHeaderToStoreInfo(file *multipart.FileHeader) (StoreInfo, error) {
	openedFile, err := file.Open()
	if err != nil {
		return StoreInfo{}, err
	}
	fileBytes := make([]byte, file.Size)
	if _, err := openedFile.Read(fileBytes); err != nil {
		return StoreInfo{}, err
	}

	return StoreInfo{
		Ext:         filepath.Ext(file.Filename),
		ContentType: file.Header.Get("Content-Type"),
		Reader:      bytes.NewReader(fileBytes),
		Size:        file.Size,
	}, nil
}
