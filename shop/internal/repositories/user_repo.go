package repositories

import (
	"database/sql"
	"errors"
	"shop/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CreateUser(username, passwordHash string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(`
		INSERT INTO users (username, password_hash, coins) 
		VALUES ($1, $2, 1000) RETURNING id, username, coins, created_at`,
		username, passwordHash).
		Scan(&user.ID, &user.Username, &user.Coins, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(`
		SELECT id, username, password_hash, coins, created_at FROM users WHERE username = $1`,
		username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	user.ID = id
	err := r.DB.QueryRow(`
		SELECT username, password_hash, coins, created_at FROM users WHERE id = $1`,
		id).Scan(&user.Username, &user.PasswordHash, &user.Coins, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserInfo(userID int) (map[string]interface{}, error) {
	var coins int
	err := r.DB.QueryRow("SELECT coins FROM users WHERE id = $1", userID).Scan(&coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	transactions := []struct {
		FromUser int `json:"from_user"`
		Amount   int `json:"amount"`
	}{}

	rows, err := r.DB.Query("SELECT sender_id, amount FROM transactions WHERE receiver_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var txn struct {
			FromUser int `json:"from_user"`
			Amount   int `json:"amount"`
		}
		if err := rows.Scan(&txn.FromUser, &txn.Amount); err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}

	return map[string]interface{}{
		"coins":                 coins,
		"transactions_received": transactions,
	}, nil
}
