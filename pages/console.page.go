package pages

import (
	"errors"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type ConsolePage struct {
	Identifier string
	NextPage   string
}

func (cp *ConsolePage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, cp.Identifier)
	if err != nil {
		return valid, err
	}

	if answer == "potato" {
		valid = true
		err = CorrectAnswer(c, participant, cp.Identifier, answer, cp.NextPage)
	} else {
		err = WrongGuess(c, participant, cp.Identifier, answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (cp *ConsolePage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"I swear I already gave you the answer",
		"Though it's not visible on the page, it's not hidden",
		"Come one, act like you're a real hacker",
		"console.log doesn't ring a bell?",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
