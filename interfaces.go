package filestores

import "io"

type (
	FileStore interface {
		Save(input Storable) (string, error)
		Delete(filepath string) error
	}

	StoreInfo struct {
		Ext, ContentType string
		Reader           io.Reader
		Size             int64
	}
	Storable interface {
		Filename() string
		GetStoreInfo() StoreInfo
	}
)
