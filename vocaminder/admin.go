package vocaminder

import (
	"encoding/json"
	"io/ioutil"

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
	}

	sendSuccessResponse(c, vocabList)
}

func deleteData(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	q := datastore.NewQuery("Vocab")
	var allVocabs []Vocab
	keys, err := q.GetAll(context, &allVocabs)

	for _, key := range keys {
		if err = datastore.Delete(context, key); err != nil {
			// Handle err
			log.Errorf(context, "%v", err)
			sendFailResponse(c, "Can't delete vocab key="+key.Encode())
			return
		}
	}

	sendSuccessResponse(c, "all the vocabs were deleted")
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
