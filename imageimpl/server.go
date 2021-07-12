package imageimpl

import (
	"context"
	"errors"
	st "github.com/nats-io/stan.go"
	"github.com/nnqq/scr-image/config"
	"github.com/nnqq/scr-image/logger"
	"github.com/nnqq/scr-image/stan"
	"github.com/nnqq/scr-proto/codegen/go/event"
	"github.com/nnqq/scr-proto/codegen/go/image"
	"google.golang.org/protobuf/encoding/protojson"
)

type server struct {
	image.UnimplementedImageServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) ConsumeCompanyNew(_m *st.Msg) {
	go func(m *st.Msg) {
		ack := func() {
			err := m.Ack()
			if err != nil {
				logger.Log.Error().Err(err).Send()
			}
		}

		const (
			willRetry = "will retry"
			notRetry  = "not retry"
		)

		if config.Env.LogLevel == "debug" {
			logger.Log.Debug().
				Str("subject", config.Env.STAN.SubjectCompanyNew).
				Str("data", string(m.Data)).
				Msg("got message")
		}

		msg := &event.CompanyNew{}
		err := protojson.UnmarshalOptions{DiscardUnknown: true, AllowPartial: true}.Unmarshal(m.Data, msg)
		if err != nil {
			logger.Log.Error().Err(err).Msg(notRetry)
			ack()
			return
		}

		if msg.GetCompanyId() == "" || msg.GetAvatarToUpload() == "" {
			ack()
			return
		}

		res, err := s.Put(context.Background(), &image.PutRequest{
			Url: msg.GetAvatarToUpload(),
		})
		if err != nil {
			if errors.Is(err, ErrS3Retryable) {
				logger.Log.Error().Err(err).Msg(willRetry)
				return
			}

			logger.Log.Error().Err(err).Msg(notRetry)
			ack()
			return
		}

		err = stan.ProduceImageUploadResult(&event.ImageUploadResult{
			CompanyId: msg.GetCompanyId(),
			AvatarUrl: res.GetS3Url(),
		})
		if err != nil {
			logger.Log.Error().Err(err).Msg(willRetry)
			return
		}

		ack()
		return
	}(_m)
}

func (s *server) ConsumeDeleteImage(_m *st.Msg) {
	go func(m *st.Msg) {
		ack := func() {
			err := m.Ack()
			if err != nil {
				logger.Log.Error().Err(err).Send()
			}
		}

		const (
			willRetry = "will retry"
			notRetry  = "not retry"
		)

		if config.Env.LogLevel == "debug" {
			logger.Log.Debug().
				Str("subject", config.Env.STAN.SubjectCompanyNew).
				Str("data", string(m.Data)).
				Msg("got message")
		}

		msg := &event.DeleteImage{}
		err := protojson.UnmarshalOptions{DiscardUnknown: true, AllowPartial: true}.Unmarshal(m.Data, msg)
		if err != nil {
			logger.Log.Error().Err(err).Msg(notRetry)
			ack()
			return
		}

		if msg.GetS3Url() == "" {
			ack()
			return
		}

		_, err = s.Remove(context.Background(), &image.RemoveRequest{
			S3Url: msg.GetS3Url(),
		})
		if err != nil {
			if errors.Is(err, ErrS3Retryable) {
				logger.Log.Error().Err(err).Msg(willRetry)
				return
			}

			logger.Log.Error().Err(err).Msg(notRetry)
			ack()
			return
		}

		ack()
		return
	}(_m)
}
