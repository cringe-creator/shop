package services

import (
	"shop/internal/repositories"
)

type InfoService struct {
	UserRepo        *repositories.UserRepository
	TransactionRepo *repositories.TransactionRepository
	PurchaseRepo    *repositories.PurchaseRepository
}

func (s *InfoService) GetUserInfo(userID int) (map[string]interface{}, error) {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil || user == nil {
		return nil, err
	}

	inventory, err := s.PurchaseRepo.GetUserInventory(userID)
	if err != nil {
		inventory = []map[string]interface{}{}
	}

	receivedTransactions, err := s.TransactionRepo.GetReceivedTransactions(userID)
	if err != nil {
		receivedTransactions = []map[string]interface{}{}
	}

	sentTransactions, err := s.TransactionRepo.GetSentTransactions(userID)
	if err != nil {
		sentTransactions = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"coins":     user.Coins,
		"inventory": inventory,
		"coinHistory": map[string]interface{}{
			"received": receivedTransactions,
			"sent":     sentTransactions,
		},
	}, nil
}
