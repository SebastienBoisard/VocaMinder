package vocaminder

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func loadData(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	file, e := ioutil.ReadFile("vocaminder/data/vocab_list.json")
	if e != nil {
		log.Errorf(context, "File error: %v\n", e)
		sendFailResponse(c, "Can't load vocab data")
		return
	}

	vocabList := []Vocab{}
	json.Unmarshal(file, &vocabList)

	for _, vocab := range vocabList {

		keyVocab := datastore.NewKey(context, "Vocab", vocab.Word, 0, nil)
		if _, err := datastore.Put(context, keyVocab, &vocab); err != nil {
			// Handle err
			log.Errorf(context, "%v", err)
			sendFailResponse(c, "Can't add new vocab")
			return
		}

		s := &Scores{
			Word: vocab.Word,
			Results: []struct {
				Date  time.Time
				Score int
			}{},
		}

		scoreKey := datastore.NewKey(context, "Scores", vocab.Word, 0, nil)

		if err := datastore.Get(context, scoreKey, s); err != nil {
			log.Debugf(context, "Scores for %s doesn't exist", s.Word)
			if _, err := datastore.Put(context, scoreKey, s); err != nil {
				// Handle err
				log.Errorf(context, "%v", err)
				sendFailResponse(c, "Can't add score for word '"+s.Word+"'")
				return
			}
		}
		log.Debugf(context, "Scores for %s exists", s.Word)

	}

	sendSuccessResponse(c, vocabList)
}

func deleteData(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	q1 := datastore.NewQuery("Vocab")
	var allVocabs []Vocab
	keys, err := q1.GetAll(context, &allVocabs)

	for _, key := range keys {
		if err = datastore.Delete(context, key); err != nil {
			// Handle err
			log.Errorf(context, "%v", err)
			sendFailResponse(c, "Can't delete vocab key="+key.Encode())
			return
		}
	}

	q2 := datastore.NewQuery("Scores")
	var allScores []Scores
	keys, err = q2.GetAll(context, &allScores)

	for _, key := range keys {
		if err = datastore.Delete(context, key); err != nil {
			// Handle err
			log.Errorf(context, "%v", err)
			sendFailResponse(c, "Can't delete score key="+key.Encode())
			return
		}
	}

	sendSuccessResponse(c, "all the data from the datastore were deleted")
}

func downloadData(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	q := datastore.NewQuery("Vocab")
	var allVocab []Vocab
	_, err := q.GetAll(context, &allVocab)

	if err != nil {
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Error while downloading all the vocabs")
		return
	}
	sendSuccessResponse(c, allVocab)
}
