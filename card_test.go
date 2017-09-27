package main

import (
	"testing"
)


func TestSortRank(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "D"},
		Card{CLUBS, "8"},
	}

	sr := sortRank(cards)

	if len(sr) != len(cards) {
		t.Errorf("ERROR IN SORTRANK")
	}
}



func TestSortRank2(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{SPADE, "A"},
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{CARO, "K"},
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
		Card{SPADE, "A"},
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{CARO, "K"},
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "10"},
	}

	shCards := Shuffle(cards)

	sr := sortRankSpecial(shCards, []string{"A", "K", "D", "J", "10" ,"9", "8", "7"})

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
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{SPADE, "8"},
	}

	cardToRemove := Card{CLUBS, "J"}
	newhand := remove(hand, cardToRemove)
	if in(newhand, cardToRemove) {
		t.Errorf("First Card not removed")
	}

}

func TestRemove2(t *testing.T) {
	hand := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{SPADE, "8"},
	}

	cardToRemove := Card{SPADE, "8"}
	newhand := remove(hand, cardToRemove)
	if in(newhand, cardToRemove) {
		t.Errorf("Last Card not removed")
	}
}

func TestRemove3(t *testing.T) {
	hand := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{SPADE, "8"},
	}

	i := 0
	debugTacticsLog("cs[:i] = %v, cs[i+1:] = %v\n", hand[:i], hand[i+1:])

	debugTacticsLog("List: %v\n",hand)
	cardToRemove := Card{CLUBS, "J"}
	newhand := remove(hand, cardToRemove)
	debugTacticsLog("Remove first: %v %v\n",hand, newhand)
	if in(newhand, cardToRemove) {
		t.Errorf("First Card not removed")
	}

	cardToRemove = Card{SPADE, "8"}
	newhand2 := remove(newhand, cardToRemove)
	debugTacticsLog("Remove last: %v %v %v\n",hand, newhand, newhand2)
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
		Card{CLUBS, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "K"},
		Card{CLUBS, "D"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
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
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{CLUBS, "10"},
		Card{CLUBS, "D"},
		Card{CLUBS, "7"},
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{HEART, "9"},
		Card{HEART, "7"},
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
	if in(equiv, Card{SPADE, "J"}) || in(equiv, Card{SPADE, "10"})  {
		t.Errorf("Wrong similar cards %v", equiv)
	}
}

func TestSimilar3(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{CARO, "J"},
		Card{CLUBS, "10"},
		Card{CLUBS, "D"},
		Card{CLUBS, "7"},
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{HEART, "9"},
		Card{HEART, "7"},
	}

	s := makeSuitState()
	s.trump = CLUBS
	s.cardsPlayed = []Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{HEART, "8"},
	}

	equiv := similar(&s, cards)

	if len(equiv) != 7 {
		t.Errorf("Wrong similar cards %v", equiv)
	}
	if !in(equiv, Card{CLUBS, "J"}, Card{CLUBS, "10"}, Card{CLUBS, "D"}, Card{CLUBS, "7"}, Card{SPADE, "A"}, Card{SPADE, "K"}, Card{HEART, "9"}) {
		t.Errorf("Wrong similar cards %v", equiv)
	}
	if in(equiv, Card{CARO, "J"}) || in(equiv, Card{SPADE, "10"}, Card{HEART, "7"})  {
		t.Errorf("Wrong similar cards %v", equiv)
	}
}

func TestSimilar4(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
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
	if in(equiv, Card{CLUBS, "8"}) || in(equiv, Card{CLUBS, "7"})  {
		t.Errorf("Wrong similar cards %v", equiv)
	}
}

func TestEquivalent(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
		Card{CLUBS, "10"},
		Card{CLUBS, "K"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
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
	if in(equiv, Card{CLUBS, "9"}, Card{CLUBS, "8"})  {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}
	if in(equiv, Card{CLUBS, "9"}, Card{CLUBS, "7"})  {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}
	if in(equiv, Card{CLUBS, "8"}, Card{CLUBS, "7"})  {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}

	s.cardsPlayed = []Card{Card{SPADE, "J"}}
	equiv = equivalent(&s, cards)

	if len(equiv) != 4 {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}
	if in(equiv, Card{CLUBS, "J"}, Card{HEART, "J"}) {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}

	s.cardsPlayed = []Card{Card{SPADE, "J"}}
	s.trick = []Card{Card{SPADE, "J"}}
	equiv = equivalent(&s, cards)

	if len(equiv) != 5 {
		t.Errorf("Wrong equivalent cards %v", equiv)
	}

	debugTacticsLog("Equivalent: %v, %v\n", cards, equiv)

}