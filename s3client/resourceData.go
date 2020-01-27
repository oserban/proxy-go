package s3Client

import (
	"io"
)

type ResourceData struct {
	Name string
	Resource string
	ContentType string
	Etag string
	ModifiedDate string
	Length int64
	Data io.Reader
}