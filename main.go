package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nschimek/nice-fixture-feeder/cmd"
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/spf13/viper"
)

const (
	defaultConfig = "./config/default.yaml"
)

func main() {
	if viper.GetBool("serverless") {
		lambda.Start(LambdaHandler)
	} else {
		cmd.Execute()
	}
}

func init() {
	core.SetupViper()
}

func LambdaHandler() (string, error) {
	if err := cmd.Execute(); err != nil {
		return "", err
	} else {
		return "success", nil
	}
}