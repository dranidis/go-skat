package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"fmt"
	// "time"
	// "github.com/dranidis/go-skat/game"
	// "github.com/dranidis/go-skat/game/minimax"
	// "github.com/dranidis/go-skat/game/mcts"

)

func TestMain(m *testing.M) {
	debugTacticsLogFlag = true
	gameLogFlag = true
	delayMs = 0

	code := m.Run()
	os.Exit(code)
}

func mState(trump, follow string) SuitState {
	s := makeSuitState()
	s.trump = trump
	s.follow = follow
	return s
}

func TestGreater(t *testing.T) {
	clubsHeart := mState(CLUBS, HEART)
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
		t.Errorf("TRUMP : %s  FOLLOW: %s. %v  should be greater than %v", s.trump, s.follow, card1, card2)
	}
	if s.greater(card2, card1) {
		t.Errorf("TRUMP : %s  FOLLOW: %s. %v  should NOT be greater than %v", s.trump, s.follow, card1, card2)
	}
}

func TestGetSuite(t *testing.T) {
	s := mState(CLUBS, HEART)
	card := Card{SPADE, "J"}
	if getSuit(s.trump, card) != CLUBS {
		t.Errorf("TRUMP : %s  FOLLOW: %s. %v  should be CLUBS", s.trump, s.follow, card)
	}
}

func TestValidPlay(t *testing.T) {
	clubsHeart := mState(CLUBS, HEART)

	playerHand := []Card{Card{SPADE, "J"}, Card{HEART, "A"}, Card{CARO, "7"}}
	validAux(t, clubsHeart, playerHand, Card{HEART, "A"})
	notValidAux(t, clubsHeart, playerHand, Card{CARO, "7"})
	notValidAux(t, clubsHeart, playerHand, Card{SPADE, "J"})

	validAux(t, mState(CLUBS, SPADE), playerHand, Card{HEART, "A"})
	validAux(t, mState(CLUBS, SPADE), playerHand, Card{CARO, "7"})
	validAux(t, mState(CLUBS, SPADE), playerHand, Card{SPADE, "J"})

	notValidAux(t, mState(CLUBS, CLUBS), playerHand, Card{HEART, "A"})
	notValidAux(t, mState(CLUBS, CLUBS), playerHand, Card{CARO, "7"})
	validAux(t, mState(CLUBS, CLUBS), playerHand, Card{SPADE, "J"})

	validAux(t, mState(CLUBS, ""), playerHand, Card{HEART, "A"})
	validAux(t, mState(CLUBS, ""), playerHand, Card{CARO, "7"})
	validAux(t, mState(CLUBS, ""), playerHand, Card{SPADE, "J"})

}

func validAux(t *testing.T, s SuitState, cards []Card, card Card) {
	if !valid(s.follow, s.trump, cards, card) {
		t.Errorf("TRUMP : %s  FOLLOW: %s. %v should be valid play. HAND: %v", s.trump, s.follow, card, cards)
	}
}

func notValidAux(t *testing.T, s SuitState, cards []Card, card Card) {
	if valid(s.follow, s.trump, cards, card) {
		t.Errorf("TRUMP : %s  FOLLOW: %s. %v should NOT be valid play. HAND: %v", s.trump, s.follow, card, cards)
	}
}

func TestValidCards(t *testing.T) {
	clubsHeart := mState(CLUBS, HEART)

	playerHand := []Card{Card{SPADE, "J"}, Card{HEART, "A"}, Card{CARO, "7"}}
	cards := validCards(clubsHeart, playerHand)
	compareLists(t, cards, []Card{Card{HEART, "A"}})

	cards = validCards(mState(CLUBS, SPADE), playerHand)
	compareLists(t, cards, playerHand)
}

func compareLists(t *testing.T, returned, expected []Card) {
	if len(returned) != len(expected) {
		t.Errorf("Expected: %v, found %v", expected, returned)
	}
	for i, c := range expected {
		if !c.equals(returned[i]) {
			t.Errorf("Expected: %v, found %v", expected, returned)
		}
	}
}

func TestSetNextTrickOrder(t *testing.T) {
	state := mState(CLUBS, "")
	state.trick = []Card{Card{SPADE, "J"}, Card{HEART, "A"}, Card{HEART, "K"}}

	player1 := makePlayer([]Card{})
	player2 := makePlayer([]Card{})
	player3 := makePlayer([]Card{})

	players := []PlayerI{&player1, &player2, &player3}
	newPlayers := setNextTrickOrder(&state, players)
	comparePlayers(t, players, newPlayers)
	checkScore(t, player1, 17)
	checkScore(t, player2, 0)
	checkScore(t, player3, 0)

	if player1.schwarz {
		t.Errorf("OUT OF SCHWARZ")
	}

	state.trick = []Card{Card{HEART, "D"}, Card{SPADE, "J"}, Card{HEART, "10"}}
	newPlayers = setNextTrickOrder(&state, players)
	expected := []PlayerI{&player2, &player3, &player1}
	comparePlayers(t, expected, newPlayers)
	checkScore(t, player1, 17)
	checkScore(t, player2, 15)
	checkScore(t, player3, 0)
	if player2.schwarz {
		t.Errorf("OUT OF SCHWARZ")
	}

	state.trick = []Card{Card{HEART, "9"}, Card{HEART, "8"}, Card{SPADE, "J"}}
	newPlayers = setNextTrickOrder(&state, players)
	expected = []PlayerI{&player3, &player1, &player2}
	comparePlayers(t, expected, newPlayers)
	checkScore(t, player1, 17)
	checkScore(t, player2, 15)
	checkScore(t, player3, 2)
	if player3.schwarz {
		t.Errorf("OUT OF SCHWARZ")
	}

}

func checkScore(t *testing.T, p Player, s int) {
	if p.score != s {
		t.Errorf("Expected score %d, got %d", s, p.score)
	}
}

