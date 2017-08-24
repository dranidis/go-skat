package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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
		t.Errorf("TRUMP :" + s.trump + " FOLLOW :" + s.follow + " - " + card1.Suit + " " + card1.Rank + " should be greater than " + card2.Suit + " " + card2.Rank)
	}
	if s.greater(card2, card1) {
		t.Errorf("TRUMP :" + s.trump + " FOLLOW :" + s.follow + " - " + card2.Suit + " " + card2.Rank + " should NOT be greater than " + card1.Suit + " " + card1.Rank)
	}
}

func TestGetSuite(t *testing.T) {
	s := mState(CLUBS, HEART)
	card := Card{SPADE, "J"}
	if getSuit(s.trump, card) != CLUBS {
		t.Errorf("TRUMP :" + s.trump + " FOLLOW :" + s.follow + " - " + fmt.Sprintf("%v", card) + " should be CLUBS")
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
	clubsHeart := mState(CLUBS, HEART)

	playerHand := []Card{Card{SPADE, "J"}, Card{HEART, "A"}, Card{CARO, "7"}}
	cards := validCards(clubsHeart, playerHand)
	compareLists(t, cards, []Card{Card{HEART, "A"}})

	cards = validCards(mState(CLUBS, SPADE), playerHand)
	compareLists(t, cards, playerHand)
}

func compareLists(t *testing.T, returned, expected []Card) {
	if len(returned) != len(expected) {
		t.Errorf("Expected: " + fmt.Sprintf("%v", expected) + " found: " + fmt.Sprintf("%v", returned))
	}
	for i, c := range expected {
		if c.Suit != returned[i].Suit || c.Rank != returned[i].Rank {
			t.Errorf("Expected: " + fmt.Sprintf("%v", expected) + " found: " + fmt.Sprintf("%v", returned))
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
// 		t.Errorf("Expected: player1 len hand 2: " + fmt.Sprintf("%v", player1.hand))
// 	}
// 	if len(player2.hand) != 2 {
// 		t.Errorf("Expected: player2 len hand 2" + fmt.Sprintf("%v", player2.hand))
// 	}
// 	if len(player3.hand) != 2 {
// 		t.Errorf("Expected: player3 len hand 2" + fmt.Sprintf("%v", player3.hand))
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
	// from two A-suits prefer the strongest
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
	if most != HEART {
		t.Errorf("In hand %v, Two suits equal length. Expected %s, got %s", player.hand, HEART, most)
	}

	player = makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},

		Card{HEART, "A"},
		Card{HEART, "7"},
		Card{HEART, "8"},
		Card{HEART, "9"},

		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{CARO, "8"},
	})
	most = mostCardsSuit(player.hand)
	if most != CARO {
		t.Errorf("In hand %v, Two suits equal length. Expected %s, got %s", player.hand, CARO, most)
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

func TestPickUpSkat(t *testing.T) {
	player := makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{HEART, "D"},
		Card{CARO, "9"},
		Card{SPADE, "8"},
	})

	skat := []Card{
		Card{CARO, "D"},
		Card{CLUBS, "A"},
	}

	player.pickUpSkat(skat)

	cc1 := skat[1].Suit != SPADE || skat[1].Rank != "8"
	cc2 := skat[0].Suit != HEART || skat[0].Rank != "D"
	if cc1 || cc2 {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
	}

	if len(player.hand) != 10 {
		t.Errorf("Wrong hand size after skat change: %d", len(player.hand))
	}
}

func TestPickUpSkat1(t *testing.T) {
	player := makePlayer([]Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CLUBS, "7"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{SPADE, "D"},
		Card{SPADE, "8"},
		Card{HEART, "9"},
		Card{CARO, "K"},
		Card{CARO, "D"},
	})

	skat := []Card{
		Card{CARO, "10"},
		Card{HEART, "K"},
	}
	// fmt.Println(player.hand)
	// fmt.Println(skat)
	player.pickUpSkat(skat)
	cc1 := skat[0].Suit != CLUBS || skat[0].Rank != "7"
	cc2 := skat[1].Suit != HEART || skat[1].Rank != "K"
	if cc1 || cc2 {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
	}

	if len(player.hand) != 10 {
		t.Errorf("Wrong hand size after skat change: %d", len(player.hand))
	}
}

func TestPickUpSkat2(t *testing.T) {
	player := makePlayer([]Card{
		Card{CLUBS, "J"},

		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "9"},

		Card{SPADE, "A"},
		Card{SPADE, "K"},
		Card{SPADE, "8"},

		Card{HEART, "D"},
		Card{HEART, "7"},
		Card{CARO, "A"},
	})

	skat := []Card{
		Card{HEART, "8"},
		Card{CARO, "8"},
	}
	// fmt.Println("TestPickUpSkat2")
	// fmt.Println(player.hand)
	// fmt.Println(skat)
	player.pickUpSkat(skat)
	// fmt.Println(sort(player.hand))
	// fmt.Println(skat)
	cc1 := skat[0].Suit != HEART || skat[0].Rank != "D"
	cc2 := skat[1].Suit != HEART || skat[1].Rank != "7"
	if cc1 || cc2 {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
	}

	if len(player.hand) != 10 {
		t.Errorf("Wrong hand size after skat change: %d", len(player.hand))
	}
}

func TestPickUpSkat3(t *testing.T) {
	player := makePlayer([]Card{
		Card{CLUBS, "9"},

		Card{SPADE, "A"},
		Card{SPADE, "D"},
		Card{SPADE, "9"},
		Card{SPADE, "8"},
		Card{SPADE, "7"},

		Card{HEART, "A"},
		Card{HEART, "10"},

		Card{CARO, "A"},
		Card{CARO, "D"},
	})

	skat := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "J"},
	}
	// fmt.Println("TestPickUpSkat3")
	// fmt.Println(player.hand)
	// fmt.Println(skat)
	player.pickUpSkat(skat)
	// fmt.Println(sort(player.hand))
	// fmt.Println(skat)
	cc1 := skat[1].Suit != CLUBS || skat[1].Rank != "9"
	cc2 := skat[0].Suit != CARO || skat[0].Rank != "D"
	if cc1 || cc2 {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
	}

	if len(player.hand) != 10 {
		t.Errorf("Wrong hand size after skat change: %d", len(player.hand))
	}
}

func TestPickUpSkat4(t *testing.T) {
	player := makePlayer([]Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},

		Card{CLUBS, "10"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},

		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},

		Card{HEART, "K"},

		Card{CARO, "K"},
	})

	skat := []Card{
		Card{CLUBS, "K"},
		Card{HEART, "7"},
	}
	// fmt.Println("TestPickUpSkat4")
	// fmt.Println(player.hand)
	// fmt.Println(skat)
	player.pickUpSkat(skat)
	// fmt.Println(sort(player.hand))
	// fmt.Println(skat)
	cc1 := skat[1].Suit != HEART || skat[1].Rank != "K"
	cc2 := skat[0].Suit != CARO || skat[0].Rank != "K"
	if cc1 || cc2 {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
	}

	if len(player.hand) != 10 {
		t.Errorf("Wrong hand size after skat change: %d", len(player.hand))
	}
}

