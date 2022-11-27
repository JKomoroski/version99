package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", ":8080", "Host:port on which to listen")
)

// --------------------------------------------------------------------

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
