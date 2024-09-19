package pages

import (
	"errors"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type Cookie2Page struct {
	Identifier string
	NextPage   string
}

func (cp *Cookie2Page) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, cp.Identifier)
	if err != nil {
		return valid, err
	}

	if answer == "Llanfairpwllgwyngyll" {
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

func (cp *Cookie2Page) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"You probably thought I made a mistake and did the same round twice. Got you.",
		"This time it has nothing to do with cookies",
		"The answer is in the image",
		"Though it doesn't look very scenic, this image was taken at a pretty nice location",
		"It's a funny village name",
		"... in Wales",
		"so download the image, inspect the EXIF data, look up the coordinates and be surprised by the town name at that location. Taking this hint has cost you so many points",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