func comparePlayers(t *testing.T, expected, returned []PlayerI) {
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

// based on firstCardTactic
// func TestTrick(t *testing.T) {
// 	state := mState(CLUBS, "")
// 	firstPlayerHand := []Card{Card{SPADE, "J"}, Card{HEART, "D"}, Card{CARO, "7"}}
// 	secondPlayerHand := []Card{Card{CLUBS, "J"}, Card{HEART, "10"}, Card{CARO, "8"}}
// 	thirdPlayerHand := []Card{Card{HEART, "A"}, Card{HEART, "K"}, Card{CLUBS, "10"}}

// 	player1 := makePlayer(firstPlayerHand)
// 	player2 := makePlayer(secondPlayerHand)
// 	player3 := makePlayer(thirdPlayerHand)

// 	players := []*Player{&player1, &player2, &player3}
// 	players = round(&state, players)

// 	if len(player1.hand) != 2 {
// 		t.Errorf("Expected: player1 len hand 2: "  fmt.Sprintf("%v", player1.hand))
// 	}
// 	if len(player2.hand) != 2 {
// 		t.Errorf("Expected: player2 len hand 2"  fmt.Sprintf("%v", player2.hand))
// 	}
// 	if len(player3.hand) != 2 {
// 		t.Errorf("Expected: player3 len hand 2"  fmt.Sprintf("%v", player3.hand))
// 	}

// 	expected := []*Player{&player2, &player3, &player1}
// 	comparePlayers(t, expected, players)
// 	checkScore(t, player1, 0)
// 	checkScore(t, player2, 14)
// 	checkScore(t, player3, 0)

// 	players = round(&state, players)

// 	expected = []*Player{&player3, &player1, &player2}
// 	comparePlayers(t, expected, players)
// 	checkScore(t, player1, 0)
// 	checkScore(t, player2, 14)
// 	checkScore(t, player3, 24)

// 	players = round(&state, players)

// 	expected = []*Player{&player3, &player1, &player2}
// 	comparePlayers(t, expected, players)
// 	checkScore(t, player1, 0)
// 	checkScore(t, player2, 14)
// 	checkScore(t, player3, 28)
// }

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

	//state := SuitState{CLUBS, ""}
	state := CLUBS

	player := makePlayer([]Card{})
	check(-10, matadors(state, player.hand))

	player.hand = append(player.hand, Card{CLUBS, "8"})
	check(-9, matadors(state, player.hand))
	player.hand = append(player.hand, Card{CLUBS, "9"})
	check(-8, matadors(state, player.hand))
	player.hand = append(player.hand, Card{CLUBS, "D"})
	check(-7, matadors(state, player.hand))
	player.hand = append(player.hand, Card{CLUBS, "K"})
	check(-6, matadors(state, player.hand))
	player.hand = append(player.hand, Card{CLUBS, "10"})
	check(-5, matadors(state, player.hand))
	player.hand = append(player.hand, Card{CLUBS, "A"})
	check(-4, matadors(state, player.hand))
	player.hand = append(player.hand, Card{CARO, "J"})
	check(-3, matadors(state, player.hand))
	player.hand = append(player.hand, Card{HEART, "J"})
	check(-2, matadors(state, player.hand))
	player.hand = append(player.hand, Card{SPADE, "J"})
	check(-1, matadors(state, player.hand))

	player.hand = []Card{}

	cards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{state, "A"},
		Card{state, "10"},
		Card{state, "K"},
		Card{state, "D"},
		Card{state, "9"},
		Card{state, "8"},
	}
	m := 0
	for _, card := range cards {
		player.hand = append(player.hand, card)
		m++
		check(m, matadors(state, player.hand))
	}
}

func TestBidding(t *testing.T) {
	makeP := func(high int) Player {
		player := makePlayer([]Card{})
		player.highestBid = high
		return player
	}
	player1 := makeP(0)
	player2 := makeP(18)
	player3 := makeP(0)
	players := []PlayerI{&player1, &player2, &player3}

	bidIndex, declarer := bid(players)
	if bidIndex != bids[0] {
		t.Errorf("Expected %d, Got %d", 18, bidIndex)
	}
	if declarer != &player2 {
		t.Errorf("Wrong declarer")
	}

	/*
		scenario 2
	*/
	player1 = makeP(23)
	player2 = makeP(20)
	player3 = makeP(24)
	players = []PlayerI{&player1, &player2, &player3}

	bidIndex, declarer = bid(players)
	if bidIndex != 24 {
		t.Errorf("Expected %d, Got %d", 24, bidIndex)
	}
	if declarer != &player3 {
		t.Errorf("Wrong declarer")
	}

	/*
		scenario 3
	*/
	player1 = makeP(18)
	player2 = makeP(0)
	player3 = makeP(18)
	players = []PlayerI{&player1, &player2, &player3}

	bidIndex, declarer = bid(players)
	if bidIndex != 18 {
		t.Errorf("Expected %d, Got %d", 18, bidIndex)
	}
	if declarer != &player1 {
		t.Errorf("Wrong declarer")
	}

	/*
		scenario 4
	*/
	player1 = makeP(0)
	player2 = makeP(18)
	player3 = makeP(20)
	players = []PlayerI{&player1, &player2, &player3}

	bidIndex, declarer = bid(players)
	if bidIndex != 20 {
		t.Errorf("Expected %d, Got %d", 20, bidIndex)
	}
	if declarer != &player3 {
		t.Errorf("Wrong declarer")
	}

	/*
		scenario 5
	*/
	player1 = makeP(0)
	player2 = makeP(0)
	player3 = makeP(0)
	players = []PlayerI{&player1, &player2, &player3}

	bidIndex, declarer = bid(players)
	if bidIndex != 0 {
		t.Errorf("Expected %d, Got %d", 0, bidIndex)
	}
	if declarer != nil {
		t.Errorf("Wrong declarer. Everybody passed")
	}

	/*
		scenario 6
	*/
	player1 = makeP(18)
	player2 = makeP(0)
	player3 = makeP(0)
	players = []PlayerI{&player1, &player2, &player3}

	bidIndex, declarer = bid(players)
	if bidIndex != 18 {
		t.Errorf("Expected %d, Got %d", 18, bidIndex)
	}
	if declarer != &player1 {
		t.Errorf("Wrong declarer")
	}

	/*
		scenario 7
	*/
	player1 = makeP(0)
	player2 = makeP(0)
	player3 = makeP(18)
	players = []PlayerI{&player1, &player2, &player3}

	bidIndex, declarer = bid(players)
	if bidIndex != 18 {
		t.Errorf("Expected %d, Got %d", 18, bidIndex)
	}
	if declarer != &player3 {
		t.Errorf("Wrong declarer")
	}

}

