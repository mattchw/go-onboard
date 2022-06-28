package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FirstName string             `bson:"firstName" json:"firstName" validate:"required"`
	LastName  string             `bson:"lastName" json:"lastName" validate:"required"`
	Bio       string             `json:"bio,omitempty"`
	Age       int                `json:"age,omitempty"`
	Gender    string             `json:"gender,omitempty"`
}

func (user User) ValidateUser() error {
	err := validation.ValidateStruct(&user,
		validation.Field(&user.FirstName, validation.Required),
		validation.Field(&user.LastName, validation.Required),
		validation.Field(&user.Age, validation.Min(1)),
		validation.Field(&user.Gender, validation.In("Male", "Female", "Others")),
	)
	return err
}
