package imageimpl

import (
	"bytes"
	"context"
	"errors"
	userAgent "github.com/EDDYCJY/fake-useragent"
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
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	if req.GetUrl() == "" {
		err = errors.New("url: cannot be blank")
		return
	}

	log := logger.Log.With().Str("url", req.GetUrl()).Logger()

	client := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		ReadTimeout:              10 * time.Second,
		WriteTimeout:             10 * time.Second,
		MaxConnWaitTimeout:       10 * time.Second,
		MaxResponseBodySize:      10 * 1024 * 1024,
	}

	httpReq := fasthttp.AcquireRequest()
	httpReq.SetRequestURI(req.GetUrl())
	httpReq.Header.SetUserAgent(userAgent.Random())
	httpRes := fasthttp.AcquireResponse()
	err = client.DoRedirects(httpReq, httpRes, 3)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	i, err := imaging.Decode(bytes.NewReader(httpRes.Body()))
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	croppedImg := imaging.Fit(i, 200, 200, imaging.Box)
	buf := &bytes.Buffer{}
	err = imaging.Encode(buf, croppedImg, imaging.JPEG)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	object, err := minio.Client.PutObject(
		ctx,
		config.BucketName,
		strings.Join([]string{uuid.New().String(), "jpg"}, "."),
		buf,
		-1,
		m.PutObjectOptions{},
	)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	if object.Location == "" {
		err = errors.New("object.Location is empty")
		return
	}

	res = &image.PutResponse{
		S3Url: object.Location,
	}
	return
}
