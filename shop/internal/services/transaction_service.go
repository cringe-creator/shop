package services

import (
	"errors"
	"shop/internal/repositories"
)

type TransactionService struct {
	TransactionRepo *repositories.TransactionRepository
	UserRepo        *repositories.UserRepository
}

func (s *TransactionService) SendCoins(senderID, receiverID, amount int) error {
	if senderID == receiverID {
		return errors.New("cannot send coins to yourself")
	}

	receiver, err := s.UserRepo.GetUserByID(receiverID)
	if err != nil || receiver == nil {
		return errors.New("receiver does not exist")
	}

	return s.TransactionRepo.TransferCoins(senderID, receiver.ID, amount)
}
