package main

import (
	"database/sql"
	_ "embed"
	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "/Users/olivermcinnes/GolandProjects/pegasus14/backend/database.sqlite")
	if err != nil {
		log.Fatal("failed to connect to db", "err", err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping db", "err", err)
	}

	router := mux.NewRouter().StrictSlash(true)

	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("failed to start http server", "err", err)
	}
}
