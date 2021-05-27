package main

import (
	"flag"
	"log"
	"net/http"
)

var siteDir = flag.String("siteDir", "site", "Directory to serve")
var listenAddr = flag.String("listenAddr", ":8080", "Server listen address")

func main() {
	flag.Parse()
	log.Printf("Listening on %s\n", *listenAddr)
	fileServer := http.FileServer(http.Dir(*siteDir))
	err := http.ListenAndServe(*listenAddr, fileServer)
	if err != nil {
		log.Fatal(err)
	}
}
