package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

type User struct {
	ID      primitive.ObjectID `bson:"_id"`
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

type Config struct {
	Hosts      []string
	Database   string
	Collection string
}

var Module = fx.Option(
	fx.Provide(GetMongoClient, GetUserCollection),
)

type DB interface {
	Find(context.Context) ([]User, error)
	FindById(context.Context, string) (*User, error)
	AddUser(context.Context, *User) (*User, error)
	UpdateUser(context.Context, *User) (*User, error)
	DeleteUser(context.Context, string) error
}

type dbCollection struct {
	c *mongo.Collection
}

func GetMongoClient(config *Config) (*mongo.Client, error) {
	clientOptions := options.Client()
	clientOptions.SetHosts(config.Hosts)
	return mongo.NewClient(clientOptions)
}

func GetUserCollection(client *mongo.Client, config *Config) *mongo.Collection {
	return client.Database(config.Database).Collection(config.Collection)
}

func (client *dbCollection) Find(ctx context.Context) ([]User, error) {
	var users []User
	cursor, err := client.c.Find(ctx, bson.M{"active": true}, nil)
	if err != nil {
		return nil, err
	}
	err = cursor.Decode(&users)
	return users, err
}

func (client *dbCollection) FindById(ctx context.Context, id string) (*User, error) {
	var user User
	result := client.c.FindOne(ctx, bson.M{"active": true, "_id": id}, nil)
	err := result.Decode(&user)
	return &user, err
}

func (client *dbCollection) AddUser(ctx context.Context, user *User) (*User, error) {
	user.ID = primitive.NewObjectID()
	_, err := client.c.InsertOne(ctx, user)
	return user, err
}

func (client *dbCollection) UpdateUser(ctx context.Context, user *User) (*User, error) {
	var userFromDB *User
	err := client.c.FindOneAndUpdate(ctx, bson.M{"active": true, "_id": user.ID}, user).Decode(userFromDB)
	return userFromDB, err
}

func (client *dbCollection) DeleteUser(ctx context.Context, id string) error {
	var user *User
	return client.c.FindOneAndUpdate(ctx, bson.M{"active": true, "_id": id}, bson.M{"active": false}).Decode(user)
}
