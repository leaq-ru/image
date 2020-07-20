package main

import (
	"github.com/nnqq/scr-image/config"
	"github.com/nnqq/scr-image/imageimpl"
	"github.com/nnqq/scr-image/logger"
	graceful "github.com/nnqq/scr-lib-graceful"
	"github.com/nnqq/scr-proto/codegen/go/image"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"strings"
)

func main() {
	srv := grpc.NewServer()
	go graceful.HandleSignals(srv.GracefulStop)

	lis, err := net.Listen("tcp", strings.Join([]string{
		"0.0.0.0",
		config.Env.Grpc.Port,
	}, ":"))
	if err != nil {
		logger.Log.Error().Err(err).Send()
		return
	}

	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	image.RegisterImageServer(srv, imageimpl.NewServer())
	logger.Err(srv.Serve(lis))
}
