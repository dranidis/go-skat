package main

import (
	"fmt"
	"testing"
)

func TestGreater(t *testing.T) {
	clubsHeart := SuitState{CLUBS, HEART}
	greaterAux(t, clubsHeart, Card{HEART, "A"}, Card{HEART, "10"})
	greaterAux(t, clubsHeart, Card{HEART, "A"}, Card{HEART, "K"})

	greaterAux(t, clubsHeart, Card{CLUBS, "7"}, Card{HEART, "A"})
	greaterAux(t, clubsHeart, Card{CLUBS, "10"}, Card{CLUBS, "9"})
	greaterAux(t, clubsHeart, Card{CLUBS, "J"}, Card{CLUBS, "A"})
	greaterAux(t, clubsHeart, Card{CLUBS, "8"}, Card{CLUBS, "7"})
	greaterAux(t, clubsHeart, Card{CLUBS, "J"}, Card{SPADE, "J"})
	greaterAux(t, clubsHeart, Card{SPADE, "J"}, Card{HEART, "J"})

	greaterAux(t, clubsHeart, Card{HEART, "7"}, Card{CARO, "A"})
	greaterAux(t, clubsHeart, Card{CLUBS, "A"}, Card{HEART, "7"})
	greaterAux(t, clubsHeart, Card{SPADE, "A"}, Card{SPADE, "7"})
}

func greaterAux(t *testing.T, s SuitState, card1, card2 Card) {
	if !s.greater(card1, card2) {
		t.Errorf("TRUMP :" + s.trump + " FOLLOW :" + s.follow + " - " + card1.suit + " " + card1.rank + " should be greater than " + card2.suit + " " + card2.rank)
	}
	if s.greater(card2, card1) {
		t.Errorf("TRUMP :" + s.trump + " FOLLOW :" + s.follow + " - " + card2.suit + " " + card2.rank + " should NOT be greater than " + card1.suit + " " + card1.rank)
	}
}

func TestValidPlay(t *testing.T) {
	clubsHeart := SuitState{CLUBS, HEART}

	playerHand := []Card{Card{SPADE, "J"}, Card{HEART, "A"}, Card{CARO, "7"}}
	validAux(t, clubsHeart, playerHand, Card{HEART, "A"})
	notValidAux(t, clubsHeart, playerHand, Card{CARO, "7"})
	notValidAux(t, clubsHeart, playerHand, Card{SPADE, "J"})

	validAux(t, SuitState{CLUBS, SPADE}, playerHand, Card{HEART, "A"})
	validAux(t, SuitState{CLUBS, SPADE}, playerHand, Card{CARO, "7"})
	validAux(t, SuitState{CLUBS, SPADE}, playerHand, Card{SPADE, "J"})

	notValidAux(t, SuitState{CLUBS, CLUBS}, playerHand, Card{HEART, "A"})
	notValidAux(t, SuitState{CLUBS, CLUBS}, playerHand, Card{CARO, "7"})
	validAux(t, SuitState{CLUBS, CLUBS}, playerHand, Card{SPADE, "J"})

	validAux(t, SuitState{CLUBS, ""}, playerHand, Card{HEART, "A"})
	validAux(t, SuitState{CLUBS, ""}, playerHand, Card{CARO, "7"})
	validAux(t, SuitState{CLUBS, ""}, playerHand, Card{SPADE, "J"})

}

func validAux(t *testing.T, s SuitState, cards []Card, card Card) {
	if !s.valid(cards, card) {
		t.Errorf("TRUMP :" + s.trump + " FOLLOW :" + s.follow + " - " + fmt.Sprintf("%v", card) + " should be valid play. HAND:" + fmt.Sprintf("%v", cards))
	}
}

func notValidAux(t *testing.T, s SuitState, cards []Card, card Card) {
	if s.valid(cards, card) {
		t.Errorf("TRUMP :" + s.trump + " FOLLOW :" + s.follow + " - " + fmt.Sprintf("%v", card) + " should NOT be valid play. HAND:" + fmt.Sprintf("%v", cards))
	}
}

func TestValidCards(t *testing.T) {
	clubsHeart := SuitState{CLUBS, HEART}

	playerHand := []Card{Card{SPADE, "J"}, Card{HEART, "A"}, Card{CARO, "7"}}
	cards := validCards(clubsHeart, playerHand)
	compareLists(t, cards, []Card{Card{HEART, "A"}})

	cards = validCards(SuitState{CLUBS, SPADE}, playerHand)
	compareLists(t, cards, playerHand)
}

func compareLists(t *testing.T, returned, expected []Card) {
	if len(returned) != len(expected) {
		t.Errorf("Expected: " + fmt.Sprintf("%v", expected) + " found: " + fmt.Sprintf("%v", returned))
	}
	for i, c := range expected {
		if c.suit != returned[i].suit || c.rank != returned[i].rank {
			t.Errorf("Expected: " + fmt.Sprintf("%v", expected) + " found: " + fmt.Sprintf("%v", returned))
		}
	}
}
