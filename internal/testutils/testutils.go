package testutils

import (
	"log"
	"os"

	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func MapToItem(m map[string]string) map[string]ddbtypes.AttributeValue {
	item := map[string]ddbtypes.AttributeValue{}

	for key, value := range m {
		item[key] = &ddbtypes.AttributeValueMemberS{Value: value}
	}

	return item
}

func ItemToMap(item map[string]ddbtypes.AttributeValue) map[string]string {
	m := map[string]string{}

	for key, value := range item {
		switch v := value.(type) {
		case *ddbtypes.AttributeValueMemberS:
			m[key] = v.Value
		default:
			panic("WAT")
		}
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
