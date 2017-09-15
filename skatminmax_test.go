package main

import (
	"testing"
	"time"
	"github.com/dranidis/go-skat/game"
	"github.com/dranidis/go-skat/game/minimax"
)

func TestFindNextStateEndOfTrick(t *testing.T) {
	dist := make([][]Card, 3)
	dist[0] = []Card {
		// Card{SPADE, "J"},
		Card{SPADE, "A"},
		Card{SPADE, "K"},
		Card{SPADE, "9"},
	}
	dist[1] = []Card {
		// Card{SPADE, "8"},
		Card{CLUBS, "A"},
		Card{HEART, "8"},
		Card{CLUBS, "7"},
	}
	dist[2] = []Card {
		Card{SPADE, "10"},
		Card{SPADE, "7"},
		Card{CLUBS, "8"},
		Card{CLUBS, "9"},
	}

	skatState := SkatState{
		CARO, // trump
		dist, 			
		[]Card{Card{SPADE, "J"}, Card{SPADE, "8"}}, // trick 
		0, // declarer 
		2, // who's turn is it
		40, 
		42,
		false,
	}

	action := SkatAction{Card{SPADE, "10"}}

	var skatStateP game.State
	skatStateP = &skatState
	var a game.Action
	a = action

	skatStateN := skatStateP.FindNextState(a)
	newState := skatStateN.(*SkatState)
	debugTacticsLog("State: %v\n", newState)

	if newState.declScore != skatState.declScore + 12 {
		t.Errorf("Wrong declScore score. Is: %d", newState.declScore)
	}	
	if newState.oppScore != skatState.oppScore {
		t.Errorf("Wrong oppScore score. Is: %d", newState.oppScore)
	}	
	if newState.turn != 0 {
		t.Errorf("Wrong winner. Is: %d", newState.turn)
	}	
	

}

func TestFindNextStateEndOfTrick2(t *testing.T) {
	dist := make([][]Card, 3)
	dist[0] = []Card {
		// Card{SPADE, "J"},
		// Card{SPADE, "A"},
		Card{SPADE, "K"},
		Card{SPADE, "9"},
	}
	dist[1] = []Card {
		// Card{SPADE, "8"},
		// Card{CLUBS, "A"},
		Card{HEART, "8"},
		Card{CLUBS, "7"},
	}
	dist[2] = []Card {
		Card{SPADE, "10"},
		Card{SPADE, "7"},
		Card{CLUBS, "8"},
		Card{CLUBS, "9"},
	}

	skatState := SkatState{
		CARO, // trump
		dist, 			
		[]Card{{SPADE, "A"}, Card{CLUBS, "A"}}, // trick 
		0, // declarer 
		2, // who's turn is it
		52, 
		42,
		false,
	}

	action := SkatAction{Card{SPADE, "7"}}

	var skatStateP game.State
	skatStateP = &skatState
	var a game.Action
	a = action

	skatStateN := skatStateP.FindNextState(a)
	newState := skatStateN.(*SkatState)
	debugTacticsLog("State: %v\n", newState)

	if newState.declScore != skatState.declScore + 22 {
		t.Errorf("Wrong declScore score. Is: %d", newState.declScore)
	}	
	if newState.oppScore != skatState.oppScore {
		t.Errorf("Wrong oppScore score. Is: %d", newState.oppScore)
	}	
	if newState.turn != 0 {
		t.Errorf("Wrong winner. Is: %d", newState.turn)
	}	
}

func TestFindNextStateEndOfTrick3(t *testing.T) {
	dist := make([][]Card, 3)
	dist[0] = []Card {
		// Card{SPADE, "J"},
		// Card{SPADE, "A"},
		// Card{SPADE, "K"},
		Card{SPADE, "9"},
	}
	dist[1] = []Card {
		// Card{SPADE, "8"},
		// Card{CLUBS, "A"},
		// Card{HEART, "8"},
		Card{CLUBS, "7"},
	}
	dist[2] = []Card {
		Card{SPADE, "10"},
		// Card{SPADE, "7"},
		// Card{CLUBS, "8"},
		Card{CLUBS, "9"},
	}

	skatState := SkatState{
		CARO, // trump
		dist, 			
		[]Card{{SPADE, "K"}, Card{SPADE, "8"}}, // trick 
		0, // declarer 
		2, // who's turn is it
		52, 
		42,
		false,
	}

	action := SkatAction{Card{SPADE, "10"}}

	var skatStateP game.State
	skatStateP = &skatState
	var a game.Action
	a = action

	skatStateN := skatStateP.FindNextState(a)
	newState := skatStateN.(*SkatState)
	debugTacticsLog("State: %v\n", newState)

	if newState.declScore != skatState.declScore {
		t.Errorf("Wrong declScore score. Is: %d", newState.declScore)
	}	
	if newState.oppScore != skatState.oppScore + 14{
		t.Errorf("Wrong oppScore score. Is: %d", newState.oppScore)
	}	
	if newState.turn != 2 {
		t.Errorf("Wrong winner. Is: %d", newState.turn)
	}	
}

