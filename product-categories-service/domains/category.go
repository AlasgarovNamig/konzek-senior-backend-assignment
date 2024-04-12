package domains

import "time"

type Category struct {
	Id        uint      `gorm:"primary_key:auto_increment" json:"id"`
	CreatedAt time.Time ` json:"createdAt"`
	Name      string    `gorm:"uniqueIndex;type:varchar(20)" json:"name"`
}
