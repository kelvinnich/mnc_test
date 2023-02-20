package usecase

import (
	"mnc/model"
	"mnc/repository"
)

// merchant usecase interface
type MerchantUsecase interface {
	AddMerchant(merchant *model.Merchant) error
}

type merchantUsecase struct {
	merchantRepo repository.MerchantRepository
}

// implementasi merchant usecase
func NewMerchantUsecase(merchantRepo repository.MerchantRepository) MerchantUsecase {
	return &merchantUsecase{merchantRepo}
}

func (uc *merchantUsecase) AddMerchant(merchant *model.Merchant) error {
	err := uc.merchantRepo.AddMerchant(merchant)
	if err != nil {
		return err
	}
	return nil
}
