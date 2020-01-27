package main

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/sou-chon/proxy-go/s3Client"
	"fmt"
)

func main() {
	s3config := s3Client.S3ClientConfig {
		Endpoint: "localhost:9000",
		AccessKeyID: "admin",
		SecretAccessKey: "password",
		UseSSL: false,
	}

	s3ClientInstance := s3Client.CreateConnectedClient(s3config)

	router := routing.New()
	
	router.Get("/<store>/<project>/<resource:[^ ]+>", HandleGetResourceRequestEnv(&s3ClientInstance))

	router.Get("/", func(ctx *routing.Context) error {
		fmt.Fprintf(ctx, "Go to /<store>/<project>/<resource> for objects.")
		return nil
	})
	
	panic(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}