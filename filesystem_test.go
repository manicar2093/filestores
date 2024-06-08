package filestores_test

import (
	"os"

	"github.com/manicar2093/filestores"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileSystem", Ordered, func() {

	var (
		gopherFile       *os.File
		systemPath       string
		expectedFileUrl  string
		expectedFilePath string
		store            filestores.FileStore
	)

	BeforeAll(func() {
		gopherFile = Must(os.Open("./gopher.png"))
		systemPath = "./testing"
		expectedFileUrl = "/testing/gophers/uuid/gopher_saved.png"
		expectedFilePath = "./testing/gophers/uuid/gopher_saved.png"
		store = filestores.NewFileSystem(systemPath)
	})

	AfterAll(func() {
		os.RemoveAll(systemPath)
		gopherFile.Close()
	})

	Describe("Save", func() {
		It("stores the file into systemPath", func() {
			var (
				input = SaveableFile{
					File: gopherFile,
				}
			)

			got, err := store.Save(input)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(Equal(expectedFileUrl))
			Expect(expectedFilePath).To(BeAnExistingFile())
		})
	})

	Describe("Get", func() {
		It("retreives a file from file system", func() {
			var (
				expectedContentType = "image/png"
				expectedExtension   = ".png"
			)

			got, err := store.Get(expectedFileUrl)

			Expect(err).ToNot(HaveOccurred())
			Expect(got.ContentType).To(Equal(expectedContentType))
			Expect(got.Ext).To(Equal(expectedExtension))
			Expect(got.Size).ToNot(BeZero())
			Expect(got.Reader).ToNot(BeNil())
		})
	})

	Describe("Delete", func() {
		It("removes file from system", func() {
			err := store.Delete(expectedFileUrl)

			Expect(err).ToNot(HaveOccurred())
			Expect(expectedFilePath).ToNot(BeAnExistingFile())

		})
	})

	Describe("all crud process", func() {
		It("complete all without failing", func() {

			var (
				file       = Must(os.Open("./gopher.png"))
				systemPath = "./testing/nested_testing/nested_testing_2"
				// expectedFileUrl  = "/testing/nested_testing/nested_testing/gophers/uuid/gopher_saved.png"
				// expectedFilePath = "./testing/nested_testing/nested_testing/gophers/uuid/gopher_saved.png"
				internalStore = filestores.NewFileSystem(systemPath)
				input         = SaveableFile{
					File: file,
				}
			)

			gotUrl, err := internalStore.Save(input)

			Expect(err).ToNot(HaveOccurred())

			var (
				expectedContentType = "image/png"
				expectedExtension   = ".png"
			)

			gotFile, err := internalStore.Get(gotUrl)

			Expect(err).ToNot(HaveOccurred())
			Expect(gotFile.ContentType).To(Equal(expectedContentType))
			Expect(gotFile.Ext).To(Equal(expectedExtension))
			Expect(gotFile.Size).ToNot(BeZero())
			Expect(gotFile.Reader).ToNot(BeNil())

			err = internalStore.Delete(gotUrl)

			Expect(err).ToNot(HaveOccurred())

		})
	})
})
