package sense

import (
	"encoding/hex"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// KeyStore
type KeyStore interface {
	Get(senseId string) ([]byte, error)
}

var (
	AesKeyAttributeName   = "aes_key"
	DeviceIdAttributeName = "device_id"
)

// DynamoDBKeyStore implements the KeyStore interface
type DynamoDBKeyStore struct {
	awsConfig *aws.Config
	srv       *dynamodb.DynamoDB
	tableName string
}

// Creates a new instance of DynamoDBKeyStore
func NewDynamoDBKeyStore(tableName string, config *aws.Config) *DynamoDBKeyStore {
	return &DynamoDBKeyStore{
		tableName: tableName,
		srv:       dynamodb.New(session.New(), config),
	}
}

// Get attempts to retrieve the associated aes_key from DynamoDB Item at key senseId
func (k *DynamoDBKeyStore) Get(senseId string) ([]byte, error) {
	empty := make([]byte, 0)

	attrs := make(map[string]*dynamodb.AttributeValue)
	attrs[DeviceIdAttributeName] = &dynamodb.AttributeValue{
		S: aws.String(senseId),
	}

	item := &dynamodb.GetItemInput{
		Key:       attrs,
		TableName: aws.String(k.tableName),
		AttributesToGet: []*string{
			aws.String(AesKeyAttributeName),
		},
	}

	out, err := k.srv.GetItem(item)
	if err != nil {
		return empty, err
	}

	key, found := out.Item[AesKeyAttributeName]
	if !found {
		return empty, errors.New("missing aes_key attribute")
	}

	keyBytes, err := hex.DecodeString(*key.S)
	return keyBytes, err
}

type FakeKeyStore struct {
	SenseId string
	Key     []byte
}

func (k *FakeKeyStore) Get(senseId string) ([]byte, error) {
	return k.Key, nil
}
