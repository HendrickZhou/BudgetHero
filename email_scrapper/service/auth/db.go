package auth

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri             string = "mongodb://localhost:27017"
	collection_name string = "userauth"
	db_name         string = "db"
)

const ( // extentable
	mail_gmail   = 1
	mail_outlook = 2
)

type UserAuth struct {
	Email    string `bson:"omitempty"`
	Token    string `bson:"omitempty"`
	Provider int    `bson:"omitempty"`
}

func Connect(client *mongo.Client) (*mongo.Database, error) {
	// todo add authentication for real business
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	db := client.Database(db_name)
	return db, err
}

func Disconnect(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func useDoc(db *mongo.Database) *mongo.Collection {
	return db.Collection(collection_name)
}

///// CRUD

func db_saveUser(coll *mongo.Collection, email string, provider int) (primitive.ObjectID, error) {
	doc := UserAuth{Email: email, Provider: provider}
	if result, err := coll.InsertOne(context.TODO(), doc); err == nil {
		log.Print(result)
		objId := result.InsertedID.(primitive.ObjectID)
		return objId, nil
	} else {
		log.Printf("insert user fail! %s", err)
		return primitive.NilObjectID, err
	}
}

func db_saveToken(coll *mongo.Collection, oid *primitive.ObjectID, token string) error {
	doc := bson.D{{Key: "$set", Value: bson.D{{Key: "token", Value: token}}}}
	if result, err := coll.UpdateByID(context.TODO(), oid, doc); err == nil {
		log.Print(result)
		return nil
	} else {
		log.Printf("insert user token fail! %s", err)
		return err
	}
}
