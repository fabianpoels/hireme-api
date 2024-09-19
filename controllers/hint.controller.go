package controllers

import (
	"hireme-api/db"
	"hireme-api/middleware"
	"hireme-api/models"
	"hireme-api/pages"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HintRequest struct {
	SessionId string `json:"sessionId" binding:"required"`
	Page      string `json:"page" binding:"required"`
}

type HintController struct {
}

func (h HintController) GetHints(c *gin.Context) {
	participant, ok := middleware.GetParticipantFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Unauthorized": "We didn't find your session hombre"})
		return
	}

	var hr HintRequest
	if err := c.ShouldBindBodyWith(&hr, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you provided me with a garbage request"})
		return
	}

	if participant.Page != hr.Page {
		c.JSON(http.StatusBadRequest, gin.H{"error": "what are you trying to accomplish here"})
		return
	}

	page, err := pages.GetPage(participant.Page)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// create the page in the db if it doesn't exist
	err = pages.EnsurePage(c, participant, participant.Page)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := page.GetHintsForPage(participant.Pages[participant.Page])
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h HintController) Hint(c *gin.Context) {
	participant, ok := middleware.GetParticipantFromContext(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Unauthorized": "We didn't find your session hombre"})
		return
	}

	var hr HintRequest
	if err := c.ShouldBindBodyWith(&hr, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you provided me with a garbage request"})
		return
	}

	if participant.Page != hr.Page {
		c.JSON(http.StatusBadRequest, gin.H{"error": "what are you trying to accomplish here"})
		return
	}

	page, err := pages.GetPage(participant.Page)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// create the page in the db if it doesn't exist
	err = pages.EnsurePage(c, participant, participant.Page)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	mongoClient := db.GetDbClient()
	hints := participant.Pages[participant.Page].Hints + 1
	penalty := hints * 50
	update := bson.M{
		"$set": bson.M{
			"pages." + participant.Page + ".hints": hints,
			"score":                                participant.Score - penalty,
			"updatedAt":                            time.Now(),
		},
	}
	opts := options.Update().SetUpsert(false)
	_, err = models.GetParticipantCollection(*mongoClient).UpdateOne(c, bson.M{"_id": participant.Id}, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = models.GetParticipantCollection(*mongoClient).FindOne(c, bson.D{{"sessionId", participant.SessionId}}).Decode(&participant)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "barely looked and still didn't find anything"})
		return
	}

	resp, err := page.GetHintsForPage(participant.Pages[participant.Page])
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"hints": resp, "participant": participant})
}
