package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"unicode"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var tmpl *template.Template
var db *sql.DB

func main() {
	tmpl, _ = tmpl.ParseGlob("templates/*.html")
	var err error
	db, err = sql.Open("mysql", "pavlyysh:password@tcp(localhost:3306)/grow_adept")
	if err != nil {
		panic(err.Error)
	}
	defer db.Close()

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)

	http.ListenAndServe(":8080", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----registerHandler is running----")
	tmpl.ExecuteTemplate(w, "register.html", nil)
}

func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----registerAuthHandler is running----")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	if !IsValidUsernameLenght(username) || !IsAlphaNumericUsername(username) || !IsStrongPassword(password) {
		tmpl.ExecuteTemplate(w, "register.html", "please check username and password criteria")
		return
	}

	stmt := "SELECT userID FROM grow_adept.bcrypt WHERE username = ?"
	row := db.QueryRow(stmt, username)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		fmt.Println("username already exists")
		tmpl.ExecuteTemplate(w, "register.html", "username already taken")
		return
	}

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt error:", err)
		tmpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	fmt.Println("hash:", hash)
	fmt.Println("string(hash):", string(hash))

	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO grow_adept.bcrypt (username, Hash) VALUES (?, ?);")
	if err != nil {
		fmt.Println("error preparint statement:", err)
		tmpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	defer insertStmt.Close()

	var result sql.Result
	result, err = insertStmt.Exec(username, hash)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAffected:", rowsAff)
	fmt.Println("lastInsertID:", lastIns)
	if err != nil {
		fmt.Println("error inserting new user")
		tmpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	fmt.Fprint(w, "congrats, your account has been successfully created")
}

func IsValidUsernameLenght(username string) bool {
	if len(username) >= 5 && len(username) <= 50 {
		return true
	}
	return false
}

func IsAlphaNumericUsername(username string) bool {
	for _, char := range username {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

func IsStrongPassword(pswd string) bool {
	var upperLetter, lowerLetter, number, symbol bool
	for _, char := range pswd {
		switch {
		case unicode.IsUpper(char):
			upperLetter = true
		case unicode.IsLower(char):
			lowerLetter = true
		case unicode.IsNumber(char):
			number = true
		case unicode.IsSymbol(char) || unicode.IsPunct(char):
			symbol = true
		case unicode.IsSpace(char):
			return false
		}
	}

	if len(pswd) <= 11 || len(pswd) >= 60 {
		return false
	}

	if upperLetter && lowerLetter && number && symbol {
		return true
	}

	return false
}
