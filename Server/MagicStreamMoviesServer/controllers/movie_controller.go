package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/go-playground/validator/v10"

	model "github.com/Phc01/MagicStreamMovies/Server/MagicStreamMoviesServer/models"

	db "github.com/Phc01/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
)

var movieCollection *mongo.Collection = db.OpenCollection("movies")

var validate = validator.New()


// busca todos os filmes
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)

		defer cancel()
		var movies []model.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to fetch movies"})
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to decode movies"})
		}

		c.JSON(http.StatusOK, movies)
	}
}

// busca o filme de acordo com o id do imdb
func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context)  {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error":"Movie ID is required"})
			return 
		}
		var movie model.Movie

		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error":"Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie model.Movie
		if err := c.ShouldBindJSON(&movie); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid input"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":"Validation failed", "details": err.Error()})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to add movie"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}