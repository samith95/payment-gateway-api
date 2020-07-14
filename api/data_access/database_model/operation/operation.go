package operation

import (
	"github.com/jinzhu/gorm"
)

//Operation represents the table definition of the Operations table in the db
type Operation struct {
	gorm.Model
	AuthID   string `gorm:"column:auth_id"`
	Name     string
	Amount   float32
	Currency string
}
