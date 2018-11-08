package web

import (
	"fmt"
	"html/template"
	"net/http"
)

func Start() {
	tmpl := template.Must(template.ParseFiles("web/tpl/users.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	// static
	http.Handle("/static/", http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("web/static"))))

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
