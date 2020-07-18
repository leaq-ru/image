package config

import (
	"github.com/kelseyhightower/envconfig"
)

const BucketName = "image"

type c struct {
	Grpc     grpc
	S3       s3
	LogLevel string `envconfig:"LOGLEVEL"`
}

type grpc struct {
	Port string `envconfig:"GRPC_PORT"`
}

type s3 struct {
	Endpoint        string `envconfig:"S3_ENDPOINT"`
	AccessKeyID     string `envconfig:"S3_ACCESSKEYID"`
	SecretAccessKey string `envconfig:"S3_SECRETACCESSKEY"`
	Secure          string `envconfig:"S3_SECURE"`
	Region          string `envconfig:"S3_REGION"`
}

var Env c

func init() {
	envconfig.MustProcess("", &Env)
}
