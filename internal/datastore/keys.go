package datastore

import "fmt"

var (
	userPrefix  = "USER"
	orderPrefix = "ORDER"

	metaSK    = "META"
	requestSK = "REQUEST"

	statusWaiting = "WAITING"
	statusDone    = "DONE"
)

func genKey(prefix, id string) string {
	return fmt.Sprintf("%s#%s", prefix, id)
}
