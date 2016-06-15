package main

import (
	"log"
	"runtime"

	"github.com/guilhermebr/pirat.as/backend/shortener"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	log.Printf("CPU's %d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	if err := shortener.ConnectDB("shortener"); err != nil {
		log.Fatal(err)
	}

	n := negroni.Classic()
	r := mux.NewRouter()

	r.HandleFunc("/{key}", shortener.Redir).Methods("GET")
	r.HandleFunc("/enc/", shortener.Encode).Methods("GET")

	n.UseHandler(r)
	n.Run(":3000")
}
