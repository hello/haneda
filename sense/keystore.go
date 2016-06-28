package sense

import (
	"encoding/hex"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type KeyStore interface {
	Get(senseId string) ([]byte, error)
}

type DynamoDBKeyStore struct {
	awsConfig *aws.Config
	srv       *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBKeyStore(tableName string, config *aws.Config) KeyStore {
	return DynamoDBKeyStore{
		tableName: tableName,
		srv:       dynamodb.New(session.New(), config),
	}
}

func (k DynamoDBKeyStore) Get(senseId string) ([]byte, error) {
	empty := make([]byte, 0)

	attrs := make(map[string]*dynamodb.AttributeValue)
	attrs["device_id"] = &dynamodb.AttributeValue{
		S: aws.String(senseId),
	}

	item := &dynamodb.GetItemInput{
		Key:       attrs,
		TableName: aws.String(k.tableName),
		AttributesToGet: []*string{
			aws.String("aes_key"),
		},
	}

	out, err := k.srv.GetItem(item)
	if err != nil {
		return empty, err
	}

	key, found := out.Item["aes_key"]
	if !found {
		return empty, errors.New("missing aes_key attribute")
	}
	keyBytes, err := hex.DecodeString(*key.S)
	return keyBytes, err
}
