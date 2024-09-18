package controllers

import (
	"fmt"
	"hireme-api/config"
	"hireme-api/db"
	"hireme-api/models"
	"hireme-api/pages"
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

type AnswerRequest struct {
	SessionId string `json:"sessionId" binding:"required"`
	Answer    string `json:"answer" binding:"required"`
}

func getPage(pageType string) (pages.Page, error) {
	switch pageType {
	case "zero":
		return &pages.ZeroPage{}, nil
	case "info":
		return &pages.InfoPage{}, nil
	case "email":
		return &pages.EmailPage{}, nil
	case "otp":
		return &pages.OtpPage{}, nil
	default:
		return nil, fmt.Errorf("unknown page type: %s", pageType)
	}
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
	mongoClient := db.GetDbClient()

	var ar AnswerRequest

	if err := c.ShouldBindJSON(&ar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "get your shit together"})
		return
	}

	var participant models.Participant
	err := models.GetParticipantCollection(*mongoClient).FindOne(c, bson.D{{"sessionId", ar.SessionId}}).Decode(&participant)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "barely looked and still didn't find anything"})
		return
	}

	page, err := getPage(participant.Page)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something strange happened and I didn't write enough logic to handle the case properly"})
		return
	}

	valid, err := page.ProvideAnswer(ar.Answer, participant, c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something strange happened and I didn't write enough logic to handle the case properly"})
		return
	}

	err = models.GetParticipantCollection(*mongoClient).FindOne(c, bson.D{{"sessionId", ar.SessionId}}).Decode(&participant)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "barely looked and still didn't find anything"})
		return
	}

	if !valid {
		c.JSON(http.StatusBadRequest, participant)
	}

	c.JSON(http.StatusOK, participant)
}

func (p PublicController) Hint(c *gin.Context) {
}
