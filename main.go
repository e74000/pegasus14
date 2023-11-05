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
	"strconv"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)

	dbPath := ""
	webPath := ""
	flag.StringVar(&dbPath, "p", "data/database.sqlite", "the path to the database")
	flag.StringVar(&webPath, "w", "pages/", "the path to the homepage")
	flag.Parse()

	log.Info("connecting to database")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("failed to connect to db", "err", err)
	}

	defer func() {
		log.Info("closed database connection")
		_ = db.Close()
	}()

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping db", "err", err)
	}

	log.Info("registering handlers")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/product/{sku}/", func(w http.ResponseWriter, r *http.Request) {
		skuString := mux.Vars(r)["sku"]

		log.Debug("querying sku", "sku", skuString)

		sku, err := strconv.ParseInt(skuString, 10, 64)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("select * from Products where sku = ?", sku)
		if err != nil {
			log.Error("error running query", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		products := make([]Product, 0)

		for rows.Next() {
			var product Product
			err = rows.Scan(&product.SKU, &product.Title, &product.Img, &product.Description, &product.Price)

			log.Debug("got product", "title", product.Title)

			products = append(products, product)
		}

		_ = rows.Close()

		data, err := json.Marshal(products)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(data)
	})

	router.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
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

		_ = rows.Close()

		data, err := json.Marshal(skus)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusFound)
		_, _ = w.Write(data)
	})

	router.HandleFunc("/register/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("got registration request")

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Error("failed to decode json", "err", err)
			http.Error(w, "failed to decode request", http.StatusBadRequest)
			return
		}

		log.Info("got user details", "email", user.Email, "password", user.Password)

		rows, err := db.Query("select count(id) from Users where email = ?", user.Email)
		if err != nil {
			log.Error("failed to query db", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var count int
		rows.Next()

		if err = rows.Scan(&count); err != nil {
			log.Error("failed read response db", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		_ = rows.Close()

		if count != 0 {
			log.Error("user already exists")
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate password hash", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		log.Info("got password hash", "hash", string(hashed))

		_, err = db.Exec("INSERT INTO Users (id, email, password_hash) VALUES (COALESCE((SELECT max(id)+1 FROM Users), 1), ?, ?)", user.Email, string(hashed))
		if err != nil {
			log.Error("failed to query db", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		log.Info("user successfully created")

		w.WriteHeader(http.StatusCreated)
		log.Info("created user", "email", user.Email, "password", user.Password)
	})

	router.HandleFunc("/validate/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("validating user")
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Error("failed to decode request", "err", err)
			http.Error(w, "could not decode request", http.StatusBadRequest)
			return
		}

		log.Info("got user details", "email", user.Email, "password", user.Password)

		rows, err := db.Query("select email, password_hash from Users where email = ?", user.Email)
		if err != nil {
			log.Error("error running query", "err", err)
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
			log.Error("error fetching email", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		_ = rows.Close()

		log.Info("comparing passwords", "expected_hash", expect.Password, "password", user.Password)

		err = bcrypt.CompareHashAndPassword([]byte(expect.Password), []byte(user.Password))
		if err != nil {
			log.Error("password does not match")
			http.Error(w, "forbidden", http.StatusUnauthorized)
			return
		}

		log.Info("signing token")

		token, err := SignClaim(user.Email, time.Now().Add(time.Hour*24))
		if err != nil {
			log.Error("error creating token", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		log.Info("issued token", "token", token)

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, token)
	})

	router.HandleFunc("/validate_token/", func(w http.ResponseWriter, r *http.Request) {
		var claim Claim
		err := json.NewDecoder(r.Body).Decode(&claim)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok, err := VerifyClaim(claim)
		if err != nil || !ok {
			http.Error(w, "unauthorised", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	})

	router.HandleFunc("/impression/", func(w http.ResponseWriter, r *http.Request) {
		var impression Impression
		err := json.NewDecoder(r.Body).Decode(&impression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if impression.SKU == 0 {
			w.WriteHeader(http.StatusOK)
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

		if impression.Claim.Email != impression.Email {
			http.Error(w, "not authorised", http.StatusUnauthorized)
			return
		}

		_, err = db.Exec("insert into Impressions (email, sku, swipe) VALUES (?, ?, ?)", impression.Email, impression.SKU, impression.Swipe)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	router.HandleFunc("/suggest/{email}/", func(w http.ResponseWriter, r *http.Request) {
		email := mux.Vars(r)["email"]

		log.Info("getting suggestions", "email", email)

		rows, err := db.Query("SELECT p.sku FROM main.Products AS p LEFT JOIN main.Impressions AS i ON p.sku = i.sku AND i.email = ? WHERE i.sku IS NULL LIMIT 10", email)
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

		_ = rows.Close()

		log.Info("got skus", "skus", skus)

		data, err := json.Marshal(skus)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})

	router.HandleFunc("/basket/{email}/", func(w http.ResponseWriter, r *http.Request) {
		email := mux.Vars(r)["email"]

		rows, err := db.Query("SELECT sku FROM Impressions WHERE email = ? AND swipe = 0", email)
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

	router.HandleFunc("/basket/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/basket.html")
	})

	router.HandleFunc("/swipes/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT count(swipe) FROM Impressions")
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if !rows.Next() {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var count int
		err = rows.Scan(&count)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(count)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/home.html")
	})

	router.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/login.html")
	})

	router.HandleFunc("/app/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/cards_page.html")
	})

	router.HandleFunc("/terms_and_conditions/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/terms_and_conditions.html")
	})

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Info("starting server")

	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("failed to start http server", "err", err)
	}
}
