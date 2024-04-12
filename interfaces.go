package filestores

type (
	FileStore interface {
		Save(input Storable) (string, error)
		Delete(filepath string) error
		Get(filepath string) (ObjectInfo, error)
	}

	Storable interface {
		Filename() string // returns the filename without trailing slash and extension like my/path/to/file
		GetStoreInfo() ObjectInfo
	}
)
