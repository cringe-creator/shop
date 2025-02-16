package repositories

import (
	"database/sql"
	"errors"
	"shop/internal/models"
)

type TransactionRepository struct {
	DB *sql.DB
}

func (r *TransactionRepository) TransferCoins(senderID, receiverID, amount int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	var senderBalance int
	err = tx.QueryRow("SELECT coins FROM users WHERE id = $1", senderID).Scan(&senderBalance)
	if err != nil {
		tx.Rollback()
		return err
	}
	if senderBalance < amount {
		tx.Rollback()
		return errors.New("insufficient balance")
	}

	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", amount, senderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE users SET coins = coins + $1 WHERE id = $2", amount, receiverID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO transactions (sender_id, receiver_id, amount) VALUES ($1, $2, $3)",
		senderID, receiverID, amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *TransactionRepository) GetTransactionHistory(userID int) ([]models.Transaction, error) {
	rows, err := r.DB.Query(`
		SELECT sender_id, receiver_id, amount, created_at 
		FROM transactions 
		WHERE sender_id = $1 OR receiver_id = $1 
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *TransactionRepository) GetReceivedTransactions(userID int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(`
		SELECT sender_id, amount 
		FROM transactions 
		WHERE receiver_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []map[string]interface{}
	for rows.Next() {
		var txn struct {
			ToUserID int `json:"fromUserID"`
			Amount   int `json:"amount"`
		}

		if err := rows.Scan(&txn.ToUserID, &txn.Amount); err != nil {
			return nil, err
		}

		row := r.DB.QueryRow(`SELECT username
			FROM users 
			WHERE id = $1`, txn.ToUserID)
		var fromUserName string
		err := row.Scan(&fromUserName)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, map[string]interface{}{
			"fromUser": fromUserName,
			"amount":   txn.Amount,
		})
	}
	return transactions, nil
}

func (r *TransactionRepository) GetSentTransactions(userID int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(`
		SELECT receiver_id, amount 
		FROM transactions 
		WHERE sender_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []map[string]interface{}
	for rows.Next() {
		var txn struct {
			ToUserID int `json:"toUserID"`
			Amount   int `json:"amount"`
		}

		if err := rows.Scan(&txn.ToUserID, &txn.Amount); err != nil {
			return nil, err
		}

		row := r.DB.QueryRow(`SELECT username
			FROM users 
			WHERE id = $1`, txn.ToUserID)
		var toUserName string
		err := row.Scan(&toUserName)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, map[string]interface{}{
			"toUser": toUserName,
			"amount": txn.Amount,
		})
	}
	return transactions, nil
}
