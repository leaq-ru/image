package imageimpl

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	m "github.com/minio/minio-go/v7"
	"github.com/nnqq/scr-image/config"
	"github.com/nnqq/scr-image/logger"
	"github.com/nnqq/scr-image/minio"
	"github.com/nnqq/scr-proto/codegen/go/image"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func (s *server) PutBase64(ctx context.Context, req *image.PutBase64Request) (res *image.PutResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if req.GetBase64() == "" {
		err = errors.New("base64 required")
		return
	}

	imgBytes, err := base64.StdEncoding.DecodeString(req.GetBase64())
	if err != nil {
		logger.Log.Error().Err(err).Send()
		return
	}

	img, err := bimg.NewImage(imgBytes).SmartCrop(200, 200)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	buf := &bytes.Buffer{}
	_, err = buf.Write(img)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	object, err := minio.Client.PutObject(
		ctx,
		config.Env.S3.ImageBucketName,
		strings.Join([]string{uuid.New().String(), "png"}, "."),
		buf,
		int64(len(img)),
		m.PutObjectOptions{},
	)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	if object.Key == "" {
		err = errors.New("object.Key is empty")
		return
	}

	res = &image.PutResponse{
		S3Url: strings.Join([]string{config.Env.S3.Alias, object.Key}, "/"),
	}
	return
}
