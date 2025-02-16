package repositories

import (
	"database/sql"
	"shop/internal/models"
)

type PurchaseRepository struct {
	DB *sql.DB
}

func (r *PurchaseRepository) CreatePurchase(userID int, itemName string, price int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET coins = coins - $1 WHERE id = $2", price, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO purchases (user_id, item_name, price) VALUES ($1, $2, $3)", userID, itemName, price)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *PurchaseRepository) GetPurchaseHistory(userID int) ([]models.Purchase, error) {
	rows, err := r.DB.Query(`
		SELECT item_name, price, created_at 
		FROM purchases 
		WHERE user_id = $1 
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []models.Purchase
	for rows.Next() {
		var purchase models.Purchase
		if err := rows.Scan(&purchase.ItemName, &purchase.Price, &purchase.CreatedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}
	return purchases, nil
}

func (r *PurchaseRepository) GetUserInventory(userID int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(`
		SELECT item_name, COUNT(*) as quantity
		FROM purchases 
		WHERE user_id = $1 
		GROUP BY item_name
		ORDER BY item_name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows == nil {
		return []map[string]interface{}{}, nil
	}

	var inventory []map[string]interface{}

	for rows.Next() {
		var item struct {
			Type     string `json:"type"`
			Quantity int64  `json:"quantity"`
		}
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, err
		}
		inventory = append(inventory, map[string]interface{}{
			"type":     item.Type,
			"quantity": item.Quantity,
		})
	}

	if len(inventory) == 0 {
		return []map[string]interface{}{}, nil
	}

	return inventory, nil
}
