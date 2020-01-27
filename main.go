package main

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func main() {
	router := routing.New()
	
	router.Get("/<store>/<project>/<resource>", HandleGetResourceRequest)
	
	panic(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}