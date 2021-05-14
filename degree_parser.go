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

func DegreeParser(link string) (subjects []Subject, err error) {
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

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

			subjects = append(subjects, Subject{
				Name: selection.Text(),
				Link: href,
			})
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err = c.Visit(link)

	return
}
