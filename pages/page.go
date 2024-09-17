package pages

import (
	"fmt"
	"hireme-api/db"
	"hireme-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Page interface {
	// PreAnswerHook(*gin.Context) error
	ProvideAnswer(string, models.Participant, *gin.Context) (bool, error)
	// PostAnswerHook(*gin.Context) error
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

	_, err := models.GetParticipantCollection(*mongoClient).UpdateOne(
		c,
		initFilter,
		initUpdate,
	)
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
	opts := options.Update().SetUpsert(true)
	_, err = models.GetParticipantCollection(*mongoClient).UpdateOne(
		c,
		bson.M{
			"_id": participant.Id,
			"$or": []bson.M{
				{"pages": bson.M{"$exists": false}},
				{"pages": nil},
				{"pages." + pageKey: bson.M{"$exists": false}},
				{"pages." + pageKey: nil},
			},
		},
		update,
		opts,
	)

	if err != nil {
		return fmt.Errorf("failed to ensure page %s: %v", pageKey, err)
	}
	return nil
}