func TestPickUpSkat5(t *testing.T) {
	player := makePlayer([]Card{
		Card{HEART, "J"},

		Card{CLUBS, "A"},
		Card{CLUBS, "K"},
		Card{CLUBS, "D"},

		Card{SPADE, "A"},

		Card{HEART, "9"},

		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{CARO, "9"},
		Card{CARO, "8"},
	})

	skat := []Card{
		Card{SPADE, "J"},
		Card{SPADE, "9"},
	}
	player.trumpToDeclare = CARO
	// fmt.Println("TestPickUpSkat5")
	// fmt.Println(player.hand)
	// fmt.Println(skat)
	player.pickUpSkat(skat)
	// fmt.Println(sort(player.hand))
	// fmt.Println(skat)
	cc1 := skat[1].Suit != SPADE || skat[1].Rank != "9"
	cc2 := skat[0].Suit != HEART || skat[0].Rank != "9"
	if cc1 || cc2 {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
	}

	if len(player.hand) != 10 {
		t.Errorf("Wrong hand size after skat change: %d", len(player.hand))
	}
}

func TestPickUpSkat6(t *testing.T) {
	player := makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},

		Card{SPADE, "A"},
		Card{SPADE, "9"},

		Card{HEART, "A"},

		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{CARO, "D"},
		Card{CARO, "9"},
	})

	skat := []Card{
		Card{CARO, "7"},
		Card{SPADE, "J"},
	}
	player.trumpToDeclare = CARO
	// fmt.Println("TestPickUpSkat6")
	// fmt.Println(player.hand)
	// fmt.Println(skat)
	player.pickUpSkat(skat)
	// fmt.Println(sortSuit("", player.hand))
	// fmt.Println(skat)
	cc1 := skat[1].Suit != HEART || skat[1].Rank != "A"
	cc2 := skat[0].Suit != SPADE || skat[0].Rank != "9"
	if cc1 || cc2 {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
	}

	if len(player.hand) != 10 {
		t.Errorf("Wrong hand size after skat change: %d", len(player.hand))
	}
}

func TestPickUpSkat7(t *testing.T) {
	player := makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
		Card{SPADE, "J"},
		Card{CARO, "J"},

		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{CARO, "D"},
		Card{CARO, "9"},
		Card{CARO, "8"},
	})

	skat := []Card{
		Card{HEART, "A"},
		Card{SPADE, "A"},
	}
	// fmt.Println("TestPickUpSkat7")
	// fmt.Println(player.hand)
	// fmt.Println(skat)
	player.pickUpSkat(skat)
	// fmt.Println(sortSuit("", player.hand))
	// fmt.Println(skat)
	cc1 := skat[1].Suit != SPADE || skat[1].Rank != "A"
	cc2 := skat[0].Suit != HEART || skat[0].Rank != "A"
	if cc1 || cc2 {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
	}

	if len(player.hand) != 10 {
		t.Errorf("Wrong hand size after skat change: %d", len(player.hand))
	}
}

func TestCalculateHighestBid(t *testing.T) {

	player := makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{HEART, "D"},// loser
		Card{HEART, "9"},// loser
		Card{SPADE, "8"},// loser
	})
	player.risky = true
	player.calculateHighestBid()

	act := player.highestBid
	exp := 5 * 24
	if act != exp {
		t.Errorf("Expected high bid %d, got %d", exp, act)
	}

}
func TestHandEstimation(t *testing.T) {
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
		Card{SPADE, "8"},
	})
	_ = player
	//fmt.Println(player.handEstimation())
}

func TestSum(t *testing.T) {
	c := makeDeck()
	act := sum(c)
	if act != 120 {
		t.Errorf("Sum of deck is %d", act)
	}
}

