package service

import (
	"net/http"

	"github.com/izaakdale/lib/router"
)

func Router() http.Handler {
	return router.New(
		router.WithRoute(http.MethodPost, "/", createOrder),
		router.WithRoute(http.MethodGet, "/order/:id", getOrder),
	)
}
