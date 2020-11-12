package imageimpl

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	m "github.com/minio/minio-go/v7"
	"github.com/nnqq/scr-image/config"
	"github.com/nnqq/scr-image/logger"
	"github.com/nnqq/scr-image/minio"
	"github.com/nnqq/scr-proto/codegen/go/image"
	"strings"
	"time"
)

func (s *server) Remove(ctx context.Context, req *image.RemoveRequest) (res *empty.Empty, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if req.GetS3Url() == "" {
		err = errors.New("s3Url: cannot be blank")
		return
	}

	err = minio.Client.RemoveObject(
		ctx,
		config.Env.S3.ImageBucketName,
		strings.TrimPrefix(req.GetS3Url(), config.Env.S3.Alias+"/"),
		m.RemoveObjectOptions{},
	)
	if err != nil {
		logger.Log.Error().Err(err).Send()
		return
	}

	res = &empty.Empty{}
	return
}
