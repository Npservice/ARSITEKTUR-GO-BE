package order

// Repository interface tinggal di domain layer (dependency inversion);
// implementasinya ada di infrastructure layer.
type Repository interface {
	NextID() string
	FindByID(id string) (*Order, error)
	Save(order *Order) error
}
