package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Email string             `bson:"email,omitempty"`
}

func Insert(name string, email string) {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}
	mongo_uri := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("main")
	userDB := database.Collection("users")
	newUser := User{
		Name:  name,
		Email: email,
	}
	insertUser, err := userDB.InsertOne(ctx, newUser)
	if err != nil {
		panic(err)
	}
	fmt.Println(insertUser.InsertedID)
}
