package parser

import (
	"fmt"
)

func Parser(link string) {
	degreeData, err := DegreeParser(link)
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

			if err := SaveDataToJson(degreeName, subjectName, topicName, mcqData.MCQs); err != nil {
				fmt.Println("Couldn't save mcq data")
			}
		}
	}
}
