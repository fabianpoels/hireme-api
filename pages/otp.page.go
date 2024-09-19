package pages

import (
	"hireme-api/db"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type OtpPage struct {
	Identifier string
	NextPage   string
}

func (op *OtpPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, op.Identifier)
	if err != nil {
		return valid, err
	}

	cacheClient := db.GetCacheClient()
	otp, err := cacheClient.Get(c, participant.Id.Hex()).Result()
	if err != nil {
		log.Println(err)
		return valid, err
	}

	if otp == answer {
		valid = true

		guesses := participant.Pages["email"].Guesses
		email := guesses[len(guesses)-1]
		filter := bson.M{"_id": participant.Id}
		update := bson.M{"$set": bson.M{"email": email}}
		mongoClient := db.GetDbClient()
		_, err := models.GetParticipantCollection(*mongoClient).UpdateOne(c, filter, update)
		if err != nil {
			log.Println(err)
			return false, err
		}

		err = CorrectAnswer(c, participant, op.Identifier, answer, op.NextPage)
	} else {
		err = WrongGuess(c, participant, op.Identifier, answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (op *OtpPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hr.Hints = []string{}
	hr.HasHintsLeft = false
	return hr, nil
}
