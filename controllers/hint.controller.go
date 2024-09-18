package controllers

import (
	"hireme-api/db"
	"hireme-api/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
)

type HintRequest struct {
	SessionId string `json:"sessionId" binding:"required"`
	Page      string `json:"page" binding:"required"`
}

type HintController struct {
}

func (h HintController) GetHints(c *gin.Context) {

}

func (h HintController) Hint(c *gin.Context) {
	mongoClient := db.GetDbClient()

	var hr HintRequest

	if err := c.ShouldBindBodyWith(&hr, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "get your shit together"})
		return
	}

	var participant models.Participant
	err := models.GetParticipantCollection(*mongoClient).FindOne(c, bson.D{{"sessionId", hr.SessionId}}).Decode(&participant)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "barely looked and still didn't find anything"})
		return
	}

}
