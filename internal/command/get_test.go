package command

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/golang/mock/gomock"
	"github.com/kgaughan/gcredstash/internal"
	"github.com/kgaughan/gcredstash/internal/mockaws"
	"github.com/kgaughan/gcredstash/internal/testutils"
)

func TestGetCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)

	name := "test.key"
	table := "credential-store"

	item := map[string]string{
		"contents": "eBtO1lgLxIe6Yw==",
		"hmac":     "b23a3efafd4795e50ca87afd7d764f263e9ae456499a8d40eece70a63ed5da27",
		"key":      "CiDY1vsR456LEdoL3+0p+PrTCleoqi/sutbDfJZNiUSpphLLAQEBAQB42Nb7EeOeixHaC9/tKfj60wpXqKov7LrWw3yWTYlEqaYAAACiMIGfBgkqhkiG9w0BBwaggZEwgY4CAQAwgYgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMy/Oc2pOJsR0y9nbhAgEQgFsHECqku7QZiRjLmmeGyhcsgWdWvi7Op3luJu4soi5sP0pqcsjTrBJqOXHLazgyBS9wb6deP8zpXa/41WT0ZpNY9at4gw7+XRtbz8f4Rlh8WnyFnK5RZ7i0mOlD",
		"name":     "test.key",
		"version":  "0000000000000000002",
	}

	mddb.EXPECT().Query(&dynamodb.QueryInput{
		TableName:                aws.String(table),
		Limit:                    aws.Int64(1),
		ConsistentRead:           aws.Bool(true),
		ScanIndexForward:         aws.Bool(false),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {S: aws.String(name)},
		},
	}).Return(&dynamodb.QueryOutput{
		Count: aws.Int64(1),
		Items: []map[string]*dynamodb.AttributeValue{testutils.MapToItem(item)},
	}, nil)

	mkms.EXPECT().Decrypt(&kms.DecryptInput{
		CiphertextBlob: internal.B64Decode(item["key"]),
	}).Return(&kms.DecryptOutput{
		Plaintext: []byte{188, 163, 172, 238, 203, 68, 210, 84, 58, 152, 145, 235, 42, 23, 204, 164, 62, 139, 115, 220, 63, 85, 98, 228, 48, 229, 82, 62, 72, 86, 255, 162, 53, 75, 177, 91, 204, 232, 206, 127, 200, 23, 43, 148, 246, 221, 240, 247, 94, 72, 147, 211, 60, 139, 50, 150, 18, 100, 28, 24, 240, 2, 199, 121},
	}, nil)

	driver := &internal.Driver{Ddb: mddb, Kms: mkms}
	cmd, out := testutils.NewDummyCommand()

	args := []string{name}
	if err := getImpl(cmd, args, driver); err != nil {
		t.Errorf("\nexpected: %v\ngot: %q\n", nil, err)
	}

	expected := "test.value\n"
	txt := out.String()
	if expected != txt {
		t.Errorf("\nexpected: %q\ngot: %q\n", expected, txt)
	}
}

