package service

import (
	"log"
	"net/http"

	"github.com/izaakdale/lib/response"
	"github.com/izaakdale/service-order/internal/datastore"
	"github.com/julienschmidt/httprouter"
)

func getOrder(w http.ResponseWriter, r *http.Request) {
	log.Printf("hitting get")
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	_, err := datastore.Fetch(id)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, err)
	}
}
