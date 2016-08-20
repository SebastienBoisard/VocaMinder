package vocaminder

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

// Scores contains all the scores of a user for a word
type Scores struct {
	Word    string
	Results []struct {
		Date  int
		Score int
	}
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

	r.POST("/vocab", addVocab)
	r.PUT("/vocab", updateVocab)
	r.DELETE("/vocab/:word", deleteVocab)

	r.GET("/vocab/id/:word", getVocabID)
	r.POST("/score/new", setVocabScore)

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

func getVocabID(c *gin.Context) {

	word := c.Param("word")

	context := appengine.NewContext(c.Request)

	vocabKey := datastore.NewKey(context, "Vocab", word, 0, nil)

	var v Vocab

	err := datastore.Get(context, vocabKey, &v)

	if err != nil {
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Word ''"+word+"'' not found")
		return
	}

	response := map[string]string{
		"id": v.Word,
	}

	sendSuccessResponse(c, response)
}

func setVocabScore(c *gin.Context) {
	context := appengine.NewContext(c.Request)
	log.Errorf(context, "setVocabScore")

	// @TODO: get the user id
	// u := user.Current(context)
	// u.ID or u.ClientID
	// https://cloud.google.com/appengine/docs/go/users/reference#Current
	// https://cloud.google.com/appengine/docs/go/users/
	// @TODO: make a query to find a score with key=word and user=u.ID

	word := c.PostForm("word")
	newScore, _ := strconv.Atoi(c.PostForm("score"))

	keyScore := datastore.NewKey(context, "Scores", word, 0, nil)

	var s Scores

	err := datastore.Get(context, keyScore, &s)

	if err != nil {
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Score for word ''"+word+"'' not found")
		return
	}
	log.Errorf(context, "setVocabScore", s)

	s.Results = append(s.Results, struct {
		Date  int
		Score int
	}{
		Date:  20160819,
		Score: newScore,
	})

	log.Errorf(context, "setVocabScore", s)
	if _, err := datastore.Put(context, keyScore, &s); err != nil {
		// Handle err
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Can't add score")
		return
	}

	response := map[string]string{
		"message": "new score for word '" + word + "' added to the database",
	}

	sendSuccessResponse(c, response)
}
