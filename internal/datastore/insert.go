package datastore

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/izaakdale/lib/publisher"
	"github.com/izaakdale/service-order/schema/event"
	"github.com/izaakdale/service-order/schema/order"
)

func Insert(o *order.Order) (string, error) {
	id := uuid.NewString()
	log.Printf("handling order: %s", id)
	pk := genKey(orderPrefix, id)

	iRec := orderRecord{
		PK: pk,
		SK: itemsSK,

		Items: o.Items,
		Type:  itemsType,
	}
	recMap, err := attributevalue.MarshalMap(iRec)
	if err != nil {
		return "", err
	}

	mRec := orderRecord{
		PK: pk,
		SK: metaSK,

		Meta: o.Metadata,
		Type: metaType,
	}
	metaMap, err := attributevalue.MarshalMap(mRec)
	if err != nil {
		return "", err
	}

	dRec := orderRecord{
		PK: pk,
		SK: deliverySK,

		Delivery: o.DeliveryAddress,
		Type:     deliveryType,
	}
	deliveryMap, err := attributevalue.MarshalMap(dRec)
	if err != nil {
		return "", err
	}

	_, err = client.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			client.tableName: {
				{PutRequest: &types.PutRequest{Item: recMap}},
				{PutRequest: &types.PutRequest{Item: metaMap}},
				{PutRequest: &types.PutRequest{Item: deliveryMap}},
			},
		},
	})
	if err != nil {
		log.Printf("%+v\n", err)
		return "", err
	}

	var e = event.OrderCreated{ID: id}
	eBytes, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	publisher.Publish(string(eBytes))

	return id, nil
}
