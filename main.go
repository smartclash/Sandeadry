package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"net/url"
	"strings"
)

func main() {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs:          []string{"https://www.sanfoundry.com/computer-science-questions-answers/"},
		ParseFunc:          subjectParse,
		Exporters:          []export.Exporter{&export.JSON{}},
		ConcurrentRequests: 100,
	}).Start()
}

func stringCleaner(text string) string {
	replacers := []string{
		" MCQs with Answers - Sanfoundry",
		" MCQs - Sanfoundry",
		"50000+ ",
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

		g.Do(req, topicParse)
	})
}

func topicParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("table tbody tr td li a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		subject := fmt.Sprint(r.Request.Meta["subject"])

		req, _ := client.NewRequest("GET", href, nil)
		req.Meta["subject"] = subject
		req.Meta["topic"] = s.Text()

		g.Do(req, mcqParse)
	})
}

type MCQ struct {
	Subject string
	Topic   string
	Result  []interface{}
}

func mcqParse(g *geziyor.Geziyor, r *client.Response) {
	var answers []interface{}
	dom := r.HTMLDoc.Find("div.entry-content")

	dom.Find("div.collapseomatic_content").Each(func(i int, s *goquery.Selection) {
		answerRaw := strings.Split(s.Text(), "\n")
		answer := strings.ReplaceAll(answerRaw[0], "Answer: ", "")
		explanation := strings.ReplaceAll(answerRaw[1], "Explanation: ", "")

		finalAnswer := map[string]string{
			"Answer":      answer,
			"Explanation": explanation,
		}

		answers = append(answers, finalAnswer)
	})

	g.Exports <- MCQ{
		Subject: fmt.Sprint(r.Request.Meta["subject"]),
		Topic:   fmt.Sprint(r.Request.Meta["topic"]),
		Result:  answers,
	}
}
