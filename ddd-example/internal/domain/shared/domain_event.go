package shared

import "time"

type DomainEvent interface {
	OccurredAt() time.Time
	EventName() string
}
