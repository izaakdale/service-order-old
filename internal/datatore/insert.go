package datatore

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/izaakdale/service-order/schema/order"
)

func Insert(o *order.Order) (string, error) {
	id := uuid.NewString()
	pk := genPK(orderPKPrefix, id)

	meta := OrderRecord{
		PK:   pk,
		SK:   metaSK,
		Meta: o.GetMeta(),
	}
	request := OrderRecord{
		PK:      pk,
		SK:      requestSK,
		Request: o.Items,
	}

	metaMap, err := attributevalue.MarshalMap(meta)
	if err != nil {
		return "", err
	}
	requestMap, err := attributevalue.MarshalMap(request)
	if err != nil {
		return "", err
	}

	_, err = client.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			client.tableName: {
				{PutRequest: &types.PutRequest{Item: metaMap}},
				{PutRequest: &types.PutRequest{Item: requestMap}},
			},
		},
	})
	if err != nil {
		log.Printf("%+v\n", err)
		return "", err
	}

	return id, nil
}