func TestFindNextStateNewTrick(t *testing.T) {
	dist := make([][]Card, 3)
	dist[0] = []Card {
		// Card{SPADE, "J"},
		// Card{SPADE, "A"},
		Card{SPADE, "K"},
		Card{SPADE, "9"},
	}
	dist[1] = []Card {
		// Card{SPADE, "8"},
		// Card{CLUBS, "A"},
		Card{HEART, "8"},
		Card{CLUBS, "7"},
	}
	dist[2] = []Card {
		Card{SPADE, "10"},
		// Card{SPADE, "7"},
		Card{CLUBS, "8"},
		Card{CLUBS, "9"},
	}

	skatState := SkatState{
		CARO, // trump
		dist, 			
		[]Card{{SPADE, "A"}, Card{CLUBS, "A"}, Card{SPADE, "7"}}, // trick 
		0, // declarer 
		0, // who's turn is it
		52, 
		42,
		false,
	}

	action := SkatAction{Card{SPADE, "K"}}

	var skatStateP game.State
	skatStateP = &skatState
	var a game.Action
	a = action

	skatStateN := skatStateP.FindNextState(a)
	newState := skatStateN.(*SkatState)
	debugTacticsLog("State: %v\n", newState)

	if len(newState.trick) != 1 {
		t.Errorf("Wrong trick size. Is: %d %v", len(newState.trick), newState.trick)
	}	
	if !in(newState.trick, Card{SPADE, "K"}) {
		t.Errorf("Wrong trick. Is: %v", newState.trick)
	}	
	if newState.turn != 1 {
		t.Errorf("Wrong turn. Is: %d", newState.turn)
	}	
}

func TestIsOpponentTurn(t *testing.T) {
	dist := make([][]Card, 3)

	s := SkatState{
		CARO, // trump
		dist, 			
		[]Card{}, // trick 
		0, // declarer 
		0, // who's turn is it
		40, 
		42,
		false,
	}	

	if s.IsOpponentTurn() { // YOU
		t.Errorf("Error opponent")
	}
	s.turn = 1 // OPP1
	if !s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.turn = 2 // OPP2
	if !s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.declarer = 1
	s.turn = 1
	if !s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.turn = 0 // YOU
	if s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.turn = 2 // YOUR PARTNER
	if s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.declarer = 2
	s.turn = 2
	if !s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.turn = 0 // YOU 
	if s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.turn = 1 // YOUR PARTNET
	if s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
}


func TestCopySkatState(t *testing.T) {
	dist := make([][]Card, 3)
	dist[0] = []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{HEART, "A"},
	}
	dist[1] = []Card{
		Card{HEART, "10"},
		Card{CARO, "8"},
		Card{SPADE, "A"},
	}
	dist[2] = []Card{
		Card{CLUBS, "9"},
		Card{CARO, "10"},
		Card{CARO, "K"},
	}

	skatState := SkatState{
		CLUBS,
		dist, 			
		[]Card{}, 
		0, 
		0, 
		30, 
		45,
		false,
	}

	ss := copySkatState(skatState)

	// fmt.Printf("%v\n", skatState)
	// fmt.Printf("%v\n", ss)
	if len(ss.playerHand[0]) != 3 && len(ss.playerHand[1]) != 3 && len(ss.playerHand[2]) != 3 {
		t.Errorf("ERROR Copy: %v", ss)
	}
	if !in(ss.playerHand[0], skatState.playerHand[0]...) {
		t.Errorf("ERROR Copy: %v %v", ss.playerHand[0], skatState.playerHand[0])
	}
	if !in(ss.playerHand[1], skatState.playerHand[1]...) {
		t.Errorf("ERROR Copy: %v %v", ss.playerHand[1], skatState.playerHand[1])
	}
	if !in(ss.playerHand[2], skatState.playerHand[2]...) {
		t.Errorf("ERROR Copy: %v %v", ss.playerHand[2], skatState.playerHand[2])
	}
}

