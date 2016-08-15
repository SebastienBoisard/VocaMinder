package vocaminder

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

// Vocab is a structure for a word
type Vocab struct {
	//	id         int
	Word       string
	Phonetics  string
	Definition string
}

var tpl = template.Must(template.ParseGlob("vocaminder/templates/*.html"))

// init is called before the application starts.
func init() {
	// Register a handler for / URLs.
	http.HandleFunc("/", handleHomePage)
	http.HandleFunc("/add_vocab", handleAddVocab)
	http.HandleFunc("/new_vocab", handleNewVocab)
}

// getHomePage is an HTTP handler that prints "HomePage"
func handleHomePage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	context := appengine.NewContext(r)
	u := user.Current(context)

	loginURL, _ := user.LoginURL(context, "/")
	logoutURL, _ := user.LogoutURL(context, "/")

	d := struct {
		Data        interface{}
		AuthEnabled bool
		LoginURL    string
		LogoutURL   string
	}{
		Data:        nil,
		AuthEnabled: u != nil,
		LoginURL:    loginURL,
		LogoutURL:   logoutURL,
	}

	if err := tpl.ExecuteTemplate(w, "homepage.html", d); err != nil {
		log.Errorf(context, "%v", err)
	}

	/*
		if u == nil {
			url, _ := user.LoginURL(context, "/")
			fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
			return
		}
		url, _ := user.LogoutURL(context, "/")
	*/
	//	fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}

func handleNewVocab(w http.ResponseWriter, r *http.Request) {

	context := appengine.NewContext(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "add_vocab.html", nil); err != nil {
		log.Errorf(context, "%v", err)
	}
}

func handlNewVocab(w http.ResponseWriter, r *http.Request) {

	context := appengine.NewContext(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "addvocab.html", nil); err != nil {
		log.Errorf(context, "%v", err)
	}
}

func handleAddVocab(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
		return
	}
	context := appengine.NewContext(r)

	vocab := &Vocab{
		Word:       r.FormValue("word"),
		Phonetics:  r.FormValue("phonetics"),
		Definition: r.FormValue("definition"),
	}

	key := datastore.NewIncompleteKey(context, "Vocab", nil)
	if _, err := datastore.Put(context, key, vocab); err != nil {
		// Handle err
		log.Errorf(context, "%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect with 303 which causes the subsequent request to use GET.
	//http.Redirect(w, r, "/", http.StatusSeeOther)
}