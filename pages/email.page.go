package pages

import (
	"crypto/rand"
	"errors"
	"fmt"
	"hireme-api/db"
	"hireme-api/models"
	"log"
	"math/big"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	Identifier string
	NextPage   string
}

func (ep *EmailPage) ProvideAnswer(answer string, participant models.Participant, c *gin.Context) (valid bool, err error) {
	// create the page in the db if it doesn't exist
	err = EnsurePage(c, participant, ep.Identifier)
	if err != nil {
		return valid, err
	}

	mongoClient := db.GetDbClient()
	filter := bson.M{"email": answer}
	count, err := models.GetParticipantCollection(*mongoClient).CountDocuments(c, filter)
	if err != nil {
		log.Println(err)
		return valid, err
	}

	if validateEmail(answer) && count < 1 {
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

		err = CorrectAnswer(c, participant, ep.Identifier, answer, ep.NextPage)
	} else {
		err = WrongGuess(c, participant, ep.Identifier, answer)
	}

	if err != nil {
		log.Println(err)
		return valid, err
	}

	return valid, nil
}

func (ep *EmailPage) GetHintsForPage(page models.Page) (hr HintsResponse, err error) {
	hints := []string{
		"Electronic mail (email or e-mail) is a method of transmitting and receiving messages using electronic devices. It was conceived in the late–20th century as the digital version of, or counterpart to, mail (hence e- + mail). Email is a ubiquitous and very widely used communication medium; in current use, an email address is often treated as a basic and necessary part of many processes in business, commerce, government, education, entertainment, and other spheres of daily life in most countries.",
		"I obviously don't allow these alias e-mails from e.g. Google which include a '+",
		"... and like I said, you can only use the same address once",
	}

	if page.Hints < 0 || page.Hints > len(hints) {
		return hr, errors.New("the amount of hints taken does not make any sense")
	}

	hr.Hints = hints[:page.Hints]
	hr.HasHintsLeft = page.Hints < len(hints)
	return hr, nil
}
