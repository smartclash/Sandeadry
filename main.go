package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"github.com/samber/lo"
	"net/url"
	"strings"
	"sync"
)

type MCQ struct {
	Question    string
	OptionOne   string
	OptionTwo   string
	OptionThree string
	OptionFour  string
	Answer      string
	Explanation string
}

type MCQWrapper struct {
	Subject string
	Topic   string
	MCQs    []MCQ
}

var wg sync.WaitGroup

func main() {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{
			"https://www.sanfoundry.com/programming-questions-answers/",
			"https://www.sanfoundry.com/computer-science-questions-answers/",
			"https://www.sanfoundry.com/information-technology-questions-answers/",
			"https://www.sanfoundry.com/information-science-questions-answers/",
			"https://www.sanfoundry.com/electronics-communication-engineering-questions-answers/",
			"https://www.sanfoundry.com/electrical-electronics-engineering-questions-answers/",
			"https://www.sanfoundry.com/electrical-engineering-questions-answers/",
			"https://www.sanfoundry.com/civil-engineering-questions-answers/",
			"https://www.sanfoundry.com/mechanical-engineering-questions-answers/",
			"https://www.sanfoundry.com/chemical-engineering-questions-answers/",
			"https://www.sanfoundry.com/metallurgical-engineering-questions-answers/",
			"https://www.sanfoundry.com/mining-engineering-questions-answers/",
			"https://www.sanfoundry.com/instrumentation-engineering-questions-answers/",
			"https://www.sanfoundry.com/aerospace-engineering-questions-answers/",
			"https://www.sanfoundry.com/aeronautical-engineering-questions-answers/",
			"https://www.sanfoundry.com/biotechnology-questions-answers/",
			"https://www.sanfoundry.com/agricultural-engineering-questions-answers/",
			"https://www.sanfoundry.com/marine-engineering-questions-answers/",
			"https://www.sanfoundry.com/mechatronics-engineering-questions-answers/",
		},
		ParseFunc:          subjectParse,
		Exporters:          []export.Exporter{&export.JSON{}},
		ConcurrentRequests: 100,
	}).Start()

	wg.Wait()
}

func stringCleaner(text string) string {
	replacers := []string{
		" MCQs with Answers - Sanfoundry",
		" MCQs - Sanfoundry",
		"50000+ ",
		"10000+ ",
		" Questions and Answers - Sanfoundry",
		" MCQ Questions and Answers - Sanfoundry",
		" MCQ Questions with Answers - Sanfoundry",
		" MCQ (Multiple Choice Questions) - Sanfoundry",
		" MCQ with Answers - Sanfoundry",
	}

	for _, replacer := range replacers {
		text = strings.ReplaceAll(text, replacer, "")
	}

	return text
}

func subjectParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("table tbody tr td li a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		link, _ := url.Parse(href)
		if link.Host == "rank.sanfoundry.com" {
			return
		}

		rawTitle := r.HTMLDoc.Find("title").Text()
		title := stringCleaner(rawTitle)

		req, _ := client.NewRequest("GET", href, nil)
		req.Meta["subject"] = title

		wg.Add(1)
		go g.Do(req, topicParse)
	})
}

func topicParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("table tbody tr td li a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		subject := fmt.Sprint(r.Request.Meta["subject"])

		req, _ := client.NewRequest("GET", href, nil)
		req.Meta["subject"] = subject
		req.Meta["topic"] = s.Text()

		wg.Add(1)
		go g.Do(req, mcqParse)
	})

	wg.Done()
}

func mcqParse(g *geziyor.Geziyor, r *client.Response) {
	var MCQs []MCQ
	dom := r.HTMLDoc.Find("div.entry-content")

	dom.Find("p").Each(func(i int, s *goquery.Selection) {
		if !strings.Contains(s.Text(), "View Answer") {
			return
		}

		theText := strings.Split(strings.ReplaceAll(s.Text(), "View Answer", ""), "\n")
		if len(theText) <= 5 {
			return
		}

		theText = lo.Reject(theText, func(line string, _ int) bool {
			return line == ""
		})

		rawResult := s.Next().Text()
		if !strings.Contains(rawResult, "Answer") {
			return
		}

		splitResults := strings.Split(rawResult, "\n")
		answer := strings.ReplaceAll(splitResults[0], "Answer: ", "")
		explanation := strings.ReplaceAll(splitResults[1], "Explanation: ", "")

		var options []string
		var question string
		starters := []string{"a)", "b)", "c)", "d)"}
		for _, line := range theText {
			if len(line) <= 1 {
				return
			}

			if lo.Contains(starters, line[0:2]) {
				options = append(options, line)
			} else {
				question += "\n" + line
			}
		}

		MCQs = append(MCQs, MCQ{
			Question:    question,
			OptionOne:   options[0],
			OptionTwo:   options[1],
			OptionThree: options[2],
			OptionFour:  options[3],
			Answer:      answer,
			Explanation: explanation,
		})
	})

	g.Exports <- MCQWrapper{
		Subject: fmt.Sprint(r.Request.Meta["subject"]),
		Topic:   fmt.Sprint(r.Request.Meta["topic"]),
		MCQs:    MCQs,
	}
	wg.Done()
}
