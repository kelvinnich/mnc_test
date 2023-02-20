package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mnc/model"
)

type MerchantRepository interface {
	FindMerchantByID(id string) (*model.Merchant, error)
	AddMerchant(merchant *model.Merchant) error
}

type mDB struct {
	db *model.Database
}

func NewMerchantRepository(db *model.Database)MerchantRepository{
	return &mDB{
		db: db,
	}
}

func (db *mDB) FindMerchantByID(id string) (*model.Merchant, error) {
	// Read merchant data from JSON file
	merchantsData, err := ioutil.ReadFile("merchants.json")
	if err != nil {
		return nil, err
	}

	// Parse merchant data from JSON file to Merchant slice
	var merchants []*model.Merchant
	err = json.Unmarshal(merchantsData, &merchants)
	if err != nil {
		return nil, err
	}

	// Find merchant with matching ID
	for _, merchant := range merchants {
		if merchant.ID == id {
			return merchant, nil
		}
	}

	// If not found, return nil and error
	return nil, fmt.Errorf("merchant with ID %s not found", id)
}



func (db *mDB) AddMerchant(merchant *model.Merchant) error {
	// Baca data merchants dari file JSON
	merchantsData, err := ioutil.ReadFile("merchants.json")
	if err != nil {
		return err
	}

	// Parsing data merchants dari file JSON ke dalam slice Merchant
	var merchants []*model.Merchant
	err = json.Unmarshal(merchantsData, &merchants)
	if err != nil {
		return err
	}

	// Cek apakah merchant dengan ID yang sama sudah ada
	for _, m := range merchants {
		if m.ID == merchant.ID {
			return fmt.Errorf("merchant with ID %s already exists", merchant.ID)
		}
	}

	// Tambahkan merchant ke slice Merchant
	merchants = append(merchants, merchant)

	// Tulis kembali data merchants ke file JSON
	merchantsData, err = json.Marshal(merchants)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("merchants.json", merchantsData, 0644)
	if err != nil {
		return err
	}

	return nil
}

