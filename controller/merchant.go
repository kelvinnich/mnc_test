package controller

import (
	"encoding/json"
	"fmt"
	"mnc/model"
	"mnc/usecase"
	"net/http"
)

type MerchantController interface {
	CreateMerchant(w http.ResponseWriter, r *http.Request)
}

type merchantController struct {
	MerchantUsecase usecase.MerchantUsecase
}

func NewMerchantController(mcu usecase.MerchantUsecase)MerchantController{
	return &merchantController{
		MerchantUsecase: mcu,
	}
}

// create merchant handler
func (mc *merchantController) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	// parse request body
	var req struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create new merchant
	merchant := &model.Merchant{
		ID:      req.ID,
		Name:    req.Name,
	}

	// add merchant
	if err := mc.MerchantUsecase.AddMerchant(merchant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "success add merchant")
}
