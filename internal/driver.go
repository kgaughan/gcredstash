package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

var (
	ErrItemNotFound   = errors.New("item couldn't be found")
	ErrNeedContext    = errors.New("could not decrypt HMAC key with KMS: the credential may require that an encryption context be provided to decrypt it")
	ErrCredNotMatched = errors.New("could not decrypt HMAC key with KMS: the encryption context provided may not match the one used when the credential was stored")
	ErrBadHMAC        = errors.New("computed HMAC does not match stored HMAC")
	ErrVersionExists  = errors.New("version already in the credential store - use the -v flag to specify a new version")
	ErrBadType        = errors.New("unexpected entry type")
)

type DynamoDB interface {
	CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	ListTables(ctx context.Context, params *dynamodb.ListTablesInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ListTablesOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type Kms interface {
	Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
	GenerateDataKey(ctx context.Context, params *kms.GenerateDataKeyInput, optFns ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error)
}

type Driver struct {
	Ddb DynamoDB
	Kms Kms
}

func NewDriver(ctx context.Context) (*Driver, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't load AWS SDK configuration: %w", err)
	}
	driver := &Driver{
		Ddb: dynamodb.NewFromConfig(cfg),
		Kms: kms.NewFromConfig(cfg),
	}
	return driver, nil
}

func (driver *Driver) GetMaterialWithoutVersion(ctx context.Context, name, table string) (map[string]ddbtypes.AttributeValue, error) {
	params := &dynamodb.QueryInput{
		TableName:                aws.String(table),
		Limit:                    aws.Int32(1),
		ConsistentRead:           aws.Bool(true),
		ScanIndexForward:         aws.Bool(false),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]string{"#name": "name"},
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":name": &ddbtypes.AttributeValueMemberS{Value: name},
		},
	}

	resp, err := driver.Ddb.Query(ctx, params)
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	if resp.Count == 0 {
		return nil, fmt.Errorf(`%w: {"name": %q}`, ErrItemNotFound, name)
	}

	return resp.Items[0], nil
}

func (driver *Driver) GetMaterialWithVersion(ctx context.Context, name, version, table string) (map[string]ddbtypes.AttributeValue, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]ddbtypes.AttributeValue{
			"name":    &ddbtypes.AttributeValueMemberS{Value: name},
			"version": &ddbtypes.AttributeValueMemberS{Value: version},
		},
	}

	resp, err := driver.Ddb.GetItem(ctx, params)
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	if resp.Item == nil {
		return nil, fmt.Errorf(`%w: {"name": %q}`, ErrItemNotFound, name)
	}

	return resp.Item, nil
}

func (driver *Driver) DecryptMaterial(ctx context.Context, name string, material map[string]ddbtypes.AttributeValue, context map[string]string) (string, error) {
	var data []byte
	switch v := material["key"].(type) {
	case *ddbtypes.AttributeValueMemberS:
		data = B64Decode(v.Value)
	default:
		return "", ErrBadType
	}
	dataKey, hmacKey, err := KmsDecrypt(ctx, driver.Kms, data, context)
	if err != nil {
		if strings.Contains(err.Error(), "InvalidCiphertextException") {
			if len(context) < 1 {
				return "", fmt.Errorf("%s: %w", name, ErrNeedContext)
			}
			return "", fmt.Errorf("%s: %w", name, ErrCredNotMatched)
		}
		return "", err
	}

	var hmac []byte
	switch v := material["hmac"].(type) {
	case *ddbtypes.AttributeValueMemberB:
		hmac = HexDecode(string(v.Value))
	case *ddbtypes.AttributeValueMemberS:
		hmac = HexDecode(v.Value)
	default:
		return "", ErrBadType
	}

	var contents []byte
	switch v := material["contents"].(type) {
	case *ddbtypes.AttributeValueMemberS:
		contents = B64Decode(v.Value)
	default:
		return "", ErrBadType
	}
	if !ValidateHMAC(contents, hmac, hmacKey) {
		return "", fmt.Errorf("%s: %w", name, ErrBadHMAC)
	}

	decrypted := Crypt(contents, dataKey)

	return string(decrypted), nil
}

