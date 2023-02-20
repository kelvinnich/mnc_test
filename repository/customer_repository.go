package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mnc/model"
)

type CustomersRepository interface {
	FindCustomerByID(id string) (*model.Customer,error)
	AuthenticateCustomer(id, password string) bool
}

type Database struct {
	db *model.Database
}

func NewCustomerRepository(db *model.Database) CustomersRepository{
	return &Database{
		db: db,
	}
}


func (db *Database) FindCustomerByID(id string) (*model.Customer, error) {
	// Baca data customers dari file JSON
	customersData, err := ioutil.ReadFile("customers.json")
	if err != nil {
			return nil, err
	}

	// Parsing data customers dari file JSON ke dalam slice Customer
	var customers []*model.Customer
	err = json.Unmarshal(customersData, &customers)
	if err != nil {
			return nil, err
	}

	// Cari customer dengan ID yang cocok
	for _, customer := range customers {
			if customer.ID == id {
					return customer, nil
			}
	}

	// Jika tidak ditemukan, return nil dan error
	return nil, fmt.Errorf("customer with ID %s not found", id)
}



func (db *Database) AuthenticateCustomer(id, password string) bool {
	customer,_ := db.FindCustomerByID(id)
	return customer != nil && customer.Password == password
}
