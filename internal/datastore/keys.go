package datastore

import "fmt"

var (
	orderPrefix = "ORDER"

	metaSK     = "META"
	deliverySK = "DELIVERY"
	itemsSK    = "ITEMS"

	statusReceived   = "RECEIVED"
	statusProcessing = "WAITING"
	statusDone       = "DONE"
)

func genKey(prefix, id string) string {
	return fmt.Sprintf("%s#%s", prefix, id)
}
