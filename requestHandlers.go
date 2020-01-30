package main

import (
	"github.com/qiangxue/fasthttp-routing"
	"strings"
	"github.com/valyala/fasthttp"
	"github.com/sou-chon/proxy-go/s3Client"
	"github.com/sou-chon/proxy-go/definitions"
	"fmt"
)

func HandleGetResourceRequestEnv(s3ClientPool *s3Client.S3ClientPool, whitelistTokens []string) ( func(*routing.Context) error ) {
	/* This handler checks permission + formats. The s3client function checks the existence of objects */
	return func(ctx *routing.Context) error {
		store := ctx.Param("store")
		project := ctx.Param("project")
		resource := ctx.Param("resource")

		if (strings.HasSuffix(resource, ".ovemeta")) {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return fmt.Errorf("%v: Cannot get resource \"%v\"", definitions.ACCESS_DENIED_ERROR, resource)
		}

		/* check credentials */
		queryParams := ctx.QueryArgs()
		token := string(queryParams.Peek("token"))
		for k, v := range whitelistTokens {
			if (v == token) {
				break
			}
			if (k == len(whitelistTokens) - 1) {
				return fmt.Errorf("%v: Cannot get resource \"%v\". Invalid token. Please provide token in url as ?token={{token}}.", definitions.ACCESS_DENIED_ERROR, resource)
			}
		}


		/* getting the file */
		object, err, statusCode := s3ClientPool.DownloadFile(store, project, resource)
		if err != nil {
			ctx.SetStatusCode(statusCode)
			return fmt.Errorf("%v", err)
		}

		ctx.SetContentType("foo/bar")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBodyStream(object.Data, int(object.Length))

		return nil
	}
}