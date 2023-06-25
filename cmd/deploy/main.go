package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

const (
	configFilename = "deploy_config.yaml"
)

type lambdaConfig struct {
	Name string `yaml:"name"`
	File string `yaml:"file"`
}

var deployConfig struct {
	AWSRegion          string         `yaml:"aws_region"`
	AWSAccessKeyId     string         `yaml:"aws_access_key_id"`
	AWSSecretAccessKey string         `yaml:"aws_secret_access_key"`
	Lambda             []lambdaConfig `yaml:"lambda"`
}

func main() {
	fmt.Println("open config file", configFilename)
	file, err := os.Open(configFilename)
	if os.IsNotExist(err) {
		fmt.Println("not exists init config file", configFilename)
		file, err = os.Create(configFilename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		deployConfig.Lambda = append(deployConfig.Lambda, lambdaConfig{
			Name: "test",
			File: "./file.zip",
		})

		yaml.NewEncoder(file).Encode(deployConfig)
		//enc := json.NewEncoder(file)
		//enc.SetIndent("", "\t")
		//err = enc.Encode(deployConfig)
		if err != nil {
			panic(err)
		}

		fmt.Println("init done")
		return
	}

	defer file.Close()

	fmt.Println("read config")
	err = yaml.NewDecoder(file).Decode(&deployConfig)
	if err != nil {
		panic(err)
	}

	cfg := aws.Config{
		Region: deployConfig.AWSRegion,
		Credentials: credentials.NewStaticCredentialsProvider(
			deployConfig.AWSAccessKeyId,     // key id
			deployConfig.AWSSecretAccessKey, // secret
			"",
		),
	}

	cli := lambda.NewFromConfig(cfg)
	ctx := context.Background()

	fmt.Println("deploying...")
	for i := range deployConfig.Lambda {
		err = deploy(ctx, cli, deployConfig.Lambda[i])
		if err != nil {
			panic(err)
		}
	}
}

func deploy(ctx context.Context, cli *lambda.Client, lc lambdaConfig) error {
	fmt.Println("open file", lc.File)
	file, err := os.Open(lc.File)
	if err != nil {
		fmt.Println("failed to open file", lc.File)
		return err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fmt.Println("deploy lambda", lc.Name)
	_, err = cli.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
		FunctionName: &lc.Name,
		ZipFile:      data,
	})
	if err != nil {
		fmt.Println("failed to deploy lambda", lc.Name)
	}
	return err
}
