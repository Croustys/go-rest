package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
	Tags     []Tag              `bson:"tags"`
}
type Tag struct {
	Id      int    `bson:"id,omitempty"`
	TagType string `bson:"tag_type,omitempty"`
}

func Insert(name string, pwd string, email string) bool {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}
	mongo_uri := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	if err != nil {
		log.Fatal(err)
	}

	password := []byte(pwd)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("main")
	userDB := database.Collection("users")
	newUser := User{
		Name:     name,
		Password: string(hashedPassword),
		Email:    email,
	}
	insertUser, err := userDB.InsertOne(ctx, newUser)
	if err != nil {
		panic(err)
	}
	fmt.Println(insertUser.InsertedID)
	return true
}

func Login(email string, pwd string) bool {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}

	mongo_uri := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database("main").Collection("users")

	filter := bson.D{primitive.E{Key: "email", Value: email}}

	var result User
	if err = coll.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		panic(err)
	}

	hashedPassword := []byte(result.Password)

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(pwd))
	if err != nil {
		panic(err)
	}
	return true
}
