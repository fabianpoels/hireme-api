package pages

import (
	"errors"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type CookiePage struct {
	Identifier string
	NextPage   string
}

func (cp *CookiePage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, cp.Identifier)
	if err != nil {
		return valid, err
	}

	if answer == "isBlowingInTheWind" {
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

func (cp *CookiePage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"After the previous round, I'd though I give an easier challenge",
		"As if the image wasn't enough of a hint, it has to do with cookies this time",
		"If you made it this far, you won't need another hint, do you?",
		"Just go and inspect the cookies for this page",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
