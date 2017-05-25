package main

import (
	"log"

	"github.com/guilhermebr/pirat.as/backend/shortener"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	if err := shortener.ConnectDB("shortener"); err != nil {
		log.Fatal(err)
	}
	defer shortener.CloseDB()

	n := negroni.Classic()
	r := mux.NewRouter()

	r.HandleFunc("/enc", shortener.Encode).Methods("GET")
	r.HandleFunc("/{key}", shortener.Redir).Methods("GET")

	n.UseHandler(r)
	n.Run(":3000")
}
