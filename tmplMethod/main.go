package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Price float64

// Canadian cash price
func (p Price) CanCashPrice() string {
	remainder := int(p*100) % 5
	quotiant := int(p*100) / 5

	if remainder < 3 {
		pr := float64(quotiant*5) / 100
		s := fmt.Sprintf("%.2f", pr)
		return s
	}
	pr := (float64(quotiant*5) + 5) / 100
	s := fmt.Sprintf("%.2f", pr)
	return s
}

var tmpl *template.Template

var p Price

func main() {
	tmpl, _ = tmpl.ParseFiles("index.html")
	p = 4.97

	http.HandleFunc("/price", priceHandler)
	http.ListenAndServe(":8080", nil)
}

func priceHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", p)
}
