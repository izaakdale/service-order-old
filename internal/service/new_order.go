package service

import (
	"encoding/json"
	"net/http"

	"github.com/izaakdale/lib/response"
	"github.com/izaakdale/service-order/internal/datatore"
	"github.com/izaakdale/service-order/schema/order"
)

func createOrder(w http.ResponseWriter, r *http.Request) {
	var incoming order.Order
	err := json.NewDecoder(r.Body).Decode(&incoming)
	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, err)
		return
	}

	id, err := datatore.Insert(&incoming)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, err)
		return
	}
	response.WriteJson(w, http.StatusCreated, map[string]string{
		"order_id": id,
	})
}
