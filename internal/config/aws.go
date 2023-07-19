package config

import (
	"context"
	"d.kin-app/internal/typex"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	AWS = typex.ByLazy(func() aws.Config {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			panic(err)
		}

		return cfg
	})
)
