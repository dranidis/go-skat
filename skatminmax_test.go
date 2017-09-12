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