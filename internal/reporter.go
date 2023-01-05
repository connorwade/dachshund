package internal

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/gocolly/colly"
)

const (
	reportFilename = "report.json"
)

// type HtmlReport struct {
// 	Failures     []Page
// 	RatioOfFails float32
// 	NumOfLinks   int
// 	NumOfImages  int
// }

// func ReadIntoHtmlReport() *HtmlReport {
// 	report, err := ReadReport()
// 	if err != nil {
// 		log.Fatalln("Error reading report: ", err)
// 	}

// 	rof := float32(len(report.Failures)) / float32(len(report.Pages))

// 	var linkSum int
// 	var imageSum int

// 	for i := range report.Pages {
// 		linkSum += len(report.Pages[i].Html.Links)
// 		imageSum += len(report.Pages[i].Html.Images)
// 	}

// 	h := HtmlReport{
// 		report.Failures,
// 		rof,
// 		linkSum,
// 		imageSum,
// 	}

// 	return &h
// }

type Reporter struct {
	Pages []Page `json:"allPages"`
	// Failures []Page `json:"failures"`
}

/*
	Refactor code so that it is a better representation
	of a page.

	Pages should have their Responses, Url, Origin, and Html attached.

	Html should just be representive of the Html
*/

type Page struct {
	Origin   string          `json:"origin"`
	Url      string          `json:"url"`
	Html     Html            `json:"html"`
	Response *colly.Response `json:"response"`
}

type Html struct {
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

func (r *Reporter) AddPage(org string, url string, html Html, res *colly.Response) *Page {
	p := &Page{org, url, html, res}
	r.Pages = append(r.Pages, *p)
	return p
}

func (r *Reporter) AddLinkToPage(pageURL string, el string, href string) *Link {
	var p *Page
	for i := range r.Pages {
		if r.Pages[i].Url == pageURL {
			p = &r.Pages[i]
			break
		}
	}
	l := &Link{el, pageURL, href}
	p.Html.Links = append(p.Html.Links, *l)
	return l
}

func (r *Reporter) AddImageToPage(pageURL string, el string, src string) *Image {
	var p *Page
	for i := range r.Pages {
		if r.Pages[i].Url == pageURL {
			p = &r.Pages[i]
			break
		}
	}
	l := &Image{el, pageURL, src}
	p.Html.Images = append(p.Html.Images, *l)
	return l
}

func (r *Reporter) AddContentToPage(pageURL string, el string, content string) *ContentEl {
	var p *Page
	for i := range r.Pages {
		if r.Pages[i].Url == pageURL {
			p = &r.Pages[i]
			break
		}
	}
	l := &ContentEl{el, pageURL, content}
	p.Html.Content = append(p.Html.Content, *l)
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

func WriteContentReport() error {
	report, err := ReadReport()
	if err != nil {
		return err
	}

	// report content

	var rows = [][]string{
		{"Url", "Html Element", "Content"},
	}

	for i := range report.Pages {
		page := report.Pages[i]
		if page.Response.StatusCode < 200 && page.Response.StatusCode > 299 && len(page.Html.Content) == 0 {
			continue
		}

		for j := range page.Html.Content {
			a := []string{page.Url, page.Html.Content[j].Element, page.Html.Content[j].Text}
			rows = append(rows, a)
		}
	}

	writeToCSV(rows, "csv-content-report.csv")

	return nil
}

func WriteFailureReport() error {
	// report, err := ReadReport()
	// if err != nil {
	// 	return err
	// }

	// reconcile failures

	var brokenElements = [][]string{
		{"Origin", "Url", "Html Element", "Response Code"},
	}

	// for i := range report.Pages {
	// 	page := report.Pages[i]
	// 	if page.Response.StatusCode < 200 && page.Response.StatusCode > 299 {
	// 		continue
	// 	}
	// 	//Check link failures
	// 	for j := range page.Html.Links {
	// 		link := report.Pages[i].Html.Links[j]
	// 		for k := range report.Failures {
	// 			fail := report.Failures[k]
	// 			if link.HREF == fail.Url {
	// 				a := []string{page.Url, fail.Url, link.Element, fmt.Sprintf("%x", fail.Response.StatusCode)}
	// 				brokenElements = append(brokenElements, a)
	// 			}
	// 		}
	// 	}
	// 	//Check image failures
	// 	for j := range page.Html.Images {
	// 		link := report.Pages[i].Html.Images[j]
	// 		for k := range report.Failures {
	// 			fail := report.Failures[k]
	// 			if link.Src == fail.Url {
	// 				a := []string{page.Url, fail.Url, link.Element, fmt.Sprintf("%x", fail.Response.StatusCode)}
	// 				brokenElements = append(brokenElements, a)
	// 			}
	// 		}
	// 	}
	// }

	writeToCSV(brokenElements, "csv-error-report.csv")

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
