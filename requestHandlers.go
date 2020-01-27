package main

import (
	"github.com/qiangxue/fasthttp-routing"
	"strings"
	"github.com/valyala/fasthttp"
	"github.com/sou-chon/proxy-go/s3Client"
	"github.com/sou-chon/proxy-go/definitions"
	"fmt"
)

func HandleGetResourceRequestEnv()

func HandleGetResourceRequest(ctx *routing.Context) error {
	store := ctx.Param("store")
	project := ctx.Param("project")
	resource := ctx.Param("resource")

	// need to check for nil / "" ?

	if (strings.HasSuffix(resource, ".ovemeta")) {
		fmt.Fprintf(ctx, "%v: Cannot get resource \"%v\"", definitions.ACCESS_DENIED_ERROR, resource)
		return fmt.Errorf("%v: Cannot get resource \"%v\"", definitions.ACCESS_DENIED_ERROR, resource)
	}

	/* check credentials */

	/* getting the file */
	object, err := s3Client.DownloadFile(store, project, resource)
	if (err != nil) {
		fmt.Fprintf(ctx, "%v: Cannot get resource \"%v\"", definitions.SERVER_ERROR, resource)
		return fmt.Errorf("%v: Cannot get resource \"%v\"", definitions.SERVER_ERROR, resource)
	}

	ctx.SetContentType("foo/bar")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyStream(object.data, object.length)

	// then update status code
	// ctx.SetStatusCode(fasthttp.StatusNotFound)
}