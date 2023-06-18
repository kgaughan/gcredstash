package gcredstash

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

func KmsDecrypt(svc kmsiface.KMSAPI, blob []byte, context map[string]string) ([]byte, []byte, error) {
	params := &kms.DecryptInput{
		CiphertextBlob: blob,
	}

	if len(context) > 0 {
		ctx := map[string]*string{}

		for key, value := range context {
			ctx[key] = aws.String(value)
		}

		params.EncryptionContext = ctx
	}

	resp, err := svc.Decrypt(params)
	if err != nil {
		return nil, nil, fmt.Errorf("can't decrypt values with KMS: %w", err)
	}

	dataKey := resp.Plaintext[:32]
	hmacKey := resp.Plaintext[32:]

	return dataKey, hmacKey, nil
}

func KmsGenerateDataKey(svc kmsiface.KMSAPI, keyID string, context map[string]string) ([]byte, []byte, []byte, error) {
	params := &kms.GenerateDataKeyInput{
		KeyId:         aws.String(keyID),
		NumberOfBytes: aws.Int64(64),
	}

	if len(context) > 0 {
		ctx := map[string]*string{}

		for key, value := range context {
			ctx[key] = aws.String(value)
		}

		params.EncryptionContext = ctx
	}

	resp, err := svc.GenerateDataKey(params)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't generate data key: %w", err)
	}

	dataKey := resp.Plaintext[:32]
	hmacKey := resp.Plaintext[32:]
	wrappedKey := resp.CiphertextBlob

	return dataKey, hmacKey, wrappedKey, nil
}
