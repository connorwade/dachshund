package internal

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gocolly/colly"
)

const (
	reportFilename = "report.json"
)

type Reporter struct {
	Pages    []Page `json:"allPages"`
	Failures []Page `json:"failures"`
}

type Page struct {
	URL      string          `json:"url"`
	Links    []Link          `json:"links"`
	Response *colly.Response `json:"response"`
}

type Link struct {
	Element string `json:"element"`
	Origin  string `json:"origin"`
	HREF    string `json:"href"`
}

func InitReporter() *Reporter {
	return &Reporter{}
}

func (r *Reporter) AddPage(url string, ls []Link, res *colly.Response) *Page {
	p := &Page{url, ls, res}
	r.Pages = append(r.Pages, *p)
	return p
}

func (r *Reporter) AddFailures(url string, ls []Link, res *colly.Response) *Page {
	p := &Page{url, ls, res}
	r.Pages = append(r.Failures, *p)
	return p
}

func (r *Reporter) AddLinkToPage(pageURL string, el string, href string) *Link {
	var p *Page
	for i := range r.Pages {
		if r.Pages[i].URL == pageURL {
			p = &r.Pages[i]
			break
		}
	}
	l := &Link{el, pageURL, href}
	p.Links = append(p.Links, *l)
	return l
}

func (r *Reporter) WriteReport() error {
	file, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(reportFilename, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func WriteCSVReport() error {
	//read report
	content, err := os.ReadFile(reportFilename)
	if err != nil {
		return err
	}

	var jsonRep Reporter

	err = json.Unmarshal(content, &jsonRep)
	if err != nil {
		return err
	}

	//reconcile failures

	var brokenLinks = [][]string{
		{"Origin", "URL", "HTML Element", "Response Code"},
	}

	for i := range jsonRep.Pages {
		page := jsonRep.Pages[i]
		if page.Response.StatusCode < 200 && page.Response.StatusCode > 299 {
			continue
		}
		for j := range page.Links {
			link := jsonRep.Pages[i].Links[j]
			for k := range jsonRep.Failures {
				fail := jsonRep.Failures[k]
				if link.HREF == fail.URL {
					a := []string{page.URL, fail.URL, link.Element, fmt.Sprintf("%x", fail.Response.StatusCode)}
					brokenLinks = append(brokenLinks, a)
				}
			}
		}
	}

	writeToCSV(brokenLinks, "csv-report.csv")

	return nil
}

func writeToCSV(data [][]string, fn string) {
	path := "results"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			log.Fatalln("failed to create results dir", err)
		}
	}
	f, err := os.Create("results/" + fn)
	if err != nil {
		log.Fatalln("failed to open file: ", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.WriteAll(data)
	if err != nil {
		log.Fatalln("failed to write file: ", err)
	}
}
