VERSION:=unknown

all: gcredstash

gcredstash:
	CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w -X github.com/kgaughan/gcredstash/internal.Version=$(VERSION)" -tags netgo -o $@ ./cmd/$@

test:
	go test -cover -v ./...

clean:
	rm -rf gcredstash *.gz *.zip dist/

mock:
	go install go.uber.org/mock/mockgen@v0.4.0
	mockgen -package mockaws -destination internal/mockaws/dynamodbmock.go github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface DynamoDBAPI
	mockgen -package mockaws -destination internal/mockaws/kmsmock.go github.com/aws/aws-sdk-go/service/kms/kmsiface KMSAPI

.PHONY: all test clean mock
