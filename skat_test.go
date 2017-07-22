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
	greaterAux(t, clubsHeart, Card{HEART, "J"}, Card{CARO, "J"})
	greaterAux(t, clubsHeart, Card{CARO, "J"}, Card{CLUBS, "A"})

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

func TestGetSuite(t *testing.T) {
	s := SuitState{CLUBS, HEART}
	card := Card{SPADE, "J"}
	if getSuite(s, card) != CLUBS {
		t.Errorf("TRUMP :" + s.trump + " FOLLOW :" + s.follow + " - " + fmt.Sprintf("%v", card) + " should be CLUBS")
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

func TestSetNextTrickOrder(t *testing.T) {
	state := SuitState{CLUBS, ""}
	trick := []Card{Card{SPADE, "J"}, Card{HEART, "A"}, Card{HEART, "K"}}

	player1 := makePlayer([]Card{})
	player2 := makePlayer([]Card{})
	player3 := makePlayer([]Card{})

	players := []*Player{&player1, &player2, &player3}
	newPlayers := setNextTrickOrder(&state, players, trick)
	comparePlayers(t, players, newPlayers)
	checkScore(t, player1, 17)
	checkScore(t, player2, 0)
	checkScore(t, player3, 0)

	trick = []Card{Card{HEART, "D"}, Card{SPADE, "J"}, Card{HEART, "10"}}
	newPlayers = setNextTrickOrder(&state, players, trick)
	expected := []*Player{&player2, &player3, &player1}
	comparePlayers(t, expected, newPlayers)
	checkScore(t, player1, 17)
	checkScore(t, player2, 15)
	checkScore(t, player3, 0)

	trick = []Card{Card{HEART, "9"}, Card{HEART, "8"}, Card{SPADE, "J"}}
	newPlayers = setNextTrickOrder(&state, players, trick)
	expected = []*Player{&player3, &player1, &player2}
	comparePlayers(t, expected, newPlayers)
	checkScore(t, player1, 17)
	checkScore(t, player2, 15)
	checkScore(t, player3, 2)
}

func checkScore(t *testing.T, p Player, s int) {
	if p.score != s {
		t.Errorf("Expected score %d, got %d", s, p.score)
	}
}

func comparePlayers(t *testing.T, expected, returned []*Player) {
	wrong := false
	for i, p := range expected {
		if p != returned[i] {
			wrong = true
		}
	}
	if wrong {
		t.Errorf("Wrong order of players: Expected: %v, Got: %v", expected, returned)
	}
}

func TestTrick(t *testing.T) {
	state := SuitState{CLUBS, ""}
	firstPlayerHand := []Card{Card{SPADE, "J"}, Card{HEART, "D"}, Card{CARO, "7"}}
	secondPlayerHand := []Card{Card{CLUBS, "J"}, Card{HEART, "10"}, Card{CARO, "8"}}
	thirdPlayerHand := []Card{Card{HEART, "A"}, Card{HEART, "K"}, Card{CLUBS, "10"}}

	player1 := makePlayer(firstPlayerHand)
	player2 := makePlayer(secondPlayerHand)
	player3 := makePlayer(thirdPlayerHand)

	players := []*Player{&player1, &player2, &player3}
	players = round(&state, players)

	if len(player1.hand) != 2 {
		t.Errorf("Expected: player1 len hand 2: " + fmt.Sprintf("%v", player1.hand))
	}
	if len(player2.hand) != 2 {
		t.Errorf("Expected: player2 len hand 2" + fmt.Sprintf("%v", player2.hand))
	}
	if len(player3.hand) != 2 {
		t.Errorf("Expected: player3 len hand 2" + fmt.Sprintf("%v", player3.hand))
	}

	expected := []*Player{&player2, &player3, &player1}
	comparePlayers(t, expected, players)
	checkScore(t, player1, 0)
	checkScore(t, player2, 14)
	checkScore(t, player3, 0)

	players = round(&state, players)

	expected = []*Player{&player3, &player1, &player2}
	comparePlayers(t, expected, players)
	checkScore(t, player1, 0)
	checkScore(t, player2, 14)
	checkScore(t, player3, 24)

	players = round(&state, players)

	expected = []*Player{&player3, &player1, &player2}
	comparePlayers(t, expected, players)
	checkScore(t, player1, 0)
	checkScore(t, player2, 14)
	checkScore(t, player3, 28)
}

func TestMakeDeck(t *testing.T) {
	cards := makeDeck()
	if len(cards) != 32 {
		t.Errorf("Not 32 cards")
	}
}

func TestMatadors(t *testing.T) {
	check := func(exp, act int) {
		if act != exp {
			t.Errorf("Expected %d matadors, got %d", exp, act)
		}
	}

	state := SuitState{CLUBS, ""}

	player := makePlayer([]Card{})
	check(-10, matadors(state, player))

	player.hand = append(player.hand, Card{CLUBS, "8"})
	check(-9, matadors(state, player))
	player.hand = append(player.hand, Card{CLUBS, "9"})
	check(-8, matadors(state, player))
	player.hand = append(player.hand, Card{CLUBS, "D"})
	check(-7, matadors(state, player))
	player.hand = append(player.hand, Card{CLUBS, "K"})
	check(-6, matadors(state, player))
	player.hand = append(player.hand, Card{CLUBS, "10"})
	check(-5, matadors(state, player))
	player.hand = append(player.hand, Card{CLUBS, "A"})
	check(-4, matadors(state, player))
	player.hand = append(player.hand, Card{CARO, "J"})
	check(-3, matadors(state, player))
	player.hand = append(player.hand, Card{HEART, "J"})
	check(-2, matadors(state, player))
	player.hand = append(player.hand, Card{SPADE, "J"})
	check(-1, matadors(state, player))

	player.hand = []Card{}

	cards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{state.trump, "A"},
		Card{state.trump, "10"},
		Card{state.trump, "K"},
		Card{state.trump, "D"},
		Card{state.trump, "9"},
		Card{state.trump, "8"},
	} 
	m := 0
	for _, card := range cards {
		player.hand = append(player.hand, card)
		m++
		check(m, matadors(state, player))		
	}	
}
