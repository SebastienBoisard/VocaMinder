package vocaminder

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Card contains the data of a card to learn
type Card struct {
	Word         string
	vocabToLearn Vocab
	score        int
	learned      bool
}

func initCard(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	log.Debugf(context, "initCard")

	q := datastore.NewQuery("Scores")
	var allScores []Scores
	_, err := q.GetAll(context, &allScores)

	if err != nil {
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Error while listing all scores")
		return
	}

	var newCards []Card
	var repeatCards []Card

	log.Debugf(context, "initCard %v", allScores)

	for _, score := range allScores {
		nbResults := len(score.Results)

		log.Debugf(context, "nbResults=%d", nbResults)

		if nbResults == 0 {
			newCards = append(newCards, Card{
				Word: score.Word,
			})
			continue
		}

		log.Debugf(context, "score.Results[nbResults].Score=%d", score.Results[nbResults-1].Score)

		if score.Results[nbResults-1].Score != goodScore {
			repeatCards = append(repeatCards, Card{
				Word: score.Word,
			})
			continue
		}

		if nbResults == 1 {
			repeatCards = append(repeatCards, Card{
				Word: score.Word,
			})
			continue
		}

		if score.Results[nbResults-2].Score != goodScore {
			repeatCards = append(repeatCards, Card{
				Word: score.Word,
			})
			continue
		}

		interval := score.Results[nbResults-1].Date - score.Results[nbResults-2].Date

		if interval == 1 {
			repeatCards = append(repeatCards, Card{
				Word: score.Word,
			})
			continue
		}

		now := 1

		if score.Results[nbResults-1].Date+interval+1 > now {
			continue
		}

		repeatCards = append(repeatCards, Card{
			Word: score.Word,
		})

	}

	sendSuccessResponse(c, repeatCards)
}

func getCard(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	word := "first"

	var v Vocab

	vocabKey := datastore.NewKey(context, "Vocab", word, 0, nil)

	err := datastore.Get(context, vocabKey, &v)
	if err != nil {
		log.Errorf(context, "%v", err)
		sendFailResponse(c, "Word ''"+word+"'' not found")
		return
	}

	sendSuccessResponse(c, v)
}
