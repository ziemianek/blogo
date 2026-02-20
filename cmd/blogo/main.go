package main

import (
	"html/template"
	"log"
	"os"
)

const Title string = "Test strona"

type config struct {
	Title string
}

func main() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.Create("output.html")
	check(err)
	defer file.Close()

	tpl, err := template.ParseFiles("./web/static/html/index.tmpl")
	check(err)

	tpl.Execute(file, &config{Title: Title})
	check(err)
}
