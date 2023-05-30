package filestores

import (
	"os"
)

type (
	FileStore interface {
		Save(input Storable) (string, error)
		Delete(filepath string) error
	}
	Storable interface {
		Nameable
		GetFile() *os.File
	}
	Nameable interface {
		Filename() string
	}
)
