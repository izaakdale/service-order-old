package datastore

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/izaakdale/service-order/schema/order"
)

func Fetch(username, id string) (*order.Order, error) {
	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &client.tableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + username},
			"SK": &types.AttributeValueMemberS{Value: "ORDER#" + id},
		},
	})
	if err != nil {
		return nil, err
	}
	var rec orderRecord
	err = attributevalue.UnmarshalMap(out.Item, &rec)
	if err != nil {
		return nil, err
	}

	log.Printf("rec: %+v", rec)

	return nil, nil
}
