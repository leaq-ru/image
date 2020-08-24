package imageimpl

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	m "github.com/minio/minio-go/v7"
	"github.com/nnqq/scr-image/config"
	"github.com/nnqq/scr-image/logger"
	"github.com/nnqq/scr-image/minio"
	"github.com/nnqq/scr-proto/codegen/go/image"
	"net/url"
	"strings"
	"time"
)

func (s *server) Remove(ctx context.Context, req *image.RemoveRequest) (res *empty.Empty, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsedURL, err := url.Parse(req.GetS3Url())
	if err != nil {
		logger.Log.Error().Err(err).Send()
		return
	}

	err = minio.Client.RemoveObject(
		ctx,
		config.Env.S3.ImageBucketName,
		strings.TrimPrefix(parsedURL.Path, "/"),
		m.RemoveObjectOptions{},
	)
	if err != nil {
		logger.Log.Error().Err(err).Send()
		return
	}

	res = &empty.Empty{}
	return
}
