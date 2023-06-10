VERSION:=unknown
SRC:=$(wildcard *.go src/*/*.go src/*/*/*.go)
TEST_SRC:=$(wildcard src/gcredstash/*_test.go)
CMD_TEST_SRC:=$(wildcard src/gcredstash/command/*_test.go)

all: gcredstash

gcredstash: $(SRC)
	CGO_ENABLED=0 go build -v -trimpath -ldflags "-s -w -X main.Version=$(VERSION)" -tags netgo -o gcredstash

test: $(TEST_SRC) $(CMD_TEST_SRC)
	go test -v $(TEST_SRC)
	go test -v $(CMD_TEST_SRC)

clean:
	rm -f gcredstash{,.exe} *.gz *.zip
	rm -f pkg/*
	rm -f debian/gcredstash.debhelper.log
	rm -f debian/gcredstash.substvars
	rm -f debian/files
	rm -rf debian/gcredstash/

mock:
	go install github.com/golang/mock/mockgen@v1.6.0
	mockgen -package mockaws -destination src/mockaws/dynamodbmock.go github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface DynamoDBAPI
	mockgen -package mockaws -destination src/mockaws/kmsmock.go github.com/aws/aws-sdk-go/service/kms/kmsiface KMSAPI