func TestDeck(t *testing.T) {
	for i := 0; i < 10; i++ {
		auxDeck(t)
	}
}
func auxDeck(t *testing.T) {
	cards := Shuffle(makeDeck())
	//fmt.Printf("CARDS %d %v\n", -1, cards)
	hand1 := cards[:10]
	hand2 := cards[10:20]
	hand3 := cards[20:30]
	skat := cards[30:32]

	sumAll := sum(hand1) + sum(hand2) + sum(hand3) + sum(skat)
	if sumAll != 120 {
		t.Errorf("DEAL PROBLEM: %d", sumAll)
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

	cardToRemove := Card{CLUBS, "J"}
	newhand := remove(hand, cardToRemove)
	if in(newhand, cardToRemove) {
		t.Errorf("First Card not removed")
	}

	cardToRemove = Card{SPADE, "8"}
	newhand = remove(hand, cardToRemove)
	if in(newhand, cardToRemove) {
		t.Errorf("Last Card not removed")
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

func TestGame(t *testing.T) {
	// player1 is declared globally
	player := makePlayer([]Card{})
	player1 = &player

	player2 := makePlayer([]Card{})
	player3 := makePlayer([]Card{})
	player3.firstCardPlay = true
	player1.setName("NAME")
	if player1.getName() != "NAME" {
		t.Errorf("Error in set/get name")
	}
	if player1.getTotalScore() != 0 {
		t.Errorf("Error in get total score")
	}
	players := []PlayerI{player1, &player2, &player3}
	for i := 0; i < 10; i++ {
		_ = game(players)
	}
}

func TestOpponentTacticFOREFollowPreviousSuit3(t *testing.T) {

	// if you have a card with suit played in a previous trick
	// started from you or your partner continue with it.

	otherPlayer := makePlayer([]Card{})
	//teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer
	s.trump = CLUBS
	s.trick = []Card{}

	player.previousSuit = CARO

	validCards := []Card{Card{CARO, "8"},
		Card{HEART, "9"},
		Card{HEART, "10"},
		Card{HEART, "A"},
	}
	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and valid %v, expected to play %v to follow previously played suit, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticFOREFollowPreviousSuit4(t *testing.T) {

	// if you have a card with suit played in a previous trick
	// started from you or your partner continue with it.

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer

	s.trump = CLUBS
	s.trick = []Card{}
	teamMate.previousSuit = CARO
	validCards := []Card{Card{CARO, "8"},
		Card{HEART, "9"},
		Card{HEART, "10"},
		Card{HEART, "A"},
	}
	//
	s.opp1 = &player
	s.opp2 = &teamMate
	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and valid %v, expected to play %v to follow previously played suit, played %v",
			s.trick, validCards, exp, card)
	}
	//
	s.opp2 = &player
	s.opp1 = &teamMate
	card = player.playerTactic(&s, validCards)
	exp = Card{CARO, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and valid %v, expected to play %v to follow previously played suit, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticFOREFollowPreviousSuit5(t *testing.T) {

	// if you have a card with suit played in a previous trick
	// started from you or your partner continue with it.

	// unless your card is a J

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer

	s.trump = CLUBS
	s.trick = []Card{}
	teamMate.previousSuit = CARO
	validCards := []Card{Card{CARO, "J"},
		Card{HEART, "9"},
		Card{HEART, "10"},
		Card{HEART, "A"},
	}
	//
	s.opp1 = &player
	s.opp2 = &teamMate
	card := player.playerTactic(&s, validCards)
	unexp := Card{CARO, "J"}
	if card.equals(unexp) {
		t.Errorf("In trick %v and valid %v, not expected to play %v to follow previously played suit, played %v",
			s.trick, validCards, unexp, card)
	}
}

func TestOpponentTacticMIDTrump1(t *testing.T) {

	// if player has J caro and A trump and is required to play a trump
	// in a losing trick he should play J
	// in a winning trick he should play A

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.trump = CLUBS
	s.follow = CLUBS

	validCards := []Card{Card{CARO, "J"}, Card{CLUBS, "A"}}

	s.trick = []Card{Card{CLUBS, "J"}}

	s.declarer = &otherPlayer

	s.leader = &otherPlayer
	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "J"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by opponent, valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}

	s.leader = &teamMate
	card = player.playerTactic(&s, validCards)
	exp = Card{CLUBS, "A"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by teammate, valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}

	s.trick = []Card{Card{CLUBS, "10"}}

	s.leader = &otherPlayer
	card = player.playerTactic(&s, validCards)

	exp = Card{CLUBS, "A"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by opponent, and valid cards: %v expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}

	s.leader = &teamMate
	card = player.playerTactic(&s, validCards)

	exp = Card{CARO, "J"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, expected to play %v, played %v", s.trick, exp, card)
	}

	//////////////
	s.trick = []Card{Card{CLUBS, "J"}}
	validCards = []Card{Card{CARO, "J"}, Card{CLUBS, "A"}, Card{CLUBS, "D"}, Card{CLUBS, "9"}}
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	card = player.playerTactic(&s, validCards)
	exp = Card{CLUBS, "9"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by opponent, and valid cards: %v expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}

	//////////////
	s.trick = []Card{Card{CARO, "A"}, Card{CARO, "7"}}
	validCards = []Card{Card{CARO, "K"}, Card{CARO, "10"}, Card{CARO, "7"}}
	s.leader = &teamMate
	s.declarer = &otherPlayer
	card = player.playerTactic(&s, validCards)
	exp = Card{CARO, "10"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by teammate, and valid cards: %v expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}

	//////////////
	s.trump = SPADE
	s.follow = SPADE
	s.trick = []Card{Card{SPADE, "J"}, Card{HEART, "J"}}
	validCards = []Card{Card{CARO, "J"}, Card{SPADE, "D"}}
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	card = player.playerTactic(&s, validCards)
	exp = Card{SPADE, "D"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by declarer, and valid cards: %v expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticMIDTrump2(t *testing.T) {

	// if declarer leads with a low trump
	// to not waste your high trumps
	otherPlayer := makePlayer([]Card{})
	//teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.trump = CLUBS
	s.follow = CLUBS
	s.trick = []Card{Card{CLUBS, "8"}}

	validCards := []Card{Card{CARO, "J"}, Card{CLUBS, "9"}}
	card := player.playerTactic(&s, validCards)
	exp := Card{CLUBS, "9"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by opponent, and valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticMIDTrump6(t *testing.T) {
	// MIDDLEHAND

	// if declarer leads a low trump, and there are still higher trumps
	// smear it with a high value

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

	s.trump = CLUBS
	s.trick = []Card{Card{CLUBS, "D"}}
	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
	}

	validCards := []Card{
		Card{SPADE, "9"},
		Card{HEART, "9"},
		Card{HEART, "10"},
		Card{HEART, "A"},
	}
	//

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, with trumps in game: %v and valid %v, expected to smear %v, played %v",
			s.trick, s.trumpsInGame, validCards, exp, card)
	}
}

func TestOpponentTacticMIDPlayerLeadsLosingCard_Smear(t *testing.T) {
	// MIDDLEHAND

	// if declarer leads a losing card (there are higher cards in game), SMEAR

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

	s.trump = CLUBS
	s.trick = []Card{Card{HEART, "D"}}

	validCards := []Card{
		Card{SPADE, "A"},
		Card{SPADE, "K"},
		Card{SPADE, "9"},
	}
	//

	card := player.playerTactic(&s, validCards)
	exp := Card{SPADE, "A"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, with cards played %v and valid %v, expected to smear %v, played %v",
			s.trick, s.cardsPlayed, validCards, exp, card)
	}
}

func TestOpponentTacticMIDPlayerLeadsLosingCard_Smear1(t *testing.T) {
	// MIDDLEHAND

	// if declarer leads a losing card (there are higher cards in game), SMEAR

	validCards := []Card{
		Card{HEART, "D"},
		Card{HEART, "K"},
		Card{HEART, "9"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{CARO, "8"},
	}
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

	s.trump = CLUBS
	s.trick = []Card{Card{SPADE, "7"}}

	//
	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "10"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, with cards played %v and valid %v, expected to smear %v, played %v",
			s.trick, s.cardsPlayed, validCards, exp, card)
	}
}

func TestOpponentTacticMIDTrump7(t *testing.T) {
	// MIDDLEHAND

	// if declarer leads a trump, and there no higher trumps in game
	// do not smear it with a high value

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

	s.trump = CLUBS
	s.trick = []Card{Card{CLUBS, "D"}}
	s.trumpsInGame = []Card{
		Card{CLUBS, "9"},
	}
	s.cardsPlayed = []Card{
		Card{CLUBS, "J"},		
		Card{SPADE, "J"},		
		Card{HEART, "J"},		
		Card{CARO, "J"},		
		Card{CLUBS, "A"},		
		Card{CLUBS, "10"},	
		Card{CLUBS, "K"},		
		Card{CLUBS, "8"},		
		Card{CLUBS, "7"},		
	}

	validCards := []Card{
		Card{SPADE, "9"},
		Card{HEART, "9"},
		Card{HEART, "10"},
		Card{HEART, "A"},
	}
	//

	card := player.playerTactic(&s, validCards)
	unexp := Card{HEART, "A"}
	if card.equals(unexp) {
		t.Errorf("In trick %v, with trumps in game: %v and valid %v, it is NOT expected to smear %v, played %v",
			s.trick, s.trumpsInGame, validCards, unexp, card)
	}
}

func TestOpponentTacticMID(t *testing.T) {
	// MIDDLEHAND

	// if partner leads a very low card
	// don't SMEAR the trick

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer

	s.trump = CLUBS
	s.trick = []Card{Card{SPADE, "7"}}
	s.follow = SPADE
	s.cardsPlayed = []Card{
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{SPADE, "D"},
		Card{SPADE, "8"},
		//	Card{SPADE, "9"},
	}

	validCards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
		Card{HEART, "A"},
		Card{CARO, "8"},
		Card{CARO, "10"},
	}
	//

	card := player.playerTactic(&s, validCards)
	exp := Card{CLUBS, "7"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by teammate, and VOID hand, and SPADES still in game, and valid %v, it is expected to trump, expected: %v, played %v",
			s.trick, validCards, exp, card)
	}

	s.cardsPlayed = append(s.cardsPlayed, Card{SPADE, "9"})
	card = player.playerTactic(&s, validCards)
	exp = Card{CARO, "10"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by teammate, and VOID hand, all SPADE played, and valid %v, it is expected to Increase the value of the trick for the declarer to trump, expected: %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticMIDFollow(t *testing.T) {
	// MIDDLEHAND

	// if partner leads a very low card
	// and you cannot win it
	// slightly increase the value

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer

	s.trump = CLUBS
	s.trick = []Card{Card{SPADE, "7"}}
	s.follow = SPADE
	s.cardsPlayed = []Card{}

	validCards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{SPADE, "D"},
		Card{SPADE, "8"},
	}
	//

	card := player.playerTactic(&s, validCards)
	exp := Card{SPADE, "D"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by teammate, and A SPADE still in game, and valid %v, expected: %v, played %v",
			s.trick, validCards, exp, card)
	}

	s.cardsPlayed = append(s.cardsPlayed, Card{SPADE, "A"})
	card = player.playerTactic(&s, validCards)
	exp = Card{SPADE, "10"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by teammate, A SPADE played, and valid %v, expected: %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticMIDPartnerLeadsVoidSuit_Trump(t *testing.T) {
	// MIDDLEHAND

	// if partner leads a suit
	// that you are void and
	// cards still in play, trump it.

	validCards := []Card{
		Card{CARO, "K"},
		Card{CARO, "8"},
		Card{SPADE, "A"},
		Card{SPADE, "K"},
		Card{SPADE, "D"},
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{HEART, "8"},
	}

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer

	s.trump = CARO
	s.trick = []Card{Card{CLUBS, "K"}}
	s.follow = CLUBS
	s.cardsPlayed = []Card{
		Card{CLUBS, "A"},		
		Card{CLUBS, "7"},		
		Card{CLUBS, "8"},		
	}

	//

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "K"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by teammate, and A CLUBS still in game, and valid %v, expected: %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticFORE_short_long(t *testing.T) {
	// FOREHAND

	// if declarer BACK short
	// if declarer MID long

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer

	teamMate.previousSuit = ""
	player.previousSuit = ""
	s.trump = CLUBS
	s.trick = []Card{}
	_ = teamMate

	validCards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
		Card{HEART, "A"},
		Card{CARO, "8"},
		Card{CARO, "K"},
		Card{CARO, "D"},
	}
	// declarer MID
	s.opp2 = &player
	s.opp1 = &teamMate

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "K"}
	if !card.equals(exp) {
		t.Errorf("FOREHAND, DECLARER MID, valid %v, expected: %v, played %v",
			validCards, exp, card)
	}
	// declarer BACK
	s.opp1 = &player
	s.opp2 = &teamMate
	teamMate.previousSuit = ""
	player.previousSuit = ""

	card = player.playerTactic(&s, validCards)
	exp = Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("FOREHAND, DECLARER BACK, valid %v, expected: %v, played %v",
			validCards, exp, card)
	}
}

