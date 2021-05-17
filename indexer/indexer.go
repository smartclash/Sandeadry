package indexer

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/smartclash/Sandeadry/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type MeiliMCQ struct {
	ID          uint
	Question    string
	Explanation string
	Answer      string
	OptionOne   string
	OptionTwo   string
	OptionThree string
	OptionFour  string
	Topic       string
	Subject     string
	Degrees     []string
}

func Init(database string) (err error) {
	meili := meilisearch.NewClient(meilisearch.Config{
		Host:   os.Getenv("MEILISEARCH_HOST"),
		APIKey: os.Getenv("MEILISEARCH_KEY"),
	})
	db, err := gorm.Open(sqlite.Open(database+".sqlite"), &gorm.Config{})
	if err != nil {
		return
	}

	err = IndexDocuments(db, meili)
	if err != nil {
		return err
	}

	_, err = meili.Settings("mcqs").UpdateAttributesForFaceting([]string{"Degrees"})
	if err != nil {
		return err
	}

	_, err = meili.Settings("mcqs").UpdateSearchableAttributes([]string{
		"Question", "Explanation", "OptionOne", "OptionTwo",
		"OptionThree", "OptionFour", "Answer", "Topic",
	})
	if err != nil {
		return err
	}

	return
}

func IndexDocuments(db *gorm.DB, meili meilisearch.ClientInterface) (err error) {
	var MCQs []storage.MCQ

	db.Model(&MCQs).
		Preload("Topic").Preload("Topic.Subject").
		Preload("Topic.Subject.Degrees").
		FindInBatches(&MCQs, 5000, func(tx *gorm.DB, batch int) error {
			var meiliMCQs []MeiliMCQ

			for _, mcq := range MCQs {
				var theDegrees []string
				for _, degree := range mcq.Topic.Subject.Degrees {
					theDegree := degree.Name
					theDegrees = append(theDegrees, theDegree)
				}

				meiliMCQs = append(meiliMCQs, MeiliMCQ{
					ID:          mcq.ID,
					Question:    mcq.Question,
					Explanation: mcq.Explanation,
					Answer:      mcq.Answer,
					OptionOne:   mcq.OptionOne,
					OptionTwo:   mcq.OptionTwo,
					OptionThree: mcq.OptionThree,
					OptionFour:  mcq.OptionFour,
					Topic:       mcq.Topic.Name,
					Subject:     mcq.Topic.Subject.Name,
					Degrees:     theDegrees,
				})
			}

			if _, err = meili.Documents("mcqs").AddOrReplace(meiliMCQs); err != nil {
				return err
			}

			return nil
		})

	return
}
