package handler

import (
	"net/http"

	"github.com/goincremental/negroni-sessions"
	"github.com/snsubuga/library/database"
	"github.com/snsubuga/library/models"

	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
)

type loginPage struct {
	Error string
}

//TODO: seperate GET and POST Login and Register Handlers

//LoginHandler allows a user to login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	dbmap := database.GetDbMap()
	var p loginPage
	username := r.FormValue("username")
	password := r.FormValue("password")

	//register the user
	if r.FormValue("register") != "" {

		key, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := models.User{Username: username, Password: string(key), ID: uuid.New()}
		if err := dbmap.Insert(&user); err != nil {
			p.Error = err.Error()

		} else {
			sessions.GetSession(r).Set("User", user.Username)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

	} else if r.FormValue("login") != "" {
		user, err := dbmap.Get(models.User{}, username)
		if err != nil {
			p.Error = err.Error()
		} else if user == nil {
			p.Error = "Wrong username"
		} else {
			u := user.(*models.User)
			if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
				p.Error = err.Error()
			} else {
				sessions.GetSession(r).Set("User", u.Username)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}

	}

	if err := templates.ExecuteTemplate(w, "login.html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//LogoutHandler logs a user out of the app
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessions.GetSession(r).Set("User", nil)
	sessions.GetSession(r).Set("Filter", nil)

	http.Redirect(w, r, "/login", http.StatusFound)
	return
}
