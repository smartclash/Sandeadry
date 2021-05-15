package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
	link := flag.String("l", "", "Link to the degree you want to parse MCQs")
	flag.Parse()

	if *link == "" {
		flag.PrintDefaults()
		return
	}

	parse, err := url.Parse(*link)
	if err != nil {
		fmt.Println("Please enter a proper link")
		flag.PrintDefaults()
		return
	}

	if !strings.EqualFold(parse.Hostname(), "www.sanfoundry.com") {
		fmt.Println("Enter only sanfoundry links")
		flag.PrintDefaults()
		return
	}

	degreeData, err := DegreeParser(*link)
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