func (driver *Driver) GetHighestVersion(ctx context.Context, name, table string) (int, error) {
	params := &dynamodb.QueryInput{
		TableName:                aws.String(table),
		Limit:                    aws.Int32(1),
		ConsistentRead:           aws.Bool(true),
		ScanIndexForward:         aws.Bool(false),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]string{"#name": "name"},
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":name": &ddbtypes.AttributeValueMemberS{Value: name},
		},
		ProjectionExpression: aws.String("version"),
	}

	resp, err := driver.Ddb.Query(ctx, params)
	if err != nil {
		return -1, fmt.Errorf("can't query version: %w", err)
	}

	if resp.Count == 0 {
		return 0, nil
	}

	var version int
	switch v := resp.Items[0]["version"].(type) {
	case *ddbtypes.AttributeValueMemberS:
		version = Atoi(v.Value)
	default:
		return 0, ErrBadType
	}

	return version, nil
}

func (driver *Driver) PutItem(ctx context.Context, name, version string, key, contents, hmac []byte, table string) error {
	b64key := B64Encode(key)
	b64contents := B64Encode(contents)
	hexHmac := HexEncode(hmac)

	params := &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]ddbtypes.AttributeValue{
			"name":     &ddbtypes.AttributeValueMemberS{Value: name},
			"version":  &ddbtypes.AttributeValueMemberS{Value: version},
			"key":      &ddbtypes.AttributeValueMemberS{Value: b64key},
			"contents": &ddbtypes.AttributeValueMemberS{Value: b64contents},
			"hmac":     &ddbtypes.AttributeValueMemberS{Value: hexHmac},
		},
		ConditionExpression:      aws.String("attribute_not_exists(#name)"),
		ExpressionAttributeNames: map[string]string{"#name": "name"},
	}

	_, err := driver.Ddb.PutItem(ctx, params)
	if err != nil {
		return fmt.Errorf("can't store secret: %w", err)
	}

	return nil
}

func (driver *Driver) GetDeleteTargetWithoutVersion(ctx context.Context, name, table string) (map[string]string, error) {
	items := map[string]string{}

	params := &dynamodb.QueryInput{
		TableName:                aws.String(table),
		ConsistentRead:           aws.Bool(true),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]string{"#name": "name"},
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":name": &ddbtypes.AttributeValueMemberS{Value: name},
		},
	}

	resp, err := driver.Ddb.Query(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("can't find deletion target: %w", err)
	}

	if resp.Count == 0 {
		return nil, fmt.Errorf(`%w: {"name": %q}`, ErrItemNotFound, name)
	}

	for _, i := range resp.Items {
		var name string
		switch v := i["name"].(type) {
		case *ddbtypes.AttributeValueMemberS:
			name = v.Value
		default:
			return nil, ErrBadType
		}

		var version string
		switch v := i["version"].(type) {
		case *ddbtypes.AttributeValueMemberS:
			version = v.Value
		default:
			return nil, ErrBadType
		}

		items[name] = version
	}

	return items, nil
}

func (driver *Driver) GetDeleteTargetWithVersion(ctx context.Context, name, version, table string) (map[string]string, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]ddbtypes.AttributeValue{
			"name":    &ddbtypes.AttributeValueMemberS{Value: name},
			"version": &ddbtypes.AttributeValueMemberS{Value: version},
		},
	}

	resp, err := driver.Ddb.GetItem(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("can't find deletion target: %w", err)
	}

	if resp.Item == nil {
		versionNum := Atoi(version)
		return nil, fmt.Errorf(`%w: {"name": %q, "version": %d}`, ErrItemNotFound, name, versionNum)
	}

	items := map[string]string{}

	{
		var name string
		switch v := resp.Item["name"].(type) {
		case *ddbtypes.AttributeValueMemberS:
			name = v.Value
		default:
			return nil, ErrBadType
		}

		var version string
		switch v := resp.Item["version"].(type) {
		case *ddbtypes.AttributeValueMemberS:
			version = v.Value
		default:
			return nil, ErrBadType
		}

		items[name] = version
	}

	return items, nil
}

