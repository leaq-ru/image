package main

import (
	"github.com/nnqq/scr-image/config"
	"github.com/nnqq/scr-image/imageimpl"
	"github.com/nnqq/scr-image/logger"
	"github.com/nnqq/scr-proto/codegen/go/image"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	srv := grpc.NewServer()

	done := make(chan struct{}, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signals
		srv.GracefulStop()
		done <- struct{}{}
	}()

	lis, err := net.Listen("tcp", strings.Join([]string{
		"0.0.0.0",
		config.Env.Grpc.Port,
	}, ":"))
	logger.Must(err)

	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	image.RegisterImageServer(srv, imageimpl.NewServer())
	err = srv.Serve(lis)
	logger.Must(err)

	<-done
}
