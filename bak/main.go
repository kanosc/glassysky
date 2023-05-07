package main

import (
	"html/template"
	"log"
	"net/http"
)

func process(w http.ResponseWriter, r *http.Request) {
	log.Println("recive request")
	t, err := template.ParseFiles("pages/index.html")
	if err != nil {
		log.Println(err.Error())
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("response success")
	}
}

func main() {
	server := http.Server{
		Addr: ":80",
	}
	log.Println("website start")
	http.HandleFunc("/", process)
	fs := http.FileServer(http.Dir("pages"))
        http.Handle("/pages/", http.StripPrefix("/pages/", fs))
	err :=server.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
	}
}
