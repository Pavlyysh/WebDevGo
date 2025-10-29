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

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)

	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---loginHandler---")
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---loginAuthHandler---")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	var hash string
	stmt := "SELECT hash FROM grow_adept.login_table WHERE username = ?"
	row := db.QueryRow(stmt, username)
	err := row.Scan(&hash)
	if err != nil {
		fmt.Println("error selecting hash in db by username:", err)
		tmpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	fmt.Println("hash from db:", hash)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		fmt.Fprint(w, "You have successfully logged in")
		return
	}

	fmt.Println("incorrect password")
	tmpl.ExecuteTemplate(w, "login.html", "check username and password")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---registerHandler is running---")
	tmpl.ExecuteTemplate(w, "register.html", nil)
}

func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---registerAuthHandler---")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	if !IsValidUsernameLenght(username) || !IsAlphaNumericUsername(username) || !IsStrongPassword(password) {
		tmpl.ExecuteTemplate(w, "register.html", "please enter correct username and password")
		return
	}

	stmt := "SELECT userID FROM grow_adept.login_table WHERE username = ?;"
	row := db.QueryRow(stmt, username)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		fmt.Println("username already taken")
		tmpl.ExecuteTemplate(w, "register.html", "username already taken")
		return
	}

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("hashing error:", err)
		tmpl.ExecuteTemplate(w, "register.html", "server error, please try again later")
		return
	}

	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO grow_adept.login_table (username, hash) VALUES (?, ?);")
	if err != nil {
		fmt.Println("db preparing error:", err)
		tmpl.ExecuteTemplate(w, "register.html", "server error, please try again later")
		return
	}
	defer insertStmt.Close()

	var sqlResult sql.Result
	sqlResult, err = insertStmt.Exec(username, hash)
	rowsAff, _ := sqlResult.RowsAffected()
	insIds, _ := sqlResult.LastInsertId()
	fmt.Println("Rows Affected:", rowsAff)
	fmt.Println("Last Insert ID:", insIds)
	if err != nil {
		fmt.Println("error inserting user")
		tmpl.ExecuteTemplate(w, "register.html", "server error, please ty again later")
		return
	}

	fmt.Fprint(w, "Thank you for registration!\nYour account has been successfully created!")

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
