package pages

import (
	"errors"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type TimerPage struct {
	Identifier string
	NextPage   string
}

func (thePage *TimerPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, thePage.Identifier)
	if err != nil {
		return valid, err
	}

	if answer == "13" {
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

func (thePage *TimerPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"This is all skill baby",
		"or is it?",
		"Maybe you could hack it with some javascript in the console?",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
