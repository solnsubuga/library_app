package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	gmux "github.com/gorilla/mux"
	"github.com/snsubuga/library/database"
	"github.com/snsubuga/library/handler"
	"github.com/urfave/negroni"
)

func verifyDatabase(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	db := database.GetDbConnection()
	if err := db.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	next(w, r)
}

func main() {
	//initialize database
	database.DB, _ = sql.Open("mysql", "root:@/librarystoredb")
	database.InitDb()

	mux := gmux.NewRouter().StrictSlash(false)

	//serve static files like CSS, JS
	mux.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	//Application routes
	mux.HandleFunc("/", handler.HomeHandler).Methods("GET")
	mux.HandleFunc("/search", handler.SearchHandler).Methods("POST")
	mux.HandleFunc("/books", handler.AddBookHandler).Methods("PUT")
	mux.HandleFunc("/books", handler.SortBooksHandler).Methods("GET").Queries("sortBy", "{sortBy:title|author|classification}")
	mux.HandleFunc("/books", handler.FilterBooksHandler).Methods("GET").Queries("filter", "{filter:all|fiction|nonfiction}")
	mux.HandleFunc("/books/{pk}", handler.DeleteBookHandler).Methods("DELETE")

	//Authentication
	mux.HandleFunc("/login", handler.LoginHandler)
	mux.HandleFunc("/logout", handler.LogoutHandler)

	n := negroni.Classic()

	store := cookiestore.New([]byte("secret123"))
	n.Use(sessions.Sessions("library_session", store))

	n.Use(negroni.HandlerFunc(verifyDatabase))
	n.Use(negroni.HandlerFunc(handler.VerifyUser))

	//recover
	n.Use(negroni.NewRecovery())

	n.UseHandler(mux)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        n,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
