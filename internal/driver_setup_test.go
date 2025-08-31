package internal

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kgaughan/gcredstash/internal/mockaws"
	"go.uber.org/mock/gomock"
)

func TestCreateTable(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDB(ctrl)
	mkms := mockaws.NewMockKms(ctrl)
	table := "credential-store"

	mddb.EXPECT().CreateTable(
		ctx,
		&dynamodb.CreateTableInput{
			TableName: aws.String(table),
			KeySchema: []ddbtypes.KeySchemaElement{
				{
					AttributeName: aws.String("name"),
					KeyType:       ddbtypes.KeyTypeHash,
				},
				{
					AttributeName: aws.String("version"),
					KeyType:       ddbtypes.KeyTypeRange,
				},
			},
			AttributeDefinitions: []ddbtypes.AttributeDefinition{
				{
					AttributeName: aws.String("name"),
					AttributeType: ddbtypes.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("version"),
					AttributeType: ddbtypes.ScalarAttributeTypeS,
				},
			},
			ProvisionedThroughput: &ddbtypes.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		},
	).Return(nil, nil)

	driver := &Driver{
		Ddb: mddb,
		Kms: mkms,
	}

	err := driver.CreateTable(ctx, table)
	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestWaitUntilTableExists(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDB(ctrl)
	mkms := mockaws.NewMockKms(ctrl)
	table := "credential-store"

	mddb.EXPECT().DescribeTable(
		ctx,
		&dynamodb.DescribeTableInput{
			TableName: aws.String(table),
		},
	).Return(&dynamodb.DescribeTableOutput{
		Table: &ddbtypes.TableDescription{
			TableStatus: ddbtypes.TableStatusActive,
		},
	}, nil)

	driver := &Driver{
		Ddb: mddb,
		Kms: mkms,
	}

	err := driver.WaitUntilTableExists(ctx, table)
	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestIsTableExists(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDB(ctrl)
	mkms := mockaws.NewMockKms(ctrl)
	table := "credential-store"

	driver := &Driver{
		Ddb: mddb,
		Kms: mkms,
	}

	mddb.EXPECT().ListTables(
		ctx,
		&dynamodb.ListTablesInput{},
		gomock.Any(),
	).Return(&dynamodb.ListTablesOutput{}, nil)

	isExist, err := driver.IsTableExists(ctx, table)

	if isExist {
		t.Errorf("\nexpected: %v\ngot: %v\n", false, isExist)
	}

	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}

func TestCreateDdbTable(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mddb := mockaws.NewMockDynamoDB(ctrl)
	mkms := mockaws.NewMockKms(ctrl)
	table := "credential-store"

	mddb.EXPECT().ListTables(
		ctx,
		&dynamodb.ListTablesInput{},
		gomock.Any(),
	).Return(&dynamodb.ListTablesOutput{}, nil)

	mddb.EXPECT().CreateTable(
		ctx,
		&dynamodb.CreateTableInput{
			TableName: aws.String(table),
			KeySchema: []ddbtypes.KeySchemaElement{
				{
					AttributeName: aws.String("name"),
					KeyType:       ddbtypes.KeyTypeHash,
				},
				{
					AttributeName: aws.String("version"),
					KeyType:       ddbtypes.KeyTypeRange,
				},
			},
			AttributeDefinitions: []ddbtypes.AttributeDefinition{
				{
					AttributeName: aws.String("name"),
					AttributeType: ddbtypes.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("version"),
					AttributeType: ddbtypes.ScalarAttributeTypeS,
				},
			},
			ProvisionedThroughput: &ddbtypes.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		},
	).Return(nil, nil)

	mddb.EXPECT().DescribeTable(
		ctx,
		&dynamodb.DescribeTableInput{
			TableName: aws.String(table),
		},
	).Return(&dynamodb.DescribeTableOutput{
		Table: &ddbtypes.TableDescription{
			TableStatus: ddbtypes.TableStatusActive,
		},
	}, nil)

	driver := &Driver{
		Ddb: mddb,
		Kms: mkms,
	}

	err := driver.CreateDdbTable(ctx, table)
	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}
}