// TODO
// func TestOpponentTacticFORE_long_Not_Full(t *testing.T) {
// 	// FOREHAND

// 	otherPlayer := makePlayer([]Card{})
// 	teamMate := makePlayer([]Card{})
// 	player := makePlayer([]Card{})
// 	s := makeSuitState()
// 	s.leader = &player
// 	s.declarer = &otherPlayer

// 	teamMate.previousSuit = ""
// 	player.previousSuit = ""
// 	s.trump = CLUBS
// 	s.trick = []Card{}
// 	_ = teamMate

// 	validCards := []Card{
// 		Card{CLUBS, "J"},
// 		Card{CLUBS, "8"},
// 		Card{CLUBS, "7"},
// 		Card{HEART, "A"},
// 		Card{CARO, "8"},
// 		Card{CARO, "K"},
// 		Card{CARO, "D"},
// 		Card{CARO, "A"},
// 	}
// 	// declarer MID
// 	s.opp2 = &player
// 	s.opp1 = &teamMate

// 	card := player.playerTactic(&s, validCards)
// 	exp := Card{CARO, "K"}
// 	if !card.equals(exp) {
// 		t.Errorf("FOREHAND, DECLARER MID, valid %v, expected: %v, played a full one %v",
// 			validCards, exp, card)
// 	}
// }

func TestOpponentTacticFORE_short_TOD_SUENDE_1_1(t *testing.T) {
	// FOREHAND

	// if declarer BACK short

	// never 2 numbers, or D number suit
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer

	teamMate.previousSuit = ""
	player.previousSuit = ""
	s.trump = CLUBS
	s.trick = []Card{}
	_ = teamMate

	validCards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},

		Card{HEART, "9"},
		Card{HEART, "8"},

		Card{SPADE, "A"},
		Card{SPADE, "10"},

		Card{CARO, "D"},
		Card{CARO, "9"},

	}

	// declarer BACK
	s.opp1 = &player
	s.opp2 = &teamMate
	teamMate.previousSuit = ""
	player.previousSuit = ""

	card := player.playerTactic(&s, validCards)
	notExpSuit1 := HEART
	if card.Suit == notExpSuit1 {
		t.Errorf("FOREHAND, DECLARER BACK, not expected suit of 2 numbers: %v, played %v",
			validCards, card)
	}
	notExpSuit2 := CARO
	if card.Suit == notExpSuit2 {
		t.Errorf("FOREHAND, DECLARER BACK, not expected suit of D-number: %v, played %v",
			validCards, card)
	}
}

func TestOpponentTacticFORE_short_TOD_SUENDE_1_2(t *testing.T) {
	// FOREHAND

	// if declarer BACK short

	// never 2 numbers, or D number suit
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer

	teamMate.previousSuit = ""
	player.previousSuit = ""
	s.trump = CLUBS
	s.trick = []Card{}
	_ = teamMate

	validCards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},

		Card{HEART, "9"},
		Card{HEART, "8"},

		Card{SPADE, "A"},
		Card{SPADE, "K"},
		Card{SPADE, "9"},
	}

	// declarer BACK
	s.opp1 = &player
	s.opp2 = &teamMate
	teamMate.previousSuit = ""
	player.previousSuit = ""

	card := player.playerTactic(&s, validCards)
	notExpSuit := HEART
	if card.Suit == notExpSuit {
		t.Errorf("FOREHAND, DECLARER BACK, not expected suit of 2 numbers: %v, played %v",
			validCards, card)
	}
}

func TestOpponentTacticFORE_short_TOD_SUENDE_1_not_a_choice(t *testing.T) {
	// FOREHAND

	// if declarer BACK short

	// never 2 numbers, or D number suit
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer

	teamMate.previousSuit = ""
	player.previousSuit = ""
	s.trump = CLUBS
	s.trick = []Card{}
	_ = teamMate

	validCards := []Card{
		Card{HEART, "D"},
		Card{HEART, "8"},

		Card{SPADE, "D"},
		Card{SPADE, "7"},
	}

	// declarer BACK
	s.opp1 = &player
	s.opp2 = &teamMate
	teamMate.previousSuit = ""
	player.previousSuit = ""

	card := player.playerTactic(&s, validCards)
	if card.Suit == "" && card.Rank == "" {
		t.Errorf("Error validCards: %v, played %v",
			validCards, card)
	}
}

func TestOpponentTacticBACK1(t *testing.T) {

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer

	s.trump = CARO
	s.trick = []Card{Card{CLUBS, "9"}, Card{CARO, "A"}}
	s.follow = CLUBS

	validCards := []Card{
		Card{CARO, "D"},
		Card{CARO, "9"},
		Card{CARO, "7"},
		Card{SPADE, "D"},
		Card{HEART, "10"},
		Card{HEART, "D"},
		Card{HEART, "8"},
	}
	//

	card := player.playerTactic(&s, validCards)
	if getSuit(s.trump, card) == s.trump {
		t.Errorf("In trick %v, TRUMPS CARO, and valid %v, NOT expected a trump: %v",
			s.trick, validCards, card)
	}
}

