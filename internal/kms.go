package internal

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func KmsDecrypt(ctx context.Context, svc Kms, blob []byte, encCtx map[string]string) ([]byte, []byte, error) {
	params := &kms.DecryptInput{
		CiphertextBlob: blob,
	}

	if len(encCtx) > 0 {
		params.EncryptionContext = encCtx
	}

	resp, err := svc.Decrypt(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("can't decrypt values with KMS: %w", err)
	}

	dataKey := resp.Plaintext[:32]
	hmacKey := resp.Plaintext[32:]

	return dataKey, hmacKey, nil
}

func KmsGenerateDataKey(ctx context.Context, svc Kms, keyID string, encCtx map[string]string) ([]byte, []byte, []byte, error) {
	params := &kms.GenerateDataKeyInput{
		KeyId:         aws.String(keyID),
		NumberOfBytes: aws.Int32(64),
	}

	if len(encCtx) > 0 {
		params.EncryptionContext = encCtx
	}

	resp, err := svc.GenerateDataKey(ctx, params)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't generate data key: %w", err)
	}

	dataKey := resp.Plaintext[:32]
	hmacKey := resp.Plaintext[32:]
	wrappedKey := resp.CiphertextBlob

	return dataKey, hmacKey, wrappedKey, nil
}
