package pages

import (
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type ZeroPage struct {
}

func (z *ZeroPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, "zero")
	if err != nil {
		return valid, err
	}

	if answer == "i'm a moron" {
		valid = true
		err = CorrectAnswer(c, participant, "zero", answer, "info")
	} else {
		err = WrongGuess(c, participant, "zero", answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}
