package s3Client

import (
	"strings"
	"fmt"
	"log"
	"github.com/sou-chon/proxy-go/definitions"
	"github.com/valyala/fasthttp"
	"github.com/minio/minio-go/v6"
)

type S3ClientConfig struct {
	Endpoint string
	AccessKeyID string
	SecretAccessKey string
	UseSSL bool
}

type S3Client struct {
	minioClient *minio.Client
}

func CreateConnectedClient(config S3ClientConfig) S3Client {
	minioClient, err := minio.New(config.Endpoint, config.AccessKeyID, config.SecretAccessKey, config.UseSSL)
	if err != nil {
		log.Fatalln(err)
		panic("Cannot connect to minio.")
	}

	buckets, err := minioClient.ListBuckets()
	if err != nil {
		fmt.Println(err)
		panic("Cannot list buckets.")
	}
	fmt.Println("Connected to object store.")
	fmt.Println("Buckets:")
	for _, bucket := range buckets {
		fmt.Println(bucket)
	}
	return S3Client { minioClient }
}

func (client *S3Client)DownloadFile(store, project, resource string) (ResourceData, error, int) {
	var data ResourceData

	if (strings.HasSuffix(resource, ".ovemeta")) {
		return data, fmt.Errorf("%v: Cannot get resource \"%v\".", definitions.ACCESS_DENIED_ERROR, resource), fasthttp.StatusUnauthorized 
	}

	/* check if store exists */

	/* check if project (bucket) exists */
	if found, _ := client.minioClient.BucketExists(project); !found {
		return data, fmt.Errorf("%v: Project not found \"%v\".", definitions.RESOURCE_NOT_FOUND_ERROR, project), fasthttp.StatusNotFound
	}

	/* get the object (resource) and check if it exists */
	objStats, statsErr := client.minioClient.StatObject(project, resource, minio.StatObjectOptions{})
	if (statsErr != nil) {
		return data, fmt.Errorf("%v: Cannot get object metadata for \"%v\". It might not exist. Minio error: \"%v\"", definitions.RESOURCE_NOT_FOUND_ERROR, resource, statsErr), fasthttp.StatusNotFound
	}

	obj, err := client.minioClient.GetObject(project, resource, minio.GetObjectOptions{})
	if (err != nil) {
		return data, fmt.Errorf("%v: Cannot get object data for \"%v\". Minio error: \"%v\"", definitions.SERVER_ERROR, resource, err), fasthttp.StatusInternalServerError
	}

	/* return data struct with embedded data stream as io.Reader */
	data = ResourceData {
		Name: resource,
		Resource: resource,
		ContentType: objStats.ContentType,
		Etag: "",
		ModifiedDate: objStats.LastModified.String(),
		Length: objStats.Size,
		Data: obj,
	}

	return data, nil, fasthttp.StatusOK
}


