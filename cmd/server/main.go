package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <db file>\n", os.Args[0])
		os.Exit(1)
	}

	db, err := sqlx.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatalf("Could not open db file:\n%s", err)
	}

	apiHandler := &apiHandler{db: db}

	sm := http.NewServeMux()
	sm.Handle("/api/dinos", apiHandler)
	sm.Handle("/", http.FileServer(http.Dir("static")))

	addr := "[::]:8090"
	log.Printf("Listening on: %s", addr)
	log.Fatal(http.ListenAndServe(addr, loggingWrapper(sm)))
}
