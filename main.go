package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/samber/lo"
)

var c = colly.NewCollector(
	colly.AllowedDomains("sanfoundry.com", "www.sanfoundry.com"),
	colly.Async(true),
	colly.CacheDir("./sanfoundry_cache"),
)

type Subject struct {
	Name       string
	Link       string
	DegreeName string
}

type Topic struct {
	Name        string
	Link        string
	SubjectName string
}

type Answer struct {
	Answer      string
	Explanation string
}

type MCQ struct {
	Question string
	Options  []string

	Answer
	TopicName string
}

func main() {
	go func() { log.Fatal(http.ListenAndServe(":4000", nil)) }()

	var theURL string

	flag.StringVar(&theURL, "u", "", "Specify the degree url to parse")

	flag.Parse()
	c.SetRequestTimeout(time.Minute * 2)

	if len(theURL) <= 0 {
		fmt.Println("Please enter a degree link")
		return
	}

	subRes := make(chan Subject, 10)
	topicRes := make(chan Topic, 10)
	mcqRes := make(chan []MCQ, 10)

	go degreeParser(theURL, subRes)
	go subjectWorker(subRes, topicRes)
	go topicWorker(topicRes, mcqRes)

	file, err := os.Create("lmao/mcq.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	tf, err := os.OpenFile("lmao/mcq.json", os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Print("cry")
	}
	defer tf.Close()

	for mcqs := range mcqRes {
		res, err := json.Marshal(mcqs)
		if err != nil {
			fmt.Println("bruh")
		}

		tf.Write(res)
	}

	c.Wait()
}

func subjectWorker(subRes chan Subject, topicRes chan Topic) {
	for sub := range subRes {
		go subjectParser(sub, topicRes)
	}
}

func topicWorker(topicRes chan Topic, mcqRes chan []MCQ) {
	for topic := range topicRes {
		go topicParser(topic, mcqRes)
	}
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
				//DegreeName:
				Name: selection.Text(),
				Link: href,
			}
		})
	})

	c.OnHTML("title", func(element *colly.HTMLElement) {
		//fmt.Println("Scrapped subject", element.Text)
	})

	if err := c.Visit(theURL); err != nil {
		//fmt.Println("Couldn't visit degree site", err.Error())
	}
}

func subjectParser(sub Subject, topicRes chan<- Topic) {
	c.OnHTML("div.sf-section", func(e *colly.HTMLElement) {
		e.ForEach("div.sf-section table tbody tr td li a", func(bruh int, selection *colly.HTMLElement) {
			href := selection.Attr("href")
			topic := Topic{
				Name: selection.Text,
				Link: href,
			}

			topicRes <- topic
		})
	})

	c.OnScraped(func(res *colly.Response) {
		//fmt.Println("Scrapped topic", sub.Link)
	})

	if err := c.Visit(sub.Link); err != nil {
		//fmt.Println("Couldn't visit subject site", err.Error())
	}
}

func topicParser(topic Topic, mcqRes chan<- []MCQ) {
	offset := 0
	var mcqs []MCQ
	var answers []Answer
	var skippedQuestions []int

	c.OnHTML("div.entry-content", func(e *colly.HTMLElement) {
		e.ForEach("div.collapseomatic_content", func(_ int, ans *colly.HTMLElement) {
			if len(answers) <= 0 {
				return
			}

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
				if len(line) <= 1 {
					return
				}

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
		//fmt.Println("Scrapped MCQ", element.Text)
	})

	if err := c.Visit(topic.Link); err != nil {
		//fmt.Println("Couldn't visit topic site", err.Error())
	}
}
