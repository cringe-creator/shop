package services

import (
	"shop/internal/repositories"
)

type HistoryService struct {
	TransactionRepo *repositories.TransactionRepository
	PurchaseRepo    *repositories.PurchaseRepository
}

func (s *HistoryService) GetUserHistory(userID int) (map[string]interface{}, error) {
	transactions, err := s.TransactionRepo.GetTransactionHistory(userID)
	if err != nil {
		return nil, err
	}

	purchases, err := s.PurchaseRepo.GetPurchaseHistory(userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"transactions": transactions,
		"purchases":    purchases,
	}, nil
}
