package datastore

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/izaakdale/service-order/schema/order"
)

func Fetch(id string) (*order.Order, error) {
	log.Printf("fetching item %s\n", id)

	keyCond := expression.Key("PK").Equal(expression.Value(genKey(orderPrefix, id)))
	proj := expression.NamesList(
		expression.Name("items"),
		expression.Name("meta"),
		expression.Name("delivery"),
		expression.Name("type"),
	)
	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
		WithProjection(proj).
		Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 &client.tableName,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
	}

	out, err := client.Query(context.Background(), input)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, errors.New("empty dynamo response")
	}

	orderRec := order.Order{}

	for _, item := range out.Items {
		var rec orderRecord
		if err := attributevalue.UnmarshalMap(item, &rec); err != nil {
			return nil, err
		}

		switch rec.Type {
		case itemsType:
			orderItems := []*order.Item{}
			err := attributevalue.Unmarshal(item[itemsType], &orderItems)
			if err != nil {
				return nil, err
			}
			orderRec.Items = orderItems
		case deliveryType:
			del := order.Delivery{}
			err := attributevalue.Unmarshal(item[deliveryType], &del)
			if err != nil {
				return nil, err
			}
			orderRec.DeliveryAddress = &del
		case metaType:
			meta := order.MetaData{}
			err := attributevalue.Unmarshal(item[metaType], &meta)
			if err != nil {
				return nil, err
			}
			orderRec.Metadata = &meta
		default:
			log.Printf("hit default switch while fetching %s\n", id)
		}
	}

	return &orderRec, nil
}
