package person

import "time"

type Person struct {
	Name        string    `json:"name" bson:"name"`
	Email       string    `json:"email" bson:"email"`
	PhoneNumber string    `json:"phoneNumber" bson:"phoneNumber"`
	Address     string    `json:"address" bson:"address"`
	Company     string    `json:"company" bson:"company"`
	CreatedOn   time.Time `json:"createdOn" bson:"createdon"`
}
