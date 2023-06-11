VERSION:=unknown
SRC=$(wildcard *.go src/*/*.go src/*/*/*.go)

all: gcredstash

gcredstash: $(SRC)
	CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w -X main.Version=$(VERSION)" -tags netgo -o gcredstash

test:
	go test -v ./...

clean:
	rm -f gcredstash{,.exe} *.gz *.zip dist/

mock:
	go install github.com/golang/mock/mockgen@v1.6.0
	mockgen -package mockaws -destination src/mockaws/dynamodbmock.go github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface DynamoDBAPI
	mockgen -package mockaws -destination src/mockaws/kmsmock.go github.com/aws/aws-sdk-go/service/kms/kmsiface KMSAPI

.PHONY: all test clean mock
