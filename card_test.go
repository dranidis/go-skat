package main

import (
	"testing"
)

func TestSortRank(t *testing.T) {
	cards := []Card{
		{CLUBS, "J"},
		{CLUBS, "A"},
		{CLUBS, "D"},
		{CLUBS, "8"},
	}

	sr := sortRank(cards)

	if len(sr) != len(cards) {
		t.Errorf("ERROR IN SORTRANK")
	}
}

func TestSortRank2(t *testing.T) {
	cards := []Card{
		{CLUBS, "J"},
		{SPADE, "J"},
		{SPADE, "A"},
		{HEART, "A"},
		{HEART, "K"},
		{CARO, "K"},
	}

	shCards := Shuffle(cards)

	sr := sortRank(shCards)

	if len(sr) != len(cards) {
		t.Errorf("ERROR IN SORTRANK")
	}
	for i, c := range sr {
		if !c.equals(cards[i]) {
			t.Errorf("Wrong ordering: %d, %v, %v", i, c, cards[i])
		}
	}
}

func TestSortRankSpecial1(t *testing.T) {
	cards := []Card{
		{SPADE, "A"},
		{HEART, "A"},
		{HEART, "K"},
		{CARO, "K"},
		{CLUBS, "J"},
		{SPADE, "J"},
		{HEART, "10"},
	}

	shCards := Shuffle(cards)

	sr := sortRankSpecial(shCards, []string{"A", "K", "D", "J", "10", "9", "8", "7"})

	if len(sr) != len(cards) {
		t.Errorf("ERROR IN SORTRANKspecial")
	}
	for i, c := range sr {
		if !c.equals(cards[i]) {
			t.Errorf("Wrong ordering: %d, %v, %v", i, c, cards[i])
		}
	}
}

func TestRemove1(t *testing.T) {
	hand := []Card{
		{CLUBS, "J"},
		{SPADE, "J"},
		{HEART, "J"},
		{CARO, "J"},
		{CARO, "A"},
		{CARO, "10"},
		{CARO, "K"},
		{HEART, "D"},
		{HEART, "9"},
		{SPADE, "8"},
	}

	cardToRemove := Card{CLUBS, "J"}
	newhand := remove(hand, cardToRemove)
	if in(newhand, cardToRemove) {
		t.Errorf("First Card not removed")
	}

}

func TestRemove2(t *testing.T) {
	hand := []Card{
		{CLUBS, "J"},
		{SPADE, "J"},
		{HEART, "J"},
		{CARO, "J"},
		{CARO, "A"},
		{CARO, "10"},
		{CARO, "K"},
		{HEART, "D"},
		{HEART, "9"},
		{SPADE, "8"},
	}

	cardToRemove := Card{SPADE, "8"}
	newhand := remove(hand, cardToRemove)
	if in(newhand, cardToRemove) {
		t.Errorf("Last Card not removed")
	}
}

func TestRemove3(t *testing.T) {
	hand := []Card{
		{CLUBS, "J"},
		{SPADE, "J"},
		{HEART, "J"},
		{CARO, "J"},
		{CARO, "A"},
		{CARO, "10"},
		{CARO, "K"},
		{HEART, "D"},
		{HEART, "9"},
		{SPADE, "8"},
	}

	i := 0
	debugTacticsLog("cs[:i] = %v, cs[i+1:] = %v\n", hand[:i], hand[i+1:])

	debugTacticsLog("List: %v\n", hand)
	cardToRemove := Card{CLUBS, "J"}
	newhand := remove(hand, cardToRemove)
	debugTacticsLog("Remove first: %v %v\n", hand, newhand)
	if in(newhand, cardToRemove) {
		t.Errorf("First Card not removed")
	}

	cardToRemove = Card{SPADE, "8"}
	newhand2 := remove(newhand, cardToRemove)
	debugTacticsLog("Remove last: %v %v %v\n", hand, newhand, newhand2)
	if in(newhand2, cardToRemove) {
		t.Errorf("Last Card not removed")
	}
}

func TestNextCard(t *testing.T) {
	p := Card{CLUBS, "J"}
	c := nextCard(CARO, Card{CLUBS, "J"})
	debugTacticsLog("next of %v is %v\n", p, c)
	if !c.equals(Card{SPADE, "J"}) {
		t.Errorf("Wrong next card: %v", c)
	}
}

func TestSimilar1(t *testing.T) {
	cards := []Card{
		{CLUBS, "J"},
		{HEART, "J"},
		{CARO, "J"},
		{CLUBS, "A"},
		{CLUBS, "10"},
		{CLUBS, "K"},
		{CLUBS, "D"},
		{CLUBS, "8"},
		{CLUBS, "7"},
	}
	s := makeSuitState()
	s.trump = CLUBS

	equiv := similar(&s, cards)

	if len(equiv) != 5 {
		t.Errorf("Wrong similar cards %v", equiv)
	}
	if !in(equiv, Card{CLUBS, "J"}, Card{HEART, "J"}, Card{CLUBS, "A"}, Card{CLUBS, "K"}, Card{CLUBS, "8"}) {
		t.Errorf("Wrong similar cards %v", equiv)
	}
	if in(equiv, Card{CARO, "J"}) || in(equiv, Card{CLUBS, "10"}) || in(equiv, Card{CLUBS, "D"}) || in(equiv, Card{CLUBS, "7"}) {
		t.Errorf("Wrong similar cards %v", equiv)
	}
}

