package internal

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	ErrAttemptsExceeded = errors.New("timeout while creating table")
	ErrTableExists      = errors.New("credential store table already exists")
)

func (driver *Driver) IsTableExists(table string) (bool, error) {
	params := &dynamodb.ListTablesInput{}
	isExist := false

	err := driver.Ddb.ListTablesPages(params, func(page *dynamodb.ListTablesOutput, lastPage bool) bool {
		for _, tableName := range page.TableNames {
			if *tableName == table {
				isExist = true
				return false
			}
		}

		return true
	})
	if err != nil {
		return false, fmt.Errorf("can't check if %q exists: %w", table, err)
	}

	return isExist, nil
}

func (driver *Driver) CreateTable(table string) error {
	params := &dynamodb.CreateTableInput{
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
	}

	if _, err := driver.Ddb.CreateTable(params); err != nil {
		return fmt.Errorf("can't create %q: %w", table, err)
	}
	return nil
}

func (driver *Driver) WaitUntilTableExists(table string) error {
	delay := 20 * time.Second
	maxAttempts := 25
	isCreated := false

	params := &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}

	for i := 0; i < maxAttempts; i++ {
		resp, err := driver.Ddb.DescribeTable(params)
		if err != nil {
			return fmt.Errorf("can't describe %q: %w", table, err)
		}

		if *resp.Table.TableStatus == "ACTIVE" {
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

func (driver *Driver) CreateDdbTable(table string) error {
	if tableIsExist, err := driver.IsTableExists(table); err != nil {
		return err
	} else if tableIsExist {
		return fmt.Errorf("%w: %s", ErrTableExists, table)
	}

	if err := driver.CreateTable(table); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Creating table...")
	fmt.Fprintln(os.Stderr, "Waiting for table to be created...")

	if err := driver.WaitUntilTableExists(table); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Table has been created. Go read the README about how to create your KMS key")

	return nil
}