func TestOpponentTacticBACK2(t *testing.T) {

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer

	s.trump = CARO
	s.trick = []Card{Card{CLUBS, "10"}, Card{CLUBS, "7"}}
	s.follow = CLUBS

	validCards := []Card{
		Card{CARO, "10"},
		Card{SPADE, "D"},
		Card{HEART, "K"},
		Card{HEART, "8"},
	}
	//

	card := player.playerTactic(&s, validCards)
	if getSuit(s.trump, card) == s.trump {
		t.Errorf("In trick %v, TRUMPS CARO, and valid %v, NOT expected a trump: %v",
			s.trick, validCards, card)
	}

	exp := Card{HEART, "K"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, TRUMPS CARO, and valid %v, expected : %v, got %v",
			s.trick, validCards, exp, card)
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

func TestDeclarerTacticFORE0(t *testing.T) {
	// don't play your A-10 trumps if Js still there
	validCards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CLUBS, "K"},
		Card{CLUBS, "D"},
	}

	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{}


	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "K"},
		Card{CLUBS, "D"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
	}

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "J"}
	if !card.equals(exp) {
		t.Errorf("Trump: CLUBS, In trick %v and valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticFORE1(t *testing.T) {
	// don't play your A-10 trumps if Js still there

	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{}
	validCards := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "7"},
		Card{HEART, "A"},
	}

	s.trumpsInGame = []Card{Card{CLUBS, "J"}}

	card := player.playerTactic(&s, validCards)
	unexp1 := Card{CLUBS, "A"}
	unexp2 := Card{CLUBS, "10"}
	if card.equals(unexp1) || card.equals(unexp2) {
		t.Errorf("Trump: CLUBS, In trick %v and valid %v, not expected to play %v since still in game are: %v",
			s.trick, validCards, card, s.trumpsInGame)
	}
}

func TestDeclarerTacticFORE2(t *testing.T) {
	// don't play a value trump if higher trumps
	// in game

	validCards := []Card{
		Card{CLUBS, "D"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
		Card{SPADE, "A"},
		Card{SPADE, "9"},
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{HEART, "7"},
	}
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{}


	s.trumpsInGame = []Card{Card{SPADE, "J"}}

	card := player.playerTactic(&s, validCards)
	unexp1 := Card{CLUBS, "D"}
	if card.equals(unexp1) {
		t.Errorf("Trump: CLUBS, In trick %v and valid %v, not expected to play %v since still in game are: %v",
			s.trick, validCards, card, s.trumpsInGame)
	}
}

func TestDeclarerTacticFORE3(t *testing.T) {
	// BUT play your A-10 trumps if Js ARE NOT still there

	validCards := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "7"},
	}
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{}

	s.trumpsInGame = []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "7"},
		Card{CLUBS, "K"},
	}

	card := player.playerTactic(&s, validCards)
	exp := Card{CLUBS, "A"}
	if !card.equals(exp) {
		t.Errorf("Trump: CLUBS, In trick %v and valid %v, was expected to play %v since still in game are: %v",
			s.trick, validCards, exp, s.trumpsInGame)
	}
}

func TestDeclarerTacticFORE4(t *testing.T) {
	// BUT play your A-10 trumps if Js ARE NOT still there

	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{}
	validCards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "A"},
		Card{HEART, "A"},
		Card{HEART, "7"},
	}
	player.hand = validCards

	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "A"},
	}

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("Trump: CLUBS, In trick %v and valid %v, was expected to play %v since still in game are: %v. Played %v",
			s.trick, validCards, exp, s.trumpsInGame, card)
	}
}

func TestDeclarerTacticFORE5(t *testing.T) {
	// don't play a trump if opponents have many

	validCards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{SPADE, "7"},
		Card{HEART, "8"},
		Card{CARO, "D"},
	}
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = GRAND
	s.trick = []Card{}

	s.trumpsInGame = []Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
	}

	card := player.playerTactic(&s, validCards)
	notexp := Card{CLUBS, "J"}
	if card.equals(notexp) {
		t.Errorf("Trump: CLUBS, In trick %v and valid %v, was NOT expected to play %v since still in game are: %v. Played %v",
			s.trick, validCards, notexp, s.trumpsInGame, card)
	}
}

func TestDeclarerTacticFORE_LowTrump(t *testing.T) {
	// don't play a high trump if you are not strong

	validCards := []Card{
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{SPADE, "10"},
		Card{SPADE, "8"},
		Card{SPADE, "7"},
		Card{HEART, "A"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "D"},
	}
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = SPADE
	s.trick = []Card{}

	s.trumpsInGame = makeTrumpDeck(SPADE)

	card := player.playerTactic(&s, validCards)
	exp := Card{SPADE, "7"}
	if !card.equals(exp) {
		t.Errorf("Trump: CLUBS, In trick %v and valid %v, was expected to play %v. Played %v",
			s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticAKX(t *testing.T) {
	// BUT play your A-10 trumps if Js ARE NOT still there

	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{}
	player.hand = []Card{
		Card{CLUBS, "J"},
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{HEART, "7"},
	}

	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
	}

	card := player.playerTactic(&s, player.hand)
	exp := Card{HEART, "7"}
	if !card.equals(exp) {
		t.Errorf("A-K-X tactic. Trump: CLUBS, In trick %v and hand %v, was expected to play %v since still in game are: %v. Played %v",
			s.trick, player.hand, exp, s.trumpsInGame, card)
	}

	s.cardsPlayed = append(s.cardsPlayed, Card{HEART, "10"})

	card = player.playerTactic(&s, player.hand)
	exp = Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("A-K-X tactic. 10 already player. Trump: CLUBS, In trick %v and hand %v, was expected to play %v since still in game are: %v. Played %v",
			s.trick, player.hand, exp, s.trumpsInGame, card)
	}
}

func TestDeclarerTacticKX(t *testing.T) {

	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{}
	player.hand = []Card{
		Card{CLUBS, "J"},
		Card{HEART, "K"},
		Card{HEART, "7"},
	}

	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
	}

	card := player.playerTactic(&s, player.hand)
	exp := Card{HEART, "7"}
	if !card.equals(exp) {
		t.Errorf("K-X tactic. Trump: CLUBS exhausted, In trick %v and hand %v, was expected to play %v since still in game are: A 10. Played %v",
			s.trick, player.hand, exp, card)
	}

	s.cardsPlayed = append(s.cardsPlayed, Card{HEART, "10"}, Card{HEART, "A"})

	card = player.playerTactic(&s, player.hand)
	exp = Card{HEART, "K"}
	if !card.equals(exp) {
		t.Errorf("K-X tactic. Trump: CLUBS, In trick %v and hand %v, was expected to play %v since still in game are: %v. Played %v",
			s.trick, player.hand, exp, s.trumpsInGame, card)
	}
}

func TestDeclarerTacticKXLessValuableLoser(t *testing.T) {

	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{}
	player.hand = []Card{
		Card{CLUBS, "J"},
		Card{HEART, "K"},
		Card{HEART, "D"},
		Card{CARO, "10"},
		Card{CARO, "7"},
	}

	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
	}

	card := player.playerTactic(&s, player.hand)
	exp := Card{CARO, "7"}
	if !card.equals(exp) {
		t.Errorf("K-X tactic. Trump: CLUBS exhausted, In trick %v and hand %v, was expected to play %v since still in game are: A 10. Played %v",
			s.trick, player.hand, exp, card)
	}
}

