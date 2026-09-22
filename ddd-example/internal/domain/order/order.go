package order

import (
	"errors"
	"time"

	"ddd-example/internal/domain/shared"
)

type Status string

const (
	StatusDraft     Status = "DRAFT"
	StatusPlaced    Status = "PLACED"
	StatusCancelled Status = "CANCELLED"
)

// Aggregate Root: satu-satunya titik masuk untuk memodifikasi Order dan item-itemnya.
// Bertanggung jawab menjaga invariant bisnis (mis. tidak boleh place order tanpa item).
type Order struct {
	id       string
	customer string
	items    []OrderItem
	status   Status
	events   []shared.DomainEvent
}

func NewOrder(id, customerID string) (*Order, error) {
	if id == "" {
		return nil, errors.New("order id is required")
	}
	if customerID == "" {
		return nil, errors.New("customer id is required")
	}
	return &Order{id: id, customer: customerID, status: StatusDraft}, nil
}

func (o *Order) ID() string     { return o.id }
func (o *Order) Status() Status { return o.status }
func (o *Order) Items() []OrderItem {
	return append([]OrderItem{}, o.items...)
}

func (o *Order) AddItem(productID string, quantity int, unitPrice shared.Money) error {
	if o.status != StatusDraft {
		return errors.New("cannot add item to a non-draft order")
	}
	item, err := NewOrderItem(productID, quantity, unitPrice)
	if err != nil {
		return err
	}
	o.items = append(o.items, item)
	return nil
}

// Place menegakkan invariant: order harus punya minimal satu item sebelum dipesan.
func (o *Order) Place() error {
	if o.status != StatusDraft {
		return errors.New("order already placed or cancelled")
	}
	if len(o.items) == 0 {
		return errors.New("cannot place an order without items")
	}
	o.status = StatusPlaced
	o.record(OrderPlaced{OrderID: o.id, At: time.Now()})
	return nil
}

// Cancel menegakkan invariant: hanya order yang sudah PLACED yang boleh
// dibatalkan (draft cukup dibuang, tidak perlu "dibatalkan"), dan order
// yang sudah CANCELLED tidak boleh dibatalkan lagi.
func (o *Order) Cancel(reason string) error {
	if o.status == StatusDraft {
		return errors.New("cannot cancel a draft order")
	}
	if o.status == StatusCancelled {
		return errors.New("order already cancelled")
	}
	if reason == "" {
		return errors.New("cancellation reason is required")
	}
	o.status = StatusCancelled
	o.record(OrderCancelled{OrderID: o.id, Reason: reason, At: time.Now()})
	return nil
}

func (o *Order) Total() (shared.Money, error) {
	if len(o.items) == 0 {
		return shared.NewMoney(0, "IDR")
	}
	total := o.items[0].Subtotal()
	for _, item := range o.items[1:] {
		var err error
		total, err = total.Add(item.Subtotal())
		if err != nil {
			return shared.Money{}, err
		}
	}
	return total, nil
}

func (o *Order) record(event shared.DomainEvent) {
	o.events = append(o.events, event)
}

// PullEvents mengambil dan mengosongkan domain event yang belum dipublikasikan.
func (o *Order) PullEvents() []shared.DomainEvent {
	events := o.events
	o.events = nil
	return events
}
