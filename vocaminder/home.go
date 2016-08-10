package vocaminder

import (
	"fmt"
	"net/http"
)

// init is called before the application starts.
func init() {
	// Register a handler for / URLs.
	http.HandleFunc("/", getHomePage)
}

// getHomePage is an HTTP handler that prints "HomePage"
func getHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "HomePage")
}
