package model


type Database struct {
	Customers []*Customer
	Merchants []*Merchant
	History   []*History
}

// NewDatabase creates a new in-memory database
func NewDatabase() *Database {
	return &Database{
		Customers: []*Customer{},
		Merchants: []*Merchant{},
		History:   []*History{},
	}
}
