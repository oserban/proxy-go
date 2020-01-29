package configType 

import (
	"github.com/sou-chon/proxy-go/s3Client"
)

type Config struct {
	Stores []s3Client.S3ClientConfig `json:"Stores"`
	AccessTokens []string `json:"AccessTokens"`
	NumberOfClientsPerHost int `json:"NumberOfClientsPerHost"`
}