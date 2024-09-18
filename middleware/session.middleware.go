package middleware

import (
	"hireme-api/db"
	"hireme-api/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
)

type SessionRequest struct {
	SessionId string `json:"sessionId" binding:"required"`
}

func LoadSession() gin.HandlerFunc {
	var mongoClient = db.GetDbClient()

	return func(c *gin.Context) {
		var sr SessionRequest

		if err := c.ShouldBindBodyWith(&sr, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "are you trying to hack me?"})
			return
		}

		var participant models.Participant
		err := models.GetParticipantCollection(*mongoClient).FindOne(c, bson.D{{"sessionId", sr.SessionId}}).Decode(&participant)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "session not found"})
			return
		}

		c.Set("participant", participant)
		c.Next()
	}
}

func GetParticipantFromContext(c *gin.Context) (models.Participant, bool) {
	participantInterface, exists := c.Get("participant")
	if !exists {
		return models.Participant{}, false
	}
	participant, ok := participantInterface.(models.Participant)
	return participant, ok
}
