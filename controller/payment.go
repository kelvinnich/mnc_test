package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"mnc/model"
	"mnc/repository"
	"mnc/usecase"
	"net/http"
)


type PaymentController interface {
	PaymentHandler(w http.ResponseWriter, r *http.Request)
}

type payment struct {
	db *model.Database
	customerRepo repository.CustomersRepository
	merchantRepo repository.MerchantRepository
	logTransaction repository.LogTransactionRepository
	data usecase.DataUseCase
}

func NewPaymentController(db *model.Database, crp repository.CustomersRepository, mrp repository.MerchantRepository, ltp repository.LogTransactionRepository, data usecase.DataUseCase)PaymentController {
	return &payment{
		db: db,
		customerRepo: crp,
		merchantRepo: mrp,
		logTransaction: ltp,
		data: data,
	}
}

func (p *payment) PaymentHandler(w http.ResponseWriter, r *http.Request) {
	// read session cookie
	session, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	// find customer
	customer, err := p.customerRepo.FindCustomerByID(session.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// parse request body
	var req struct {
		MerchantID string `json:"merchant_id"`
		Amount     int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// find merchant
	merchant, err := p.merchantRepo.FindMerchantByID(req.MerchantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// check customer balance
	if customer.Balance < req.Amount {
		http.Error(w, "insufficient balance", http.StatusBadRequest)
		return
	}

	// perform transaction
	if err := p.processTransaction(customer.ID, merchant.ID, req.Amount, *p.db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "payment success")
}

func (p *payment) processTransaction(id string, merchantId string, amount int, db model.Database) error {
	
	cust,_ := p.customerRepo.FindCustomerByID(id)
	merch,_ := p.merchantRepo.FindMerchantByID(merchantId)
	
	// check customer balance
	if cust.Balance < amount {
		return errors.New("insufficient balance")
	}

	// log transaction
	p.logTransaction.LogTransaction(cust.ID, merch.ID, amount, "history.json")

	// update customer balance
	cust.Balance -= amount

	// save data
	if err := p.data.SaveData(db); err != nil {
		return err
	}

	if err := p.data.UpdateCustomerBalance(cust.ID, cust.Balance); err != nil {
		return err
	}

	// save transaction history
	if err := p.data.SaveHistory(cust.ID, merch.ID, amount); err != nil {
		return err
	}

	return nil
}

