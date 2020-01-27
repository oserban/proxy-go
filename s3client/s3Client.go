package s3Client

import (
	"strings"
	"fmt"
	"log"
	"github.com/sou-chon/proxy-go/definitions"
	"github.com/minio/minio-go/v6"
)

type s3ClientConfig struct {
	endpoint string
	accessKeyID string
	secretAccessKey string
	useSSL bool
}

type S3Client struct {
	minioClient *minio.Client
}

func CreateConnectedClient(config s3ClientConfig) S3Client {
	minioClient, err := minio.New(config.endpoint, config.accessKeyID, config.secretAccessKey, config.useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient)
	return S3Client { minioClient }
}

func (client *S3Client)DownloadFile(store, project, resource string) (ResourceData, error) {
	var data ResourceData

	if (strings.HasSuffix(resource, ".ovemeta")) {
		return data, fmt.Errorf("%v: Cannot get resource \"%v\".", definitions.ACCESS_DENIED_ERROR, resource)
	}

	/* check if store exists */

	/* check if project (bucket) exists */
	if found, _ := client.minioClient.BucketExists(resource); !found {
		return data, fmt.Errorf("%v: Resource not found \"%v\".", definitions.RESOURCE_NOT_FOUND_ERROR, resource)
	}

	/* get the object */
	// minio.GetObject(bucketName, objectName string, opts GetObjectOptions) (*Object, error)
	// minio.Object represents object reader. It implements io.Reader, io.Seeker, io.ReaderAt and io.Closer interfaces.
	objStats, statsErr := client.minioClient.StatObject(project, resource, minio.StatObjectOptions{})
	if (statsErr != nil) {
		return data, fmt.Errorf("%v: Cannot get object metadata for \"%v\".", definitions.SERVER_ERROR, resource)
	}

	obj, err := client.minioClient.GetObject(project, resource, minio.GetObjectOptions{})
	if (err != nil) {
		return data, fmt.Errorf("%v: Cannot get object data for \"%v\".", definitions.SERVER_ERROR, resource)
	}

	data = ResourceData {
		name: resource,
		resource: resource,
		contentType: objStats.ContentType,
		etag: "",
		modifiedDate: objStats.LastModified.String(),
		length: objStats.Size,
		data: obj,
	}

	return data, nil
}


