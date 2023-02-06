package datatore

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

type OrderRecord struct {
	PK string `dynamodbav:"PK" json:"-"`
	SK string `dynamodbav:"SK" json:"-"`

	Items           any    `dynamodbav:"items"`
	Status          string `dynamodbav:"status"`
	CreatedAt       int64  `dynamodbav:"created_at"`
	DeliveryAddress any    `dynamodbav:"delivery_address"`
}

type ProfileRecord struct {
	PK string `dynamodbav:"PK" json:"-"`
	SK string `dynamodbav:"SK" json:"-"`

	Username  string `dynamodbav:"username"`
	FullName  string `dynamodbav:"full_name"`
	Email     string `dynamodbav:"email"`
	CreatedAt int64  `dynamodbav:"created_at"`
	Addresses any    `dynamodbav:"addresses"`
}

type Address struct {
	NameOrNumber string
	Street       string
	Postcode     string
}

type Item struct {
	ID       string
	Quantity int32
}
