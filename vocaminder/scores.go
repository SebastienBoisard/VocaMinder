package vocaminder

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Scores contains all the scores of a user for a word
type Scores struct {
	Word    string
	Results []struct {
		Date  int
		Score int
	}
}

func addScore(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	scores := &Scores{
		Word: c.PostForm("word"),
		Results: []struct {
			Date  int
			Score int
		}{},
	}

	keyScores := datastore.NewKey(context, "Scores", scores.Word, 0, nil)
	if _, err := datastore.Put(context, keyScores, scores); err != nil {
		// Handle err
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Can't add score for word '"+scores.Word+"'")
		return
	}

	response := map[string]string{
		"message": "score for word '" + scores.Word + "' added to the datastore",
	}

	sendSuccessResponse(c, response)
}

func updateScore(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	word := c.PostForm("word")

	newScore, _ := strconv.Atoi(c.PostForm("score"))

	var s Scores

	keyScore := datastore.NewKey(context, "Scores", word, 0, nil)

	err := datastore.Get(context, keyScore, &s)

	if err != nil {
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Score for word ''"+word+"'' not found")
		return
	}

	// Get the current time
	t := time.Now()
	// Store the date of this score as an integer as YYYYMMDDHH
	currentDate, _ := strconv.Atoi(t.Format("2006010215"))

	s.Results = append(s.Results, struct {
		Date  int
		Score int
	}{
		Date:  currentDate,
		Score: newScore,
	})

	if _, err := datastore.Put(context, keyScore, &s); err != nil {
		// Handle err
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Can't update score for word '"+word+"'")
		return
	}

	response := map[string]string{
		"message": "new score for word '" + word + "' updated to the database",
	}

	sendSuccessResponse(c, response)
}