func (driver *Driver) DeleteItem(ctx context.Context, name, version, table string) error {
	svc := driver.Ddb

	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key: map[string]ddbtypes.AttributeValue{
			"name":    &ddbtypes.AttributeValueMemberS{Value: name},
			"version": &ddbtypes.AttributeValueMemberS{Value: version},
		},
	}

	if _, err := svc.DeleteItem(ctx, params); err != nil {
		return fmt.Errorf("can't delete secret %q (%v): %w", name, version, err)
	}

	return nil
}

func (driver *Driver) DeleteSecrets(ctx context.Context, name, version, table string) error {
	var items map[string]string
	var err error

	if version == "" {
		items, err = driver.GetDeleteTargetWithoutVersion(ctx, name, table)
	} else {
		items, err = driver.GetDeleteTargetWithVersion(ctx, name, version, table)
	}

	if err != nil {
		return err
	}

	for name, version := range items {
		err := driver.DeleteItem(ctx, name, version, table)
		if err != nil {
			return err
		}

		versionNum := Atoi(version)
		fmt.Fprintf(os.Stderr, "Deleting %s -- version %d\n", name, versionNum)
	}

	return nil
}

func (driver *Driver) PutSecret(ctx context.Context, name, secret, version, kmsKey, table string, context map[string]string) error {
	dataKey, hmacKey, wrappedKey, err := KmsGenerateDataKey(ctx, driver.Kms, kmsKey, context)
	if err != nil {
		return fmt.Errorf("could not generate key using KMS key(%s): %w", kmsKey, err)
	}

	cipherText := Crypt([]byte(secret), dataKey)
	hmac := Digest(cipherText, hmacKey)

	err = driver.PutItem(ctx, name, version, wrappedKey, cipherText, hmac, table)
	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			latestVersion, err := driver.GetHighestVersion(ctx, name, table)
			if err != nil {
				//nolint:wrapcheck
				return err
			}

			return fmt.Errorf("%w (name: %q, version: %d)", ErrVersionExists, name, latestVersion)
		}
		return err
	}

	return nil
}

func (driver *Driver) GetSecret(ctx context.Context, name, version, table string, context map[string]string) (string, error) {
	var material map[string]ddbtypes.AttributeValue
	var err error

	if version == "" {
		material, err = driver.GetMaterialWithoutVersion(ctx, name, table)
	} else {
		material, err = driver.GetMaterialWithVersion(ctx, name, version, table)
	}

	if err != nil {
		return "", fmt.Errorf("can't fetch secret: %w", err)
	}

	value, err := driver.DecryptMaterial(ctx, name, material, context)
	if err != nil {
		return "", fmt.Errorf("can't decrypt secret: %w", err)
	}

	return value, nil
}

func (driver *Driver) ListSecrets(ctx context.Context, table string) (map[string]string, error) {
	svc := driver.Ddb

	params := &dynamodb.ScanInput{
		TableName:                aws.String(table),
		ProjectionExpression:     aws.String("#name,version"),
		ExpressionAttributeNames: map[string]string{"#name": "name"},
	}

	resp, err := svc.Scan(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("can't list secrets: %w", err)
	}

	items := map[string]string{}

	for _, i := range resp.Items {
		var name string
		switch v := i["name"].(type) {
		case *ddbtypes.AttributeValueMemberS:
			name = v.Value
		default:
			return nil, ErrBadType
		}

		var version string
		switch v := i["version"].(type) {
		case *ddbtypes.AttributeValueMemberS:
			version = v.Value
		default:
			return nil, ErrBadType
		}

		items[name] = version
	}

	return items, nil
}
