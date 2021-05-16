package parser

import (
	"fmt"
	"github.com/smartclash/Sandeadry/utils"
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

			degreeName := utils.Parameterize(degreeData.Degree)
			subjectName := utils.Parameterize(subject.Name)
			topicName := utils.Parameterize(topic.Name)

			if err := utils.SaveDataToJson(degreeName, subjectName, topicName, mcqData.MCQs); err != nil {
				fmt.Println("Couldn't save mcq data")
			}
		}
	}
}
