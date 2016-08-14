package vocaminder

import (
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
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
	fmt.Fprint(w, "HomePage")
}

func handleNewVocab(w http.ResponseWriter, r *http.Request) {

	context := appengine.NewContext(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "addvocab.html", nil); err != nil {
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
