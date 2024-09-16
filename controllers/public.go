package controllers

import (
	"hireme-api/config"
	"hireme-api/db"
	"hireme-api/models"
	"hireme-api/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type PublicController struct {
}

type StatusRequest struct {
	SessionId string `json:"sessionId" binding:"required"`
}

func (p PublicController) Status(c *gin.Context) {
	mongoClient := db.GetDbClient()

	var sr StatusRequest

	if err := c.ShouldBindJSON(&sr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "get your shit together"})
		return
	}

	var participant models.Participant
	err := models.GetParticipantCollection(*mongoClient).FindOne(c, bson.D{{"sessionId", sr.SessionId}}).Decode(&participant)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "barely looked and still didn't find anything"})
		return
	}

	c.JSON(http.StatusOK, participant)
}

func (p PublicController) Init(c *gin.Context) {
	mongoClient := db.GetDbClient()
	startPage := config.GetEnv("STARTPAGE")
	randomId := utils.GenerateRandomString()
	score, err := strconv.Atoi(config.GetEnv("STARTINGSCORE"))
	if err != nil {
		score = 0
	}

	var participant models.Participant

	participant.CreatedAt = time.Now()
	participant.UpdatedAt = time.Now()
	participant.SessionId = randomId
	participant.Page = startPage
	participant.Score = score

	result, err := models.GetParticipantCollection(*mongoClient).InsertOne(c, participant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var createdParticipant models.Participant
	err = models.GetParticipantCollection(*mongoClient).FindOne(c, bson.D{{"_id", result.InsertedID}}).Decode(&createdParticipant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdParticipant)
}