func TestMostCardsSuit(t *testing.T) {

	player := makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{HEART, "K"},
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{SPADE, "8"},
	})

	act := len(trumpCards(CLUBS, player.hand))
	if act != 4 {
		t.Errorf("Expected %d, got %d", 4, act)
	}
	act = len(trumpCards(SPADE, player.hand))
	if act != 5 {
		t.Errorf("Expected %d, got %d", 5, act)
	}
	act = len(trumpCards(CARO, player.hand))
	if act != 6 {
		t.Errorf("Expected %d, got %d", 6, act)
	}
	act = len(trumpCards(HEART, player.hand))
	if act != 7 {
		t.Errorf("Expected %d, got %d", 7, act)
	}

	most := mostCardsSuit(player.hand)
	if most != HEART {
		t.Errorf("Expected %s, got %s", HEART, most)
	}

	player = makePlayer([]Card{
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
	})
	most = mostCardsSuit(player.hand)
	if most != CARO {
		t.Errorf("Expected %s, got %s", CARO, most)
	}

	player = makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{SPADE, "8"},
	})
	most = mostCardsSuit(player.hand)
	if most != SPADE {
		t.Errorf("Expected %s, got %s", SPADE, most)
	}
	player = makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{SPADE, "10"},
		Card{CLUBS, "K"},
		Card{CLUBS, "D"},
		Card{HEART, "9"},
		Card{CLUBS, "8"},
	})
	most = mostCardsSuit(player.hand)
	if most != CLUBS {
		t.Errorf("Expected %s, got %s", CLUBS, most)
	}
}

func TestMostCardsSuitA(t *testing.T) {
	// if two suits have the same length, choose the non-A suit
	player := makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},

		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{HEART, "8"},
	})
	most := mostCardsSuit(player.hand)
	if most != HEART {
		t.Errorf("In hand %v, Two suits equal length. Expected %s, got %s", player.hand, HEART, most)
	}
}

func TestMostCardsSuit1(t *testing.T) {
	// from two A-suits prefer the weakest
	player := makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},

		Card{CARO, "A"},
		Card{CARO, "7"},
		Card{CARO, "8"},
		Card{CARO, "9"},

		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "K"},
		Card{HEART, "8"},
	})
	most := mostCardsSuit(player.hand)
	exp := CARO
	if most != exp {
		t.Errorf("In hand %v, Two suits equal length. Expected %s, got %s", player.hand, exp, most)
	}

	player = makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},

		Card{HEART, "A"},
		Card{HEART, "7"},
		Card{HEART, "8"},
		Card{HEART, "9"},

		Card{CARO, "A"}, // keep as suit
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{CARO, "8"},
	})
	most = mostCardsSuit(player.hand)
	if most != HEART {
		t.Errorf("In hand %v, Two suits equal length. Expected %s, got %s", player.hand, HEART, most)
	}
}

func TestMostCardsSuitTwoAsuitsAndAWeakSuit(t *testing.T) {
	// from two A-suits prefer the weakest
	player := makePlayer([]Card{
		Card{CLUBS, "J"},

		Card{CARO, "A"},
		Card{CARO, "7"},
		Card{CARO, "8"},
		Card{CARO, "9"},

		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "K"},
		Card{HEART, "8"},

		Card{SPADE, "10"},
	})
	most := mostCardsSuit(player.hand)
	exp := CARO
	if most != exp {
		t.Errorf("In hand %v, Two suits equal length. Expected %s, got %s", player.hand, exp, most)
	}
}

func TestMostCardsSuitTwoAsuitsAndAWeakSuit2(t *testing.T) {
	// from two A-suits prefer the weakest
	player := makePlayer([]Card{
		Card{CLUBS, "J"},

		Card{CARO, "A"},
		Card{CARO, "7"},
		Card{CARO, "8"},

		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "K"},

		Card{SPADE, "10"},

		Card{CLUBS, "K"},
		Card{CLUBS, "D"},
	})
	most := mostCardsSuit(player.hand)
	exp := CARO
	if most != exp {
		t.Errorf("In hand %v, Two suits equal length. Expected %s, got %s", player.hand, exp, most)
	}
}


func gameScore1(trump string, cs []Card, s int, bid int, decS, oppS bool, hg bool) int {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})
	d.setScore(s)
	d.setSchwarz(decS)
	o1.setSchwarz(oppS)
	o2.setSchwarz(oppS)
	testState = makeSuitState()
	testState.trump = trump
	testState.declarer = &d
	testState.opp1 = &o1
	testState.opp2 = &o2
	testState.declarer.setDeclaredBid(bid)


	gs := gameScore(testState, cs, hg)
	return gs.GameScore
}

var testState SuitState

func TestGameScore(t *testing.T) {
	declarerCards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{HEART, "K"},
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{SPADE, "8"},
	}

	act := gameScore1(CARO, declarerCards, 61, 63, false, false, false)
	exp := 63
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

	act = gameScore1(CARO, declarerCards, 60, 63, false, false, false)
	exp = -126
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

	act = gameScore1(HEART, declarerCards, 61, 50, false, false, false)
	exp = 50
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

	act = gameScore1(CLUBS, declarerCards, 61, 50, false, false, false)
	exp = 60
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

	act = gameScore1(SPADE, declarerCards, 61, 50, false, false, false)
	exp = 55
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

	act = gameScore1(SPADE, declarerCards, 61, 50, false, false, true)
	exp = 66
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

	// hand is 50, OVERBID
	act = gameScore1(HEART, declarerCards, 61, 51, false, false, false)
	exp = -120
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

	declarerCards = []Card{
		//Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{HEART, "K"},
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{SPADE, "8"},
	}
	act = gameScore1(CARO, declarerCards, 61, 18, false, false, false)
	exp = 18
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

	// schneider winner
	act = gameScore1(CARO, declarerCards, 90, 18, false, false, false)
	exp = 27
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}
	// schneider loss
	act = gameScore1(CARO, declarerCards, 30, 18, false, false, false)
	exp = -54
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

	// schwarz winner
	act = gameScore1(CARO, declarerCards, 120, 18, false, true, false)
	exp = 36
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

	// schwarz loss
	act = gameScore1(CARO, declarerCards, 0, 18, true, false, false)
	exp = -72
	if act != exp {
		t.Errorf("Expected GAME SCORE %d, got %d", exp, act)
	}

}


func TestFindBlankCards(t *testing.T) {
	player := makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{SPADE, "10"},
	})
	cards := findBlankCards(player.hand)

	if len(cards) != 1 {
		t.Errorf("Expected 1. Found blank cards: %d", len(cards))
	} else {
		c := cards[0]
		if c.Suit != SPADE && c.Rank != "10" {
			t.Errorf("Found wrong blank card %v", c)
		}
	}

	player.hand = append(player.hand, Card{CLUBS, "7"})

	cards = findBlankCards(player.hand)

	if len(cards) != 2 {
		t.Errorf("Expected 2. Found blank cards: %d", len(cards))
	} else {
		c1 := cards[0]
		c2 := cards[1]
		if c1.Suit != SPADE && c1.Rank != "10" {
			t.Errorf("Blank Cards in wrong order %v, %v", c1, c2)
		}
	}
}

