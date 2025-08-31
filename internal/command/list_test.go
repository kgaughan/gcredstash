package command

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kgaughan/gcredstash/internal"
	"github.com/kgaughan/gcredstash/internal/mockaws"
	"github.com/kgaughan/gcredstash/internal/testutils"
	"go.uber.org/mock/gomock"
)

func TestListCommand(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDB(ctrl)
	mkms := mockaws.NewMockKms(ctrl)

	table := "credential-store"
	name := "test.key"
	version := "0000000000000000002"

	item := map[string]string{
		"contents": "eBtO1lgLxIe6Yw==",
		"hmac":     "b23a3efafd4795e50ca87afd7d764f263e9ae456499a8d40eece70a63ed5da27",
		"key":      "CiDY1vsR456LEdoL3+0p+PrTCleoqi/sutbDfJZNiUSpphLLAQEBAQB42Nb7EeOeixHaC9/tKfj60wpXqKov7LrWw3yWTYlEqaYAAACiMIGfBgkqhkiG9w0BBwaggZEwgY4CAQAwgYgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMy/Oc2pOJsR0y9nbhAgEQgFsHECqku7QZiRjLmmeGyhcsgWdWvi7Op3luJu4soi5sP0pqcsjTrBJqOXHLazgyBS9wb6deP8zpXa/41WT0ZpNY9at4gw7+XRtbz8f4Rlh8WnyFnK5RZ7i0mOlD",
		"name":     name,
		"version":  version,
	}

	mddb.EXPECT().Scan(
		ctx,
		&dynamodb.ScanInput{
			TableName:                aws.String(table),
			ProjectionExpression:     aws.String("#name,version"),
			ExpressionAttributeNames: map[string]string{"#name": "name"},
		},
	).Return(&dynamodb.ScanOutput{
		Items: []map[string]ddbtypes.AttributeValue{testutils.MapToItem(item)},
	}, nil)

	driver := &internal.Driver{Ddb: mddb, Kms: mkms}
	cmd, out := testutils.NewDummyCommand(ctx)

	args := []string{}
	if err := listImpl(cmd, args, driver, out); err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}

	expected := fmt.Sprintf("%s -- version: %d\n", name, internal.Atoi(version))
	txt := out.String()
	if expected != txt {
		t.Errorf("\nexpected: %q\ngot: %q\n", expected, txt)
	}
}
