package controllers

import (
	"hireme-api/db"
	"hireme-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScoreController struct {
}

func (s ScoreController) Scores(c *gin.Context) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"page": "score"}}},
		{{Key: "$project", Value: bson.M{
			"sessionId": 1,
			"username":  1,
			"score":     1,
			"guesses": bson.M{"$sum": bson.M{
				"$map": bson.M{
					"input": bson.M{"$objectToArray": "$pages"},
					"as":    "page",
					"in":    bson.M{"$size": "$$page.v.guesses"},
				},
			}},
			"hints": bson.M{"$sum": bson.M{
				"$map": bson.M{
					"input": bson.M{"$objectToArray": "$pages"},
					"as":    "page",
					"in":    "$$page.v.hints",
				},
			}},
		}}},
	}

	mongoClient := db.GetDbClient()
	cursor, err := models.GetParticipantCollection(*mongoClient).Aggregate(c, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer cursor.Close(c)

	var scores []models.Score
	if err = cursor.All(c, &scores); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, scores)
}
