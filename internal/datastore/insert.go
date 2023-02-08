package datastore

import (
	"context"
	"encoding/json"
	"log"
	"time"

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

	items := make(map[productID]quantity, len(o.Items))
	for _, v := range o.Items {
		items[productID(v.Id)] = quantity(v.Quantity)
	}
	iRec := itemsRecord{
		PK: pk,
		SK: itemsSK,

		Items:    items,
		Subtotal: &o.Subtotal,
		Tax:      &o.Tax,
	}
	recMap, err := attributevalue.MarshalMap(iRec)
	if err != nil {
		return "", err
	}

	mRec := metaRecord{
		PK: pk,
		SK: metaSK,

		CreatedAt:     time.Now().Unix(),
		IpAddress:     &o.Metadata.IpAddress,
		PaymentMethod: &o.Metadata.PaymentMethod,
	}
	metaMap, err := attributevalue.MarshalMap(mRec)
	if err != nil {
		return "", err
	}

	dRec := deliveryRecord{
		PK: pk,
		SK: deliverySK,

		Name:        &o.DeliveryAddress.Name,
		HouseNumber: &o.DeliveryAddress.HouseNumber,
		Street:      &o.DeliveryAddress.Street,
		Postcode:    &o.DeliveryAddress.Postcode,
		Phone:       &o.DeliveryAddress.Phone,
		Status:      &statusReceived,
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
