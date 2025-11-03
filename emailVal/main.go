package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/joho/godotenv"
)

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func main() {
	err := godotenv.Load("emailVal.env")
	if err != nil {
		panic(err)
	}

	emails := []string{
		"!#$%&'*+-/=?^_`{|}~0123456789@mail.com", // allowed characters in local part !#$%&'*+-/=?^_`{|}~
		"janedoe@email.com",                      // uppercase and lowercase Latin letters A to Z and a to z allowd in local part
		"JaneDoe@email.com",                      // // uppercase and lowercase Latin letters A to Z and a to z allowd in local part
		"Jane.Doe@email.com",                     // . allowed in local part as long as it's not the first, last, or there are consectutives
		"Jane-Doe@email.com",                     // allowed characters in local part !#$%&'*+-/=?^_`{|}~
		"a@email.com",
		// "abcdedfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzazzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz@mail.com", // local longer than 64 characters
		"JaneDoeemail.com",           // must have @
		"Jane@Doe@email.com",         // cannot have more than one @
		"Jane:Doe@email.com",         // (),:;<>@[\] not allowed outsided of qoutes
		"janedoe@email.-com",         // - cannot be first or last character in domain
		"janedoe@email.com-",         // - cannot be first or last character in domain
		"janedoe@email.)com",         // domain may only contain a-z, A-Z, 0-9, or hyphen
		"janedoe@mysite_email.!com",  // domain may only contain a-z, A-Z, 0-9, or hyphen
		"janedoe@mysite+account.com", // domain may only contain a-z, A-Z, 0-9, or hyphen
	}

	w := tabwriter.NewWriter(os.Stdout, 10, 0, 3, ' ', tabwriter.Debug)

	fmt.Fprint(w, "Pass Regex\tPass Length\tEmail Address\t\n")
	for _, v := range emails {
		rg := emailRegexp.MatchString(v)
		ln := len(v) > 3 && len(v) < 254
		fmt.Fprintln(w, rg, "\t", ln, "\t", v, "\t")
	}
	w.Flush()

	addr := os.Getenv("ToEmailAddr")
	// addr := "johndoe@nothisisnotadomain.com"
	// func Split(s, sep string) []string
	i := strings.Index(addr, "@")
	fmt.Println("i:", i)
	host := addr[i+1:]
	fmt.Println(host)
	// func LookupMX(name string) ([]*MX, error)
	_, err = net.LookupMX(host)
	if err != nil {
		fmt.Println("invalid receiver email domain:", err)
	}
	fmt.Println(err)
}
