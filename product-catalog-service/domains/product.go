package domains

import (
	"time"
)

type Product struct {
	ID         uint      `gorm:"primary_key:auto_increment" json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	Name       string    `gorm:"uniqueIndex;type:varchar(20)" json:"name"`
	Price      float64   `gorm:"type:decimal(10,2);" json:"price"`
	CategoryID uint      `json:"category_id"`
}
