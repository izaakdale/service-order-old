package datatore

import "fmt"

var (
	orderPKPrefix = "ORDER#"

	metaSK    = "META"
	requestSK = "REQUEST"

	statusWaiting = "WAITING"
	statusDone    = "DONE"
)

func genPK(prefix, id string) string {
	return fmt.Sprintf("%s%s", prefix, id)
}
