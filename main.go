package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func saveDataToJson(degree string, subject string, topic string, mcqs []MCQ) (err error) {
	thePath := filepath.Join("datas", degree, subject)
	err = os.MkdirAll(thePath, os.ModePerm)

	res, err := json.Marshal(mcqs)
	if err != nil {
		return
	}

	theFilePath := filepath.Join(thePath, topic+".json")
	err = ioutil.WriteFile(theFilePath, res, os.ModePerm)
	if err != nil {
		return
	}

	return
}

func main() {
	degreeData, err := DegreeParser("https://www.sanfoundry.com/computer-science-questions-answers/")
	if err != nil {
		fmt.Println("Couldn't visit degree link")
		return
	}

	for _, subject := range degreeData.Subjects {
		subjectData, err := SubjectParser(subject.Name, subject.Link)
		if err != nil {
			fmt.Println("Couldn't visit subject link")
			return
		}

		for _, topic := range subjectData.Topics {
			mcqData := MCQParser(topic.Name, topic.Link)

			degreeName := Parameterize(degreeData.Degree)
			subjectName := Parameterize(subject.Name)
			topicName := Parameterize(topic.Name)

			if err := saveDataToJson(degreeName, subjectName, topicName, mcqData.MCQs); err != nil {
				fmt.Println("Couldn't save mcq data")
			}
		}
	}
}
