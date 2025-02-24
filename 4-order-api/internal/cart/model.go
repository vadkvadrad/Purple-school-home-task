package cart

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	Phone string `json:"phone"`
	Products pq.StringArray `json:"products" gorm:"type:text"`
	Date datatypes.Date `json:"date"`
}