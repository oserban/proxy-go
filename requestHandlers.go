package main

import (
	"github.com/qiangxue/fasthttp-routing"
	"strings"
	"github.com/valyala/fasthttp"
	"github.com/sou-chon/proxy-go/s3Client"
	"github.com/sou-chon/proxy-go/definitions"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/sou-chon/proxy-go/mongoClient"
	"fmt"
	"context"
	"time"
)

func HandleGetResourceRequestEnv(s3ClientPool *s3Client.S3ClientPool) ( func(*routing.Context) error ) {
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


		/* getting the file */
		object, err, statusCode := s3ClientPool.DownloadFile(store, project, resource)
		if (err != nil) {
			ctx.SetStatusCode(statusCode)
			return fmt.Errorf("%v", err)
		}

		ctx.SetContentType("foo/bar")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBodyStream(object.Data, int(object.Length))

		return nil
	}
}



func HandleMongoResourceRequestEnv(mongoClient *mongoClient.MongoClient) ( func(*routing.Context) error ) {
	/* This handler checks permission + formats. The s3client function checks the existence of objects */
	return func(ctx *routing.Context) error {
		collection := mongoClient.MongoClient.Database("ukbiobank").Collection("STUDY_COLLECTION")

		type study struct {
			Id string
			Name string
		}
		var result []*study
		filter := bson.M{}
		ctxmongo, _ := context.WithTimeout(context.Background(), 5*time.Second)

		projection := bson.D{
			{"_id", 0},
			{"id", 1 },
			{"name", 1},
		}

		cur, err := collection.Find(ctxmongo, filter, options.Find().SetProjection(projection))

		for cur.Next(context.Background()) {
			var entry study
			err = cur.Decode(&entry)
			if err != nil {
				panic("no!!!!!")
			}
			result = append(result, &entry)
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
		fmt.Fprintf(ctx, "Go to /<store>/<project>/<resource> for objects.")
		return nil
	}
}