func TestGetCommandWithWildcard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)

	name := "test.key"
	table := "credential-store"

	item := map[string]string{
		"contents": "eBtO1lgLxIe6Yw==",
		"hmac":     "b23a3efafd4795e50ca87afd7d764f263e9ae456499a8d40eece70a63ed5da27",
		"key":      "CiDY1vsR456LEdoL3+0p+PrTCleoqi/sutbDfJZNiUSpphLLAQEBAQB42Nb7EeOeixHaC9/tKfj60wpXqKov7LrWw3yWTYlEqaYAAACiMIGfBgkqhkiG9w0BBwaggZEwgY4CAQAwgYgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMy/Oc2pOJsR0y9nbhAgEQgFsHECqku7QZiRjLmmeGyhcsgWdWvi7Op3luJu4soi5sP0pqcsjTrBJqOXHLazgyBS9wb6deP8zpXa/41WT0ZpNY9at4gw7+XRtbz8f4Rlh8WnyFnK5RZ7i0mOlD",
		"name":     "test.key",
		"version":  "0000000000000000002",
	}

	mddb.EXPECT().Scan(&dynamodb.ScanInput{
		TableName:                aws.String(table),
		ProjectionExpression:     aws.String("#name,version"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
	}).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{testutils.MapToItem(item)},
	}, nil)

	mddb.EXPECT().Query(&dynamodb.QueryInput{
		TableName:                aws.String(table),
		Limit:                    aws.Int64(1),
		ConsistentRead:           aws.Bool(true),
		ScanIndexForward:         aws.Bool(false),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {S: aws.String(name)},
		},
	}).Return(&dynamodb.QueryOutput{
		Count: aws.Int64(1),
		Items: []map[string]*dynamodb.AttributeValue{testutils.MapToItem(item)},
	}, nil)

	mkms.EXPECT().Decrypt(&kms.DecryptInput{
		CiphertextBlob: internal.B64Decode(item["key"]),
	}).Return(&kms.DecryptOutput{
		Plaintext: []byte{188, 163, 172, 238, 203, 68, 210, 84, 58, 152, 145, 235, 42, 23, 204, 164, 62, 139, 115, 220, 63, 85, 98, 228, 48, 229, 82, 62, 72, 86, 255, 162, 53, 75, 177, 91, 204, 232, 206, 127, 200, 23, 43, 148, 246, 221, 240, 247, 94, 72, 147, 211, 60, 139, 50, 150, 18, 100, 28, 24, 240, 2, 199, 121},
	}, nil)

	driver := &internal.Driver{Ddb: mddb, Kms: mkms}
	cmd, out := testutils.NewDummyCommand()

	args := []string{"test.*"}
	if err := getImpl(cmd, args, driver); err != nil {
		t.Errorf("\nexpected: %v\ngot: %q\n", nil, err)
	}

	expected := `{
  "test.key": "test.value"
}
`
	txt := out.String()
	if expected != txt {
		t.Errorf("\nexpected: %q\ngot: %q\n", expected, txt)
	}
}

