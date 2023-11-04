package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"flag"
	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)

func main() {
	dbPath := ""
	flag.StringVar(&dbPath, "p", "database.sqlite", "the path to the database")
	flag.Parse()

	db, err := sql.Open("sqlite3", dbPath)
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

	getStmt, err := db.Prepare("select * from Products where sku = ?")
	if err != nil {
		log.Fatal("failed to prepare get statement", "err", err)
		return
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/product/{sku}", func(w http.ResponseWriter, r *http.Request) {
		skuString := mux.Vars(r)["sku"]

		sku, err := strconv.ParseInt(skuString, 10, 64)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		rows, err := getStmt.Query(sku)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(ParseRows(rows))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(data)
	}).Methods("GET")

	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("select sku from Products")
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		skus := make([]int, 0)
		for rows.Next() {
			sku := 0
			err := rows.Scan(&sku)
			if err != nil {
				continue
			}

			skus = append(skus, sku)
		}

		data, err := json.Marshal(skus)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(data)
	}).Methods("GET")

	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("failed to start http server", "err", err)
	}
}
