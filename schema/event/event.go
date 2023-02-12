package event

type OrderCreatedPayload struct {
	ID string `json:"id,omitempty"`
}

var (
	TypeOrderCreated = "order.created"
	TypeOrderUpdated = "order.updated"
)
