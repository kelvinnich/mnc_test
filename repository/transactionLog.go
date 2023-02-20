package repository

import (
	"encoding/json"
	"io/ioutil"
	"mnc/model"
	"time"
)

type LogTransactionRepository interface {
	LogTransaction(customerID, merchantID string, amount int, fileName string) error
}

type DB struct {
	db *model.Database
}

func NewLogTransaction(db *model.Database) LogTransactionRepository{
	return &DB{
		db: db,
	}
}

func (db *DB) LogTransaction(customerID, merchantID string, amount int, fileName string) error {
	// create new transaction
	transaction := &model.History{
		CustomerID: customerID,
		MerchantID: merchantID,
		Amount:     amount,
		Time:       time.Now(),
	}

	// add transaction to the history slice
	db.db.History = append(db.db.History, transaction)

	// encode history slice to JSON
	historyData, err := json.MarshalIndent(db.db.History, "", "  ")
	if err != nil {
		return err
	}

	// write history data to file
	if err := ioutil.WriteFile(fileName, historyData, 0644); err != nil {
		return err
	}

	return nil
}

