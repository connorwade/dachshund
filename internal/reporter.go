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
	OriginEl DomNode         `json:"originEl"`
	Url      string          `json:"url"`
	Html     Html            `json:"html"`
	Response *colly.Response `json:"response"`
}

type Html struct {
	Links   []Link      `json:"links"`
	Content []ContentEl `json:"content"`
}

type ContentEl struct {
	Element DomNode `json:"element"`
}

type Link struct {
	Element DomNode `json:"element"`
	Href    string  `json:"href"`
}

type DomNode struct {
	Tag     string `json:"tag"`
	Text    string `json:"text"`
	Classes string `json:"classes"`
	ID      string `json:"id"`
	Inner   string `json:"innerHtml"`
}

func InitReporter() *Reporter {
	return &Reporter{}
}

func (r *Reporter) AddPage(org string, orgEl DomNode, url string, html Html, res *colly.Response) *Page {
	p := &Page{org, orgEl, url, html, res}
	r.Pages = append(r.Pages, *p)
	return p
}

func (r *Reporter) AddLinkToPage(pageURL string, el DomNode, link string) *Link {
	var p *Page
	for i := range r.Pages {
		if r.Pages[i].Url == pageURL {
			p = &r.Pages[i]
			break
		}
	}
	l := &Link{el, link}
	p.Html.Links = append(p.Html.Links, *l)
	return l
}

func (r *Reporter) AddContentToPage(pageURL string, el DomNode) *ContentEl {
	var p *Page
	for i := range r.Pages {
		if r.Pages[i].Url == pageURL {
			p = &r.Pages[i]
			break
		}
	}
	l := &ContentEl{el}
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

func WriteReports(failRep bool, contRep bool) error {
	report, err := ReadReport()
	if err != nil {
		return err
	}

	var brokenElements = [][]string{
		{"Origin", "Url", "Response Code", "Origin Html Element", "Element Text", "Element Classes", "Element Id", "Inner Html"},
	}

	var contRows = [][]string{
		{"Url", "Html Element", "Element Text", "Element Classes", "Element Id", "Inner Html"},
	}

	for _, page := range report.Pages {

		if contRep {
			for _, content := range page.Html.Content {
				element := content.Element
				a := []string{
					page.Url,
					element.Tag,
					element.Text,
					element.Classes,
					element.ID,
					element.Inner,
				}
				contRows = append(contRows, a)
			}
		}

		if failRep {
			// skip pages with no failures or redirects
			if page.Response.StatusCode >= 200 && page.Response.StatusCode < 300 {
				continue
			}

			b := []string{
				page.OriginEl.Tag,
				page.OriginEl.Text,
				page.OriginEl.Classes,
				page.OriginEl.ID,
				page.OriginEl.Inner,
			}
			a := []string{page.Origin, page.Url, fmt.Sprintf("%d", page.Response.StatusCode)}
			a = append(a, b...)
			brokenElements = append(brokenElements, a)
		}

	}

	if failRep {
		writeToCSV(brokenElements, "csv-failure-report.csv")
	}
	if contRep {
		writeToCSV(contRows, "csv-content-report.csv")
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
