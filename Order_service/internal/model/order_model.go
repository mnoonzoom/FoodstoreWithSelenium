package model

import "time"

type Order struct {
	ID         string    `bson:"_id,omitempty"`
	UserID     string    `bson:"user_id"`
	ItemIDs    []string  `bson:"item_ids"`
	TotalPrice float64   `bson:"total_price"`
	Status     string    `bson:"status"`
	CreatedAt  time.Time `bson:"created_at"`
}
