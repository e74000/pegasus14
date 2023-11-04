package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"flag"
	"github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func main() {
	dbPath := ""
	webPath := ""
	flag.StringVar(&dbPath, "p", "database.sqlite", "the path to the database")
	flag.StringVar(&webPath, "w", "../frontend/homepage.html", "the path to the homepage")
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

		w.WriteHeader(http.StatusFound)
		_, _ = w.Write(data)
	}).Methods("GET")

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		re := regexp.MustCompile(`^((?!\.)[\w\-_.]*[^.])(@\w+)(\.\w+(\.\w+)?[^.\W])$`)
		if !re.MatchString(user.Email) {
			http.Error(w, "invalid email", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("select count(id) from Users where email = ?", user.Email)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var count int
		rows.Next()

		if err = rows.Scan(&count); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if count != 0 {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		_, err = db.Query("insert into Users (id, email, password_hash) VALUES ((select max(id)+1 from Users), ?, ?)", user.Email, string(hashed))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Info("created user", "email", user.Email, "password", user.Password)
	}).Methods("POST")

	router.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rows, err := db.Query("select * from Users where email = ?", user.Email)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if !rows.Next() {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		var expect User
		err = rows.Scan(&expect.Email, &expect.Password)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(expect.Password), []byte(user.Password))
		if err != nil {
			http.Error(w, "forbidden", http.StatusUnauthorized)
			return
		}

		token, err := SignClaim(user.Email, time.Now().Add(time.Hour*24))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, token)
	}).Methods("POST")

	router.HandleFunc("/impression", func(w http.ResponseWriter, r *http.Request) {
		var impression Impression
		err := json.NewDecoder(r.Body).Decode(&impression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok, err := VerifyClaim(impression.Claim)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !ok {
			http.Error(w, "not authorised", http.StatusUnauthorized)
			return
		}

		if impression.Claim.Email != impression.User {
			http.Error(w, "not authorised", http.StatusUnauthorized)
			return
		}

		_, err = db.Query("insert into Impressions (user, product, liked, view_seconds) VALUES ((select id from Users where user = ?), ?, ?, ?)")
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	router.HandleFunc("/suggest/{email}", func(w http.ResponseWriter, r *http.Request) {
		email := mux.Vars(r)["email"]

		rows, err := db.Query("SELECT p.sku FROM main.Products AS p LEFT JOIN main.Impressions AS i ON p.sku = i.product AND i.user = ? WHERE i.product IS NULL LIMIT 10", email)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		skus := make([]int, 0)

		for rows.Next() {
			var sku int
			err = rows.Scan(&sku)
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

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, webPath)
	})

	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("failed to start http server", "err", err)
	}
}
