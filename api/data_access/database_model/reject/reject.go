package reject

import "github.com/jinzhu/gorm"

type Reject struct {
	gorm.Model
	CardNumber string
	Operation  string
}
