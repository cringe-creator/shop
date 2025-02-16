package services

import (
	"errors"
	"shop/internal/repositories"
)

var MerchItems = map[string]int{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}

type PurchaseService struct {
	UserRepo     *repositories.UserRepository
	PurchaseRepo *repositories.PurchaseRepository
}

func (s *PurchaseService) BuyItem(userID int, itemName string) error {
	price, exists := MerchItems[itemName]
	if !exists {
		return errors.New("item not found")
	}

	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	if user.Coins < price {
		return errors.New("insufficient balance")
	}

	err = s.PurchaseRepo.CreatePurchase(userID, itemName, price)
	if err != nil {
		return err
	}

	return nil
}