// func TestGame(t *testing.T) {
// 	// player1 is declared globally
// 	player := makePlayer([]Card{})
// 	player1 := &player

// 	player2 := makePlayer([]Card{})
// 	player3 := makePlayer([]Card{})
// 	player3.firstCardPlay = true
// 	player1.setName("NAME")
// 	if player1.getName() != "NAME" {
// 		t.Errorf("Error in set/get name")
// 	}
// 	if player1.getTotalScore() != 0 {
// 		t.Errorf("Error in get total score")
// 	}
// 	gamePlayers = []PlayerI{player1, &player2, &player3}
// 	for i := 0; i < 20; i++ {
// 		_ = skatGame()
// 	}
// }

func TestGame2(t *testing.T) {
	// player1 is declared globally
	player := makePlayer([]Card{})
	player1 = &player

	player2 := makeMinMaxPlayer([]Card{})
	player3 := makePlayer([]Card{})
	player3.firstCardPlay = true
	player1.setName("NAME")
	player2.setName("MINMAX")
	if player1.getName() != "NAME" {
		t.Errorf("Error in set/get name")
	}
	if player1.getTotalScore() != 0 {
		t.Errorf("Error in get total score")
	}
	gamePlayers = []PlayerI{player1, &player2, &player3}
	for i := 0; i < 10; i++ {
		// _ = skatGame()
	}
}

func TestLongestNonTrumpSuit(t *testing.T) {
	cards := []Card{
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{CARO, "7"},
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "K"},
		Card{SPADE, "8"},
		Card{SPADE, "7"},
	}
	suit := LongestNonTrumpSuit(CARO, cards)
	if suit == CARO {
		t.Errorf("CARO is the trump")
	}

	cards = []Card{
		Card{SPADE, "K"},
		Card{CLUBS, "10"},
		Card{CLUBS, "9"},
		Card{CLUBS, "7"},
		Card{HEART, "K"},
		Card{HEART, "8"},
		Card{HEART, "7"},
	}
	suit = LongestNonTrumpSuit(SPADE, cards)
	if suit == CARO {
		t.Errorf("CARO is not in the cards")
	}
}

func TestShortestNonTrumpSuit(t *testing.T) {
	cards := []Card{
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{CARO, "7"},
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "K"},
		Card{SPADE, "8"},
		Card{SPADE, "7"},
	}
	suit := ShortestNonTrumpSuit(CARO, cards)
	if suit == CARO {
		t.Errorf("CARO is the trump")
	}

	cards = []Card{
		Card{SPADE, "K"},
		Card{CLUBS, "10"},
		Card{CLUBS, "9"},
		Card{CLUBS, "7"},
		Card{HEART, "K"},
		Card{HEART, "8"},
		Card{HEART, "7"},
	}
	suit = ShortestNonTrumpSuit(CARO, cards)
	if suit == CARO {
		t.Errorf("CARO is not in the cards")
	}

	if suit != SPADE {
		t.Errorf("SPADE is the shortest: %v", cards)
	}
}

func TestHighestLong(t *testing.T) {
	cards := []Card{
		Card{SPADE, "K"},
		Card{CLUBS, "10"},
		Card{CLUBS, "9"},
		Card{CLUBS, "7"},
		Card{HEART, "K"},
		Card{HEART, "8"},
		Card{HEART, "7"},
	}

	card := HighestLong(SPADE, cards)

	if card.equals(cards[0]) {
		t.Errorf("Error in HighestLong")
	}
}

func TestHighestShort(t *testing.T) {
	cards := []Card{
		Card{SPADE, "K"},
		Card{SPADE, "10"},
		Card{SPADE, "9"},
		Card{HEART, "7"},
		Card{HEART, "K"},
		Card{HEART, "8"},
		Card{CARO, "10"},
		Card{CARO, "7"},
	}

	card := HighestShort(SPADE, cards)
	exp := Card{CARO, "10"}
	if !card.equals(exp) {
		t.Errorf("Error in HighestShort: %v, exp: %v, got %v", cards, exp, card)
	}
}


func TestOtherPlayersTrumps(t *testing.T) {
	player := makePlayer([]Card{})
	s := makeSuitState()

	s.trump = CLUBS
	player.hand = []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "A"},
		Card{HEART, "A"},
		Card{HEART, "7"},
	}

	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "A"},
	}

	other := player.otherPlayersTrumps(&s)
	//fmt.Printf("OTHER: %v\n", other)
	if len(other) != 0 {
		t.Errorf("No other trumps, since hand: %v and trumps in game %v",
			player.hand, s.trumpsInGame)
	}

	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "K"},
	}

	other = player.otherPlayersTrumps(&s)
	//fmt.Printf("OTHER: %v\n", other)

	if len(other) != 1 {
		t.Errorf("One more trump in game, since hand: %v and trumps in game %v",
			player.hand, s.trumpsInGame)
	}

}

func TestRotatePlayers(t *testing.T) {
	player1 := makePlayer([]Card{})
	player2 := makePlayer([]Card{})
	player3 := makePlayer([]Card{})
	players := []PlayerI{&player1, &player2, &player3}

	players = rotatePlayers(players)
	next := players[0]
	expNext := &player2

	if next != expNext {
		t.Errorf("Wrong order")
	}

	next = players[1]
	expNext = &player3

	if next != expNext {
		t.Errorf("Wrong order")
	}

	next = players[2]
	expNext = &player1

	if next != expNext {
		t.Errorf("Wrong order")
	}

	players = rotatePlayers(players)
	next = players[0]
	expNext = &player3

	if next != expNext {
		t.Errorf("Wrong order")
	}

	next = players[1]
	expNext = &player1

	if next != expNext {
		t.Errorf("Wrong order")
	}

	next = players[2]
	expNext = &player2

	if next != expNext {
		t.Errorf("Wrong order")
	}

	players = rotatePlayers(players)
	next = players[0]
	expNext = &player1

	if next != expNext {
		t.Errorf("Wrong order")
	}

	next = players[1]
	expNext = &player2

	if next != expNext {
		t.Errorf("Wrong order")
	}

	next = players[2]
	expNext = &player3

	if next != expNext {
		t.Errorf("Wrong order")
	}
}

