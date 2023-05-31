package filestores_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var cloudLabels = []string{"cloud", "integration"}

type SaveableFile struct {
	File *os.File
}

func (c SaveableFile) Filename() string {
	return "gophers/uuid/gopher_saved"
}

func (c SaveableFile) GetFile() *os.File {
	return c.File
}

func TestFilestores(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filestores Suite")
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
