package s3Client

import (
	"fmt"
	"github.com/sou-chon/proxy-go/definitions"
	"github.com/valyala/fasthttp"
	"math/rand"
)

type S3ClientPool struct {
	numClientsPerHost int
	clientMap map[string][]S3Client
}

func InitialiseS3ClientPool(numClientsPerHost int, listOfHosts []S3ClientConfig) S3ClientPool {
	clientMap := make(map[string][]S3Client)

	for _, host := range listOfHosts {
		var clients []S3Client
		for i := 0; i < numClientsPerHost; i++ {
			clients = append(clients, CreateConnectedClient(host))
		}
		clientMap[host.Name] = clients
	}
	return S3ClientPool{ numClientsPerHost, clientMap } 
}

func (pool *S3ClientPool)DownloadFile(store, project, resource string) (ResourceData, error, int) {
	var data ResourceData

	s3ClientsList, ok := pool.clientMap[store]
	if !ok {
		return data, fmt.Errorf("%v: Cannot get store \"%v\".", definitions.RESOURCE_NOT_FOUND_ERROR, store), fasthttp.StatusNotFound
	}

	pickedClientNumber := rand.Intn(pool.numClientsPerHost)
	fmt.Println(pickedClientNumber)

	resourceData, err, statusCode := s3ClientsList[pickedClientNumber].DownloadFile(store, project, resource)
	if err != nil {
		return data, err, statusCode 
	}

	return  resourceData, nil, fasthttp.StatusOK
}