// DOES NOT IMPROVE ???!!!!

// func TestDeclarerTacticLosingTrickLow(t *testing.T) {

// 	player := makePlayer([]Card{})
// 	s := makeSuitState()
// 	s.leader = &player
// 	s.declarer = &player

// 	s.trump = CLUBS
// 	s.trick = []Card{}
// 	player.hand = []Card{
// 		Card{CLUBS, "K"},
// 		Card{SPADE, "10"},
// 		Card{SPADE, "D"},
// 	}

// 	s.trumpsInGame = []Card{
// 		Card{CARO, "J"},
// 	}

// 	suit := SPADE
// 	s.cardsPlayed = []Card{
// 		Card{suit, "A"},
// 		Card{suit, "K"},
// 		Card{suit, "8"},
// 		Card{suit, "7"},
// 	}

// 	card := player.playerTactic(&s, player.hand)
// 	exp := Card{SPADE, "D"}
// 	if !card.equals(exp) {
// 		t.Errorf("In trick %v and hand %v, was expected to play %v. Played %v",
// 			s.trick, player.hand, exp, card)
// 	}
// }

func TestDeclarerTacticDoNotTrumpZeroValueTricks(t *testing.T) {
	// BUT play your A-10 trumps if Js ARE NOT still there

	other := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{
		Card{HEART, "7"},
		Card{HEART, "8"},
	}
	s.follow = HEART

	player.hand = []Card{
		Card{CLUBS, "J"},
		Card{CARO, "A"},
		Card{CARO, "K"},
		Card{CARO, "7"},
	}

	card := player.playerTactic(&s, player.hand)
	exp := Card{CARO, "7"}
	if !card.equals(exp) {
		t.Errorf("Zero value trick. In trick %v and hand %v, was expected to play %v . Played %v",
			s.trick, player.hand, exp, card)
	}

}

func TestDeclarerTacticKeepTheAForThe10(t *testing.T) {
	other := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{
		Card{HEART, "7"},
		Card{HEART, "K"},
	}
	s.follow = HEART

	player.hand = []Card{
		//	Card{HEART, "J"},
		Card{HEART, "A"},
		Card{HEART, "D"},
		Card{HEART, "8"},
	}

	s.cardsPlayed = []Card{}

	card := player.playerTactic(&s, player.hand)
	exp := Card{HEART, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and hand %v, and 10 still in game it was expected to play %v . Played %v",
			s.trick, player.hand, exp, card)
	}
}

func TestDeclarerTacticKeepTheAForThe10_2(t *testing.T) {
	other := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = CLUBS
	s.trick = []Card{
		Card{HEART, "7"},
		Card{HEART, "K"},
	}
	s.follow = HEART

	player.hand = []Card{
		//	Card{HEART, "J"},
		Card{HEART, "A"},
		Card{HEART, "D"},
		Card{HEART, "8"},
	}

	s.cardsPlayed = []Card{}

	s.cardsPlayed = append(s.cardsPlayed, Card{HEART, "10"})
	card := player.playerTactic(&s, player.hand)
	exp := Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and hand %v, and 10 played, it was expected to play %v . Played %v",
			s.trick, player.hand, exp, card)
	}

	s.trick = []Card{
		Card{HEART, "7"},
		Card{HEART, "10"},
	}
	s.cardsPlayed = []Card{}
	card = player.playerTactic(&s, player.hand)
	exp = Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and hand %v, it was expected to play %v . Played %v",
			s.trick, player.hand, exp, card)
	}
}

func TestDeclarerTacticKeepTheAForThe10_1(t *testing.T) {

	other := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = CLUBS
	s.follow = HEART

	player.hand = []Card{
		//	Card{HEART, "J"},
		Card{HEART, "A"},
		Card{HEART, "D"},
		Card{HEART, "8"},
	}

	s.cardsPlayed = []Card{}

	s.trick = []Card{
		Card{HEART, "7"},
		Card{HEART, "9"},
	}
	s.cardsPlayed = []Card{}
	card := player.playerTactic(&s, player.hand)
	exp := Card{HEART, "D"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and hand %v, it was expected to play %v . Played %v",
			s.trick, player.hand, exp, card)
	}

}

func TestDeclarerTacticKeepTheAForThe10_3(t *testing.T) {
	other := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = CLUBS
	s.follow = HEART

	player.hand = []Card{
		//	Card{HEART, "J"},
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{HEART, "8"},
	}

	s.cardsPlayed = []Card{}

	s.trick = []Card{
		Card{HEART, "7"},
		Card{HEART, "9"},
	}
	s.cardsPlayed = []Card{}
	card := player.playerTactic(&s, player.hand)
	exp := Card{HEART, "K"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and hand %v, it was expected to play %v . Played %v",
			s.trick, player.hand, exp, card)
	}

}

func TestDeclarerTacticGrand(t *testing.T) {
	//other := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = GRAND

	player.hand = []Card{
		Card{HEART, "J"},
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "8"},
	}

	s.cardsPlayed = []Card{}

	s.trick = []Card{}
	s.cardsPlayed = []Card{}
	card := player.playerTactic(&s, player.hand)
	exp := Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and hand %v, it was expected to play %v . Played %v",
			s.trick, player.hand, exp, card)
	}

	player.hand = []Card{
		Card{HEART, "J"},
		Card{HEART, "10"},
		Card{HEART, "8"},
	}
	s.cardsPlayed = []Card{Card{HEART, "A"}}

	card = player.playerTactic(&s, player.hand)
	exp = Card{HEART, "10"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and hand %v, A played, it was expected to play %v . Played %v",
			s.trick, player.hand, exp, card)
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

func TestDiscardInSkat(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "D"},
		Card{CLUBS, "9"},
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{HEART, "A"},
		Card{HEART, "D"},
		Card{HEART, "8"},
		Card{CARO, "A"},
	}
	skat := []Card{Card{CARO, "7"}, Card{SPADE, "D"}}
	p := makePlayer(cards)
	p.discardInSkat(skat)

	if in(skat, Card{SPADE, "A"}) || in(skat, Card{HEART, "A"}) || in(skat, Card{CARO, "A"}) {
		t.Errorf("A discarded in SKAT: %v", skat)
	}
}

func TestDiscardInSkatGRAND(t *testing.T) {
	cards := []Card{
		Card{HEART, "J"},  //L1
		Card{CLUBS, "A"},  //1
		Card{CLUBS, "9"},  //L2
		Card{SPADE, "A"},  //2
		Card{SPADE, "9"},  //DISCARD
		Card{SPADE, "8"},  //L3
		Card{HEART, "A"},  //3
		Card{HEART, "10"}, //4
		Card{HEART, "K"},
		Card{HEART, "9"},
		Card{CARO, "10"}, // DISCARD
	}
	skat := []Card{Card{SPADE, "9"}, Card{HEART, "9"}}
	p := makePlayer(cards)
	p.discardInSkat(skat)

	if in(skat, Card{HEART, "10"}) {
		t.Errorf("%v from hand: %v discarded in SKAT: %v", Card{HEART, "10"}, p.hand, skat)
	}
}

