package main

import (
	"testing"
	"github.com/dranidis/go-skat/game"

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
	if s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.turn = 0 // YOU
	if !s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.turn = 2 // YOUR PARTNER
	if !s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.declarer = 2
	s.turn = 2
	if s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.turn = 0 // YOU 
	if !s.IsOpponentTurn() {
		t.Errorf("Error opponent")
	}
	s.turn = 1 // YOUR PARTNET
	if !s.IsOpponentTurn() {
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