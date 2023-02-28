package controllers

import (
	"github.com/demola234/real-estate/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var propertyCollection *mongo.Collection = database.OpenCollection(database.Client, "properties")
func AddProperty() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

func GetProperty() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetProperties() gin.HandlerFunc {
	return func(c *gin.Context) {

	}

}

func UpdateProperty() gin.HandlerFunc {
	return func(c *gin.Context) {

	}

}

func DeleteProperty() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetPropertyByAgent() gin.HandlerFunc {
	return func(c *gin.Context) {

	}

}

func RateProperty() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetPropertyByLocation() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func FeaturedProperties() gin.HandlerFunc {
	return func(c *gin.Context) {


	}

}
