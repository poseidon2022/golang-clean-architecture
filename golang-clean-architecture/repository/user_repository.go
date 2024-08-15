package repository

import (
	"context"
	"errors"
	"golang-clean-architecture/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct  {
	Database 		*mongo.Database
	Collection 		string
}

func NewUserRepository(db *mongo.Database, collection string) domain.UserRepository {
	return &UserRepository{
		Database : db,
		Collection : collection,
	}
}

func (ur *UserRepository) Register(newUser *domain.User) error {
	collection := ur.Database.Collection(ur.Collection)
	newUser.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(context.TODO(), newUser)
	return err
}

func (ur *UserRepository) VerifyFirst(newUser *domain.User) error {
	collection := ur.Database.Collection(ur.Collection)
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if cur.Next(context.TODO()) {return errors.New("a user is found on db")}
	if err != nil {
		return errors.New("internal server error")
	}
	return nil
}

func (ur *UserRepository) UserExists(newUser *domain.User) error {
	collection := ur.Database.Collection(ur.Collection)
	var existingUser domain.User
	err := collection.FindOne(context.TODO(), bson.D{{Key : "email", Value : newUser.Email}}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {return errors.New("internal server error")}

	if existingUser != (domain.User{}) {return errors.New("user email already in use")}
	return nil
}

func (ur *UserRepository) GetUserByEmail(email string) domain.User {
	collection := ur.Database.Collection(ur.Collection)
	filter := bson.D{{Key : "email", Value : email}}

	var existingUser domain.User
	err := collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err != nil {
		return domain.User{}
	}
	return existingUser
}

func (ur *UserRepository) PromoteUser(userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}
	filter := bson.D{{Key : "_id", Value : objectID}}
	update := bson.D{{Key : "$set", Value : bson.D{{Key : "role", Value : "admin"}}}}
	collection := ur.Database.Collection(ur.Collection)
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	
	if updateResult.MatchedCount == 0 {
		return errors.New("no user with the specified id found")
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("user is already an admin")
	}

	if err != nil {
		return errors.New("internal server error")
	}
	return nil
}