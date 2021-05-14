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

type MCQ struct {
	Question string
	Options  []string
	Answer
}

type MCQParserResult struct {
	Topic string
	MCQs  []MCQ
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.EqualFold(v, str) {
			return true
		}
	}

	return false
}

func optionsParser(rawText []string) (options []string, exists bool) {
	optionStarters := []string{"a)", "b)", "c)", "d)"}

	for _, line := range rawText {
		for _, starter := range optionStarters {
			if !strings.HasPrefix(line, starter) {
				continue
			}

			options = append(options, strings.ReplaceAll(line, starter+" ", ""))
		}
	}

	if len(options) > 0 {
		exists = true
	}

	return
}

func mcqBuilder(answers []Answer, options [][]string, questions []string) (mcqs []MCQ) {
	for i, question := range questions {
		mcqs = append(mcqs, MCQ{
			Question: question,
			Options:  options[i],
			Answer:   answers[i],
		})
	}

	return
}

func MCQParser(topic string, link string) (mcqs MCQParserResult, err error) {
	var answers []Answer
	var options [][]string
	var questions []string

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

		e.ForEach("div.entry-content p", func(i int, element *colly.HTMLElement) {
			if i == 0 {
				return
			}

			questionsRaw := strings.Split(element.Text, "\n")
			if val, exists := optionsParser(questionsRaw); exists {
				options = append(options, val)
			}

			questionNotFound := true
			for questionNotFound {
				if contains(questionsRaw, "View Answer") {
					if !strings.HasPrefix(questionsRaw[0], "a)") {
						questions = append(questions, questionsRaw[0])
						questionNotFound = false
					}
				}

				if len(questionsRaw) == 1 {
					questions = append(questions, questionsRaw[0])
					questionNotFound = false
				}

				questionNotFound = false
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err = c.Visit(link)

	questions = questions[:len(questions)-2]
	mcqs = MCQParserResult{
		Topic: topic,
		MCQs:  mcqBuilder(answers, options, questions),
	}

	return
}
