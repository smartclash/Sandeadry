package storage

import (
	"encoding/json"
	"fmt"
	"github.com/smartclash/Sandeadry/parser"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)

func Init(database string) (err error) {
	db, err := gorm.Open(sqlite.Open(database+".sqlite"), &gorm.Config{})

	if err != nil {
		return
	}

	err = db.AutoMigrate(&Degree{}, &Subject{}, &Topic{}, &MCQ{})
	if err != nil {
		return
	}

	theDegree, err := CreateOrInsertDegree(db)
	if err != nil {
		return
	}

	if err = CreateOrInsertSubjects(theDegree, db); err != nil {
		return
	}

	return
}

func CreateOrInsertDegree(db *gorm.DB) (theDegree Degree, err error) {
	folders, err := ioutil.ReadDir("datas")
	if err != nil {
		return
	}

	for _, f := range folders {
		var count int64
		degree := Degree{Name: parser.Humanize(f.Name())}
		db.Model(&degree).Where(&degree).Count(&count)

		fmt.Println("Found degree", degree.Name)
		if count <= int64(0) {
			db.Create(&degree)
			fmt.Println("Created degree", degree.Name)
			if err = CreateOrInsertSubjects(degree, db); err != nil {
				return
			}
		}
	}

	return
}

func CreateOrInsertSubjects(degree Degree, db *gorm.DB) error {
	var subjects []Subject
	folders, err := ioutil.ReadDir("datas/" + parser.Parameterize(degree.Name) + "/")
	if err != nil {
		return err
	}

	err = db.Model(&degree).
		Association("Subjects").
		Clear()
	if err != nil {
		return err
	}

	for _, f := range folders {
		subject := Subject{Name: parser.Humanize(f.Name())}
		var count int64

		fmt.Println("Found subject", subject.Name)
		db.Where(subject).First(&subject).Count(&count)

		if count <= int64(0) {
			db.Create(&subject)
			fmt.Println("Created subject", subject.Name)

			if err = CreateTopics(degree, subject, db); err != nil {
				return err
			}

			subjects = append(subjects, subject)
		}

		err := db.Model(&degree).
			Association("Subjects").
			Append(&subject)
		if err != nil {
			return err
		}
	}

	return err
}

func CreateTopics(degree Degree, subject Subject, db *gorm.DB) (err error) {
	var topics []Topic
	path := "datas/" + parser.Parameterize(degree.Name) + "/" + parser.Parameterize(subject.Name)
	folders, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range folders {
		topicName := strings.ReplaceAll(f.Name(), ".json", "")
		fmt.Println("Found topic", topicName)
		topic := Topic{Name: parser.Humanize(topicName)}
		db.Create(&topic)

		err := db.Model(&subject).
			Association("Topics").
			Append(&topic)
		if err != nil {
			return err
		}

		topics = append(topics, topic)
		fmt.Println("Created topic", topic.Name)

		if err = CreateMCQs(degree, subject, topic, db); err != nil {
			return err
		}
	}

	return
}

func CreateMCQs(degree Degree, subject Subject, topic Topic, db *gorm.DB) (err error) {
	var fileMCQs []parser.MCQ
	path := "datas/" + parser.Parameterize(degree.Name) +
		"/" + parser.Parameterize(subject.Name) + "/" +
		parser.Parameterize(topic.Name) + ".json"

	res, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	rawFile := string(res)
	if rawFile == "null" {
		return nil
	}

	err = json.Unmarshal(res, &fileMCQs)
	if err != nil {
		return err
	}

	for _, fMCQ := range fileMCQs {
		mcq := MCQ{
			Question:    fMCQ.Question,
			Answer:      fMCQ.Answer.Answer,
			Explanation: fMCQ.Answer.Explanation,
			OptionOne:   getOption(fMCQ.Options, 0),
			OptionTwo:   getOption(fMCQ.Options, 1),
			OptionThree: getOption(fMCQ.Options, 2),
			OptionFour:  getOption(fMCQ.Options, 3),
		}

		db.Create(&mcq)
		err := db.Model(&topic).Association("MCQs").Append(&mcq)
		if err != nil {
			return err
		}
	}

	return
}

func getOption(s []string, index int) (option string) {
	if len(s) <= index {
		return ""
	}

	option = s[index]

	return
}