func TestIsAKX(t *testing.T) {
	hand := []Card{
		Card{CARO, "A"},
		Card{CARO, "K"},
	}

	act := isAKX(CARO, hand)
	if act {
		t.Errorf("Hand is NOT AKX (too short): %v", hand)
	}

	hand = append(hand, Card{CARO, "8"})
	act = isAKX(CARO, hand)
	if !act {
		t.Errorf("Hand is AKX: %v", hand)
	}

	hand = append(hand, Card{CARO, "10"})
	act = isAKX(CARO, hand)
	if act {
		t.Errorf("Hand is NOT AKX: %v", hand)
	}
}

func TestInMany(t *testing.T) {
	cards := []Card{
		Card{CARO, "A"},
		Card{CARO, "K"},
		Card{CARO, "7"},
	}

	act := in(cards, Card{CARO, "K"}, Card{CARO, "A"})
	if !act {
		t.Errorf("FAILED inMANY")
	}
	act = in(cards, Card{CARO, "10"}, Card{CARO, "A"})
	if act {
		t.Errorf("FAILED inMANY")
	}

	if !in(cards, Card{CARO, "A"}, Card{CARO, "K"}) {
		t.Errorf("FAILED inMANY")
	}

}


func TestNextLowestCardsStillInPlay(t *testing.T) {
	s := makeSuitState()

	s.trick = []Card{Card{SPADE, "7"}, Card{SPADE, "9"}}
	s.cardsPlayed = []Card{
		Card{SPADE, "D"},
		Card{SPADE, "8"},
		//	Card{SPADE, "9"},
	}
	w := Card{SPADE, "A"}
	followCards := []Card{Card{SPADE, "A"}, Card{SPADE, "K"}}
	still10 := true
	if nextLowestCardsStillInPlay(&s, w, followCards) != still10 {
		t.Errorf("10 not played")
	}

	s.cardsPlayed = append(s.cardsPlayed, Card{SPADE, "10"})
	still10 = false
	if nextLowestCardsStillInPlay(&s, w, followCards) != still10 {
		t.Errorf("10 played")
	}
}






func TestSingletons(t *testing.T) {
	cs := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "8"},
		Card{CLUBS, "10"},
		Card{SPADE, "9"},
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "K"},
		Card{CARO, "9"},
	}

	s := singletons(cs)
	if len(s) != 2 {
		t.Errorf("Singleton error, found: %v", s)
	}

	if !in(s, Card{CARO, "9"}, Card{SPADE, "9"}) {
		t.Errorf("Singleton error, found: %v", s)
	}
}

func TestDeclarerVoidSuits(t *testing.T) {
	cs := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "8"},
		Card{CLUBS, "10"},
		Card{SPADE, "9"},
		Card{CARO, "9"},
	}
	player := makePlayer(cs)
	opp1 := makePlayer([]Card{})
	opp2 := makePlayer([]Card{})

	s := makeSuitState()
	s.declarer = &player
	s.trump = NULL
	s.follow = HEART
	s.trick = []Card{Card{HEART, "10"}}
	s.opp1 = &opp1
	s.opp2 = &opp2

	play(&s, &player)

	if !s.declarerVoidSuit[HEART] {
		t.Errorf("declarerVoidSuits %v", s.declarerVoidSuit)
	}
}

func TestSmallerCardsInPlay(t *testing.T) {
	s := makeSuitState()
	trick := Card{HEART, "J"}
	s.cardsPlayed = []Card{Card{HEART, "8"}, Card{HEART, "10"}}
	cs := []Card{Card{HEART, "K"}, Card{HEART, "9"}}
	// {8, 10} J  {K,9} ===> 7 in play
	act := smallerCardsInPlay(&s, trick, cs)

	if !act {
		t.Errorf("Error in smallerCardsInPlay")
	}

	s.cardsPlayed = []Card{Card{HEART, "7"}, Card{HEART, "8"}, Card{HEART, "10"}}
	act = smallerCardsInPlay(&s, trick, cs)

	if act {
		t.Errorf("Error in smallerCardsInPlay")
	}

}

func TestHTMLBid1(t *testing.T) {
	router := startServer()
	var m BidData

	makeChannels()

	makeP := func(high int) Player {
		player := makePlayer([]Card{})
		player.highestBid = high
		return player
	}
	player1 := makeP(24)
	player2 := makeP(36)
	player3 := makeP(0)
	// GLOBAL
	players = []PlayerI{&player1, &player2, &player3}

	// req, _ := http.NewRequest("GET", "/start", nil)
	rr := httptest.NewRecorder()
	// router.ServeHTTP(rr, req)

	// MIDDLEHAND 18
	req, _ := http.NewRequest("GET", "/getbidvalue/1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp := 18
	act := m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}
	//FOREHAND Yes
	req, _ = http.NewRequest("GET", "/bid/0", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp = 18
	act = m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}
	expB := true
	actB := m.Accepted
	if actB != expB {
		t.Errorf("Error bid, exp: %v, found %v", expB, actB)
	}

	// MIDDLEHAND 20
	req, _ = http.NewRequest("GET", "/getbidvalue/1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp = 20
	act = m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}

	// FOREHAND (20) Yes
	req, _ = http.NewRequest("GET", "/bid/0", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp = 20
	act = m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}
	expB = true
	actB = m.Accepted
	if actB != expB {
		t.Errorf("Error bid, exp: %v, found %v", expB, actB)
	}
}

func TestHTMLBid2(t *testing.T) {
	router := startServer()
	var m BidData

	makeChannels()

	makeP := func(high int) Player {
		player := makePlayer([]Card{})
		player.highestBid = high
		return player
	}
	player1 := makeP(0)
	player2 := makeP(24)
	player3 := makeP(48)
	player1.name = "Bob"
	player2.name = "Ana"
	player3.name = "You"
	// GLOBAL
	players = []PlayerI{&player1, &player2, &player3}

	// req, _ := http.NewRequest("GET", "/start", nil)
	rr := httptest.NewRecorder()
	// router.ServeHTTP(rr, req)

	// MIDDLEHAND 18
	req, _ := http.NewRequest("GET", "/bid/1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp := 18
	act := m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}
	expB := true
	actB := m.Accepted
	if actB != expB {
		t.Errorf("Error bid, exp: %v, found %v", expB, actB)
	}
	//FOREHAND No
	req, _ = http.NewRequest("GET", "/bid/0", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp = 18
	act = m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}
	expB = false
	actB = m.Accepted
	if actB != expB {
		t.Errorf("Error bid, exp: %v, found %v", expB, actB)
	}

	// BACKHAND 20
	req, _ = http.NewRequest("GET", "/getbidvalue/2", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp = 20
	act = m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}

	// MIDDLEHAND (20) Yes
	req, _ = http.NewRequest("GET", "/bid/1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp = 20
	act = m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}
	expB = true
	actB = m.Accepted
	if actB != expB {
		t.Errorf("Error bid, exp: %v, found %v", expB, actB)
	}
}

