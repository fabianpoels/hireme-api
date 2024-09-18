package pages

import (
	"crypto/rand"
	"fmt"
	"hireme-api/db"
	"hireme-api/models"
	"log"
	"math/big"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var tempEmailProviders = []string{
	"tempmail.com",
	"throwawaymail.com",
	"10minutemail.com",
	"mailinator.com",
	"guerrillamail.com",
	"yopmail.com",
	"sharklasers.com",
	"dispostable.com",
	"mailnesia.com",
	"trashmail.com",
	"temp-mail.org",
	"getnada.com",
	"fakeinbox.com",
	"guerrillamail.org",
	"temp-mail.io",
	"deadfake.com",
	"mintemail.com",
	"mohmal.com",
	"tempmail.ninja",
	"burnermail.io",
	"mytemp.email",
	"tempmailaddress.com",
}

var emailRegex *regexp.Regexp
var regexOnce sync.Once

func compileRegex() {
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
}

func validateEmail(email string) bool {
	regexOnce.Do(compileRegex)

	if !emailRegex.MatchString(email) {
		return false
	}

	parts := strings.Split(email, "@")
	domain := parts[1]

	for _, tempDomain := range tempEmailProviders {
		if domain == tempDomain {
			return false
		}
	}

	return true
}

func generateRandomNumber() (string, error) {
	max := big.NewInt(10000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Format the number as a 4-digit string with leading zeros
	return fmt.Sprintf("%04d", n), nil
}

type EmailPage struct {
}

func (e *EmailPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, "email")
	if err != nil {
		return valid, err
	}

	if validateEmail(answer) {
		cacheClient := db.GetCacheClient()
		valid = true
		// generate confirmation code and store in cache
		randomNumber, err := generateRandomNumber()
		if err != nil {
			log.Println(err)
			return valid, err
		}

		err = cacheClient.Set(c, participant.Id.Hex(), randomNumber, 0).Err()
		if err != nil {
			log.Println(err)
			return valid, err
		}

		log.Printf("OTP set for %s : %s", answer, randomNumber)
		// TODO
		// SENT EMAIL

		err = CorrectAnswer(c, participant, "email", answer, "otp")
	} else {
		err = WrongGuess(c, participant, "email", answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (e *EmailPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hr.Hints = []string{}
	hr.HasHintsLeft = false
	return hr, nil
}
