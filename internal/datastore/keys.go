package datastore

import "fmt"

var (
	orderPrefix = "ORDER"

	metaSK       = "META"
	metaType     = "meta"
	deliverySK   = "DELIVERY"
	deliveryType = "delivery"
	itemsSK      = "ITEMS"
	itemsType    = "items"

	statusReceived   = "RECEIVED"
	statusProcessing = "WAITING"
	statusDone       = "DONE"
)

func genKey(prefix, id string) string {
	return fmt.Sprintf("%s#%s", prefix, id)
}
