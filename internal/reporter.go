package internal

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
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
	HTML     HTML            `json:"html"`
	Response *colly.Response `json:"response"`
}

type HTML struct {
	Links   []Link      `json:"links"`
	Images  []Image     `json:"images"`
	Content []ContentEl `json:"content"`
}

type Image struct {
	Element string `json:"element"`
	Origin  string `json:"origin"`
	Src     string `json:"src"`
}

type ContentEl struct {
	Element string `json:"element"`
	Origin  string `json:"origin"`
	Text    string `json:"text"`
}

type Link struct {
	Element string `json:"element"`
	Origin  string `json:"origin"`
	HREF    string `json:"href"`
}

func InitReporter() *Reporter {
	return &Reporter{}
}

func (r *Reporter) AddPage(url string, html HTML, res *colly.Response) *Page {
	p := &Page{url, html, res}
	r.Pages = append(r.Pages, *p)
	return p
}

func (r *Reporter) AddFailures(url string, html HTML, res *colly.Response) *Page {
	p := &Page{url, html, res}
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
	p.HTML.Links = append(p.HTML.Links, *l)
	return l
}

func (r *Reporter) AddImageToPage(pageURL string, el string, src string) *Image {
	var p *Page
	for i := range r.Pages {
		if r.Pages[i].URL == pageURL {
			p = &r.Pages[i]
			break
		}
	}
	l := &Image{el, pageURL, src}
	p.HTML.Images = append(p.HTML.Images, *l)
	return l
}

func (r *Reporter) AddContentToPage(pageURL string, el string, content string) *ContentEl {
	var p *Page
	for i := range r.Pages {
		if r.Pages[i].URL == pageURL {
			p = &r.Pages[i]
			break
		}
	}
	l := &ContentEl{el, pageURL, content}
	p.HTML.Content = append(p.HTML.Content, *l)
	return l
}

func (r *Reporter) WriteReport() error {
	file, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(reportFilename, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadReport() (*Reporter, error) {
	report, err := os.ReadFile(reportFilename)
	if err != nil {
		return nil, err
	}

	var jsonRep Reporter

	err = json.Unmarshal(report, &jsonRep)
	if err != nil {
		return nil, err
	}

	return &jsonRep, nil

}

func WriteContentReport() error {
	report, err := ReadReport()
	if err != nil {
		return err
	}

	// report content

	var rows = [][]string{
		{"URL", "HTML Element", "Content"},
	}

	for i := range report.Pages {
		page := report.Pages[i]
		if page.Response.StatusCode < 200 && page.Response.StatusCode > 299 && len(page.HTML.Content) == 0 {
			continue
		}

		for j := range page.HTML.Content {
			a := []string{page.URL, page.HTML.Content[j].Element, page.HTML.Content[j].Text}
			rows = append(rows, a)
		}
	}

	writeToCSV(rows, "csv-content-report.csv")

	return nil
}

func WriteFailureReport() error {
	report, err := ReadReport()
	if err != nil {
		return err
	}

	// reconcile failures

	var brokenElements = [][]string{
		{"Origin", "URL", "HTML Element", "Response Code"},
	}

	for i := range report.Pages {
		page := report.Pages[i]
		if page.Response.StatusCode < 200 && page.Response.StatusCode > 299 {
			continue
		}
		//Check link failures
		for j := range page.HTML.Links {
			link := report.Pages[i].HTML.Links[j]
			for k := range report.Failures {
				fail := report.Failures[k]
				if link.HREF == fail.URL {
					a := []string{page.URL, fail.URL, link.Element, fmt.Sprintf("%x", fail.Response.StatusCode)}
					brokenElements = append(brokenElements, a)
				}
			}
		}
		//Check image failures
		for j := range page.HTML.Images {
			link := report.Pages[i].HTML.Images[j]
			for k := range report.Failures {
				fail := report.Failures[k]
				if link.Src == fail.URL {
					a := []string{page.URL, fail.URL, link.Element, fmt.Sprintf("%x", fail.Response.StatusCode)}
					brokenElements = append(brokenElements, a)
				}
			}
		}
	}

	writeToCSV(brokenElements, "csv-error-report.csv")

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
