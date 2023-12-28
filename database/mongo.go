package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func Dial(url string, dbName string, sessionNum int) {
	client, err := mongo.Connect(context.Background())
	client.Database("")
}
