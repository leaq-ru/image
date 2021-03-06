package stan

import (
	"github.com/leaq-ru/image/config"
	"github.com/leaq-ru/proto/codegen/go/event"
	"google.golang.org/protobuf/encoding/protojson"
)

func ProduceImageUploadResult(msg *event.ImageUploadResult) error {
	b, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	return Conn.Publish(config.Env.STAN.SubjectImageUploadResult, b)
}
