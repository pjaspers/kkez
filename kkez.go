package main

import "net/http"
import "log"
import "fmt"
import "html"
// import "html/template"

func main() {
	// Serve our static assets
	// (Make sure to strip out the public, before `FileServer` sees it)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
