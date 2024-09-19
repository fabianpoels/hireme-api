package pages

import (
	"hireme-api/models"

	"github.com/gin-gonic/gin"
)

type ScorePage struct {
	Identifier string
	NextPage   string
}

func (z *ScorePage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	return valid, nil
}

func (z *ScorePage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hr.Hints = []string{}
	hr.HasHintsLeft = false
	return hr, nil
}