func TestDiscardInSkatGRANDBlank(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"}, //LOSER
		Card{CLUBS, "A"},
		Card{CLUBS, "8"}, //LOSER
		Card{CLUBS, "9"}, //LOSER
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "D"}, //LOSER
		Card{SPADE, "9"}, //LOSER
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "K"},
		Card{CARO, "9"}, //LOSER
	}
	skat := []Card{Card{CLUBS, "9"}, Card{SPADE, "D"}}
	p := makePlayer(cards)
	p.trumpToDeclare = GRAND
	p.discardInSkat(skat)

	if !in(skat, Card{SPADE, "D"}, Card{CARO, "9"}) {
		t.Errorf("hand: %v discarded in SKAT: %v", p.hand, skat)
	}
}

func TestDiscardInSkatAllTrumps(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "K"},
		Card{CLUBS, "D"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
		Card{SPADE, "A"},
	}
	skat := make([]Card, 2)
	p := makePlayer(cards)
	p.discardInSkat(skat)

	if in(skat, Card{SPADE, "A"}) || in(skat, Card{CLUBS, "A"}) {
		t.Errorf("A discarded in SKAT: %v", skat)
	}

	if in(skat, Card{CLUBS, "J"}) || in(skat, Card{SPADE, "J"}) || in(skat, Card{HEART, "J"}) || in(skat, Card{CARO, "J"}) {
		t.Errorf("J discarded in SKAT: %v", skat)
	}

	if !in(skat, Card{CLUBS, "7"}, Card{CLUBS, "8"}) {
		t.Errorf("Wrong discarded: %v", skat)
	}
}

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

func TestCanWinGRAND(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},

		Card{CLUBS, "A"},
		Card{CLUBS, "10"},

		Card{CLUBS, "9"}, // LOSER
		Card{CLUBS, "8"}, // LOSER
		Card{CLUBS, "7"}, // LOSER
		Card{SPADE, "A"},
		Card{SPADE, "10"},
	}
	p := makePlayer(cards)

	canWin := p.canWin()
	if canWin != "GRAND" {
		t.Errorf("Hand %v can win GRAND. Got: %v", p.hand, canWin)
	}
}

func TestCanWinGRAND2(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},

		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "K"},

		Card{CLUBS, "8"}, // LOSER
		Card{CLUBS, "7"}, // LOSER
		Card{SPADE, "10"},
		Card{SPADE, "9"}, // LOSER
		Card{SPADE, "D"}, // LOSER
	}
	p := makePlayer(cards)
	p.risky = true

	canWin := p.canWin()
	if canWin != "GRAND" {
		t.Errorf("Hand %v can win GRAND. Got: %v", p.hand, canWin)
	}
}

func TestGrandLosers(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "K"},
		Card{CLUBS, "9"}, //LOSER
		Card{SPADE, "A"},
		Card{SPADE, "K"},  //LOSER
		Card{HEART, "10"}, //LOSER
	}

	losers := grandLosers(cards)
	if !in(losers, Card{CLUBS, "9"}, Card{SPADE, "K"}, Card{HEART, "10"}) {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersCount(t *testing.T) {
	cs := []Card{
		Card{SPADE, "J"}, //LOSER
		Card{CLUBS, "A"},
		Card{CLUBS, "8"}, //LOSER
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "9"}, //LOSER
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "K"},
		Card{CARO, "9"}, //LOSER
	}

	losers := len(grandLosers(cs)) + jackLosers(cs)
	if losers != 4 {
		t.Errorf("Losers: %v, Expected 4", losers)
	}
}

