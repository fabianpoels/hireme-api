package pages

import (
	"errors"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type PingPage struct {
	Identifier string
	NextPage   string
}

func (pp *PingPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, pp.Identifier)
	if err != nil {
		return valid, err
	}

	if answer == "pong" {
		valid = true
		err = CorrectAnswer(c, participant, pp.Identifier, answer, pp.NextPage)
	} else {
		err = WrongGuess(c, participant, pp.Identifier, answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (pp *PingPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"...",
		"Ok if you seriously need a hint for this",
		"pong, maybe?",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