func TestFindLegals(t *testing.T) {
	dist := make([][]Card, 3)
	dist[0] = []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "8"},
		Card{HEART, "A"},
	}
	dist[1] = []Card{
		Card{HEART, "10"},
		Card{CARO, "8"},
		Card{SPADE, "A"},
	}
	dist[2] = []Card{
		Card{CARO, "10"},
		Card{CARO, "K"},
	}

	skatState := SkatState{
		CLUBS,
		dist, 			
		[]Card{
			Card{CLUBS, "9"},

		}, 
		0, // declarer
		0, // turn
		30, 
		45,
		false,
	}

	actions := skatState.FindLegals()
	cards := []Card{}
	for _, action := range actions {
		ma := action.(SkatAction)
		cards = append(cards, ma.card)
	}
	// fmt.Println(cards)
	if !in(cards, Card{CLUBS, "J"}, Card{CLUBS, "8"}) {
		t.Errorf("ERROR TestFindLegals: %v", actions)
	}
	if in(cards, Card{HEART, "A"}) {
		t.Errorf("ERROR TestFindLegals: %v", actions)
	}

	ssa := skatState.FindNextState(actions[0])

	ss := ssa.(*SkatState)
	// fmt.Println(skatState)
	// fmt.Println(ss)



	if in(ss.playerHand[0], Card{CLUBS, "J"}) {
		t.Errorf("ERROR TestFindLegals: %v", ss.playerHand[0])
	}

	if !ss.trick[1].equals(Card{CLUBS, "J"}) {
		t.Errorf("ERROR TestFindLegals: %v", ss.trick)
	}

	if ss.turn != 1 {
		t.Errorf("Error turn %d", ss.turn)
	}

	ssa1 := ss.FindNextState(SkatAction{Card{HEART, "10"}})
	// fmt.Println(ssa1)	

	ss1 := ssa1.(*SkatState)

	if ss1.declScore != 42 {
		t.Errorf("Declarer score %d", ss1.declScore)
	} 	
	if ss1.oppScore != skatState.oppScore {
		t.Errorf("Opponent score %d", ss1.declScore)
	} 
}

func TestMoveOne(t *testing.T) {
	p1Hand := []Card{
		Card{CARO, "J"},
		Card{SPADE, "A"},
		}

	p1 := makePlayer(p1Hand)
	p2 := makePlayer([]Card{
		Card{CARO, "10"},
		Card{SPADE, "10"},
		})
	p3 := makePlayer([]Card{
		Card{HEART, "D"},
		Card{CLUBS, "A"},
		})

	players := []PlayerI{&p1, &p2, &p3}

	s := makeSuitState()
	s.trump = CARO
	s.declarer = &p1
	s.leader = &p1
	s.opp1 = &p2
	s.opp2 = &p3

	s.trick = []Card{
	}

	card, newplayers := moveOne(&s, players)

	if len(s.trick) != 1 {
		t.Errorf("Error moveone")
	}

	if !in(s.trick, card) {
		t.Errorf("NOt in trick")
	}

	if newplayers[0] != players[0] {
		t.Errorf("Wrong player order")
	}

	if !in(p1Hand, card) {
		t.Errorf("Wrong card")
	}

}


func TestGetGameSuitStateAndPlayers(t *testing.T) {

	h1 := []Card{
		Card{CARO, "J"},
		Card{SPADE, "A"},
		}
	h2 := []Card{
		Card{CARO, "10"},
		Card{SPADE, "10"},
		}
	h3 := []Card{
		Card{HEART, "D"},
		Card{CLUBS, "A"},
		}

	dist := [][]Card{h1, h2, h3}
	skatState := SkatState{
		CARO, // trump
		dist, 			
		[]Card{}, // trick 
		0, // declarer 
		0, // who's turn is it
		5, 
		6,
		true,
	}

	s, players := skatState.getGameSuitStateAndPlayers()

	debugTacticsLog("Suitstate %v\n", *s)
	debugTacticsLog("Player[0] %v\n", players[0])
	// if players[0].getScore() != 5 {
	// 	t.Errorf("Wrong declarer score: ", players[0].getScore())
	// }
	if s.trump != skatState.trump {
		t.Errorf("Wrong trump")
	}

	if s.declarer != players[0] && s.opp1 != players[1] && s.opp2 != players[2] {
		t.Errorf("Wrong turn order")
	}

	// new test case
	skatState.declarer = 2
	s, players = skatState.getGameSuitStateAndPlayers()

	if s.declarer != players[2] && s.opp1 != players[0] && s.opp2 != players[1] {
		t.Errorf("Wrong turn order, when declarer is 2")
	}

	// new test case
	skatState.declarer = 1
	s, players = skatState.getGameSuitStateAndPlayers()

	if s.declarer != players[1] && s.opp1 != players[2] && s.opp2 != players[0] {
		t.Errorf("Wrong turn order, when declarer is 1")
	}

	// new test case
	skatState.trick = []Card{Card{CARO, "K"}}
	// new test case
	skatState.declarer = 0
	s, players = skatState.getGameSuitStateAndPlayers()

	if s.declarer != players[1] && s.opp1 != players[2] && s.opp2 != players[0] {
		t.Errorf("Wrong turn order, when declarer is 0 and trick len 1")
	}

	// new test case
	skatState.turn = 2
	s, players = skatState.getGameSuitStateAndPlayers()

	if s.declarer != players[2] && s.opp1 != players[0] && s.opp2 != players[1] {
		t.Errorf("Wrong turn order, when declarer is 0 and trick len 1, and turn 2")
	}

}

