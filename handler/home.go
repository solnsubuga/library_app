package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
	"github.com/snsubuga/library/api"
	"github.com/snsubuga/library/database"
	"github.com/snsubuga/library/models"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFiles("templates/index.html", "templates/login.html"))
}

//HomeHandler ...
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	p := models.Page{
		Books:  []models.Book{},
		Filter: getStringFromSession(r, "Filter"),
		User:   getStringFromSession(r, "User"),
	}
	if !getBookCollection(&p.Books, getStringFromSession(r, "SortBy"), getStringFromSession(r, "Filter"), w) {
		return
	}
	if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//SearchHandler ...
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	//handle book search

	var results []models.SearchResult
	var err error
	if results, err = api.Search(r.FormValue("search")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//AddBookHandler ...
func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	dbmap := database.GetDbMap()
	var err error
	var book models.ClassifyBookResponse
	if book, err = api.Find(r.FormValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := models.Book{
		Pk:             -1,
		Title:          book.BookData.Title,
		Author:         book.BookData.Author,
		Classification: book.Classification.MostPopular,
	}
	if err = dbmap.Insert(&b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//DeleteBookHandler deletes unwanted books from the collection
func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	dbmap := database.GetDbMap()
	pk, _ := strconv.ParseInt(mux.Vars(r)["pk"], 10, 64)
	if _, err := dbmap.Delete(&models.Book{Pk: pk}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//SortBooksHandler sorts books based on a column name
func SortBooksHandler(w http.ResponseWriter, r *http.Request) {
	var b []models.Book
	if !getBookCollection(&b, r.FormValue("sortBy"), getStringFromSession(r, "Filter"), w) {
		return
	}
	sessions.GetSession(r).Set("SortBy", r.FormValue("sortBy"))

	if err := json.NewEncoder(w).Encode(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//FilterBooksHandler filters books based on a user preference
func FilterBooksHandler(w http.ResponseWriter, r *http.Request) {
	var b []models.Book
	if !getBookCollection(&b, getStringFromSession(r, "SortBy"), r.FormValue("filter"), w) {
		return
	}

	//store the filter in a session
	sessions.GetSession(r).Set("Filter", r.FormValue("filter"))

	if err := json.NewEncoder(w).Encode(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//VerifyUser middleware
func VerifyUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//TODO: remove simple hack for middleware not to run on routes
	//For login route check if user is already logged in and redirect to main page
	if r.URL.Path == "/login" || r.URL.Path == "/static/style.css" {
		next(w, r)
		return
	}
	if username := getStringFromSession(r, "User"); username != "" {
		dmap := database.GetDbMap()
		if user, _ := dmap.Get(models.User{}, username); user != nil {
			next(w, r)
			return
		}
	}
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
