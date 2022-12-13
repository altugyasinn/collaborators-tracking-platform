package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string
	CompanyName  string
	Location     string
	PhoneNumber  int
	SecondNumber int
	IDNumber     int
	LastTonnage  int
}
