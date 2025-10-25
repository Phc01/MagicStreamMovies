package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	model "github.com/Phc01/MagicStreamMovies/Server/MagicStreamMoviesServer/models"

	db "github.com/Phc01/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
)

var movieCollection *mongo.Collection = db.OpenCollection("movies")

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)

		defer cancel()
		var movies []model.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies"})
		}

		c.JSON(http.StatusOK, movies)
	}
}
