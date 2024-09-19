package pages

import (
	"errors"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type ComplexPage struct {
	Identifier string
	NextPage   string
}

func (thePage *ComplexPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, thePage.Identifier)
	if err != nil {
		return valid, err
	}

	if answer == "" {
		valid = true
		err = CorrectAnswer(c, participant, thePage.Identifier, answer, thePage.NextPage)
	} else {
		err = WrongGuess(c, participant, thePage.Identifier, answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (thePage *ComplexPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"Nothing technologically advanced about this round",
		"Some rounds, I expected a lot from you. Not this time",
		"Less is more",
		"And nothing is even better",
		"worthless hint #5",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
