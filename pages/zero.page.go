package pages

import (
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type ZeroPage struct {
	Identifier string
	NextPage   string
}

func (z *ZeroPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, z.Identifier)
	if err != nil {
		return valid, err
	}

	if answer == "i'm a moron" {
		valid = true
		err = CorrectAnswer(c, participant, z.Identifier, answer, z.NextPage)
	} else {
		err = WrongGuess(c, participant, z.Identifier, answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (z *ZeroPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hr.Hints = []string{}
	hr.HasHintsLeft = false
	return hr, nil
}
