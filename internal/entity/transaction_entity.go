package entity

import "github.com/google/uuid"

type Transaction struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`

	AmountPaid int    `json:"amount_paid"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	InvoiceURL string `json:"invoice_url"`

	Reference string `json:"reference"` // untuk webhook

	User *User `gorm:"foreignKey:UserID"`

	Timestamp
}
