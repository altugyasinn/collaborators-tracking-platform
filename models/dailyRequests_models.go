package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DailyRequests struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CrateRequest      int
	MoneyRequest      int
	EmptyCrateRequest int
	PostScript        string
}
