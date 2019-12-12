package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID      primitive.ObjectID `bson:"_id" json:""`
	Active  bool               `bson:"active" json:"active,omitempty"`
	Address struct {
		Area    string `bson:"area" json:"area,omitempty"`
		City    string `bson:"city" json:"city,omitempty"`
		Country string `bson:"country" json:"country,omitempty"`
	} `bson:"address" json:"address,omitempty"`
	Age        int    `bson:"age" json:"age,omitempty"`
	Favourites []int  `bson:"favourites" json:"favourites,omitempty"`
	FirstName  string `bson:"first_name" json:"first_name,omitempty"`
	LastName   string `bson:"last_name" json:"last_name,omitempty"`
}
