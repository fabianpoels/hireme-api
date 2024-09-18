package controllers

import (
	"hireme-api/config"
	"hireme-api/db"
	"hireme-api/middleware"
	"hireme-api/models"
	"hireme-api/pages"
	"hireme-api/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
)

type PublicController struct {
}

type AnswerRequest struct {
	Answer string `json:"answer" binding:"required"`
}

func (p PublicController) Status(c *gin.Context) {
	participant, ok := middleware.GetParticipantFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Unauthorized": "We didn't find your session hombre"})
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
		score = 1000
	}

	var participant models.Participant

	participant.CreatedAt = time.Now()
	participant.UpdatedAt = time.Now()
	participant.SessionId = randomId
	participant.Page = startPage
	participant.Score = score
	participant.Pages = make(map[string]models.Page)

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

func (p PublicController) Answer(c *gin.Context) {
	participant, ok := middleware.GetParticipantFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Unauthorized": "We didn't find your session hombre"})
		return
	}

	page, err := pages.GetPage(participant.Page)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something strange happened and I didn't write enough logic to handle the case properly"})
		return
	}

	var ar AnswerRequest
	if err := c.ShouldBindBodyWith(&ar, binding.JSON); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "get your shit together"})
		return
	}

	valid, err := page.ProvideAnswer(ar.Answer, participant, c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something strange happened and I didn't write enough logic to handle the case properly"})
		return
	}

	mongoClient := db.GetDbClient()
	err = models.GetParticipantCollection(*mongoClient).FindOne(c, bson.D{{"sessionId", participant.SessionId}}).Decode(&participant)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "barely looked and still didn't find anything"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"participant": participant, "valid": valid})
}
