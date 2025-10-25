package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var tmpl, _ = template.New("index.html").Funcs(template.FuncMap{
	"CanCashPrice": func(p float64) string {
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
	},
	"Upper": strings.ToUpper,
}).ParseFiles("index.html")

var p float64

func main() {
	p = 3.93
	http.HandleFunc("/price", priceHandler)
	http.ListenAndServe(":8080", nil)
}

func priceHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", p)
}
