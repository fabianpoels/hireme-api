package pages

import (
	"hireme-api/db"
	"hireme-api/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ZeroPage struct {
}

func (z *ZeroPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	mongoClient := db.GetDbClient()

	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, "zero")
	if err != nil {
		return valid, err
	}

	var update = bson.M{}
	if answer == "i'm a moron" {
		valid = true
		update = bson.M{
			"$push": bson.M{
				"pages.zero.guesses": answer,
			},
			"$set": bson.M{
				"page":      "info",
				"score":     participant.Score + 1,
				"updatedAt": time.Now(),
			},
		}
	} else {
		update = bson.M{
			"$push": bson.M{
				"pages.zero.guesses": answer,
			},
			"$set": bson.M{
				"score":     participant.Score + 10,
				"updatedAt": time.Now(),
			},
		}
	}
	opts := options.Update().SetUpsert(true)
	_, err = models.GetParticipantCollection(*mongoClient).UpdateOne(c, bson.M{"_id": participant.Id}, update, opts)

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}
