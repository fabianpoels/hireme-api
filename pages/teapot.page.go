package pages

import (
	"errors"
	"fmt"
	"hireme-api/config"
	"hireme-api/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TeapotPage struct {
	Identifier string
	NextPage   string
}

func (thePage *TeapotPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, thePage.Identifier)
	if err != nil {
		return valid, err
	}

	responseCode := getResponseCode(answer)

	if responseCode == 418 {
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

func (thePage *TeapotPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	finalHint := fmt.Sprintf("and for the real losers who didn't find any publicly available API, I provided one on my side: https://%s/hireme/api/v69/teapot", config.GetEnv("DOMAIN"))
	hints := []string{
		"I will admit this is one of the harder rounds",
		"Your answer has to start with either http:// or https://",
		"If you ever developed an API, you should be aware of this quirk",
		"At some point, NodeJS was about to remove the feature, but luckily the internet rallied and prevented this huge mistake",
		"I'm not a moka",
		"HTCPCP",
		"I'm a teapot",
		"If you haven't figured it out by now: go look for a publicly available API endpoint that returns you response code 418",
		finalHint,
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}

func getResponseCode(answer string) (responseCode int) {
	responseCode = 500

	method, url, ok := splitRequestString(answer)
	if !ok {
		return responseCode
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return responseCode
	}

	req.Header.Add("content-type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return responseCode
	}
	defer res.Body.Close()

	return res.StatusCode
}

func splitRequestString(s string) (string, string, bool) {
	prefixes := []string{"GET", "POST"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix+":") {
			return prefix, strings.TrimPrefix(s, prefix+":"), true
		}
	}
	return "", "", false
}
