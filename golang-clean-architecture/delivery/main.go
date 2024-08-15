package main

import (
	"context"
	"fmt"
	routers "golang-clean-architecture/delivery/router"
	"log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database Connected")
	db := client.Database("task_management")
	router := gin.Default()
	routers.Setup(db, router)
	router.Run("localhost:8080")
}