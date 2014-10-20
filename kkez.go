package main

import "net/http"
import "log"
import "fmt"
import "html"
import "path"
import "html/template"

type LayoutData struct {
	Title string
	Data interface{}
}

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
		layout := path.Join("templates", "layout.html")
		parse := func(name string) (*template.Template, error) {
			t := template.New("") //.Funcs(funcMap)
			return t.ParseFiles(layout, path.Join("templates", name))
		}
		index, err:= parse("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var layoutData = LayoutData{Title: "Dinges"}
		layoutData.Data = Index{"Is het kaarten vandaag?", "Neeje", "Volgende is bij iemand"}
		if err:= index.ExecuteTemplate(w, "layout.html", layoutData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
