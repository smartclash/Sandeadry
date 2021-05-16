package parser

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
	"time"
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

func QuestionsParser(rawText []string) (questions []string) {
	for _, line := range rawText {
		if strings.HasPrefix(line, "a) ") ||
			strings.HasPrefix(line, "b) ") ||
			strings.HasPrefix(line, "c) ") ||
			strings.HasPrefix(line, "d) ") {
		} else {
			questions = append(questions, line)
		}
	}

	return
}

func mcqBuilder(answers []Answer, options [][]string, questions []string) (mcqs []MCQ) {
	for i, answer := range answers {
		mcqs = append(mcqs, MCQ{
			Question: questions[i],
			Options:  options[i],
			Answer:   answer,
		})
	}

	return
}

func MCQParser(topic string, link string) MCQParserResult {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in MCQParser. Error is: %v \n", r)
		}
	}()

	var answers []Answer
	var options [][]string
	var questions []string
	var rawText [][]string

	c := colly.NewCollector(
		colly.CacheDir("./sanfoundry_cache"),
	)

	if err := c.Limit(&colly.LimitRule{
		DomainGlob:  "www.sanfoundry.com/*",
		RandomDelay: 1 * time.Second,
	}); err != nil {
		return MCQParserResult{}
	}

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
			rawText = append(rawText, questionsRaw)
		})
	})

	c.OnScraped(func(r *colly.Response) {
		skip := 0
		var purifiedText [][]string

		for i, text := range rawText {
			if skip != 0 && i == skip {
				continue
			}

			if contains(text, "View Answer") {
				purifiedText = append(purifiedText, text)
				continue
			}

			if len(rawText) <= i+1 {
				continue
			}

			theText := append(text, rawText[i+1]...)
			purifiedText = append(purifiedText, theText)
			skip = i + 1
		}

		var rawQuestions []string
		purifiedText = purifiedText[:len(purifiedText)-1]
		for _, text := range purifiedText {
			if theOptions, exists := optionsParser(text); exists {
				options = append(options, theOptions)
			}

			rawQuestions = append(rawQuestions, QuestionsParser(text)...)
		}

		constructQuestion := ""
		for _, theQ := range rawQuestions {
			if strings.EqualFold(theQ, "View Answer") {
				questions = append(questions, strings.TrimSpace(constructQuestion))
				constructQuestion = ""
				continue
			}

			constructQuestion = constructQuestion + " " + theQ
		}

		fmt.Println("Scrapped", r.Request.URL)
	})

	if err := c.Visit(link); err != nil {
		return MCQParserResult{}
	}

	return MCQParserResult{
		topic,
		mcqBuilder(answers, options, questions),
	}
}
