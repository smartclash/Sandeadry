package main

import "fmt"

func main() {
	_, err := DegreeParser("https://www.sanfoundry.com/computer-science-questions-answers/")
	if err != nil {
		fmt.Println("Couldn't visit degree link")
	}

	_, err = SubjectParser("https://www.sanfoundry.com/1000-data-structure-questions-answers/")
	if err != nil {
		fmt.Println("Couldn't visit subject link")
	}
}