func TestHTMLBid3(t *testing.T) {
	router := startServer()
	var m BidData

	makeChannels()

	makeP := func(high int) Player {
		player := makePlayer([]Card{})
		player.highestBid = high
		return player
	}
	player1 := makeP(0)
	player2 := makeP(0)
	player3 := makeP(18)
	player1.name = "Ana"
	player2.name = "You"
	player3.name = "Bob"
	// GLOBAL
	players = []PlayerI{&player1, &player2, &player3}

	// req, _ := http.NewRequest("GET", "/start", nil)
	rr := httptest.NewRecorder()
	// router.ServeHTTP(rr, req)

	// MIDDLEHAND NO
	req, _ := http.NewRequest("GET", "/getbidvalue/1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp := 18
	act := m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}

	// BACKHAND 18 yes
	req, _ = http.NewRequest("GET", "/bid/2", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp = 18
	act = m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}
	expB := true
	actB := m.Accepted
	if actB != expB {
		t.Errorf("Error bid, exp: %v, found %v", expB, actB)
	}


	// FOREHAND (18) No
	req, _ = http.NewRequest("GET", "/bid/0", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp = 18
	act = m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}
	expB = false
	actB = m.Accepted
	if actB != expB {
		t.Errorf("Error bid, exp: %v, found %v", expB, actB)
	}
}

func TestHTMLBid4(t *testing.T) {
	router := startServer()
	var m BidData

	makeChannels()

	makeP := func(high int) Player {
		player := makePlayer([]Card{})
		player.highestBid = high
		return player
	}
	player1 := makeP(0)
	player2 := makeP(18)
	player3 := makeP(18)
	player1.name = "Ana"
	player2.name = "You"
	player3.name = "Bob"
	// GLOBAL
	players = []PlayerI{&player1, &player2, &player3}

	// req, _ := http.NewRequest("GET", "/start", nil)
	rr := httptest.NewRecorder()
	// router.ServeHTTP(rr, req)

	// MIDDLEHAND yes
	req, _ := http.NewRequest("GET", "/getbidvalue/1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp := 18
	act := m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}

	//FOREHAND No
	req, _ = http.NewRequest("GET", "/bid/0", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp = 18
	act = m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}
	expB := false
	actB := m.Accepted
	if actB != expB {
		t.Errorf("Error bid, exp: %v, found %v", expB, actB)
	}

	// BACKHAND 20 no
	req, _ = http.NewRequest("GET", "/bid/2", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	json.Unmarshal(rr.Body.Bytes(), &m)
	exp = 20
	act = m.Bid
	if act != exp {
		t.Errorf("Error bid, exp: %v, found %v", exp, act)
	}
	expB = false
	actB = m.Accepted
	if actB != expB {
		t.Errorf("Error bid, exp: %v, found %v", expB, actB)
	}


}

func TestWinnerCards(t *testing.T) {
	cards := []Card{
		Card{HEART, "D"},
		Card{HEART, "K"},
		Card{HEART, "9"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{CARO, "8"},
	}
	s := makeSuitState()
	s.trump = CLUBS
	s.follow = SPADE
	s.trick = []Card{Card{SPADE, "7"}}

	winners := winnerCards(&s, cards)

	if len(winners) > 0 {
		t.Errorf("Expected 0 winners, got %v", winners)
	}
}


func TestGreaterTrump(t *testing.T) {
	s := makeSuitState()
	s.skat = []Card{}
	s.trump = CLUBS
	card := Card{CLUBS, "J"}
	c := Card{CARO, "J"}
	act := s.greater(card, c) 
	exp := true
	if act != exp {
		t.Errorf("Error in higher")
	}
}



func TestDealCardsMID(t *testing.T) {
	p := makeMinMaxPlayer([]Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{SPADE, "A"},
	})
	p.maxHandSize = 6

	s := makeSuitState()
	s.declarer = &p
	s.skat = []Card{Card{CLUBS, "10"}, Card{SPADE, "K"}}

	s.trick = []Card{Card{CLUBS, "K"}}

	notPlayedYet := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "D"},
		Card{SPADE, "10"},
		Card{CARO, "10"},
		Card{HEART, "10"},
	}

	s.cardsPlayed = makeDeck()
	s.cardsPlayed = remove(s.cardsPlayed, s.skat...)
	s.cardsPlayed = remove(s.cardsPlayed, s.trick...)
	s.cardsPlayed = remove(s.cardsPlayed, p.hand...)
	s.cardsPlayed = remove(s.cardsPlayed, notPlayedYet...)

	fmt.Printf("PLAYED CARDS %v\n", s.cardsPlayed)

	worlds, _ := p.dealCards(&s)
	if len(worlds) != 10  {
		t.Errorf("Expecting 10 world, found %d: %v", len(worlds), worlds)
	}
	for i := 0; i < len(worlds); i++ {
		// SET world
		p.p1Hand = worlds[i][0]
		p.p2Hand = worlds[i][1]

		fmt.Printf("DEALT: %v %v\n", p.p1Hand, p.p2Hand)

		if len(p.p2Hand) != 2 {
			t.Errorf("Wrong hand size for player who just opened the trick: %v", p.p2Hand)
		}

		if len(p.p1Hand) != 3 {
			t.Errorf("Wrong hand size for next player: %v", p.p1Hand)
		}

		if in(p.p1Hand, s.cardsPlayed...) {
			t.Errorf("Cards already played: %v", p.p1Hand)
		}
		if in(p.p2Hand, s.cardsPlayed...) {
			t.Errorf("Cards already played: %v", p.p2Hand)
		}
	}
}


