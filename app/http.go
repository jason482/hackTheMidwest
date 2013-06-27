package perfectpet4me

import (
	"appengine"
	"appengine/datastore"
	aeuser "appengine/user"
	"html/template"
	"net/http"
	"perfectpet4me/user"
)

func init() {
	http.HandleFunc("/", handler)
}

var baseTmpl = template.Must(template.ParseFiles("templates/base.html"))

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	u := aeuser.Current(c)
	if u == nil {
		url, err := aeuser.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	} else {
		k := datastore.NewKey(c, "user", u.Email, 0, nil)
		ppu := user.User{
			Email: u.Email,
		}
		if _, err := datastore.Put(c, k, &ppu); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/index.html", http.StatusMovedPermanently)
	return
}
