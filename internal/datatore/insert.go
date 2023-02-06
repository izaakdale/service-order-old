package datatore

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/izaakdale/lib/publisher"
	"github.com/izaakdale/service-order/schema/event"
	"github.com/izaakdale/service-order/schema/order"
)

func Insert(o *order.Order) (string, error) {
	id := uuid.NewString()
	log.Printf("handling order: %s", id)

	rec := OrderRecord{
		PK: genKey(userPrefix, o.Username),
		SK: genKey(orderPrefix, id),

		Items:           o.Items,
		Status:          statusWaiting,
		CreatedAt:       time.Now().Unix(),
		DeliveryAddress: o.DeliveryAddress,
	}

	recMap, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return "", err
	}

	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: &client.tableName,
		Item:      recMap,
	})
	if err != nil {
		log.Printf("%+v\n", err)
		return "", err
	}

	var e = event.OrderCreated{OrderID: id}
	eBytes, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	publisher.Publish(string(eBytes))

	return id, nil
}