func TestDealCardsMID2(t *testing.T) {
	p := makeMinMaxPlayer([]Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{SPADE, "A"},
	})
	p.maxHandSize = 6

	s := makeSuitState()
	s.declarer = &p
	s.skat = []Card{Card{CLUBS, "10"}, Card{SPADE, "K"}}

	s.trick = []Card{Card{CLUBS, "K"}}

	notPlayedYet := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "D"},
		Card{SPADE, "10"},
		Card{CARO, "10"},
		Card{HEART, "10"},
	}

	s.cardsPlayed = makeDeck()
	s.cardsPlayed = remove(s.cardsPlayed, s.skat...)
	s.cardsPlayed = remove(s.cardsPlayed, s.trick...)
	s.cardsPlayed = remove(s.cardsPlayed, p.hand...)
	s.cardsPlayed = remove(s.cardsPlayed, notPlayedYet...)

	fmt.Printf("PLAYED CARDS %v\n", s.cardsPlayed)

	// VOID SUITS
	void1 := CARO
	void2 := CLUBS
	s.opp1VoidSuit[CARO] = true
	s.opp2VoidSuit[CLUBS] = true

	worlds,_ := p.dealCards(&s)

	for i := 0; i < len(worlds); i++ {
		// SET world
		p.p1Hand = worlds[i][0]
		p.p2Hand = worlds[i][1]

		fmt.Printf("DEALT: %v %v\n", p.p1Hand, p.p2Hand)

		for _,c := range  makeSuitDeck(void1) {
			if in(p.p1Hand, c) {
				t.Errorf("Opp1 is VOID of %s: %v %v", p.p1Hand, c, void1)
			}		
		}
		for _,c := range  makeSuitDeck(void2) {
			if in(p.p2Hand, c) {
				t.Errorf("Opp2 is VOID of %s: %v %v", p.p2Hand, c, void2)
			}		
		}

		if len(p.p2Hand) != 2 {
			t.Errorf("Wrong hand size for player who just opened the trick: %v", p.p2Hand)
		}

		if len(p.p1Hand) != 3 {
			t.Errorf("Wrong hand size for next player: %v", p.p1Hand)
		}	
	}


}

func TestDealCardsLeader(t *testing.T) {
	p := makeMinMaxPlayer([]Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{SPADE, "A"},
	})
	p.maxHandSize = 6

	s := makeSuitState()
	s.declarer = &p
	s.skat = []Card{Card{CLUBS, "10"}, Card{SPADE, "K"}}

	s.trick = []Card{}

	notPlayedYet := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "D"},
		Card{SPADE, "10"},
		Card{CARO, "10"},
		Card{HEART, "10"},
//
		Card{CLUBS, "K"},
	}

	s.cardsPlayed = makeDeck()
	s.cardsPlayed = remove(s.cardsPlayed, s.skat...)
	s.cardsPlayed = remove(s.cardsPlayed, s.trick...)
	s.cardsPlayed = remove(s.cardsPlayed, p.hand...)
	s.cardsPlayed = remove(s.cardsPlayed, notPlayedYet...)

	p1HandSize := 3
	p2HandSize := 3

	fmt.Printf("PLAYED CARDS %v\n", s.cardsPlayed)

	worlds, _ := p.dealCards(&s)

	for i := 0; i < len(worlds); i++ {
		// SET world
		p.p1Hand = worlds[i][0]
		p.p2Hand = worlds[i][1]

		fmt.Printf("DEALT: %v %v\n", p.p1Hand, p.p2Hand)

		if len(p.p2Hand) != p2HandSize {
			t.Errorf("Wrong hand size for player who just opened the trick: %v", p.p2Hand)
		}

		if len(p.p1Hand) != p1HandSize {
			t.Errorf("Wrong hand size for next player: %v", p.p1Hand)
		}

		if in(p.p1Hand, s.cardsPlayed...) {
			t.Errorf("Cards already played: %v", p.p1Hand)
		}
		if in(p.p2Hand, s.cardsPlayed...) {
			t.Errorf("Cards already played: %v", p.p2Hand)
		}
	}

}

func TestDealCardsLeader3(t *testing.T) {
	p := makeMinMaxPlayer([]Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{SPADE, "A"},
	})
	p.maxHandSize = 6

	s := makeSuitState()
	s.declarer = &p
	s.skat = []Card{Card{CLUBS, "10"}, Card{SPADE, "K"}}

	s.trick = []Card{}

	notPlayedYet := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "D"},
		Card{SPADE, "10"},
		Card{CARO, "10"},
		Card{HEART, "10"},
//
		Card{CLUBS, "K"},
	}

	s.cardsPlayed = makeDeck()
	s.cardsPlayed = remove(s.cardsPlayed, s.skat...)
	s.cardsPlayed = remove(s.cardsPlayed, s.trick...)
	s.cardsPlayed = remove(s.cardsPlayed, p.hand...)
	s.cardsPlayed = remove(s.cardsPlayed, notPlayedYet...)

	p1HandSize := 3
	p2HandSize := 3

	fmt.Printf("PLAYED CARDS %v\n", s.cardsPlayed)	

	// VOID SUITS
	void1 := CARO
	void2 := CLUBS
	s.opp1VoidSuit[CARO] = true
	s.opp2VoidSuit[CLUBS] = true

	worlds, _ := p.dealCards(&s)

	for i := 0; i < len(worlds); i++ {
		// SET world
		p.p1Hand = worlds[i][0]
		p.p2Hand = worlds[i][1]

		fmt.Printf("DEALT: %v %v\n", p.p1Hand, p.p2Hand)

		for _,c := range  makeSuitDeck(void1) {
			if in(p.p1Hand, c) {
				t.Errorf("Opp1 is VOID of %s: %v %v", p.p1Hand, c, void1)
			}		
		}
		for _,c := range  makeSuitDeck(void2) {
			if in(p.p2Hand, c) {
				t.Errorf("Opp2 is VOID of %s: %v %v", p.p2Hand, c, void2)
			}		
		}

		if len(p.p2Hand) != p2HandSize {
			t.Errorf("Wrong hand size for player who just opened the trick: %v", p.p2Hand)
		}

		if len(p.p1Hand) != p1HandSize {
			t.Errorf("Wrong hand size for next player: %v", p.p1Hand)
		}	
	}


}

