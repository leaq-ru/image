package imageimpl

import (
	"bytes"
	"context"
	"errors"
	userAgent "github.com/EDDYCJY/fake-useragent"
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"github.com/leaq-ru/image/config"
	"github.com/leaq-ru/image/logger"
	"github.com/leaq-ru/image/minio"
	"github.com/leaq-ru/proto/codegen/go/image"
	m "github.com/minio/minio-go/v7"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

var ErrS3Retryable = errors.New("got err, try later")

func (s *server) Put(ctx context.Context, req *image.PutRequest) (res *image.PutResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if req.GetUrl() == "" {
		err = errors.New("url: cannot be blank")
		return
	}

	log := logger.Log.With().Str("url", req.GetUrl()).Logger()

	client := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		ReadTimeout:              5 * time.Second,
		WriteTimeout:             5 * time.Second,
		MaxConnWaitTimeout:       5 * time.Second,
		MaxResponseBodySize:      5 * 1024 * 1024,
		ReadBufferSize:           5 * 1024 * 1024,
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

	img, err := bimg.NewImage(httpRes.Body()).SmartCrop(200, 200)
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
		err = ErrS3Retryable
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
