package pages

import (
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type InfoPage struct {
}

func (z *InfoPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, "info")
	if err != nil {
		return valid, err
	}

	if answer == "i'm still a moron" {
		valid = true
		err = CorrectAnswer(c, participant, "info", answer, "email")
	} else {
		err = WrongGuess(c, participant, "info", answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}
