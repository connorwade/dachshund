package internal

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"gopkg.in/yaml.v3"
)

var mutex = &sync.Mutex{}

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

	origins := make(map[int]string)
	originEls := make(map[int]DomNode)
	uid := 1

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

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting " + r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		var org string
		var orgEl DomNode

		if r.Request.ID < 2 {
			org = "starting url"
			orgEl = DomNode{
				"starting url",
				"starting url",
				"starting url",
				"starting url",
				"starting url",
			}
		} else {
			org = origins[int(r.Request.ID)]
			orgEl = originEls[int(r.Request.ID)]
		}

		reqUrl := r.Request.URL.String()
		h := Html{
			[]Link{},
			[]ContentEl{},
		}
		rep.AddPage(org, orgEl, reqUrl, h, r)
	})

	c.OnHTML(selToChk, func(e *colly.HTMLElement) {
		tag := e.Name
		url := e.Request.URL.String()
		text := e.DOM.Text()
		classes := e.DOM.AttrOr("class", "")
		id := e.DOM.AttrOr("id", "")
		inner, err := e.DOM.Html()

		if err != nil {
			log.Println("Could not get inner HTML of element")
		}

		el := &DomNode{
			tag, text, classes, id, inner,
		}

		var link string
		if tag == "a" {
			link = e.Attr("href")
		} else if tag == "img" {
			link = e.Attr("src")
		}

		mutex.Lock()
		origins[uid] = url
		originEls[uid] = *el
		uid++
		mutex.Unlock()

		rep.AddLinkToPage(url, *el, link)
		e.Request.Visit(link)
	})

	c.OnHTML(selToGet, func(e *colly.HTMLElement) {
		tag := e.Name
		url := e.Request.URL.String()
		text := e.DOM.Text()
		classes := e.DOM.AttrOr("class", "")
		id := e.DOM.AttrOr("id", "")
		inner, err := e.DOM.Html()

		if err != nil {
			log.Println("Could not get inner HTML of element")
		}

		el := &DomNode{
			tag, text, classes, id, inner,
		}

		if err != nil {
			log.Println("Error with retrieving DOM HTML of element", err)
		}
		rep.AddContentToPage(url, *el)
	})

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode != 0 {
			var org string
			var orgEl DomNode
			if r.Request.ID < 2 {
				org = "starting url"
				orgEl = DomNode{
					"starting url",
					"starting url",
					"starting url",
					"starting url",
					"starting url",
				}
			} else {
				org = origins[int(r.Request.ID)]
				orgEl = originEls[int(r.Request.ID)]
			}

			log.Println("----------------------------------------------------")
			log.Printf("ERROR: %q\nURL: %q\nResponse Recieved: %d\n", err, r.Request.URL, r.StatusCode)
			log.Println("----------------------------------------------------")

			rep.AddPage(org, orgEl, r.Request.URL.String(), Html{nil, nil}, r)
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
