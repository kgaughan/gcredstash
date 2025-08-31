package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ErrAttemptsExceeded = errors.New("timeout while creating table")
	ErrTableExists      = errors.New("credential store table already exists")
)

func (driver *Driver) IsTableExists(ctx context.Context, table string) (bool, error) {
	params := &dynamodb.ListTablesInput{}
	isExist := false

	paginator := dynamodb.NewListTablesPaginator(driver.Ddb, params)
out:
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return false, fmt.Errorf("can't check if %q exists: %w", table, err)
		}
		for _, tableName := range page.TableNames {
			if tableName == table {
				isExist = true
				break out
			}
		}
	}
	return isExist, nil
}

func (driver *Driver) CreateTable(ctx context.Context, table string) error {
	params := &dynamodb.CreateTableInput{
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
	}

	if _, err := driver.Ddb.CreateTable(ctx, params); err != nil {
		return fmt.Errorf("can't create %q: %w", table, err)
	}
	return nil
}

func (driver *Driver) WaitUntilTableExists(ctx context.Context, table string) error {
	delay := 20 * time.Second
	maxAttempts := 25
	isCreated := false

	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}

	for i := 0; i < maxAttempts; i++ {
		resp, err := driver.Ddb.DescribeTable(ctx, params)
		if err != nil {
			return fmt.Errorf("can't describe %q: %w", table, err)
		}

		if resp.Table.TableStatus == ddbtypes.TableStatusActive {
			isCreated = true
			break
		}

		time.Sleep(delay)
	}

	if !isCreated {
		return ErrAttemptsExceeded
	}

	return nil
}

func (driver *Driver) CreateDdbTable(ctx context.Context, table string) error {
	if tableIsExist, err := driver.IsTableExists(ctx, table); err != nil {
		return err
	} else if tableIsExist {
		return fmt.Errorf("%w: %s", ErrTableExists, table)
	}

	if err := driver.CreateTable(ctx, table); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Creating table...")
	fmt.Fprintln(os.Stderr, "Waiting for table to be created...")

	if err := driver.WaitUntilTableExists(ctx, table); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Table has been created. Go read the README about how to create your KMS key")

	return nil
}
