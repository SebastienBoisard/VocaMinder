package vocaminder

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

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
	r.GET("/vocab/:word", getVocab)

	r.POST("/scores", addScore)
	r.PUT("/scores", updateScore)
	r.GET("/scores/:word", getScore)

	r.GET("/card", getCard)
	r.PUT("/card", updateCard)

	r.POST("/admin/datastore/data", loadData)
	r.DELETE("/admin/datastore/data", deleteData)
	r.GET("/admin/datastore/data", downloadData)

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
