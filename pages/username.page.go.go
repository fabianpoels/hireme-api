package pages

import (
	"errors"
	"hireme-api/db"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type UsernamePage struct {
	Identifier string
	NextPage   string
}

func (up *UsernamePage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, up.Identifier)
	if err != nil {
		return valid, err
	}

	startsWithDigit := answer[0] >= '0' && answer[0] <= '9'

	mongoClient := db.GetDbClient()
	filter := bson.M{"username": answer}
	count, err := models.GetParticipantCollection(*mongoClient).CountDocuments(c, filter)
	if err != nil {
		log.Println(err)
		return valid, err
	}

	if startsWithDigit && count < 1 {
		valid = true
		err = CorrectAnswer(c, participant, up.Identifier, answer, up.NextPage)
	} else {
		err = WrongGuess(c, participant, up.Identifier, answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (up *UsernamePage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"Obviously the username needs to have a digit as the first character",
		"And it needs to be unique, you should have guessed which ones were taken",
		"Haha fooled you, there's no more hints",
		"Or are there?",
		"No there aren't",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
