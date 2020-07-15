package reject

import "github.com/jinzhu/gorm"

//Reject represents the table definition of the Rejects table in the db
//this will contain rejected card numbers and the operation they are not
//allowed to perform
type Reject struct {
	gorm.Model
	CardNumber string
	Operation  string
}
