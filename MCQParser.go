package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

type Answer struct {
	Answer      string
	Explanation string
}

type Option struct {
	Text string
}

type MCQ struct {
	Question string
	options  []Option
	answer   []Answer
}

func MCQParser(link string) {
	var answers []Answer
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("div.entry-content", func(e *colly.HTMLElement) {
		e.ForEach("div.collapseomatic_content", func(_ int, ans *colly.HTMLElement) {
			answerRaw := strings.Split(ans.Text, "\n")
			answers = append(answers, Answer{
				Answer:      strings.ReplaceAll(answerRaw[0], "Answer: ", ""),
				Explanation: strings.ReplaceAll(answerRaw[1], "Explanation: ", ""),
			})
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	if err := c.Visit(link); err != nil {
		fmt.Println("Couldn't visit ", link)
	}

	fmt.Println(answers)
}
