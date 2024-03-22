package entity

import "time"

type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `json:"created_at"` // TODO accshualllyy, we dont need this
}
