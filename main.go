package main

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/sou-chon/proxy-go/s3Client"
	"github.com/sou-chon/proxy-go/configType"
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"
)

func main() {
	/* reading from config file */
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	/* parsing config file TO_DO: what if config file does not conform to type? */
	byteValue, _ := ioutil.ReadAll(configFile)
	var config configType.Config
	json.Unmarshal(byteValue, &config)

	/* create connections */
	clientpool := s3Client.InitialiseS3ClientPool(config.NumberOfClientsPerHost, config.Stores)

	router := routing.New()
	
	router.Get("/<store>/<project>/<resource:[^ ]+>", HandleGetResourceRequestEnv(&clientpool, config.AccessTokens))

	router.Get("/", func(ctx *routing.Context) error {
		fmt.Fprintf(ctx, "Go to /<store>/<project>/<resource> for objects.")
		return nil
	})
	
	panic(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}