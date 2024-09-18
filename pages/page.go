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
			"score":     participant.Score - 10,
			"updatedAt": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
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
	opts := options.Update().SetUpsert(true)
	_, err := models.GetParticipantCollection(*mongoClient).UpdateOne(c, bson.M{"_id": participant.Id}, update, opts)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
