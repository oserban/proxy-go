package s3Client

import (
	"io"
)

type ResourceData struct {
	name string
	resource string
	contentType string
	etag string
	modifiedDate string
	length int64
	data io.Reader
}