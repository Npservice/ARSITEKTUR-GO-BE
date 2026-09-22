package order

import "time"

type OrderPlaced struct {
	OrderID string
	At      time.Time
}

func (e OrderPlaced) OccurredAt() time.Time { return e.At }
func (e OrderPlaced) EventName() string     { return "order.placed" }

type OrderCancelled struct {
	OrderID string
	Reason  string
	At      time.Time
}

func (e OrderCancelled) OccurredAt() time.Time { return e.At }
func (e OrderCancelled) EventName() string     { return "order.cancelled" }
