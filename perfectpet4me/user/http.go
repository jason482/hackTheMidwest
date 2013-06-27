package user

import (
	"appengine"
	"appengine/datastore"
	aeuser "appengine/user"
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/user", handler)
	http.HandleFunc("/user/save", save)
}

var userTmpls = template.Must(template.ParseFiles("templates/base.html", "perfectpet4me/user/templates/user.html"))

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	aeu := aeuser.Current(c)
	k := datastore.NewKey(c, "user", aeu.Email, 0, nil)

	var u User

	if err := datastore.Get(c, k, &u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := userTmpls.Execute(w, u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func save(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	aeu := aeuser.Current(c)
	k := datastore.NewKey(c, "user", aeu.Email, 0, nil)

	var u User

	if err := datastore.Get(c, k, &u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u.Zip = r.FormValue("zip")
	datastore.Put(c, k, &u)

	http.Redirect(w, r, "/search", http.StatusMovedPermanently)
}
