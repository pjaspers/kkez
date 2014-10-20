package main

import "net/http"
import "log"
import "fmt"
import "html"
import "path"
import "html/template"

type Index struct {
	Title string
	Content string
	Subtitle string
}

func main() {
	// Serve our static assets
	// (Make sure to strip out the public, before `FileServer` sees it)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templatePath := path.Join("templates", "index.html")
		tmpl, err:= template.ParseFiles(templatePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		index := Index{"Is het kaarten vandaag?", "Neeje", "Volgende is bij iemand"}
		if err:= tmpl.Execute(w, index); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
