package web

// import (
// 	"html/template"
// 	"log"
// 	"net/http"

// 	"github.com/connorwade/dachshund/internal"
// )

// var templates = template.Must(template.ParseFiles("web/report.html"))

// func renderTemplate(w http.ResponseWriter) {
// 	report := internal.ReadIntoHtmlReport()

// 	err := templates.ExecuteTemplate(w, "report.html", report)
// 	if err != nil {
// 		log.Fatal("Couldn't render template ", err)
// 	}
// }

// func reportHandler(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w)
// }

// func Serve() {
// 	http.HandleFunc("/", reportHandler)
// 	err := http.ListenAndServe(":8083", nil)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }
