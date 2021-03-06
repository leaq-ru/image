package main

import (
	"context"
	"github.com/h2non/bimg"
	"github.com/leaq-ru/image/config"
	"github.com/leaq-ru/image/imageimpl"
	"github.com/leaq-ru/image/logger"
	"github.com/leaq-ru/image/stan"
	graceful "github.com/leaq-ru/lib-graceful"
	"github.com/leaq-ru/proto/codegen/go/image"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"strings"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	bimg.VipsCacheSetMaxMem(0)
	bimg.VipsCacheSetMax(0)

	srv := grpc.NewServer()
	img := imageimpl.NewServer()

	lis, err := net.Listen("tcp", strings.Join([]string{
		"0.0.0.0",
		config.Env.Grpc.Port,
	}, ":"))
	logger.Must(err)

	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	image.RegisterImageServer(srv, img)

	companyNew, err := stan.NewConsumer(
		logger.Log,
		stan.Conn,
		config.Env.STAN.SubjectCompanyNew,
		config.ServiceName,
		3,
		img.ConsumeCompanyNew,
	)
	logger.Must(err)

	deleteImage, err := stan.NewConsumer(
		logger.Log,
		stan.Conn,
		config.Env.STAN.SubjectDeleteImage,
		config.ServiceName,
		3,
		img.ConsumeDeleteImage,
	)
	logger.Must(err)

	var eg errgroup.Group
	eg.Go(func() error {
		graceful.HandleSignals(srv.GracefulStop, cancel)
		return nil
	})
	eg.Go(func() error {
		return srv.Serve(lis)
	})
	eg.Go(func() error {
		companyNew.Serve(ctx)
		return nil
	})
	eg.Go(func() error {
		deleteImage.Serve(ctx)
		return nil
	})
	logger.Must(eg.Wait())
}
