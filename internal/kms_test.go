package internal

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/kgaughan/gcredstash/internal/mockaws"
	"go.uber.org/mock/gomock"
)

func TestKmsDecrypt(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	blob := []byte("123")
	context := map[string]string{"foo": "bar"}
	expectedDataKey := []byte("12345678901234567890123456789012")
	expectedHmacKey := []byte("abc")

	mkms := mockaws.NewMockKms(ctrl)

	mkms.EXPECT().Decrypt(
		ctx,
		&kms.DecryptInput{
			CiphertextBlob:    blob,
			EncryptionContext: map[string]string{"foo": "bar"},
		},
	).Return(&kms.DecryptOutput{
		Plaintext: append(expectedDataKey, expectedHmacKey...),
	}, nil)

	dataKey, hmacKey, err := KmsDecrypt(ctx, mkms, blob, context)
	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}

	if !bytes.Equal(expectedDataKey, dataKey) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedDataKey, dataKey)
	}

	if !bytes.Equal(expectedHmacKey, hmacKey) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedHmacKey, hmacKey)
	}
}

func TestKmsGenerateDataKey(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	keyID := "alias/credstash"
	context := map[string]string{"foo": "bar"}

	expectedDataKey := []byte("12345678901234567890123456789012")
	expectedHmacKey := []byte("abc")
	expectedWrappedKey := []byte("blobData")

	mkms := mockaws.NewMockKms(ctrl)

	mkms.EXPECT().GenerateDataKey(
		ctx,
		&kms.GenerateDataKeyInput{
			KeyId:             aws.String(keyID),
			NumberOfBytes:     aws.Int32(64),
			EncryptionContext: map[string]string{"foo": "bar"},
		}).Return(&kms.GenerateDataKeyOutput{
		Plaintext:      append(expectedDataKey, expectedHmacKey...),
		CiphertextBlob: expectedWrappedKey,
	}, nil)

	dataKey, hmacKey, wrappedKey, err := KmsGenerateDataKey(ctx, mkms, keyID, context)
	if err != nil {
		t.Errorf("\nexpected: %v\ngot: %v\n", nil, err)
	}

	if !bytes.Equal(expectedDataKey, dataKey) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedDataKey, dataKey)
	}

	if !bytes.Equal(expectedHmacKey, hmacKey) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedHmacKey, hmacKey)
	}

	if !bytes.Equal(expectedWrappedKey, wrappedKey) {
		t.Errorf("\nexpected: %v\ngot: %v\n", expectedWrappedKey, wrappedKey)
	}
}
