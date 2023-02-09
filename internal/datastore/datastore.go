package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var client *Client

type dynamodbIface interface {
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

type Client struct {
	dynamodbIface
	tableName string
}

func Init(db dynamodbIface, tableName string) {
	client = &Client{
		db,
		tableName,
	}
}

type (
	productID string
	quantity  int

	orderRecord struct {
		PK string `dynamodbav:"PK" json:"-"`
		SK string `dynamodbav:"SK" json:"-"`

		Items    any `dynamodbav:"items" json:"items"`
		Delivery any `dynamodbav:"delivery" json:"delivery"`
		Meta     any `dynamodbav:"meta" json:"meta"`

		Type string `dynamodbav:"type" json:"type"`
	}
)
