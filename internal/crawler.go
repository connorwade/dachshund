package internal

import (
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"gopkg.in/yaml.v3"
)

func setCrawlerVars() (*CrawlerVars, error) {
	crlVars := &CrawlerVars{}
	data, err := os.ReadFile("dachshund.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(data), crlVars)
	if err != nil {
		return nil, err
	}

	return crlVars, err
}

func Crawl() {
	crlVars, err := setCrawlerVars()
	if err != nil {
		log.Fatalln("Cannot set vars: ", err)
	}

	rep := InitReporter()

	strtUrl := "https://" + crlVars.StarterURL
	allowedDomains := crlVars.AllowedDomains

	selToChk := strings.Join(crlVars.Selectors.CheckLinks, ", ")
	selToGet := strings.Join(crlVars.Selectors.GetContent, ", ")

	c := colly.NewCollector(
		colly.MaxDepth(crlVars.Colly.MaxDepth),
		colly.Async(crlVars.Colly.Async),
		colly.AllowedDomains(allowedDomains...),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	c.OnResponse(func(r *colly.Response) {
		reqUrl := r.Request.URL.String()
		h := HTML{
			[]Link{},
			[]Image{},
			[]ContentEl{},
		}
		rep.AddPage(reqUrl, h, r)
	})

	c.OnHTML(selToChk, func(e *colly.HTMLElement) {
		tag := e.Name
		url := e.Request.URL.String()
		el, err := e.DOM.Html()
		if err != nil {
			log.Println("Error with retrieving DOM HTML of element", err)
		}
		var link string
		if tag == "a" {
			link = e.Attr("href")
			rep.AddLinkToPage(url, el, link)
		} else if tag == "img" {
			link = e.Attr("src")
			rep.AddImageToPage(url, el, link)
		}

		rep.AddLinkToPage(url, el, link)
		e.Request.Visit(link)
	})

	c.OnHTML(selToGet, func(e *colly.HTMLElement) {
		content := e.Text
		url := e.Request.URL.String()
		el, err := e.DOM.Html()
		if err != nil {
			log.Println("Error with retrieving DOM HTML of element", err)
		}
		rep.AddContentToPage(url, el, content)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting " + r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode != 0 {
			log.Println("----------------------------------------------------")
			log.Printf("ERROR: %q\nURL: %q\nResponse Recieved: %d\n", err, r.Request.URL, r.StatusCode)
			log.Println("----------------------------------------------------")
			rep.AddPage(r.Request.URL.String(), HTML{nil, nil, nil}, r)
		}
	})

	log.Println("Starting crawl at " + strtUrl)
	err = c.Visit(strtUrl)
	if err != nil {
		log.Fatalln("Failed to start crawl")
	}

	c.Wait()
	err = rep.WriteReport()
	if err != nil {
		log.Fatalln("Report failed: ", err)
	}
}
