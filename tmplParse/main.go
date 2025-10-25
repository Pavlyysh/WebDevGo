package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

type Product struct {
	Name   string
	ProdID int
	Cost   float64
	Specs  ProdSpecs
}

type ProdSpecs struct {
	Size        string
	Weight      float32
	Description string
}

var prod1 Product

func main() {
	prod1 = Product{
		Name:   "iPhone",
		ProdID: 1,
		Cost:   899,
		Specs: ProdSpecs{
			Size:        "17 x 8 x 2",
			Weight:      200,
			Description: "Brand New Phone with features",
		},
	}
	tmpl, _ = tmpl.ParseGlob("templates/*.html")
	http.HandleFunc("/home", welcomePageHandler)
	http.HandleFunc("/about", aboutPageHandler)
	http.ListenAndServe(":8080", nil)
}

func welcomePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "product_info.html", prod1)
}

func aboutPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "about.html", nil)
}
