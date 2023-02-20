package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mnc/model"
	"os"
)
type DataUseCase interface{
	LoadData(customersFile, merchantsFile string) (*model.Database, error)
	SaveData(db model.Database) error
	SaveHistory(customerID, merchantID string, amount int) error
	UpdateCustomerBalance(customerID string, newBalance int) error
}

type dataUseCase struct {
	db *model.Database
}

func NewDataUseCase(hstry *model.Database)DataUseCase{
	return &dataUseCase{db: hstry}
	
}

func (d *dataUseCase) LoadData(customersFile, merchantsFile string) (*model.Database, error) {
	db := model.NewDatabase()

	// load customers
	customersData, err := ioutil.ReadFile(customersFile)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(customersData, &db.Customers); err != nil {
		return nil, err
	}

	// load merchants
	merchantsData, err := ioutil.ReadFile(merchantsFile)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(merchantsData, &db.Merchants); err != nil {
		return nil, err
	}

	return db, nil
}

func (d *dataUseCase) SaveData(db model.Database) error {
	//Save customers
	customersData, err := json.MarshalIndent(db.Customers, "", " ")
	if err != nil {
		return fmt.Errorf("error marshaling customers data: %v", err)
	}

	// Backup customers data file
	if err := backupFile("customers.json"); err != nil {
		return fmt.Errorf("error creating backup for customers file: %v", err)
	}

	// Write customers data to file
	if err := writeFile("customers.json", customersData); err != nil {
		return fmt.Errorf("error saving customers data: %v", err)
	}

	// Save merchants
	merchantsData, err := json.MarshalIndent(db.Merchants, "", " ")
	if err != nil {
		return fmt.Errorf("error marshaling merchants data: %v", err)
	}

	// Backup merchants data file
	if err := backupFile("merchants.json"); err != nil {
		return fmt.Errorf("error creating backup for merchants file: %v", err)
	}

	// Write merchants data to file
	if err := writeFile("merchants.json", merchantsData); err != nil {
		return fmt.Errorf("error saving merchants data: %v", err)
	}

	return nil
}

func backupFile(filepath string) error {
	// Check if file exists
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			// File does not exist, no need to create backup
			return nil
		} else {
			return fmt.Errorf("error checking if file exists: %v", err)
		}
	}

	// Copy file to backup file
	backupFilePath := filepath + ".bak"
	if err := copyFile(filepath, backupFilePath); err != nil {
		return fmt.Errorf("error creating backup for file: %v", err)
	}

	return nil
}

func writeFile(filepath string, data []byte) error {
	// Write data to file
	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("error saving data to file: %v", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	// Open source file for reading
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	// Create destination file for writing
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer dstFile.Close()

	// Copy data from source to destination
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("error copying data: %v", err)
	}

	return nil
}


func (p *dataUseCase) SaveHistory(customerID, merchantID string, amount int) error {
	// create new history entry
	h := &model.History{
		CustomerID: customerID,
		MerchantID: merchantID,
		Amount:     amount,
	}
	// append to payment history slice
	p.db.History = append(p.db.History, h)

	// encode history slice to JSON
	jsonData, err := json.MarshalIndent(p.db, "", " ")
	if err != nil {
		return err
	}

	// write JSON data to file
	if err := ioutil.WriteFile("history.json", jsonData, 0644); err != nil {
		return err
	}

	return nil
}


func (d *dataUseCase) UpdateCustomerBalance(customerID string, newBalance int) error {
	// Read customers data from file
	customersData, err := ioutil.ReadFile("customers.json")
	if err != nil {
		return fmt.Errorf("failed to read customers data file: %v", err)
	}

	// Unmarshal customers data
	var customers []model.Customer
	err = json.Unmarshal(customersData, &customers)
	if err != nil {
		return fmt.Errorf("failed to unmarshal customers data: %v", err)
	}

	// Find customer by ID and update balance
	var found bool
	for i := range customers {
		if customers[i].ID == customerID {
			customers[i].Balance = newBalance
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("customer with ID %s not found", customerID)
	}

	// Marshal updated customers data
	updatedCustomersData, err := json.MarshalIndent(customers, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated customers data: %v", err)
	}

	// Write updated customers data to file
	err = ioutil.WriteFile("customers.json", updatedCustomersData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated customers data to file: %v", err)
	}

	return nil
}