/*
func TestGetCommandWithTrailingNewline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)

	name := "test.key"
	table := "credential-store"

	item := map[string]string{
		"contents": "eBtO1lgLxIe6Yw==",
		"hmac":     "b23a3efafd4795e50ca87afd7d764f263e9ae456499a8d40eece70a63ed5da27",
		"key":      "CiDY1vsR456LEdoL3+0p+PrTCleoqi/sutbDfJZNiUSpphLLAQEBAQB42Nb7EeOeixHaC9/tKfj60wpXqKov7LrWw3yWTYlEqaYAAACiMIGfBgkqhkiG9w0BBwaggZEwgY4CAQAwgYgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMy/Oc2pOJsR0y9nbhAgEQgFsHECqku7QZiRjLmmeGyhcsgWdWvi7Op3luJu4soi5sP0pqcsjTrBJqOXHLazgyBS9wb6deP8zpXa/41WT0ZpNY9at4gw7+XRtbz8f4Rlh8WnyFnK5RZ7i0mOlD",
		"name":     "test.key",
		"version":  "0000000000000000002",
	}

	mddb.EXPECT().Query(&dynamodb.QueryInput{
		TableName:                aws.String(table),
		Limit:                    aws.Int64(1),
		ConsistentRead:           aws.Bool(true),
		ScanIndexForward:         aws.Bool(false),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {S: aws.String(name)},
		},
	}).Return(&dynamodb.QueryOutput{
		Count: aws.Int64(1),
		Items: []map[string]*dynamodb.AttributeValue{testutils.MapToItem(item)},
	}, nil)

	mkms.EXPECT().Decrypt(&kms.DecryptInput{
		CiphertextBlob: internal.B64Decode(item["key"]),
	}).Return(&kms.DecryptOutput{
		Plaintext: []byte{188, 163, 172, 238, 203, 68, 210, 84, 58, 152, 145, 235, 42, 23, 204, 164, 62, 139, 115, 220, 63, 85, 98, 228, 48, 229, 82, 62, 72, 86, 255, 162, 53, 75, 177, 91, 204, 232, 206, 127, 200, 23, 43, 148, 246, 221, 240, 247, 94, 72, 147, 211, 60, 139, 50, 150, 18, 100, 28, 24, 240, 2, 199, 121},
	}, nil)

	cmd := &GetCommand{
		Meta: Meta{
			Table:  table,
			KmsKey: "alias/credstash",
			Driver: &internal.Driver{Ddb: mddb, Kms: mkms},
		},
	}

	args := []string{name}
	t.Setenv("GCREDSTASH_GET_TRAILING_NEWLINE", "1")
	out, err := cmd.RunImpl(args)
	expected := "test.value"

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}

	if expected != out {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, out)
	}
}

func TestGetCommandWithN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)

	name := "test.key"
	table := "credential-store"

	item := map[string]string{
		"contents": "eBtO1lgLxIe6Yw==",
		"hmac":     "b23a3efafd4795e50ca87afd7d764f263e9ae456499a8d40eece70a63ed5da27",
		"key":      "CiDY1vsR456LEdoL3+0p+PrTCleoqi/sutbDfJZNiUSpphLLAQEBAQB42Nb7EeOeixHaC9/tKfj60wpXqKov7LrWw3yWTYlEqaYAAACiMIGfBgkqhkiG9w0BBwaggZEwgY4CAQAwgYgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMy/Oc2pOJsR0y9nbhAgEQgFsHECqku7QZiRjLmmeGyhcsgWdWvi7Op3luJu4soi5sP0pqcsjTrBJqOXHLazgyBS9wb6deP8zpXa/41WT0ZpNY9at4gw7+XRtbz8f4Rlh8WnyFnK5RZ7i0mOlD",
		"name":     "test.key",
		"version":  "0000000000000000002",
	}

	mddb.EXPECT().Query(&dynamodb.QueryInput{
		TableName:                aws.String(table),
		Limit:                    aws.Int64(1),
		ConsistentRead:           aws.Bool(true),
		ScanIndexForward:         aws.Bool(false),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {S: aws.String(name)},
		},
	}).Return(&dynamodb.QueryOutput{
		Count: aws.Int64(1),
		Items: []map[string]*dynamodb.AttributeValue{testutils.MapToItem(item)},
	}, nil)

	mkms.EXPECT().Decrypt(&kms.DecryptInput{
		CiphertextBlob: internal.B64Decode(item["key"]),
	}).Return(&kms.DecryptOutput{
		Plaintext: []byte{188, 163, 172, 238, 203, 68, 210, 84, 58, 152, 145, 235, 42, 23, 204, 164, 62, 139, 115, 220, 63, 85, 98, 228, 48, 229, 82, 62, 72, 86, 255, 162, 53, 75, 177, 91, 204, 232, 206, 127, 200, 23, 43, 148, 246, 221, 240, 247, 94, 72, 147, 211, 60, 139, 50, 150, 18, 100, 28, 24, 240, 2, 199, 121},
	}, nil)

	cmd := &GetCommand{
		Meta: Meta{
			Table:  table,
			KmsKey: "alias/credstash",
			Driver: &internal.Driver{Ddb: mddb, Kms: mkms},
		},
	}

	args := []string{"-n", name}
	out, err := cmd.RunImpl(args)
	expected := "test.value"

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}

	if expected != out {
		t.Errorf("\nexpected: %v\ngot: %v\n", expected, out)
	}
}
*/

func TestGetCommandWithoutItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)

	name := "test.key"
	table := "credential-store"

	mddb.EXPECT().Query(&dynamodb.QueryInput{
		TableName:                aws.String(table),
		Limit:                    aws.Int64(1),
		ConsistentRead:           aws.Bool(true),
		ScanIndexForward:         aws.Bool(false),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name")},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {S: aws.String(name)},
		},
	}).Return(&dynamodb.QueryOutput{
		Count: aws.Int64(0),
		Items: []map[string]*dynamodb.AttributeValue{},
	}, nil)

	driver := &internal.Driver{Ddb: mddb, Kms: mkms}
	cmd, _ := testutils.NewDummyCommand()

	args := []string{name}
	err := getImpl(cmd, args, driver)
	if err == nil {
		t.Errorf("expected error does not happen")
	}
	if !errors.Is(err, internal.ErrItemNotFound) {
		t.Errorf("Unexpected error: %v", err)
	}
}
