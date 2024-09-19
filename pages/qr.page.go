package pages

import (
	"errors"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type QrPage struct {
	Identifier string
	NextPage   string
}

func (qp *QrPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, qp.Identifier)
	if err != nil {
		return valid, err
	}

	if answer == "https://fabianpoels.com/hireme" {
		valid = true
		err = CorrectAnswer(c, participant, qp.Identifier, answer, qp.NextPage)
	} else {
		err = WrongGuess(c, participant, qp.Identifier, answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (qp *QrPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"That's right, you need to allow your camera to be used",
		"For sure I was too lazy to implement some kind of facial recognition pattern",
		"So it's probably a very common form of image recognition",
		"You guessed it right, you need to provide me with a QR code for a very specific URL",
		"Let's say this webpage has very narcissistic tendencies",
		"... just wasting some points for you on this worthless hint",
		"Which url are you currently on?",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
