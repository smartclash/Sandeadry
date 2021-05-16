package parser

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Parameterize(s string) (newString string) {
	newString = strings.ReplaceAll(s, " ", "_")
	newString = strings.ReplaceAll(newString, "/", "_OR_")

	return
}

func SaveDataToJson(degree string, subject string, topic string, mcqs []MCQ) (err error) {
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
