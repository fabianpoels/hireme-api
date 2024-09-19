package pages

import (
	"errors"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type ButtonPage struct {
	Identifier string
	NextPage   string
}

func (bp *ButtonPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, bp.Identifier)
	if err != nil {
		return valid, err
	}

	if answer == "clickedthebutton" {
		valid = true
		err = CorrectAnswer(c, participant, bp.Identifier, answer, bp.NextPage)
	} else {
		err = WrongGuess(c, participant, bp.Identifier, answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (bp *ButtonPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"Look closely",
		"It's not because things are hidden, that they are not there",
		"Like I said, there is a button on the page, you just don't see it - yet",
		"If this is already troubling you, the next rounds will be hard. Inspect the page, find the hidden button and remove the css rule that hides it.",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
