package vocaminder

import (
	"math/rand"
	"time"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const maxCards = 30

// Card contains the data of a card to learn
type Card struct {
	Word    string
	score   int
	learned bool
}

// Cards stores all the cards that had to be learned
type Cards struct {
	Date     time.Time
	CardList []Card
}

func newCards(c context.Context) (*Cards, error) {

	log.Debugf(c, "newCards")

	q := datastore.NewQuery("Scores")
	var allScores []Scores
	_, err := q.GetAll(c, &allScores)

	if err != nil {
		log.Errorf(c, "%v", err)
		return nil, err
	}

	var newCards []Card
	var repeatCards []Card

	log.Debugf(c, "initCard %v", allScores)

	for _, score := range allScores {
		nbResults := len(score.Results)

		log.Debugf(c, "nbResults=%d", nbResults)

		// Test if there's a Score for this vocab
		if nbResults == 0 {
			// There's no Score for this vocab, so it's a new one.
			// Add this vocab to the list of the new cards.
			newCards = append(newCards, Card{
				Word: score.Word,
			})
			continue
		}

		log.Debugf(c, "score.Results[nbResults].Score=%d", score.Results[nbResults-1].Score)

		// Test if the last Score is a good one.
		if score.Results[nbResults-1].Score != goodScore {
			// The last Score of the current vocab is not a good one,
			// so we have to do a rehearsal today
			repeatCards = append(repeatCards, Card{
				Word: score.Word,
			})
			continue
		}

		// Test if there is only one Score (so only one learning) for the current word
		if nbResults == 1 {
			// There is only one Score for the current word, so we have to do a rehearsal today.
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

		intervalDuration := score.Results[nbResults-1].Date.Sub(score.Results[nbResults-2].Date)

		if intervalDuration.Hours() < 25.0 {
			repeatCards = append(repeatCards, Card{
				Word: score.Word,
			})
			continue
		}

		now := time.Now()

		if now.After(score.Results[nbResults-1].Date.Add(intervalDuration).AddDate(0, 0, 1)) == true {
			continue
		}

		repeatCards = append(repeatCards, Card{
			Word: score.Word,
		})

	}

	if len(repeatCards) < maxCards {

		nbNewCardToAdd := maxCards - len(repeatCards)
		for i := 0; i < nbNewCardToAdd && i < len(newCards); i++ {
			repeatCards = append(repeatCards, newCards[i])
		}
	}
	cards := &Cards{Date: time.Now(), CardList: repeatCards}
	return cards, nil
}

func getCard(c *gin.Context) {

	context := appengine.NewContext(c.Request)

	cardsKey := datastore.NewKey(context, "Cards", "cards", 0, nil)

	var cs Cards

	err := datastore.Get(context, cardsKey, &cs)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			log.Errorf(context, "Cards data not found %v", err)
			cs2, err := newCards(context)
			if err != nil {
				sendFailResponse(c, "Can't create Cards")
				return
			}

			cardsKey := datastore.NewKey(context, "Cards", "cards", 0, nil)
			if _, err = datastore.Put(context, cardsKey, cs2); err != nil {
				// Handle err
				log.Errorf(context, "Erro while storing Cards: %v", err)
				sendFailResponse(c, "Can't store cards")
				return
			}
			cs = *cs2

		} else {
			log.Errorf(context, "Cards error: %v", err)
			sendFailResponse(c, "Error while retrieving Cards")
			return
		}
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	randGenerator := rand.New(randSource)
	for i := 0; i < len(cs.CardList); i++ {
		idx := randGenerator.Intn(len(cs.CardList))
		cd := cs.CardList[idx]
		if cd.learned == false {
			sendSuccessResponse(c, cd)
			return
		}
	}

	sendFailResponse(c, "No more card to learn")
}
