package command

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kgaughan/gcredstash/internal"
	"github.com/kgaughan/gcredstash/internal/mockaws"
	"github.com/kgaughan/gcredstash/internal/testutils"
	"go.uber.org/mock/gomock"
)

func TestSetupCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDBAPI(ctrl)
	mkms := mockaws.NewMockKMSAPI(ctrl)
	table := "credential-store"

	mddb.EXPECT().ListTablesPages(
		&dynamodb.ListTablesInput{},
		gomock.Any(),
	).Return(nil)

	mddb.EXPECT().CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(table),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       aws.String("RANGE"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: aws.String("S"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}).Return(nil, nil)

	mddb.EXPECT().DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}).Return(&dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			TableStatus: aws.String("ACTIVE"),
		},
	}, nil)

	driver := &internal.Driver{Ddb: mddb, Kms: mkms}
	cmd, out := testutils.NewDummyCommand()

	args := []string{}
	if err := setupImpl(cmd, args, driver, out); err != nil {
		t.Errorf("\nexpected: %v\ngot: %q\n", nil, err)
	}
}
