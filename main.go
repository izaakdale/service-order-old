package main

import "github.com/izaakdale/service-order/internal/service"

var (
	name = "service-order"
)

func main() {
	service.New(name).Run()
}
