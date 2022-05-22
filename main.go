package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"net/url"
	"strings"
	"time"
)

var c = colly.NewCollector(
	colly.AllowedDomains("sanfoundry.com", "www.sanfoundry.com"),
	colly.AllowURLRevisit(),
	colly.Async(true))

type Subject struct {
	Name string
	Link string
}

type Topic struct {
	Name string
	Link string
}

func main() {
	c.SetRequestTimeout(time.Minute * 2)

	subRes := make(chan Subject, 10)
	topicRes := make(chan Topic, 1000)

	go degreeParser(subRes)

	for sub := range subRes {
		go subjectParser(sub, topicRes)
	}

	for topic := range topicRes {
		fmt.Println("Topic parsed", topic.Name)
	}

	c.Wait()
}

func degreeParser(subRes chan<- Subject) {
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
		fmt.Println(strings.ReplaceAll(element.Text, " Questions and Answers - Sanfoundry", ""))
	})

	if err := c.Visit("https://www.sanfoundry.com/chemical-engineering-questions-answers/"); err != nil {
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

	if err := c.Visit(sub.Link); err != nil {
		fmt.Println("Couldn't visit the website", err.Error())
	}
}
