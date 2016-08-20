package vocaminder

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Vocab contains all the data needed by the learning card
type Vocab struct {
	Word       string
	Phonetics  string
	Definition string
	Audio      string
}

func addVocab(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	vocab := &Vocab{
		Word:       c.PostForm("word"),
		Phonetics:  c.PostForm("phonetics"),
		Definition: c.PostForm("definition"),
		Audio:      c.PostForm("audio"),
	}

	scores := &Scores{
		Word: c.PostForm("word"),
		Results: []struct {
			Date  int
			Score int
		}{},
	}

	key := datastore.NewKey(context, "Vocab", vocab.Word, 0, nil)
	if _, err := datastore.Put(context, key, vocab); err != nil {
		// Handle err
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Can't add new vocab")
		return
	}

	keyScores := datastore.NewKey(context, "Score", scores.Word, 0, nil)
	if _, err := datastore.Put(context, keyScores, scores); err != nil {
		// Handle err
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Can't add score")
		return
	}

	response := map[string]string{
		"message": "word '" + vocab.Word + "' added to the datastore",
	}

	sendSuccessResponse(c, response)
}

func updateVocab(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	vocab := &Vocab{
		Word:       c.PostForm("word"),
		Phonetics:  c.PostForm("phonetics"),
		Definition: c.PostForm("definition"),
		Audio:      c.PostForm("audio"),
	}

	key := datastore.NewKey(context, "Vocab", vocab.Word, 0, nil)
	if _, err := datastore.Put(context, key, vocab); err != nil {
		// Handle err
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Can't add new vocab")
		return
	}

	response := map[string]string{
		"message": "word '" + vocab.Word + "' updated in the datastore",
	}

	sendSuccessResponse(c, response)
}

func deleteVocab(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	word := c.Param("word")

	key := datastore.NewKey(context, "Vocab", word, 0, nil)
	if err := datastore.Delete(context, key); err != nil {
		// Handle err
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Can't delete vocab '"+word+"''")
		return
	}

	response := map[string]string{
		"message": "word '" + word + "' removed from the datastore",
	}

	sendSuccessResponse(c, response)
}

func getVocab(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	word := c.Param("word")

	vocabKey := datastore.NewKey(context, "Vocab", word, 0, nil)

	var v Vocab

	err := datastore.Get(context, vocabKey, &v)

	if err != nil {
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Word ''"+word+"'' not found")
		return
	}

	sendSuccessResponse(c, v)
}
