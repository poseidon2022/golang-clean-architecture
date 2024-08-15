// +build db

package repository_test

import (
	"context"
	"testing"
	"golang-clean-architecture/domain"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang-clean-architecture/repository"
)


type APITestSuite struct {
	suite.Suite
	client 			*mongo.Client
	db				*mongo.Database
	collection		*mongo.Collection
	repo			domain.UserRepository
}

func (suite *APITestSuite) SetupSuite() {
	//connecting to the mongo db database and setting the values for the
	//db, client, collection and repo of the struct
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		suite.T().Fatal(err)
	}

	db := client.Database("test")
	collection := db.Collection("users")

	suite.client = client
	suite.collection = collection
	suite.db = db
	suite.repo = repository.NewUserRepository(db, "users")
}

func (suite *APITestSuite) TearDownSuite() {
    // Drop the test database after all tests are run
    err := suite.db.Drop(context.TODO())
    if err != nil {
        suite.T().Fatal(err)
    }

    // Disconnect the MongoDB client
    err = suite.client.Disconnect(context.TODO())
    if err != nil {
        suite.T().Fatal(err)
    }
}

func (suite *APITestSuite) TestRegisterSuccess() {
	user := &domain.User{
		ID: primitive.NewObjectID(),
		Email: "kidusm3l@gmail.com",
		Password: "123456789hbihbi",
		Role : "admin",
	}

	err := suite.repo.Register(user)
	suite.NoError(err, "no error when registering user")
	var insertedUser domain.User
	err = suite.collection.FindOne(context.TODO(), bson.M{"_id" : user.ID}).Decode(&insertedUser)

	suite.NoError(err, "no error when retreiving user")
	suite.Equal(user.Role, insertedUser.Role)
	suite.Equal(user.Email, insertedUser.Email)
}

func (suite *APITestSuite) TestSignUp_DupEmail() {

	user := &domain.User{
		ID: primitive.NewObjectID(),
		Email: "kidus.melaku@gmail.com",
		Password: "123456789",
		Role : "user",
	}

	err := suite.repo.Register(user)
	suite.NoError(err, "no error when registering user")
	err = suite.repo.UserExists(user)
	suite.Error(err, "error user email already exists")
}

func (suite *APITestSuite) TestPromoteUser_Positive() {
	user := &domain.User{
		ID: primitive.NewObjectID(),
		Email: "kidusm33l@gmail.com",
		Password: "123456789",
		Role : "user",
	}
	err := suite.repo.Register(user)
	suite.NoError(err, "no error when registering user")
	err = suite.repo.PromoteUser(user.ID.Hex())
	suite.NoError(err, "no error while promoting user")
}

func (suite *APITestSuite) TestPromoteUser_Negative() {

	user := &domain.User{
		ID: primitive.NewObjectID(),
		Email: "kidus.another@gmail.com",
		Password: "123456789",
		Role : "admin",
	}

	err := suite.repo.Register(user)
	suite.NoError(err, "no error when registering user")
	err = suite.repo.PromoteUser(user.ID.Hex())
	suite.Error(err, "error user already an admin")
}


func TestAPITestSuite(t *testing.T) {
    suite.Run(t, new(APITestSuite))
}
