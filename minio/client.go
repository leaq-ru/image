package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nnqq/scr-image/config"
	"github.com/nnqq/scr-image/logger"
	"strconv"
	"time"
)

var Client *minio.Client

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	secure, err := strconv.ParseBool(config.Env.S3.Secure)
	logger.Must(err)

	cl, err := minio.New(config.Env.S3.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			config.Env.S3.AccessKeyID,
			config.Env.S3.SecretAccessKey,
			"",
		),
		Secure: secure,
	})
	logger.Must(err)

	// ping
	_, err = cl.ListBuckets(ctx)
	logger.Must(err)

	err = cl.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{
		Region: config.Env.S3.Region,
	})
	if err != nil {
		// ok, seems bucket exists
		logger.Log.Debug().Err(err).Send()
	}

	Client = cl
}
