package inventory

type Repository interface {
	FindByProductID(productID string) (*Inventory, error)
	Save(inventory *Inventory) error
}
