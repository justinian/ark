package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func reloadDb(path string, api *apiHandler) error {
	api.lock.Lock()
	defer api.lock.Unlock()

	api.stmt = nil
	api.db.Close()

	db, err := sqlx.Open("sqlite3", os.Args[1])
	if err != nil {
		return err
	}

	log.Print("Reloaded modified database")
	api.db = db
	return nil
}

func dbWatcher(path string, api *apiHandler) {
	for {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatalf("Error creating DB watcher: %s", err)
		}

		err = watcher.Add(path)
		if err != nil {
			log.Fatalf("Error watching DB: %s", err)
		}

		select {
		case <-watcher.Events:
			err = watcher.Close()
			if err != nil {
				log.Fatalf("Error closing watcher: %s", err)
			}
			err = reloadDb(path, api)
			if err != nil {
				log.Fatalf("Error reloading DB: %s", err)
			}

		case err := <-watcher.Errors:
			log.Fatalf("Error watching DB: %s", err)
		}
	}
}

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
	go dbWatcher(os.Args[1], apiHandler)

	sm := http.NewServeMux()
	sm.Handle("/api/dinos", apiHandler)
	sm.Handle("/", http.FileServer(http.Dir("static")))

	addr := "[::]:8090"
	log.Printf("Listening on: %s", addr)
	log.Fatal(http.ListenAndServe(addr, loggingWrapper(sm)))
}
