package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/izaakdale/lib/publisher"
	"github.com/izaakdale/lib/response"
	"github.com/izaakdale/service-order/internal/datastore"
	"github.com/izaakdale/service-order/schema/event"
	"github.com/izaakdale/service-order/schema/order"
)

func createOrder(w http.ResponseWriter, r *http.Request) {
	var incoming order.Order
	err := json.NewDecoder(r.Body).Decode(&incoming)
	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, fmt.Sprintf("error decoding: %+e", err))
		return
	}

	if incoming.Metadata == nil || incoming.DeliveryAddress == nil || incoming.Items == nil {
		response.WriteJson(w, http.StatusBadRequest, err)
		return
	}

	id, err := datastore.Insert(&incoming)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, err)
		return
	}

	var e = event.OrderCreatedPayload{ID: id}
	cEvent := cloudevents.NewEvent()
	cEvent.SetSource(name)
	cEvent.SetType(event.TypeOrderCreated)
	cEvent.SetData(cloudevents.ApplicationJSON, e)
	bytes, err := json.Marshal(cEvent)
	if err != nil {
		log.Printf("error creating event for %s\n", id)
	}
	publisher.Publish(string(bytes))

	response.WriteJson(w, http.StatusCreated, map[string]string{
		"order_id": id,
	})
}
