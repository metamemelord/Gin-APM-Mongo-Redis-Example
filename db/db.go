package db

import (
	"context"
	"fmt"
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/configuration"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

var Module = fx.Option(
	fx.Provide(getMongoClient, getMongoDatabase, GetCollection),
)

type DB interface {
	Find(context.Context) ([]*User, error)
	FindByFilters(context.Context, map[string]interface{}) ([]*User, error)
	FindById(context.Context, string) (*User, error)
	FindOneByFilters(context.Context, map[string]interface{}) (*User, error)
	AddUser(context.Context, *User) (*User, error)
	UpdateUser(context.Context, string, map[string]interface{}) (*User, error)
	DeleteUser(context.Context, string) error
	OverwriteUser(context.Context, string, *User) (*User, error)
}

type dbCollection struct {
	c *mongo.Collection
}

func getMongoClient(config *configuration.Configuration) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Database.ConnectionString))
	if err != nil {
		return nil, err
	}
	err = client.Connect(context.Background())
	return client, err
}

func getMongoDatabase(client *mongo.Client, config *configuration.Configuration) (*mongo.Database, error) {
	db := client.Database(config.Database.Database)
	if db == nil {
		return nil, fmt.Errorf("Database (%s) does not exist", config.Database.Database)
	}
	return db, nil
}

func GetCollection(db *mongo.Database, config *configuration.Configuration) (DB, error) {
	collection := db.Collection(config.Database.Collection)
	if collection == nil {
		return nil, fmt.Errorf("Collection (%s) does not exist on DB (%s)", config.Database.Collection, config.Database.Database)
	}
	return &dbCollection{c: collection}, nil
}

func (client *dbCollection) Find(ctx context.Context) ([]*User, error) {
	return client.FindByFilters(ctx, nil)
}

func (client *dbCollection) FindByFilters(ctx context.Context, filters map[string]interface{}) ([]*User, error) {
	users := []*User{}
	fields := primitive.M{"active": true}
	for k, v := range filters {
		fields[k] = v
	}
	cursor, err := client.c.Find(ctx, fields, nil)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &users)
	return users, err
}

func (client *dbCollection) FindById(ctx context.Context, id string) (*User, error) {
	return client.FindOneByFilters(ctx, map[string]interface{}{"_id": id})
}

func (client *dbCollection) FindOneByFilters(ctx context.Context, filters map[string]interface{}) (*User, error) {
	var user User
	fields := []bson.M{{"active": true}}
	for k, v := range filters {
		fields = append(fields, bson.M{k: v})
	}
	result := client.c.FindOne(ctx, fields, nil)
	err := result.Decode(&user)
	return &user, err
}

func (client *dbCollection) AddUser(ctx context.Context, user *User) (*User, error) {
	user.ID = primitive.NewObjectID()
	_, err := client.c.InsertOne(ctx, user)
	return user, err
}

func (client *dbCollection) UpdateUser(ctx context.Context, id string, fields map[string]interface{}) (*User, error) {
	var user *User
	updateFields := []bson.M{}
	for k, v := range fields {
		updateFields = append(updateFields, bson.M{k: v})
	}
	err := client.c.FindOneAndUpdate(ctx, bson.M{"active": true, "_id": id}, updateFields).Decode(user)
	return user, err
}

func (client *dbCollection) OverwriteUser(ctx context.Context, id string, user *User) (*User, error) {
	res := client.c.FindOneAndReplace(ctx, bson.M{"_id": id}, user)
	if res.Err() != nil {
		return nil, res.Err()
	}
	err := res.Decode(user)
	return user, err
}

func (client *dbCollection) DeleteUser(ctx context.Context, id string) error {
	_, err := client.UpdateUser(ctx, id, map[string]interface{}{"active": false})
	return err
}