func TestDealCardsBack(t *testing.T) {
	p := makeMinMaxPlayer([]Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{SPADE, "A"},
	})
	p.maxHandSize = 6

	s := makeSuitState()
	s.declarer = &p
	s.skat = []Card{Card{CLUBS, "10"}, Card{SPADE, "K"}}

	s.trick = []Card{Card{CLUBS, "K"},Card{CLUBS, "A"}}

	notPlayedYet := []Card{
		Card{CLUBS, "D"},
		Card{SPADE, "10"},
		Card{CARO, "10"},
		Card{HEART, "10"},
//
		
	}

	s.cardsPlayed = makeDeck()
	s.cardsPlayed = remove(s.cardsPlayed, s.skat...)
	s.cardsPlayed = remove(s.cardsPlayed, s.trick...)
	s.cardsPlayed = remove(s.cardsPlayed, p.hand...)
	s.cardsPlayed = remove(s.cardsPlayed, notPlayedYet...)

	p1HandSize := 2
	p2HandSize := 2

	fmt.Printf("PLAYED CARDS %v\n", s.cardsPlayed)

	worlds, _ := p.dealCards(&s)
	if len(worlds) != 6  {
		t.Errorf("Expecting 6 world, found %d: %v", len(worlds), worlds)
	}

	for i := 0; i < len(worlds); i++ {
		// SET world
		p.p1Hand = worlds[i][0]
		p.p2Hand = worlds[i][1]

		fmt.Printf("DEALT: %v %v\n", p.p1Hand, p.p2Hand)

		if len(p.p2Hand) != p2HandSize {
			t.Errorf("Wrong hand size for player who just opened the trick: %v", p.p2Hand)
		}

		if len(p.p1Hand) != p1HandSize {
			t.Errorf("Wrong hand size for next player: %v", p.p1Hand)
		}

		if in(p.p1Hand, s.cardsPlayed...) {
			t.Errorf("Cards already played: %v", p.p1Hand)
		}
		if in(p.p2Hand, s.cardsPlayed...) {
			t.Errorf("Cards already played: %v", p.p2Hand)
		}
	}
}

func TestDealCardsBack2(t *testing.T) {
	p := makeMinMaxPlayer([]Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{SPADE, "A"},
	})
	p.maxHandSize = 6

	s := makeSuitState()
	s.declarer = &p
	s.skat = []Card{Card{CLUBS, "10"}, Card{SPADE, "K"}}

	s.trick = []Card{Card{CLUBS, "K"},Card{CLUBS, "A"}}

	notPlayedYet := []Card{
		Card{CLUBS, "D"},
		Card{SPADE, "10"},
		Card{CARO, "10"},
		Card{HEART, "10"},
//
		
	}

	s.cardsPlayed = makeDeck()
	s.cardsPlayed = remove(s.cardsPlayed, s.skat...)
	s.cardsPlayed = remove(s.cardsPlayed, s.trick...)
	s.cardsPlayed = remove(s.cardsPlayed, p.hand...)
	s.cardsPlayed = remove(s.cardsPlayed, notPlayedYet...)

	p1HandSize := 2
	p2HandSize := 2

	fmt.Printf("PLAYED CARDS %v\n", s.cardsPlayed)

	// VOID SUITS
	void1 := CARO
	void2 := CLUBS
	s.opp1VoidSuit[CARO] = true
	s.opp2VoidSuit[CLUBS] = true

	worlds, _ := p.dealCards(&s)

	for i := 0; i < len(worlds); i++ {
		// SET world
		p.p1Hand = worlds[i][0]
		p.p2Hand = worlds[i][1]

		fmt.Printf("DEALT: %v %v\n", p.p1Hand, p.p2Hand)

		for _,c := range  makeSuitDeck(void1) {
			if in(p.p1Hand, c) {
				t.Errorf("Opp1 is VOID of %s: %v %v", p.p1Hand, c, void1)
			}		
		}
		for _,c := range  makeSuitDeck(void2) {
			if in(p.p2Hand, c) {
				t.Errorf("Opp2 is VOID of %s: %v %v", p.p2Hand, c, void2)
			}		
		}

		if len(p.p2Hand) != p2HandSize {
			t.Errorf("Wrong hand size for player who just opened the trick: %v", p.p2Hand)
		}

		if len(p.p1Hand) != p1HandSize {
			t.Errorf("Wrong hand size for next player: %v", p.p1Hand)
		}	
	}


}

func TestDealCardsLeader2(t *testing.T) {
	p := makeMinMaxPlayer([]Card{
		Card{CLUBS, "K"},
		Card{SPADE, "K"},
		Card{HEART, "8"},
		Card{CARO, "7"},
	})
	p.maxHandSize = 6

	s := makeSuitState()
	s.declarer = &p
	s.skat = []Card{Card{HEART, "10"}, Card{HEART, "9"}}
	s.trump = HEART
	s.trick = []Card{}

	notPlayedYet := []Card{
		Card{CLUBS, "10"},
		Card{CLUBS, "A"},
		Card{CLUBS, "8"},

		Card{CARO, "10"},
//
		Card{HEART, "D"},
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{HEART, "J"},
	}

	s.cardsPlayed = makeDeck()
	s.cardsPlayed = remove(s.cardsPlayed, s.skat...)
	s.cardsPlayed = remove(s.cardsPlayed, s.trick...)
	s.cardsPlayed = remove(s.cardsPlayed, p.hand...)
	s.cardsPlayed = remove(s.cardsPlayed, notPlayedYet...)

	p1HandSize := 4
	p2HandSize := 4

	fmt.Printf("PLAYED CARDS %v\n", s.cardsPlayed)

	// VOID SUITS
	void1 := CARO
	void2 := CLUBS
	s.opp1VoidSuit[void1] = true
	s.opp2VoidSuit[void2] = true


	worlds, _ := p.dealCards(&s)

	for i := 0; i < len(worlds); i++ {
		// SET world
		p.p1Hand = worlds[i][0]
		p.p2Hand = worlds[i][1]

		fmt.Printf("DEALT: %v %v\n", p.p1Hand, p.p2Hand)

		for _,c := range  makeSuitDeck(void1) {
			if in(p.p1Hand, c) {
				t.Errorf("Opp1 is VOID of %s: %v %v", p.p1Hand, c, void1)
			}		
		}
		for _,c := range  makeSuitDeck(void2) {
			if in(p.p2Hand, c) {
				t.Errorf("Opp2 is VOID of %s: %v %v", p.p2Hand, c, void2)
			}		
		}
		if len(p.p2Hand) != p2HandSize {
			t.Errorf("Wrong hand size for player who just opened the trick: %v", p.p2Hand)
		}

		if len(p.p1Hand) != p1HandSize {
			t.Errorf("Wrong hand size for next player: %v", p.p1Hand)
		}	
	}
}

