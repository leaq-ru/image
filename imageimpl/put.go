package imageimpl

import (
	"bytes"
	"context"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	m "github.com/minio/minio-go/v7"
	"github.com/nnqq/scr-image/config"
	"github.com/nnqq/scr-image/logger"
	"github.com/nnqq/scr-image/minio"
	"github.com/nnqq/scr-proto/codegen/go/image"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

func (s *server) Put(ctx context.Context, req *image.PutRequest) (res *image.PutResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		ReadTimeout:              10 * time.Second,
		WriteTimeout:             10 * time.Second,
		MaxConnWaitTimeout:       10 * time.Second,
		MaxResponseBodySize:      30 * 1024 * 1024,
	}

	httpReq := fasthttp.AcquireRequest()
	httpReq.SetRequestURI(req.GetUrl())
	httpRes := fasthttp.AcquireResponse()
	err = client.DoRedirects(httpReq, httpRes, 3)
	if err != nil {
		logger.Log.Error().Err(err).Send()
		return
	}

	i, err := imaging.Decode(bytes.NewReader(httpRes.Body()))
	if err != nil {
		logger.Log.Error().Err(err).Send()
		return
	}

	cropped := imaging.Fill(i, 200, 200, imaging.Center, imaging.Box)
	buf := &bytes.Buffer{}
	err = imaging.Encode(buf, cropped, imaging.JPEG)
	if err != nil {
		logger.Log.Error().Err(err).Send()
		return
	}

	info, err := minio.Client.PutObject(
		ctx,
		config.BucketName,
		strings.Join([]string{uuid.New().String(), "jpg"}, "."),
		buf,
		-1,
		m.PutObjectOptions{},
	)
	if err != nil {
		logger.Log.Error().Err(err).Send()
		return
	}

	res = &image.PutResponse{
		S3Url: info.Location,
	}
	return
}
