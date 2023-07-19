package config

import (
	"d.kin-app/internal/awsx/lambdax"
	"github.com/joho/godotenv"
)

func init() {
	if !lambdax.IsLambdaRuntime() {
		godotenv.Load(".env.local")
	}
}
