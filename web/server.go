package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/connorwade/dachshund/internal"
)

type PageBuilder struct {
	NumOfLinks   int        `json:"numOfLinks"`
	RatioOfFails float32    `json:"ratioOfFails"`
	Failures     [][]string `json:"failures"`
	Content      [][]string `json:"content"`
}

var templates = template.Must(template.ParseFiles("web/report.html"))

func renderTemplate(w http.ResponseWriter) {
	raw := internal.GetHtmlReport()

	numOfLinks := len(raw.Pages)
	ratioOfFails := float32(len(raw.Failures)) / float32(len(raw.Pages)) * 100

	report := PageBuilder{
		NumOfLinks:   numOfLinks,
		RatioOfFails: ratioOfFails,
		Failures:     raw.Failures,
		Content:      raw.Content,
	}

	err := templates.ExecuteTemplate(w, "report.html", report)
	if err != nil {
		log.Fatal("Couldn't render template ", err)
	}
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w)
}

func Serve() {
	http.HandleFunc("/", reportHandler)
	log.Println("Starting a report on http://localhost:8083/")
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
