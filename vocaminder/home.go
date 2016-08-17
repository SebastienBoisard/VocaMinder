package vocaminder

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

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

	// Starts a new Gin instance with no middle-ware
	r := gin.New()

	// Define your handlers
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/vocab/id/:word", getVocabID)
	r.GET("/vocab/new", addNewVocab)

	// Handle all requests using net/http
	http.Handle("/", r)
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

// Response is following the [JSend convention](https://labs.omniti.com/labs/jsend)
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func getVocabID(c *gin.Context) {
	word := c.Param("word")

	context := appengine.NewContext(c.Request)

	vocabKey := datastore.NewKey(context, "Vocab", word, 0, nil)

	var v Vocab

	err := datastore.Get(context, vocabKey, &v)

	if err != nil {

		r := &Response{
			Status: "fail",
			Data: map[string]string{
				"title": "Word ''" + word + "'' not found",
				"error": err.Error(),
			},
		}
		jsonResponse, _ := json.Marshal(r)

		log.Errorf(context, "%v", err)

		c.String(http.StatusBadRequest, string(jsonResponse))
		return
	}

	r := &Response{
		Status: "success",
		Data: map[string]string{
			"word": word,
			"id":   "1",
		},
	}
	jsonResponse, _ := json.Marshal(r)

	c.String(http.StatusOK, string(jsonResponse))
}

func addNewVocab(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	vocab := &Vocab{
		Word:       "first",        //c.PostForm("word"),
		Phonetics:  "/first/",      //c.PostForm("phonetics"),
		Definition: "def of first", //c.PostForm("definition"),
	}

	key := datastore.NewKey(context, "Vocab", vocab.Word, 0, nil)
	if _, err := datastore.Put(context, key, vocab); err != nil {
		// Handle err
		log.Errorf(context, "%v", err)
		c.String(http.StatusInternalServerError, "error="+err.Error())
		return
	}

	// Redirect with 303 which causes the subsequent request to use GET.
	//http.Redirect(w, r, "/", http.StatusSeeOther)

	c.String(http.StatusOK, "Word '"+vocab.Word+"' added")
}
