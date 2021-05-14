package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Topic struct {
	Name string
	Link string
}

func SubjectParser(link string) (topics []Topic, err error) {
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("li", func(e *colly.HTMLElement) {
		e.DOM.Find("div.sf-section table tbody tr td li a").Each(func(_ int, selection *goquery.Selection) {
			if href, exists := selection.Attr("href"); exists {
				topics = append(topics, Topic{
					Name: selection.Text(),
					Link: href,
				})
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err = c.Visit(link)

	return
}
