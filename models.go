package filestores

import "io"

type ObjectInfo struct {
	Ext, ContentType string
	Reader           io.Reader
	Size             int64
}
