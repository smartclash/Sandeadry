package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"net/url"
	"strings"
)

type Subject struct {
	Name string
	Link string
}

type DegreeParserResult struct {
	Degree   string
	Subjects []Subject
}

func DegreeParser(link string) (subjects DegreeParserResult, err error) {
	c := colly.NewCollector()

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

			subjects.Subjects = append(subjects.Subjects, Subject{
				Name: selection.Text(),
				Link: href,
			})
		})
	})

	c.OnHTML("title", func(element *colly.HTMLElement) {
		subjects.Degree = strings.ReplaceAll(element.Text, " Questions and Answers - Sanfoundry", "")
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Degree Scrapped", r.Request.URL)
	})

	err = c.Visit(link)

	return
}