func TestAlphaBetaTactics2C(t *testing.T) {

	h1 := []Card{
		Card{CARO, "J"},
		Card{SPADE, "A"},
		}
	h2 := []Card{
		Card{CARO, "10"},
		Card{SPADE, "10"},
		}
	h3 := []Card{
		Card{HEART, "D"},
		Card{CLUBS, "A"},
		}

	dist := [][]Card{h1, h2, h3}
	skatState := SkatState{
		CARO, // trump
		dist, 			
		[]Card{}, // trick 
		0, // declarer 
		0, // who's turn is it
		5, 
		15,
		true,
	}

	minimax.DEBUG = true

	var skatStateP game.State
	skatStateP = &skatState
	var a game.Action
	var v float64

	if true {
		a, v = minimax.AlphaBetaTactics(skatStateP)
	}

	debugTacticsLog("Action: %v, Value: %.4f\n", a, v)
	if false {
		t.Errorf("TEST")
	}
}

func TestAlphaBetaTactics2CDef(t *testing.T) {

	h1 := []Card{
		Card{CARO, "J"},
		Card{SPADE, "A"},
		}
	h2 := []Card{
		Card{CARO, "10"},
		Card{SPADE, "10"},
		}
	h3 := []Card{
		Card{HEART, "D"},
		Card{CLUBS, "A"},
		}

	dist := [][]Card{h1, h2, h3}
	skatState := SkatState{
		CARO, // trump
		dist, 			
		[]Card{}, // trick 
		1, // declarer 
		0, // who's turn is it
		0, 
		0,
		true,
	}

	minimax.DEBUG = true

	var skatStateP game.State
	skatStateP = &skatState
	var a game.Action
	var v float64

	if (false) {
 		a, v = minimax.AlphaBetaTactics(skatStateP)
	}
	debugTacticsLog("Action: %v, Value: %.4f\n", a, v)
	if false {
		t.Errorf("TEST")
	}
}

func TestAlphaBetaTactics10C(t *testing.T) {
	dist := make([][]Card, 3)
	dist[0] = []Card {
		Card{SPADE, "J"},
		Card{HEART, "J"},

		Card{CARO, "K"},
		Card{CARO, "10"},
		Card{CARO, "9"},
		Card{CARO, "8"},

		Card{SPADE, "A"},
		Card{SPADE, "K"},
		Card{SPADE, "9"},

		Card{CLUBS, "K"},
	}
	// dist[0] = Shuffle(dist[0])
	dist[1] = []Card {
		Card{CARO, "J"},

		Card{CARO, "A"},
		Card{CARO, "7"},

		Card{CLUBS, "A"},
		Card{CLUBS, "8"},
		Card{CLUBS, "D"},

		Card{SPADE, "8"},
		Card{SPADE, "D"},

		Card{HEART, "8"},
		Card{HEART, "10"},
	}
	// dist[1] = Shuffle(dist[1])
	dist[2] = []Card {
		Card{CLUBS, "J"},

		Card{CARO, "D"},

		Card{SPADE, "10"},
		Card{SPADE, "7"},

		Card{HEART, "K"},
		Card{HEART, "A"},

		Card{CLUBS, "D"},
		Card{CLUBS, "8"},
		Card{CLUBS, "9"},
		Card{CLUBS, "7"},
	}
	// dist[2] = Shuffle(dist[2])

	p1 := makePlayer(dist[0])
	p2 := makePlayer(dist[1])
	p3 := makePlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players = []PlayerI{&p1, &p2, &p3}

	skatState := SkatState{
		CARO, // trump
		dist, 			
		[]Card{}, // trick 
		0, // declarer 
		0, // who's turn is it
		10, // score because of skat  
		0,
		true,
	}

	var skatStateP game.State
	skatStateP = &skatState
	var a game.Action
	var v float64

	minimax.MAXDEPTH = 9


	startWhole := time.Now()

	if (false) {
 		a, v = minimax.AlphaBetaTactics(skatStateP)
	}
	
	ti := time.Now()
	elapsed := ti.Sub(startWhole)		
	debugTacticsLog("TOTAL %v\n", elapsed)

	debugTacticsLog("Action: %v, Value: %.4f\n", a, v)
	if false {
		t.Errorf("TEST")
	}
}
