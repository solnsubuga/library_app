package handler

import (
	"net/http"

	"github.com/goincremental/negroni-sessions"
	"github.com/snsubuga/library/database"
	"github.com/snsubuga/library/models"
)

/*
  Helper functions
*/

func getStringFromSession(r *http.Request, key string) string {
	var strVal string
	if val := sessions.GetSession(r).Get(key); val != nil {
		strVal = val.(string)
	}
	return strVal
}

func getBookCollection(books *[]models.Book, sortCol string, filterByClass string, w http.ResponseWriter) bool {
	dbmap := database.GetDbMap()
	if sortCol == "" {
		sortCol = "pk"
	}
	var where string
	if filterByClass == "fiction" {
		where = "WHERE classification BETWEEN '800' AND '900'"
	} else if filterByClass == "nonfiction" {
		where = "WHERE classification NOT BETWEEN '800' AND '900'"
	}
	if _, err := dbmap.Select(books, "SELECT * FROM books "+where+" ORDER BY "+sortCol); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	return true
}
