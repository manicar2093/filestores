package filestores_test

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/filestores"
)

var _ = Describe("AwsBucket", Ordered, Labels(cloudLabels), func() {
	var (
		gopherFile       *os.File
		bucketName       string
		cfg              aws.Config
		expectedFilePath string
		expectedFileUrl  string
		store            filestores.FileStore
	)

	BeforeAll(func() {
		gopherFile = Must(os.Open("./gopher.png"))
		bucketName = os.Getenv("BUCKET_NAME")
		cfg = Must(config.LoadDefaultConfig(context.Background()))
		expectedFilePath = "gophers/uuid/gopher_saved.png"
		expectedFileUrl = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, cfg.Region, expectedFilePath)
		store = filestores.NewAwsBucket(bucketName, cfg)
	})

	Describe("Save", func() {
		It("stores the file into s3", func() {
			var (
				input = SaveableFile{
					File: gopherFile,
				}
			)

			got, err := store.Save(input)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(Equal(expectedFileUrl))
		})
	})

	Describe("Get", func() {
		It("retreives a file from bucket", func() {
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
		})
	})
})
