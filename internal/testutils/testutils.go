package testutils

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func MapToItem(m map[string]string) map[string]*dynamodb.AttributeValue {
	item := map[string]*dynamodb.AttributeValue{}

	for key, value := range m {
		item[key] = &dynamodb.AttributeValue{S: aws.String(value)}
	}

	return item
}

func ItemToMap(item map[string]*dynamodb.AttributeValue) map[string]string {
	m := map[string]string{}

	for key, value := range item {
		m[key] = *value.S
	}

	return m
}

func TempFile(content string, f func(*os.File)) {
	tmpfile, err := os.CreateTemp("", "gcredstash")
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.WriteString(content); err != nil {
		log.Printf("Error writing to temporary file: %v", err)
		panic(err)
	}

	if err = tmpfile.Sync(); err != nil {
		log.Printf("Error syncing temporary file: %v", err)
		panic(err)
	}

	f(tmpfile)

	if err = tmpfile.Close(); err != nil {
		log.Printf("Error closing temporary file: %v", err)
		panic(err)
	}
}

func Setenv(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		panic(err)
	}
}
