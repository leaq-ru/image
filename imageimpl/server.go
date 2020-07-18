package imageimpl

import "github.com/nnqq/scr-proto/codegen/go/image"

type server struct {
	image.UnimplementedImageServer
}

func NewServer() *server {
	return &server{}
}
