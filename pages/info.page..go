package pages

import (
	"errors"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type InfoPage struct {
}

func (i *InfoPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
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

func (i *InfoPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"This is an info page, so you don't need hints to answer it correctly. Anyway, you lost some points now",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
