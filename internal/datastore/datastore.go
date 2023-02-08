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

	itemsRecord struct {
		PK string `dynamodbav:"PK" json:"-"`
		SK string `dynamodbav:"SK" json:"-"`

		Items    map[productID]quantity `dynamodbav:"items"`
		Subtotal *int64                 `dynamodbav:"subtotal"`
		Tax      *int64                 `dynamodbav:"tax"`
	}
	deliveryRecord struct {
		PK string `dynamodbav:"PK" json:"-"`
		SK string `dynamodbav:"SK" json:"-"`

		Name        *string `dynamodbav:"name"`
		HouseNumber *string `dynamodbav:"house_number"`
		Street      *string `dynamodbav:"street"`
		Postcode    *string `dynamodbav:"postcode"`
		Phone       *string `dynamodbav:"phone"`
		Status      *string `dynamodbav:"status"`
	}
	metaRecord struct {
		PK string `dynamodbav:"PK" json:"-"`
		SK string `dynamodbav:"SK" json:"-"`

		CreatedAt     int64   `dynamodbav:"created_at"`
		IpAddress     *string `dynamodbav:"ip_address"`
		PaymentMethod *string `dynamodbav:"payment_method"`
	}
)
