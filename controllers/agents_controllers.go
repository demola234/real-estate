package controllers

import (
	"context"
	"time"

	"github.com/demola234/real-estate/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var agentCollection *mongo.Collection = database.OpenCollection(database.Client, "agents")

func BecomeAnAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		agentCollection.InsertOne(ctx, bson.M{})
	}

}

func GetAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}

func GetAgents() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}

func UpdateAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}

func DeleteAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
	}

}
