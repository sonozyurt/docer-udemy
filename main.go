package main

import (
	"html/template"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	UserName string
	Password []byte
	First    string
	Last     string
}
type session struct {
	un           string
	lastactivity time.Time
}

var tpl *template.Template
var dbUsers = make(map[string]user)
var dbSessions = make(map[string]session)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	bs, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	dbUsers["test@test.com"] = user{"test@test.com", bs, "James", "Bond"}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", nil)

}
func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(req)
	tpl.ExecuteTemplate(w, "index.gohtml", u)
}
func bar(w http.ResponseWriter, req *http.Request) {
	u := getUser(req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}
func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")

		if _, ok := dbUsers[un]; ok {
			http.Error(w, "username already taken", http.StatusForbidden)
			return
		}
		SID := uuid.NewV4()
		cookie := &http.Cookie{
			Name:  "session",
			Value: SID.String(),
		}
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.SetCookie(w, cookie)
		u := user{un, bs, f, l}
		dbSessions[cookie.Value] = session{un, time.Now()}
		dbUsers[un] = u
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}
func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "username or password dont match", http.StatusInternalServerError)
			return
		}
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			http.Error(w, "username or password dont match", http.StatusInternalServerError)
			return
		}
		sID := uuid.NewV4()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, cookie)
		dbSessions[cookie.Value] = session{un, time.Now()}
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}
func logout(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	cookie, _ := req.Cookie("session")
	delete(dbSessions, cookie.Value)
	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}
