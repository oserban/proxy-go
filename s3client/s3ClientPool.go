package s3Client

import (
	"fmt"
	"github.com/sou-chon/proxy-go/definitions"
	"github.com/valyala/fasthttp"
)

type S3ClientPool struct {
	clientMap map[string]S3Client
}

func InitialiseS3ClientPool(listOfHosts []S3ClientConfig) S3ClientPool {
	clientMap := make(map[string]S3Client)

	for _, host := range listOfHosts {
		clientMap[host.Name] = CreateConnectedClient(host)
	}
	return S3ClientPool{ clientMap } 
}

func (pool *S3ClientPool)DownloadFile(store, project, resource string) (ResourceData, error, int) {
	var data ResourceData

	s3Client, ok := pool.clientMap[store]
	if !ok {
		return data, fmt.Errorf("%v: Cannot get store \"%v\".", definitions.RESOURCE_NOT_FOUND_ERROR, store), fasthttp.StatusNotFound
	}

	resourceData, err, statusCode := s3Client.DownloadFile(store, project, resource)
	if err != nil {
		return data, err, statusCode 
	}

	return  resourceData, nil, fasthttp.StatusOK
}