func TestSimilar2(t *testing.T) {
	cards := []Card{
		{CLUBS, "J"},
		{SPADE, "J"},
		{CLUBS, "10"},
		{CLUBS, "D"},
		{CLUBS, "7"},
		{SPADE, "A"},
		{SPADE, "10"},
		{SPADE, "K"},
		{HEART, "9"},
		{HEART, "7"},
	}

	s := makeSuitState()
	s.trump = CLUBS

	equiv := similar(&s, cards)

	if len(equiv) != 8 {
		t.Errorf("Wrong similar cards %v", equiv)
	}
	if !in(equiv, Card{CLUBS, "J"}, Card{CLUBS, "10"}, Card{CLUBS, "D"}, Card{CLUBS, "7"}, Card{SPADE, "A"}, Card{SPADE, "K"}, Card{HEART, "9"}, Card{HEART, "7"}) {
		t.Errorf("Wrong similar cards %v", equiv)
	}
	if in(equiv, Card{SPADE, "J"}) || in(equiv, Card{SPADE, "10"}) {
		t.Errorf("Wrong similar cards %v", equiv)
	}
}

func TestSimilar3(t *testing.T) {
	cards := []Card{
		{CLUBS, "J"},
		{CARO, "J"},
		{CLUBS, "10"},
		{CLUBS, "D"},
		{CLUBS, "7"},
		{SPADE, "A"},
		{SPADE, "10"},
		{SPADE, "K"},
		{HEART, "9"},
		{HEART, "7"},
	}

	s := makeSuitState()
	s.trump = CLUBS
	s.cardsPlayed = []Card{
		{SPADE, "J"},
		{HEART, "J"},
		{HEART, "8"},
	}

	equiv := similar(&s, cards)

	if len(equiv) != 7 {
		t.Errorf("Wrong similar cards %v", equiv)
	}
	if !in(equiv, Card{CLUBS, "J"}, Card{CLUBS, "10"}, Card{CLUBS, "D"}, Card{CLUBS, "7"}, Card{SPADE, "A"}, Card{SPADE, "K"}, Card{HEART, "9"}) {
		t.Errorf("Wrong similar cards %v", equiv)
	}
	if in(equiv, Card{CARO, "J"}) || in(equiv, Card{SPADE, "10"}, Card{HEART, "7"}) {
		t.Errorf("Wrong similar cards %v", equiv)
	}
}

func TestSimilar4(t *testing.T) {
	cards := []Card{
		{CLUBS, "J"},
		{CLUBS, "9"},
		{CLUBS, "8"},
		{CLUBS, "7"},
	}

	s := makeSuitState()
	s.trump = CLUBS
	s.cardsPlayed = []Card{}

	equiv := similar(&s, cards)

	if len(equiv) != 2 {
		t.Errorf("Wrong similar cards %v", equiv)
	}
	if !in(equiv, Card{CLUBS, "J"}, Card{CLUBS, "9"}) {
		t.Errorf("Wrong similar cards %v", equiv)
	}
	if in(equiv, Card{CLUBS, "8"}) || in(equiv, Card{CLUBS, "7"}) {
		t.Errorf("Wrong similar cards %v", equiv)
	}
}

func TestEquivalent(t *testing.T) {
	cards := []Card{
		{CLUBS, "J"},
		{HEART, "J"},
		{CLUBS, "10"},
		{CLUBS, "K"},
		{CLUBS, "9"},
		{CLUBS, "8"},
		{CLUBS, "7"},
	}

	s := makeSuitState()
	s.trump = CLUBS
	s.cardsPlayed = []Card{}
	s.trick = []Card{}

	equiv := equivalent(&s, cards)

	if len(equiv) != 5 {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}
	if !in(equiv, Card{CLUBS, "J"}, Card{HEART, "J"}) {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}
	if in(equiv, Card{CLUBS, "9"}, Card{CLUBS, "8"}) {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}
	if in(equiv, Card{CLUBS, "9"}, Card{CLUBS, "7"}) {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}
	if in(equiv, Card{CLUBS, "8"}, Card{CLUBS, "7"}) {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}

	s.cardsPlayed = []Card{{SPADE, "J"}}
	equiv = equivalent(&s, cards)

	if len(equiv) != 4 {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}
	if in(equiv, Card{CLUBS, "J"}, Card{HEART, "J"}) {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}

	s.cardsPlayed = []Card{{SPADE, "J"}}
	s.trick = []Card{{SPADE, "J"}}
	equiv = equivalent(&s, cards)

	if len(equiv) != 5 {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}

	debugTacticsLog("Equivalent: %v, %v\n", cards, equiv)

}

//probabilities
// func TestZuDrittProb(t *testing.T) {
// 	cards := []Card{
// 		Card{"CARO", "K"},
// 		Card{"CARO", "9"},
// 		Card{"CARO", "8"},
// 		Card{"CARO", "7"},
// 		Card{"SPADE", "7"},
// 		Card{"SPADE", "8"},
// 		Card{"SPADE", "9"},
// 		Card{"SPADE", "10"},
// 	}

// 	zuDritt := func (cs []Card) bool {
// 		caros := filter(cs, func(c Card) bool {
// 			return c.Suit == "CARO"
// 			})
// 		if len(caros) >= 3 && in(caros, Card{"CARO", "K"}) {
// 			return true
// 		}
// 		return false
// 	}
// 	zt := 0
// 	tot := 100
// 	for i := 0 ; i < tot; i++ {
// 		cards := Shuffle(cards)
// 		half1 := cards[0:4]
// 		half2 := cards[4:]
// 		// debugTacticsLog("%v %v", half1, half2)
// 		if zuDritt(half1) || zuDritt(half2) {
// 			//debugTacticsLog("\tZT\n")
// 			zt++
// 		} else {
// 			//debugTacticsLog("\n")
// 		}
// 	}
// 	debugTacticsLog("%d %d\n", zt, zt * 100/tot)

// }
