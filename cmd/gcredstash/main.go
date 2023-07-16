package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/kgaughan/gcredstash/internal"
	"github.com/kgaughan/gcredstash/internal/command"
)

var Version = "unknown"

func main() {
	awsSession, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	driver := &internal.Driver{
		Ddb: dynamodb.New(awsSession),
		Kms: kms.New(awsSession),
	}

	command.MakeRootCmd(driver, Version).Execute()
}
