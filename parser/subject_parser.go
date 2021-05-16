package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Topic struct {
	Name string
	Link string
}

type SubjectParserResult struct {
	Degree string
	Topics []Topic
}

func SubjectParser(degree string, link string) (topics SubjectParserResult, err error) {
	c := colly.NewCollector()
	topics.Degree = degree

	c.OnHTML("li", func(e *colly.HTMLElement) {
		e.DOM.Find("div.sf-section table tbody tr td li a").Each(func(_ int, selection *goquery.Selection) {
			if href, exists := selection.Attr("href"); exists {
				topics.Topics = append(topics.Topics, Topic{
					Name: selection.Text(),
					Link: href,
				})
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Subject Scrapped", r.Request.URL)
	})

	err = c.Visit(link)

	return
}
