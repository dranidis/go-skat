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
