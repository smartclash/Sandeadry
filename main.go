package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/samber/lo"
	"net/url"
	"strings"
	"time"
)

var c = colly.NewCollector(
	colly.AllowedDomains("sanfoundry.com", "www.sanfoundry.com"),
	colly.Async(true),
	colly.CacheDir("./sanfoundry_cache"))

type Subject struct {
	Name string
	Link string
}

type Topic struct {
	Name string
	Link string
}

type Answer struct {
	Answer      string
	Explanation string
}

type MCQ struct {
	Question string
	Options  []string

	Answer
}

func main() {
	var theURL string

	flag.StringVar(&theURL, "u", "", "Specify the degree url to parse")

	flag.Parse()
	c.SetRequestTimeout(time.Minute * 2)

	subRes := make(chan Subject)
	topicRes := make(chan Topic)
	mcqRes := make(chan []MCQ)

	go degreeParser(theURL, subRes)

	for sub := range subRes {
		go subjectParser(sub, topicRes)
	}

	for topic := range topicRes {
		go topicParser(topic, mcqRes)
	}

	for mcqs := range mcqRes {
		fmt.Println(mcqs)
	}

	c.Wait()
}

func degreeParser(theURL string, subRes chan<- Subject) {
	c.OnHTML("li", func(e *colly.HTMLElement) {
		e.DOM.Find("div.entry-content table tbody tr td li a").Each(func(_ int, selection *goquery.Selection) {
			href, exists := selection.Attr("href")
			if !exists {
				return
			}

			link, err := url.Parse(href)
			if err != nil {
				return
			}

			if strings.EqualFold(link.Host, "rank.sanfoundry.com") {
				return
			}

			subRes <- Subject{
				Name: selection.Text(),
				Link: href,
			}
		})
	})

	c.OnHTML("title", func(element *colly.HTMLElement) {
		fmt.Println("Scrapped degree", element.Text)
	})

	if err := c.Visit(theURL); err != nil {
		fmt.Println("Couldn't visit the website", err.Error())
	}
}

func subjectParser(sub Subject, topicRes chan<- Topic) {
	c.OnHTML("li", func(e *colly.HTMLElement) {
		e.DOM.Find("div.sf-section table tbody tr td li a").Each(func(_ int, selection *goquery.Selection) {
			if href, exists := selection.Attr("href"); exists {
				topicRes <- Topic{
					Name: selection.Text(),
					Link: href,
				}
			}
		})
	})

	c.OnHTML("title", func(element *colly.HTMLElement) {
		fmt.Println("Scrapped subject", element.Text, sub.Link)
	})

	if err := c.Visit(sub.Link); err != nil {
		fmt.Println("Couldn't visit the website", err.Error())
	}
}

func topicParser(topic Topic, mcqRes chan<- []MCQ) {
	fmt.Println("inside topic parser", topic.Name)
	offset := 0
	var mcqs []MCQ
	var answers []Answer
	var skippedQuestions []int

	c.OnHTML("div.entry-content", func(e *colly.HTMLElement) {
		e.ForEach("div.collapseomatic_content", func(_ int, ans *colly.HTMLElement) {
			answerRaw := strings.Split(ans.Text, "\n")
			answer := strings.ReplaceAll(answerRaw[0], "Answer: ", "")
			explanation := strings.ReplaceAll(answerRaw[1], "Explanation: ", "")

			answers = append(answers, Answer{
				Answer:      answer,
				Explanation: explanation,
			})
		})

		e.ForEach("div.entry-content p", func(i int, element *colly.HTMLElement) {
			if !strings.Contains(element.Text, "View Answer") {
				offset++
				return
			}

			theText := strings.Split(strings.ReplaceAll(element.Text, "View Answer", ""), "\n")
			if len(theText) <= 5 {
				skippedQuestions = append(skippedQuestions, i-offset)
				return
			}

			theText = lo.Reject(theText, func(line string, _ int) bool {
				return line == ""
			})

			mcq := MCQ{}
			starters := []string{"a)", "b)", "c)", "d)"}
			for _, line := range theText {
				if lo.Contains(starters, line[0:2]) {
					mcq.Options = append(mcq.Options, line)
				} else {
					mcq.Question = mcq.Question + "\n" + line
				}
			}

			mcqs = append(mcqs, mcq)
		})

		// Construct the MCQ with collected data
		count := 0
		for index, answer := range answers {
			questionNumber := index + 1
			if !lo.Contains(skippedQuestions, questionNumber) {
				mcqs[count].Answer = answer
				count++
			}
		}

		mcqRes <- mcqs
	})

	c.OnHTML("title", func(element *colly.HTMLElement) {
		fmt.Println("Scrapped MCQ", element.Text)
	})

	if err := c.Visit(topic.Link); err != nil {
		fmt.Println("Couldn't visit the website", err.Error())
	}
}