func TestGrandLosersJ0(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJNo(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
	}

	losers := jackLosers(cards)
	if losers != 0 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ1(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
	}

	losers := jackLosers(cards)
	if losers != 0 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ2(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ3(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ4(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"},
		Card{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 2 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ5(t *testing.T) {
	cards := []Card{
		Card{HEART, "J"},
		Card{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 2 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ6(t *testing.T) {
	cards := []Card{
		Card{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 1 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestGrandLosersJ7(t *testing.T) {
	cards := []Card{
		Card{HEART, "J"},
		Card{CARO, "J"},
	}

	losers := jackLosers(cards)
	if losers != 2 {
		t.Errorf("Losers: %v", losers)
	}
}

func TestOpponentTacticNULLBack(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	//teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.declarer = &otherPlayer
	s.trump = NULL
	s.follow = HEART
	s.trick = []Card{Card{HEART, "8"}, Card{HEART, "10"}}

	validCards := []Card{
		Card{HEART, "J"},
		Card{HEART, "D"},
		Card{HEART, "A"},
	}

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("NULL, In trick %v and valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}
}

// NOT SURE IF CORRECT
// func TestOpponentTacticNULLMIDDeclFore(t *testing.T) {
// 	otherPlayer := makePlayer([]Card{})
// 	//teamMate := makePlayer([]Card{})
// 	player := makePlayer([]Card{})
// 	s := makeSuitState()

// 	s.declarer = &otherPlayer
// 	s.leader = &otherPlayer
// 	//s.opp1 = &player

// 	s.trump = NULL
// 	s.follow = HEART
// 	s.trick = []Card{Card{HEART, "10"}}

// 	validCards := []Card{
// 		Card{HEART, "J"},
// 		Card{HEART, "D"},
// 		Card{HEART, "A"},
// 	}

// 	card := player.playerTactic(&s, validCards)
// 	exp := Card{HEART, "A"}
// 	if !card.equals(exp) {
// 		t.Errorf("NULL, In trick %v and valid %v, expected to play %v, played %v",
// 			s.trick, validCards, exp, card)
// 	}
// }

func TestOpponentTacticNULLMIDDeclBack4(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()

	s.declarer = &otherPlayer
	s.leader = &teamMate

	s.trump = NULL
	s.follow = HEART
	s.trick = []Card{Card{HEART, "J"}}

	validCards := []Card{
		Card{HEART, "8"},
		Card{HEART, "9"},
		Card{HEART, "D"},
		Card{HEART, "A"},
	}

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "D"}
	if !card.equals(exp) {
		t.Errorf("NULL, In trick %v and valid %v, expected to play (next) %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticNULLMIDDeclBack2(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()

	s.declarer = &otherPlayer
	s.leader = &teamMate
	//s.opp1 = &player

	s.cardsPlayed = []Card{}
	s.trump = NULL
	s.follow = HEART
	s.trick = []Card{Card{HEART, "9"}}

	validCards := []Card{
		Card{HEART, "10"},
		Card{HEART, "K"},
	}

	// Changed:
	// declarer surely has 7 and 8, throw off the highest of the suit

	// throw previous card same value

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "10"}
	if !card.equals(exp) {
		t.Errorf("NULL, In trick %v and valid %v, cards played: %v expected to throw off %v, played %v",
			s.trick, validCards, s.cardsPlayed, exp, card)
	}

	s.cardsPlayed = []Card{Card{HEART, "7"}, Card{HEART, "8"}}

	card = player.playerTactic(&s, validCards)
	exp = Card{HEART, "10"}
	if !card.equals(exp) {
		t.Errorf("NULL, In trick %v and valid %v, cards played: %v expected to play %v, played %v",
			s.trick, validCards, s.cardsPlayed, exp, card)
	}
}

func TestOpponentTacticNULLMIDDeclBack(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()

	s.opp2 = &player
	s.declarer = &otherPlayer
	s.leader = &teamMate
	s.trump = NULL
	s.follow = HEART
	s.trick = []Card{Card{HEART, "10"}}

	validCards := []Card{
		Card{HEART, "J"},
		Card{HEART, "D"},
		Card{HEART, "A"},
	}
	s.cardsPlayed = []Card{}

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "J"}
	if !card.equals(exp) {
		t.Errorf("NULL MID, Decl at Back, In trick %v and valid %v, with cards played: %v, expected to play %v, played %v",
			s.trick, validCards, s.cardsPlayed, exp, card)
	}
}

func TestOpponentTacticNULLMIDDeclBack1(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()

	s.opp2 = &player
	s.declarer = &otherPlayer
	s.leader = &otherPlayer
	s.trump = NULL
	s.follow = HEART
	s.trick = []Card{Card{HEART, "K"}, Card{HEART, "8"}}

	validCards := []Card{
		Card{HEART, "9"},
		Card{HEART, "A"},
	}
	s.cardsPlayed = []Card{}

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "9"}
	if !card.equals(exp) {
		t.Errorf("NULL BACK, Decl at Fore, In trick %v and valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticNULLMIDDeclBack3(t *testing.T) {
	teamMate := makePlayer([]Card{})
	otherPlayer := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()

	s.declarer = &otherPlayer
	s.leader = &teamMate
	s.trump = NULL
	s.follow = HEART
	s.trick = []Card{Card{HEART, "D"}, Card{HEART, "9"}}

	validCards := []Card{
		Card{HEART, "8"},
		Card{HEART, "10"},
		Card{HEART, "A"},
	}
	s.cardsPlayed = []Card{}

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("NULL BACK, Decl at MID, In trick %v and valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticNULLFORE(t *testing.T) {
	teamMate := makePlayer([]Card{})
	otherPlayer := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.declarer = &otherPlayer
	s.opp1 = &teamMate
	s.opp2 = &player
	s.trump = NULL

	validCards := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "8"},
		Card{CLUBS, "10"},
		Card{SPADE, "A"},
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "8"},
		Card{CARO, "9"},
	}

	teamMate.previousSuit = "HEART"
	s.trick = []Card{}

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "8"}
	if !card.equals(exp) {
		t.Errorf("NULL, VOID: %v, In trick %v and valid %v, expected to follow previous trick and play %v, played %v",
			s.declarerVoidSuit, s.trick, validCards, exp, card)
	}

	s.declarerVoidSuit[HEART] = true
	card = player.playerTactic(&s, validCards)
	exp = Card{SPADE, "A"}
	if !card.equals(exp) {
		t.Errorf("NULL, VOID: %v,In trick %v and valid %v, decl HEART void, expected to play SINGLETON %v, played %v",
			s.declarerVoidSuit, s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticNULLFORE1(t *testing.T) {
	teamMate := makePlayer([]Card{})
	otherPlayer := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.declarer = &otherPlayer
	s.leader = &player
	s.opp1 = &teamMate
	s.opp2 = &player
	s.trump = NULL

	validCards := []Card{
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "8"},
		Card{CARO, "9"},
		Card{CARO, "J"},
	}

	teamMate.previousSuit = "SPADE"
	s.trick = []Card{}

	s.cardsPlayed = []Card{
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{HEART, "D"},
		Card{HEART, "B"},
		Card{HEART, "7"},
	}

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "9"}
	if !card.equals(exp) {
		t.Errorf("NULL, In trick %v, played %v and valid %v, expected play %v, played %v",
			s.trick, s.cardsPlayed, validCards, exp, card)
	}

	if !isVoid(&s, validCards, HEART) {
		t.Errorf("Heart IS void")
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
	s := makeSuitState()
	s.declarer = &player
	s.trump = NULL
	s.follow = HEART
	s.trick = []Card{Card{HEART, "10"}}

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

func TestStrongestLowest1(t *testing.T) {
	s := makeSuitState()
	s.trump = CLUBS

	cs := []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
		Card{CLUBS, "10"},
		Card{SPADE, "9"},
		Card{CARO, "9"},
	}
	p := makePlayer(cs)

	act := p.strongestLowestNotAor10(&s, cs) 
	exp := Card{CLUBS, "J"}
	if act != exp {
		t.Errorf("Error strongestLowest, exp: %v, found %v", exp, act)
	}
}

func TestStrongestLowest2(t *testing.T) {
	s := makeSuitState()
	s.trump = CLUBS

	cs := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{CARO, "J"},
		Card{SPADE, "9"},
		Card{CARO, "9"},
	}
	p := makePlayer(cs)

	act := p.strongestLowestNotAor10(&s, cs) 
	exp := Card{SPADE, "J"}
	if act != exp {
		t.Errorf("Error strongestLowest, exp: %v, found %v", exp, act)
	}
}

func TestStrongestLowestNotA(t *testing.T) {
	s := makeSuitState()
	s.trump = CLUBS

	cs := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "K"},
	}
	p := makePlayer(cs)

	act := p.strongestLowestNotAor10(&s, cs) 
	exp := Card{CARO, "J"}
	if act != exp {
		t.Errorf("Error strongestLowest, exp: %v, found %v", exp, act)
	}
}

func TestStrongestLowestNot10(t *testing.T) {
	s := makeSuitState()
	s.trump = CLUBS

	cs := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "10"},
		Card{CLUBS, "D"},
	}
	p := makePlayer(cs)

	act := p.strongestLowestNotAor10(&s, cs) 
	exp := Card{CARO, "J"}
	if act != exp {
		t.Errorf("Error strongestLowest, exp: %v, found %v", exp, act)
	}
}

func TestStrongestLowestA_IfOnly(t *testing.T) {
	s := makeSuitState()
	s.trump = CLUBS

	cs := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "D"},
	}
	p := makePlayer(cs)

	act := p.strongestLowestNotAor10(&s, cs) 
	exp := Card{CLUBS, "A"}
	if act != exp {
		t.Errorf("Error strongestLowest, exp: %v, found %v", exp, act)
	}
}

func TestStrongestLowestinSkat(t *testing.T) {
	s := makeSuitState()
	s.trump = CLUBS
	s.skat = []Card{Card{CLUBS, "10"}}

	cs := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "K"},
		Card{CLUBS, "9"},
	}
	p := makePlayer(cs)

	act := p.strongestLowestNotAor10(&s, cs) 
	exp := Card{CLUBS, "K"}
	if act != exp {
		t.Errorf("Error strongestLowest, exp: %v, found %v", exp, act)
	}
}

func TestStrongestLowestBug1(t *testing.T) {
	s := makeSuitState()
	s.trump = CLUBS
	//s.cardsPlayed = []Card{Card{CLUBS, "8"}}
	s.trumpsInGame = []Card{Card{CLUBS, "D"}}

	cs := []Card{
		Card{CLUBS, "9"},
		Card{CLUBS, "7"},
	}
	p := makePlayer(cs)

	act := p.strongestLowestNotAor10(&s, cs) 
	exp := Card{CLUBS, "9"}
	if act != exp {
		t.Errorf("Error strongestLowest, exp: %v, found %v", exp, act)
	}
}

func TestNoHigherCard(t *testing.T) {
	s := makeSuitState()
	s.skat = []Card{}
	s.trump = CLUBS

	act := noHigherCard(&s, true, []Card{}, Card{CARO, "J"})
	exp := false

	if act != exp {
		t.Errorf("There are still other higher cards in play")
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