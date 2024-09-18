package pages

import (
	"hireme-api/db"
	"hireme-api/models"
	"log"

	"github.com/gin-gonic/gin"
)

type OtpPage struct {
}

func (z *OtpPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, "otp")
	if err != nil {
		return valid, err
	}

	cacheClient := db.GetCacheClient()
	otp, err := cacheClient.Get(c, participant.Id.Hex()).Result()
	if err != nil {
		log.Println(err)
		return valid, err
	}

	if otp == answer {
		valid = true

		err = CorrectAnswer(c, participant, "otp", answer, "ping")
	} else {
		err = WrongGuess(c, participant, "otp", answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (o *OtpPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hr.Hints = []string{}
	hr.HasHintsLeft = false
	return hr, nil
}
