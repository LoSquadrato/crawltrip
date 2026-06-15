package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CreateUserParams struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type User struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const userCollectionName = "users"
const userDatabaseName = "crawltrip"

func CreateUser(ctx context.Context, client *mongo.Client, params CreateUserParams) (*User, error) {
	log.Println("Saving user to database")
	type InsertUser struct {
		Email          string    `json:"email"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		HashedPassword string    `json:"hashed_password"`
	}
	var user = InsertUser{
		HashedPassword: params.HashedPassword,
		Email:          params.Email,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	res, err := client.Database(userDatabaseName).Collection(userCollectionName).InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	log.Printf("User saved with id: %v\n", res.InsertedID)

	return &User{
		Id:        res.InsertedID.(string),
		Email:     params.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func GetUserByEmail(ctx context.Context, client *mongo.Client, email string) (*User, error) {
	log.Printf("Getting user from database by email: %s\n", email)
	res := client.Database(userDatabaseName).Collection(userCollectionName).FindOne(ctx, bson.M{"email": email})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var user User
	if err := res.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUser(ctx context.Context, client *mongo.Client, user *User) error {
	return nil
}

func UpdateUser(ctx context.Context, client *mongo.Client, user *User) error {
	return nil
}
