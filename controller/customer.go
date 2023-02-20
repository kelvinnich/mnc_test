package controller

import (
	"encoding/json"
	"fmt"
	"mnc/model"
	"mnc/repository"
	"mnc/usecase"
	"net/http"
)

type LoginController interface {
	LoginHandler(w http.ResponseWriter, r *http.Request, db *model.Database)
	LogoutHandler(w http.ResponseWriter, r *http.Request)
	RegisterHandler(w http.ResponseWriter, r *http.Request, db *model.Database)
}

type login struct {
	custRepo repository.CustomersRepository
	data usecase.DataUseCase
}

func NewLoginController(custrepo repository.CustomersRepository, data usecase.DataUseCase)LoginController{
	return &login{
		custRepo: custrepo,
		data: data,
	}
}

func (l *login) LoginHandler(w http.ResponseWriter, r *http.Request, db *model.Database) {
	// parse request body
	var req struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// authenticate customer
	if !l.custRepo.AuthenticateCustomer(req.ID,req.Password){
		http.Error(w, "invalid login", http.StatusUnauthorized)
		return
	}

	// set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: req.ID,
	})

	// return success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "login sukses")
}


func (l *login) RegisterHandler(w http.ResponseWriter, r *http.Request, db *model.Database) {
	// parse request body
	var req struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Balance  int    `json:"balance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if user already exists
	if _, err := l.custRepo.FindCustomerByID(req.ID); err == nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	// create new customer
	customer := &model.Customer{
		ID:       req.ID,
		Name:     req.Name,
		Password: req.Password,
		Balance:  req.Balance,
	}
	db.Customers = append(db.Customers, customer)

	// save data
	if err := l.data.SaveData(*db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: customer.ID,
	})

	// return success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "register successful")
}




func (l *login)LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// delete session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	})

	// return success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "logout sukses")
}