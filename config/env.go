package config

import (
	"github.com/kelseyhightower/envconfig"
)

const ServiceName = "image"

type c struct {
	Grpc     grpc
	S3       s3
	STAN     stan
	NATS     nats
	LogLevel string `envconfig:"LOGLEVEL"`
}

type grpc struct {
	Port string `envconfig:"GRPC_PORT"`
}

type s3 struct {
	Alias           string `envconfig:"S3_ALIAS"`
	ImageBucketName string `envconfig:"S3_IMAGEBUCKETNAME"`
	Endpoint        string `envconfig:"S3_ENDPOINT"`
	AccessKeyID     string `envconfig:"S3_ACCESSKEYID"`
	SecretAccessKey string `envconfig:"S3_SECRETACCESSKEY"`
	Secure          string `envconfig:"S3_SECURE"`
	Region          string `envconfig:"S3_REGION"`
}

type stan struct {
	ClusterID                string `envconfig:"STAN_CLUSTERID"`
	SubjectCompanyNew        string `envconfig:"STAN_SUBJECTCOMPANYNEW"`
	SubjectImageUploadResult string `envconfig:"STAN_SUBJECTIMAGEUPLOADRESULT"`
}

type nats struct {
	URL string `envconfig:"NATS_URL"`
}

var Env c

func init() {
	envconfig.MustProcess("", &Env)
}
