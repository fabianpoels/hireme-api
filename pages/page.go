package pages

import (
	"fmt"
	"hireme-api/db"
	"hireme-api/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Page interface {
	ProvideAnswer(string, models.Participant, *gin.Context) (bool, error)
	GetHintsForPage(models.Page) (HintsResponse, error)
}

type HintsResponse struct {
	Hints        []string `json:"hints"`
	HasHintsLeft bool     `json:"hasHintsLeft"`
}

func GetPage(pageType string) (Page, error) {
	switch pageType {
	case "zero":
		return &ZeroPage{Identifier: "zero", NextPage: "info"}, nil
	case "info":
		return &InfoPage{Identifier: "info", NextPage: "email"}, nil
	case "email":
		return &EmailPage{Identifier: "email", NextPage: "otp"}, nil
	case "otp":
		return &OtpPage{Identifier: "otp", NextPage: "ping"}, nil
	case "ping":
		return &PingPage{Identifier: "ping", NextPage: "console"}, nil
	case "console":
		return &ConsolePage{Identifier: "console", NextPage: "username"}, nil
	case "username":
		return &UsernamePage{Identifier: "username", NextPage: "button"}, nil
	case "button":
		return &ButtonPage{Identifier: "button", NextPage: "teapot"}, nil
	case "teapot":
		return &TeapotPage{Identifier: "teapot", NextPage: "cookie"}, nil
	case "cookie":
		return &CookiePage{Identifier: "cookie", NextPage: "cookie2"}, nil
	case "cookie2":
		return &Cookie2Page{Identifier: "cookie2", NextPage: "qr"}, nil
	case "qr":
		return &QrPage{Identifier: "qr", NextPage: "complex"}, nil
	case "complex":
		return &ComplexPage{Identifier: "complex", NextPage: "complex"}, nil
	default:
		return nil, fmt.Errorf("unknown page type: %s", pageType)
	}
}

func EnsurePage(c *gin.Context, participant models.Participant, pageKey string) error {
	mongoClient := db.GetDbClient()

	initUpdate := bson.M{
		"$set": bson.M{
			"pages":     bson.M{},
			"updatedAt": time.Now(),
		},
	}

	initFilter := bson.M{
		"_id":   participant.Id,
		"pages": nil,
	}

	_, err := models.GetParticipantCollection(*mongoClient).UpdateOne(c, initFilter, initUpdate)
	if err != nil {
		return fmt.Errorf("failed to initialize pages: %v", err)
	}

	update := bson.M{
		"$set": bson.M{
			"pages." + pageKey: models.Page{
				Guesses: []string{},
				Hints:   0,
			},
			"updatedAt": time.Now(),
		},
	}

	// Perform the update
	opts := options.Update().SetUpsert(false)
	_, err = models.GetParticipantCollection(*mongoClient).UpdateOne(
		c,
		bson.M{
			"_id":              participant.Id,
			"pages." + pageKey: bson.M{"$exists": false},
		},
		update,
		opts,
	)

	if err != nil {
		return fmt.Errorf("failed to ensure page %s: %v", pageKey, err)
	}
	return nil
}

func WrongGuess(c *gin.Context, participant models.Participant, pageKey string, answer string) error {
	mongoClient := db.GetDbClient()

	update := bson.M{
		"$push": bson.M{
			"pages." + pageKey + ".guesses": answer,
		},
		"$set": bson.M{
			"score":     participant.Score - 25,
			"updatedAt": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(false)
	_, err := models.GetParticipantCollection(*mongoClient).UpdateOne(c, bson.M{"_id": participant.Id}, update, opts)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CorrectAnswer(c *gin.Context, participant models.Participant, pageKey string, answer string, nextPage string) error {
	mongoClient := db.GetDbClient()

	update := bson.M{
		"$push": bson.M{
			"pages." + pageKey + ".guesses": answer,
		},
		"$set": bson.M{
			"page":      nextPage,
			"score":     participant.Score - 1,
			"updatedAt": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(false)
	_, err := models.GetParticipantCollection(*mongoClient).UpdateOne(c, bson.M{"_id": participant.Id}, update, opts)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
