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

	p1 := makeMinMaxPlayer(dist[0])
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players := []PlayerI{&p1, &p2, &p3}

	ss := makeSuitState()
	ss.trump = CARO
	ss.declarer = &p1
	ss.opp1 = &p2
	ss.opp2 = &p3

	ss.leader = &p1

	ss.follow = CARO
	ss.trick = []Card{Card{SPADE, "J"}, Card{SPADE, "8"}}

	skatState := SkatState{
		ss,
		players,
	}
	debugTacticsLog("Players: %v\n", players)
	debugTacticsLog("Oldstate: %v\n", skatState)
	debugTacticsLog("Oldstate players: %v\n", skatState.players)

	action := SkatAction{Card{CLUBS, "J"}}

	var skatStateP game.State
	skatStateP = &skatState
	var a game.Action
	a = action

	skatStateN := skatStateP.FindNextState(a)
	newState := skatStateN.(*SkatState)
	debugTacticsLog("Oldstate: %v\nNewState: %v\n", skatState, newState)

	if newState.declarer.getScore() != skatState.declarer.getScore()  {
		t.Errorf("Wrong declScore score. Is: %d", newState.declarer.getScore())
	}	
	if newState.opp1.getScore() != skatState.opp1.getScore() {
		t.Errorf("Wrong oppScore score. Is: %d", newState.opp1.getScore())
	}	
	if newState.opp2.getScore() != skatState.opp2.getScore() + 4 {
		t.Errorf("Wrong oppScore score. Is: %d", newState.opp2.getScore())
	}	
	if newState.players[0].getName() != p3.getName() {
		t.Errorf("Wrong winner. Is: %s", newState.players[0].getName())
	}	
	if newState.players[1].getName() != p1.getName() {
		t.Errorf("Wrong 2nd. Is: %s", newState.players[1].getName())
	}	
	if newState.players[2].getName() != p2.getName() {
		t.Errorf("Wrong 3rd. Is: %s", newState.players[1].getName())
	}	
	

}

// func TestFindNextStateEndOfTrick2(t *testing.T) {
// 	dist := make([][]Card, 3)
// 	dist[0] = []Card {
// 		// Card{SPADE, "J"},
// 		// Card{SPADE, "A"},
// 		Card{SPADE, "K"},
// 		Card{SPADE, "9"},
// 	}
// 	dist[1] = []Card {
// 		// Card{SPADE, "8"},
// 		// Card{CLUBS, "A"},
// 		Card{HEART, "8"},
// 		Card{CLUBS, "7"},
// 	}
// 	dist[2] = []Card {
// 		Card{SPADE, "10"},
// 		Card{SPADE, "7"},
// 		Card{CLUBS, "8"},
// 		Card{CLUBS, "9"},
// 	}


// 	p1 := makePlayer(dist[0])
// 	p2 := makePlayer(dist[1])
// 	p3 := makePlayer(dist[2])
// 	p1.name = "Decl"
// 	p2.name = "Opp1"
// 	p3.name = "Opp2"

// 	players = []PlayerI{&p1, &p2, &p3}

// 	ss := makeSuitState()
// 	ss.trump = CARO
// 	ss.declarer = &p1
// 	ss.opp1 = &p2
// 	ss.opp2 = &p3
// 	ss.leader = &p3

// 	skatState := SkatState{
// 		ss,
// 		CARO, // trump
// 		dist, 			
// 		[]Card{{SPADE, "A"}, Card{CLUBS, "A"}}, // trick 
// 		0, // declarer 
// 		2, // who's turn is it
// 		52, 
// 		42,
// 		false,
// 	}

// 	action := SkatAction{Card{SPADE, "7"}}

// 	var skatStateP game.State
// 	skatStateP = &skatState
// 	var a game.Action
// 	a = action

// 	skatStateN := skatStateP.FindNextState(a)
// 	newState := skatStateN.(*SkatState)
// 	debugTacticsLog("State: %v\n", newState)

// 	if newState.declScore != skatState.declScore + 22 {
// 		t.Errorf("Wrong declScore score. Is: %d", newState.declScore)
// 	}	
// 	if newState.oppScore != skatState.oppScore {
// 		t.Errorf("Wrong oppScore score. Is: %d", newState.oppScore)
// 	}	
// 	if newState.turn != 0 {
// 		t.Errorf("Wrong winner. Is: %d", newState.turn)
// 	}	
// }

// func TestFindNextStateEndOfTrick3(t *testing.T) {
// 	dist := make([][]Card, 3)
// 	dist[0] = []Card {
// 		// Card{SPADE, "J"},
// 		// Card{SPADE, "A"},
// 		// Card{SPADE, "K"},
// 		Card{SPADE, "9"},
// 	}
// 	dist[1] = []Card {
// 		// Card{SPADE, "8"},
// 		// Card{CLUBS, "A"},
// 		// Card{HEART, "8"},
// 		Card{CLUBS, "7"},
// 	}
// 	dist[2] = []Card {
// 		Card{SPADE, "10"},
// 		// Card{SPADE, "7"},
// 		// Card{CLUBS, "8"},
// 		Card{CLUBS, "9"},
// 	}

// 	p1 := makePlayer(dist[0])
// 	p2 := makePlayer(dist[1])
// 	p3 := makePlayer(dist[2])
// 	p1.name = "Decl"
// 	p2.name = "Opp1"
// 	p3.name = "Opp2"

// 	players = []PlayerI{&p1, &p2, &p3}

// 	ss := makeSuitState()
// 	ss.trump = CARO
// 	ss.declarer = &p1
// 	ss.opp1 = &p2
// 	ss.opp2 = &p3
// 	ss.leader = &p3

// 	skatState := SkatState{
// 		ss,
// 		CARO, // trump
// 		dist, 			
// 		[]Card{{SPADE, "K"}, Card{SPADE, "8"}}, // trick 
// 		0, // declarer 
// 		2, // who's turn is it
// 		52, 
// 		42,
// 		false,
// 	}

// 	action := SkatAction{Card{SPADE, "10"}}

// 	var skatStateP game.State
// 	skatStateP = &skatState
// 	var a game.Action
// 	a = action

// 	skatStateN := skatStateP.FindNextState(a)
// 	newState := skatStateN.(*SkatState)
// 	debugTacticsLog("State: %v\n", newState)

// 	if newState.declScore != skatState.declScore {
// 		t.Errorf("Wrong declScore score. Is: %d", newState.declScore)
// 	}	
// 	if newState.oppScore != skatState.oppScore + 14{
// 		t.Errorf("Wrong oppScore score. Is: %d", newState.oppScore)
// 	}	
// 	if newState.turn != 2 {
// 		t.Errorf("Wrong winner. Is: %d", newState.turn)
// 	}	
// }

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

	p1 := makeMinMaxPlayer(dist[0])
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players := []PlayerI{&p1, &p2, &p3}

	ss := makeSuitState()
	ss.trump = CARO
	ss.declarer = &p1
	ss.opp1 = &p2
	ss.opp2 = &p3
	ss.leader = &p1
	ss.trick = []Card{}

	skatState := SkatState{
		ss,
		players,
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
	// if newState.turn != 1 {
	// 	t.Errorf("Wrong turn. Is: %d", newState.turn)
	// }	
}

func TestFindNextState2ndFollowsk(t *testing.T) {
	dist := make([][]Card, 3)
	dist[0] = []Card {
		// Card{SPADE, "J"},
		// Card{SPADE, "A"},
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

	p1 := makeMinMaxPlayer(dist[0])
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players := []PlayerI{&p1, &p2, &p3}

	ss := makeSuitState()
	ss.trump = CARO
	ss.declarer = &p1
	ss.opp1 = &p2
	ss.opp2 = &p3
	ss.leader = &p1
	ss.trick = []Card{Card{SPADE, "K"}}

	skatState := SkatState{
		ss,
		players,
	}

	action := SkatAction{Card{HEART, "8"}}

	var skatStateP game.State
	skatStateP = &skatState
	var a game.Action
	a = action

	skatStateN := skatStateP.FindNextState(a)
	newState := skatStateN.(*SkatState)
	debugTacticsLog("State: %v\n", newState)

	if len(newState.trick) != 2 {
		t.Errorf("Wrong trick size. Is: %d %v", len(newState.trick), newState.trick)
	}	
	if !in(newState.trick, Card{SPADE, "K"}, Card{HEART, "8"}) {
		t.Errorf("Wrong trick. Is: %v", newState.trick)
	}	
	// if newState.turn != 1 {
	// 	t.Errorf("Wrong turn. Is: %d", newState.turn)
	// }	
}

func TestCopySkatState(t *testing.T) {
	dist := make([][]Card, 3)
	dist[0] = []Card{
		
		Card{CLUBS, "8"},
		Card{HEART, "A"},
	}
	dist[1] = []Card{
		
		Card{CARO, "8"},
		Card{SPADE, "A"},
	}
	dist[2] = []Card{
		Card{CLUBS, "9"},
		Card{CARO, "10"},
		Card{CARO, "K"},
	}

	p1 := makeMinMaxPlayer(dist[0])
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players := []PlayerI{&p1, &p2, &p3}

	sst := makeSuitState()
	sst.trump = CARO
	sst.declarer = &p1
	sst.opp1 = &p2
	sst.opp2 = &p3
	sst.leader = &p1
	sst.trick = []Card{Card{CLUBS, "J"},Card{HEART, "10"}}


	skatState := SkatState{
		sst,
		players,
	}

	ss := skatState.copySkatState()

	debugTacticsLog("%v\n", skatState)
	debugTacticsLog("%v\n", ss)
	if len(ss.players[0].getHand()) != 2 && len(ss.players[1].getHand()) != 2 && len(ss.players[2].getHand()) != 3 {
		t.Errorf("ERROR Copy: %v", ss)
	}
	for i := 0; i < 3; i++ {
		if !in(ss.players[i].getHand(), skatState.players[i].getHand()...) {
			t.Errorf("ERROR Copy: %v %v", ss.players[i].getHand(), skatState.players[i].getHand())
		}
	
	}
	if len(ss.trick) != 2 {
		t.Errorf("ERROR COPY trick: %v", ss.trick)
	}
	if !in(ss.trick, skatState.trick...) {
			t.Errorf("ERROR Copy: %v %v", ss.trick, skatState.trick)
		}
	
}

// func TestFindLegals(t *testing.T) {
// 	dist := make([][]Card, 3)
// 	dist[0] = []Card{
// 		Card{CLUBS, "J"},
// 		Card{CLUBS, "8"},
// 		Card{HEART, "A"},
// 	}
// 	dist[1] = []Card{
// 		Card{HEART, "10"},
// 		Card{CARO, "8"},
// 		Card{SPADE, "A"},
// 	}
// 	dist[2] = []Card{
// 		Card{CARO, "10"},
// 		Card{CARO, "K"},
// 	}

// 	p1 := makePlayer(dist[0])
// 	p2 := makePlayer(dist[1])
// 	p3 := makePlayer(dist[2])
// 	p1.name = "Decl"
// 	p2.name = "Opp1"
// 	p3.name = "Opp2"

// 	players = []PlayerI{&p1, &p2, &p3}

// 	sst := makeSuitState()
// 	sst.trump = CARO
// 	sst.declarer = &p1
// 	sst.opp1 = &p2
// 	sst.opp2 = &p3
// 	sst.leader = &p1

// 	skatState := SkatState{
// 		sst,
// 		CLUBS,
// 		dist, 			
// 		[]Card{
// 			Card{CLUBS, "9"},

// 		}, 
// 		0, // declarer
// 		0, // turn
// 		30, 
// 		45,
// 		false,
// 	}

// 	actions := skatState.FindLegals()
// 	cards := []Card{}
// 	for _, action := range actions {
// 		ma := action.(SkatAction)
// 		cards = append(cards, ma.card)
// 	}
// 	// fmt.Println(cards)
// 	if !in(cards, Card{CLUBS, "J"}, Card{CLUBS, "8"}) {
// 		t.Errorf("ERROR TestFindLegals: %v", actions)
// 	}
// 	if in(cards, Card{HEART, "A"}) {
// 		t.Errorf("ERROR TestFindLegals: %v", actions)
// 	}

// 	ssa := skatState.FindNextState(actions[0])

// 	ss := ssa.(*SkatState)
// 	// fmt.Println(skatState)
// 	// fmt.Println(ss)



// 	if in(ss.playerHand[0], Card{CLUBS, "J"}) {
// 		t.Errorf("ERROR TestFindLegals: %v", ss.playerHand[0])
// 	}

// 	if !ss.trick[1].equals(Card{CLUBS, "J"}) {
// 		t.Errorf("ERROR TestFindLegals: %v", ss.trick)
// 	}

// 	if ss.turn != 1 {
// 		t.Errorf("Error turn %d", ss.turn)
// 	}

// 	ssa1 := ss.FindNextState(SkatAction{Card{HEART, "10"}})
// 	// fmt.Println(ssa1)	

// 	ss1 := ssa1.(*SkatState)

// 	if ss1.declScore != 42 {
// 		t.Errorf("Declarer score %d", ss1.declScore)
// 	} 	
// 	if ss1.oppScore != skatState.oppScore {
// 		t.Errorf("Opponent score %d", ss1.declScore)
// 	} 
// }

// func TestMoveOne(t *testing.T) {
// 	p1Hand := []Card{
// 		Card{CARO, "J"},
// 		Card{SPADE, "A"},
// 		}

// 	p1 := makePlayer(p1Hand)
// 	p2 := makePlayer([]Card{
// 		Card{CARO, "10"},
// 		Card{SPADE, "10"},
// 		})
// 	p3 := makePlayer([]Card{
// 		Card{HEART, "D"},
// 		Card{CLUBS, "A"},
// 		})

// 	players := []PlayerI{&p1, &p2, &p3}

// 	s := makeSuitState()
// 	s.trump = CARO
// 	s.declarer = &p1
// 	s.leader = &p1
// 	s.opp1 = &p2
// 	s.opp2 = &p3

// 	s.trick = []Card{
// 	}

// 	card, newplayers := moveOne(&s, players)

// 	if len(s.trick) != 1 {
// 		t.Errorf("Error moveone")
// 	}

// 	if !in(s.trick, card) {
// 		t.Errorf("NOt in trick")
// 	}

// 	if newplayers[0] != players[0] {
// 		t.Errorf("Wrong player order")
// 	}

// 	if !in(p1Hand, card) {
// 		t.Errorf("Wrong card")
// 	}

// }


// func TestGetGameSuitStateAndPlayers(t *testing.T) {

// 	h1 := []Card{
// 		Card{CARO, "J"},
// 		Card{SPADE, "A"},
// 		}
// 	h2 := []Card{
// 		Card{CARO, "10"},
// 		Card{SPADE, "10"},
// 		}
// 	h3 := []Card{
// 		Card{HEART, "D"},
// 		Card{CLUBS, "A"},
// 		}

// 	dist := [][]Card{h1, h2, h3}

// 	p1 := makePlayer(dist[0])
// 	p2 := makePlayer(dist[1])
// 	p3 := makePlayer(dist[2])
// 	p1.name = "Decl"
// 	p2.name = "Opp1"
// 	p3.name = "Opp2"

// 	players = []PlayerI{&p1, &p2, &p3}

// 	sst := makeSuitState()
// 	sst.trump = CARO
// 	sst.declarer = &p1
// 	sst.opp1 = &p2
// 	sst.opp2 = &p3
// 	sst.leader = &p1


// 	skatState := SkatState{
// 		sst,
// 		CARO, // trump
// 		dist, 			
// 		[]Card{}, // trick 
// 		0, // declarer 
// 		0, // who's turn is it
// 		5, 
// 		6,
// 		true,
// 	}

// 	s, players := skatState.getGameSuitStateAndPlayers()

// 	debugTacticsLog("Suitstate %v\n", *s)
// 	debugTacticsLog("Player[0] %v\n", players[0])
// 	// if players[0].getScore() != 5 {
// 	// 	t.Errorf("Wrong declarer score: ", players[0].getScore())
// 	// }
// 	if s.trump != skatState.trump {
// 		t.Errorf("Wrong trump")
// 	}

// 	if s.declarer != players[0] && s.opp1 != players[1] && s.opp2 != players[2] {
// 		t.Errorf("Wrong turn order. Expecting: Decl:%s=%s, O1:%s=%s, O2:%s=%s ",
// 			s.declarer.getName(), players[0].getName(), 
// 			s.opp1.getName(), players[1].getName(), s.opp2.getName(), players[2].getName())
// 	}

// 	// new test case
// 	skatState.declarer = 2
// 	s, players = skatState.getGameSuitStateAndPlayers()

// 	if s.declarer != players[2] && s.opp1 != players[0] && s.opp2 != players[1] {
// 		t.Errorf("Wrong turn order, when declarer is 2")
// 	}

// 	// new test case
// 	skatState.declarer = 1
// 	s, players = skatState.getGameSuitStateAndPlayers()

// 	if s.declarer != players[1] && s.opp1 != players[2] && s.opp2 != players[0] {
// 		t.Errorf("Wrong turn order, when declarer is 1")
// 	}

// 	// new test case
// 	skatState.trick = []Card{Card{CARO, "K"}}
// 	// new test case
// 	skatState.declarer = 0
// 	s, players = skatState.getGameSuitStateAndPlayers()

// 	if s.declarer != players[1] && s.opp1 != players[2] && s.opp2 != players[0] {
// 		t.Errorf("Wrong turn order, when declarer is 0 and trick len 1")
// 	}

// 	// new test case
// 	skatState.turn = 2
// 	s, players = skatState.getGameSuitStateAndPlayers()

// 	if s.declarer != players[2] && s.opp1 != players[0] && s.opp2 != players[1] {
// 		t.Errorf("Wrong turn order, when declarer is 0 and trick len 1, and turn 2")
// 	}

// }

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

	p1 := makeMinMaxPlayer(dist[0])
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players := []PlayerI{&p1, &p2, &p3}

	sst := makeSuitState()
	sst.trump = CARO
	sst.declarer = &p1
	sst.opp1 = &p2
	sst.opp2 = &p3
	sst.leader = &p1
	sst.trick = []Card{}


	skatState := SkatState{
		sst,
		players,
	}

	// minimax.DEBUG = true

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

func TestAlphaBetaTactics3C(t *testing.T) {

	h1 := []Card{
		Card{CARO, "A"},
		Card{CARO, "J"},
		Card{SPADE, "A"},
		}
	h2 := []Card{
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{SPADE, "10"},
		}
	h3 := []Card{
		Card{HEART, "D"},
		Card{CLUBS, "A"},
		Card{CLUBS, "K"},
		}

	dist := [][]Card{h1, h2, h3}

	p1 := makeMinMaxPlayer(dist[0])
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players := []PlayerI{&p1, &p2, &p3}

	sst := makeSuitState()
	sst.trump = SPADE
	sst.declarer = &p1
	sst.opp1 = &p2
	sst.opp2 = &p3
	sst.leader = &p1
	sst.trick = []Card{}


	skatState := SkatState{
		sst,
		players,
	}

	// minimax.DEBUG = true
	minimax.MAXDEPTH = 3

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


func TestAlphaBetaTactics4C(t *testing.T) {
	h1 := []Card{
		Card{SPADE, "A"},
		Card{CLUBS, "J"},
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{HEART, "9"},
		}
	h2 := []Card{
		Card{HEART, "10"},
		Card{HEART, "7"},
		Card{CLUBS, "K"},
		Card{CLUBS, "9"},
		Card{SPADE, "J"},
		}
	h3 := []Card{
		Card{CARO, "J"},
		Card{SPADE, "8"},
		Card{HEART, "D"},
		Card{CLUBS, "A"},
		Card{CLUBS, "D"},
		}

	dist := [][]Card{h1, h2, h3}
	FirstPlays(t, "abt", dist, SPADE, -120, 120)
}

func TestAlphaBeta4C(t *testing.T) {
	h1 := []Card{
		Card{SPADE, "A"},
		Card{CLUBS, "J"},
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{HEART, "9"},
		}
	h2 := []Card{
		Card{HEART, "10"},
		Card{HEART, "7"},
		Card{CLUBS, "K"},
		Card{CLUBS, "9"},
		Card{SPADE, "J"},
		}
	h3 := []Card{
		Card{CARO, "J"},
		Card{SPADE, "8"},
		Card{HEART, "D"},
		Card{CLUBS, "A"},
		Card{CLUBS, "D"},
		}

	dist := [][]Card{h1, h2, h3}
	FirstPlays(t, "ab", dist, SPADE, -120, 120)
}

func TestAlphaBeta3C(t *testing.T) {
	h1 := []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
		Card{CLUBS, "D"},
		}
	h2 := []Card{
		Card{CLUBS, "K"},
		Card{SPADE, "A"},
		Card{CLUBS, "8"},
		}
	h3 := []Card{
		Card{CARO, "J"},
		Card{SPADE, "J"},
		Card{CARO, "K"},
		}

	dist := [][]Card{h1, h2, h3}
	FirstPlays(t, "ab", dist, HEART, -120, 120)
}

func TestAlphaBeta23C(t *testing.T) {
	h1 := []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
		}
	h2 := []Card{
		Card{CLUBS, "K"},
		Card{SPADE, "A"},
		Card{CLUBS, "8"},
		}
	h3 := []Card{
		Card{CARO, "J"},
		Card{SPADE, "J"},
		Card{CARO, "K"},
		}

	dist := [][]Card{h1, h2, h3}
	SecondPlays(t, "ab", dist, Card{CLUBS, "D"}, -120, 120)
}
		

func SecondPlays(t *testing.T, alg string, dist [][]Card, first Card, a, b float64) {

	p1 := makeMinMaxPlayer(dist[0])
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players = []PlayerI{&p1, &p2, &p3}
	playersP := []PlayerI{&p1, &p2, &p3}

	sst := makeSuitState()
	sst.trump = HEART
	sst.declarer = &p1
	sst.opp1 = &p2
	sst.opp2 = &p3
	sst.leader = &p1
	sst.trick = []Card{first}
	sst.follow = getSuit(sst.trump, first)

	allcards := append([]Card{}, p1.hand...)
	allcards = append(allcards, p2.hand...)
	allcards = append(allcards, p3.hand...)
	allcards = append(allcards, first)
	sst.trumpsInGame = filter(allcards, func(c Card) bool {
			return getSuit(sst.trump, c) == sst.trump
			})

	sst.skat = []Card{
		Card{CARO, "10"},
		Card{CARO, "D"},
	}
	played := makeDeck()
	played = remove(played, p1.hand...)
	played = remove(played, p2.hand...)
	played = remove(played, p3.hand...)
	played = remove(played, sst.skat...)
	played = remove(played, first)
	sst.cardsPlayed = played

	skatState := SkatState{
		sst,
		playersP,
	}

	copy1 := make([]Card, len(p2.hand))
	copy(copy1, p2.hand)
	for _, c := range copy1 {
		card := c
		p2.hand = remove(p2.hand, card)

		nextState := skatState.FindNextState(SkatAction{card})

		// minimax.DEBUG = true
		// debugTacticsInMM = true

		minimax.MAXDEPTH = 9999


		var skatStateP game.State
		// skatStateP = &skatState
		skatStateP = nextState
		// var a game.Action
		var as []game.Action
		var v float64

		startWhole := time.Now()



		if alg == "abt" {
			_, v, as = minimax.AlphaBetaTacticsActions(skatStateP, a, b)
		} else {
			_, v, as = minimax.AlphaBetaActions(skatStateP, a, b)
		}
		ti := time.Now()
		elapsed := ti.Sub(startWhole)		

		debugTacticsLog("Action: %v, Value: %.4f, Actions: %v time: %v\n", c, v, as, elapsed)
		p2.hand = append(p2.hand, card)
	}


	if false {
		t.Errorf("TEST")
	}
}

func FirstPlays(t *testing.T, alg string, dist [][]Card, tr string, a, b float64) {


	p1 := makeMinMaxPlayer(dist[0])
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players = []PlayerI{&p1, &p2, &p3}
	playersP := []PlayerI{&p1, &p2, &p3}

	sst := makeSuitState()
	sst.trump = tr
	sst.declarer = &p1
	sst.opp1 = &p2
	sst.opp2 = &p3
	sst.leader = &p1
	sst.trick = []Card{}

	allcards := append([]Card{}, p1.hand...)
	allcards = append(allcards, p2.hand...)
	allcards = append(allcards, p3.hand...)
	sst.trumpsInGame = filter(allcards, func(c Card) bool {
			return getSuit(sst.trump, c) == sst.trump
			})

	sst.skat = []Card{
		Card{CARO, "10"},
		Card{CARO, "D"},
	}
	played := makeDeck()
	played = remove(played, p1.hand...)
	played = remove(played, p2.hand...)
	played = remove(played, p3.hand...)
	played = remove(played, sst.skat...)
	sst.cardsPlayed = played

	skatState := SkatState{
		sst,
		playersP,
	}

	copy1 := make([]Card, len(p1.hand))
	copy(copy1, p1.hand)
	for _, c := range copy1 {
		card := c
		p1.hand = remove(p1.hand, card)

		nextState := skatState.FindNextState(SkatAction{card})

		// minimax.DEBUG = true
		// debugTacticsInMM = true

		minimax.MAXDEPTH = 9999


		var skatStateP game.State
		// skatStateP = &skatState
		skatStateP = nextState
		// var a game.Action
		var as []game.Action
		var v float64

		startWhole := time.Now()



		if alg == "abt" {
			_, v, as = minimax.AlphaBetaTacticsActions(skatStateP, a, b)
		} else {
			_, v, as = minimax.AlphaBetaActions(skatStateP, a, b)
		}
		ti := time.Now()
		elapsed := ti.Sub(startWhole)		

		debugTacticsLog("Action: %v, Value: %.4f, Actions: %v time: %v\n", c, v, as, elapsed)
		p1.hand = append(p1.hand, card)
	}


	if false {
		t.Errorf("TEST")
	}
}



func TestAlphaBetaTactics7C_1(t *testing.T) {

	h1 := []Card{
		// Card{CARO, "J"},
		Card{HEART, "10"},
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{HEART, "7"},
		Card{CLUBS, "K"},
		Card{CLUBS, "8"},
		}
	h2 := []Card{
		Card{CLUBS, "A"},
		Card{CLUBS, "D"},
		Card{SPADE, "8"},
		Card{SPADE, "10"},
		Card{HEART, "K"},
		Card{HEART, "A"},
		Card{CARO, "9"},

		}
	h3 := []Card{
		Card{SPADE, "A"},
		Card{SPADE, "K"},
		Card{SPADE, "9"},
		Card{CLUBS, "9"},
		Card{CARO, "A"},
		Card{CARO, "K"},
		Card{CARO, "D"},
		}

	dist := [][]Card{h1, h2, h3}

	p1 := makeMinMaxPlayer(dist[0])
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players = []PlayerI{&p1, &p2, &p3}
	playersP := []PlayerI{&p1, &p2, &p3}

	sst := makeSuitState()
	sst.trump = GRAND
	sst.declarer = &p1
	sst.opp1 = &p2
	sst.opp2 = &p3
	sst.leader = &p1
	sst.trick = []Card{Card{CARO, "J"}}
	sst.trumpsInGame = []Card{}
	sst.cardsPlayed = []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
	}

	skatState := SkatState{
		sst,
		playersP,
	}

	// minimax.DEBUG = true
	minimax.MAXDEPTH = 3

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

func TestAlphaBetaTactics7C_2(t *testing.T) {

	h1 := []Card{
		Card{CLUBS, "J"},

		Card{CARO, "K"},
		Card{CARO, "9"},
		Card{CARO, "8"},
		Card{CLUBS, "8"},
		Card{HEART, "10"},
		Card{HEART, "D"},
		Card{HEART, "7"},
		}
	h2 := []Card{
		Card{CARO, "J"},

		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "K"},
		Card{SPADE, "9"},
		Card{SPADE, "7"},
		Card{HEART, "8"},
		}
	h3 := []Card{
		Card{CARO, "7"},

		Card{CARO, "10"},
		Card{CARO, "D"},
		Card{CLUBS, "10"},
		Card{CLUBS, "9"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{SPADE, "D"},
		}

	dist := [][]Card{h1, h2, h3}

	p1 := makeMinMaxPlayer(dist[0])
	p1.score = 37
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players = []PlayerI{&p1, &p2, &p3}
	playersP := []PlayerI{&p1, &p2, &p3}

	sst := makeSuitState()
	sst.trump = CARO
	sst.declarer = &p1
	sst.opp1 = &p2
	sst.opp2 = &p3
	sst.leader = &p1
	sst.skat = []Card{Card{CLUBS, "D"}, Card{CLUBS, "7"}}
	sst.trick = []Card{
	}

	trick1 := []Card{
		Card{SPADE, "A"},
		Card{SPADE, "8"},
		Card{CARO, "A"},
	}
	trick2 := []Card{
		Card{HEART, "A"},
		Card{HEART, "9"},
		Card{HEART, "K"},
	}
	// trick3 := []Card{
	// 	Card{CLUBS, "J"},
	// 	Card{CARO, "J"},
	// 	Card{CARO, "7"},
	// }
	sst.cardsPlayed = append(sst.cardsPlayed, trick1...)
	sst.cardsPlayed = append(sst.cardsPlayed, trick2...)
	// sst.cardsPlayed = append(sst.cardsPlayed, trick3...)
	debugTacticsLog("Played: %v\n", sst.cardsPlayed)

	sst.trumpsInGame = makeTrumpDeck(sst.trump)
	sst.trumpsInGame = remove(sst.trumpsInGame, trick1...)
	sst.trumpsInGame = remove(sst.trumpsInGame, trick2...)
	// sst.trumpsInGame = remove(sst.trumpsInGame, trick3...)
	debugTacticsLog("trumpsInGame: %v\n", sst.trumpsInGame)

	sst.declarerVoidSuit[SPADE] = true
	sst.opp2VoidCards = []Card{Card{HEART, "D"}, Card{HEART, "8"}, Card{HEART, "7"}}

	skatState := SkatState{
		sst,
		playersP,
	}

	for _, play := range p1.hand {
		debugTacticsLog("Card %v\n", play)
		var skatStateP game.State
		skatStateP = skatState.FindNextState(SkatAction{play})

		// card := p1.playerTactic(&sst, h1)
		// debugTacticsLog("Playing card %v\n", card)
		// minimax.DEBUG = true
		minimax.MAXDEPTH = 6

		// skatStateP = &skatState
		var a game.Action
		var v float64
		var as []game.Action

		if true {
			a, v, as = minimax.AlphaBetaActions(skatStateP, 45, 61)
		}

		debugTacticsLog("Action: %v, Value: %.4f, Actions: %v \n", a, v, as)	
		break
	}

	if false {
		t.Errorf("TEST")
	}
}

func TestAlphaBetaTactics8C_play(t *testing.T) {
	MINIMAX_ALG = "ab"
	h1 := []Card{
		Card{CLUBS, "J"},

		Card{CARO, "K"},
		Card{CARO, "9"},
		Card{CARO, "8"},
		Card{CLUBS, "8"},
		Card{HEART, "10"},
		Card{HEART, "D"},
		Card{HEART, "7"},
		}
	h2 := []Card{
		Card{CARO, "J"},

		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "K"},
		Card{SPADE, "9"},
		Card{SPADE, "7"},
		Card{HEART, "8"},
		}
	h3 := []Card{
		Card{CARO, "7"},

		Card{CARO, "10"},
		Card{CARO, "D"},
		Card{CLUBS, "10"},
		Card{CLUBS, "9"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{SPADE, "D"},
		}

	dist := [][]Card{h1, h2, h3}

	p1 := makeMinMaxPlayer(dist[0])
	p1.score = 37
	p1.maxWorlds = 20
	p1.timeOutMs = 30000
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players = []PlayerI{&p1, &p2, &p3}
	// playersP := []PlayerI{&p1, &p2, &p3}

	sst := makeSuitState()
	sst.trump = CARO
	sst.declarer = &p1
	sst.opp1 = &p2
	sst.opp2 = &p3
	sst.leader = &p1
	sst.skat = []Card{Card{CLUBS, "D"}, Card{CLUBS, "7"}}
	sst.trick = []Card{
	}

	trick1 := []Card{
		Card{SPADE, "A"},
		Card{SPADE, "8"},
		Card{CARO, "A"},
	}
	trick2 := []Card{
		Card{HEART, "A"},
		Card{HEART, "9"},
		Card{HEART, "K"},
	}
	// trick3 := []Card{
	// 	Card{CLUBS, "J"},
	// 	Card{CARO, "J"},
	// 	Card{CARO, "7"},
	// }
	sst.cardsPlayed = append(sst.cardsPlayed, trick1...)
	sst.cardsPlayed = append(sst.cardsPlayed, trick2...)
	// sst.cardsPlayed = append(sst.cardsPlayed, trick3...)
	debugTacticsLog("Played: %v\n", sst.cardsPlayed)

	sst.trumpsInGame = makeTrumpDeck(sst.trump)
	sst.trumpsInGame = remove(sst.trumpsInGame, trick1...)
	sst.trumpsInGame = remove(sst.trumpsInGame, trick2...)
	// sst.trumpsInGame = remove(sst.trumpsInGame, trick3...)
	debugTacticsLog("trumpsInGame: %v\n", sst.trumpsInGame)

	sst.declarerVoidSuit[SPADE] = true
	sst.opp2VoidCards = []Card{Card{HEART, "D"}, Card{HEART, "8"}, Card{HEART, "7"}}

	// skatState := SkatState{
	// 	sst,
	// 	playersP,
	// }

	// debugTacticsLog("Card %v\n", h1[0])
	// var skatStateP game.State
	// skatStateP = skatState.FindNextState(SkatAction{h1[0]})

	//p1.depth = 9

	card := p1.playerTactic(&sst, p1.hand)
	debugTacticsLog("Playing card %v\n", card)
	// minimax.DEBUG = true
	// minimax.MAXDEPTH = 3

	// skatStateP = &skatState
	// var a game.Action
	// var v float64

	// if true {
	// 	a, v = minimax.AlphaBeta(skatStateP)
	// }

	// debugTacticsLog("Action: %v, Value: %.4f\n", a, v)
	// if false {
	// 	t.Errorf("TEST")
	// }
}

func TestAlphaBetaTactics7C_play(t *testing.T) {
	MINIMAX_ALG = "ab"
	h1 := []Card{
		// Card{CLUBS, "J"},

		Card{CARO, "K"},
		Card{CARO, "9"},
		Card{CARO, "8"},
		Card{CLUBS, "8"},
		Card{HEART, "10"},
		Card{HEART, "D"},
		Card{HEART, "7"},
		}
	h2 := []Card{
		// Card{CARO, "J"},

		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CLUBS, "A"},
		Card{CLUBS, "K"},
		Card{SPADE, "9"},
		Card{SPADE, "7"},
		Card{HEART, "8"},
		}
	h3 := []Card{
		// Card{CARO, "7"},

		Card{CARO, "10"},
		Card{CARO, "D"},
		Card{CLUBS, "10"},
		Card{CLUBS, "9"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{SPADE, "D"},
		}

	dist := [][]Card{h1, h2, h3}

	p1 := makeMinMaxPlayer(dist[0])
	p1.score = 41
	p1.maxWorlds = 20
	p1.timeOutMs = 30000
	p2 := makeMinMaxPlayer(dist[1])
	p3 := makeMinMaxPlayer(dist[2])
	p1.name = "Decl"
	p2.name = "Opp1"
	p3.name = "Opp2"

	players = []PlayerI{&p1, &p2, &p3}
	// playersP := []PlayerI{&p1, &p2, &p3}

	sst := makeSuitState()
	sst.trump = CARO
	sst.declarer = &p1
	sst.opp1 = &p2
	sst.opp2 = &p3
	sst.leader = &p1
	sst.skat = []Card{Card{CLUBS, "D"}, Card{CLUBS, "7"}}
	sst.trick = []Card{
	}

	trick1 := []Card{
		Card{SPADE, "A"},
		Card{SPADE, "8"},
		Card{CARO, "A"},
	}
	trick2 := []Card{
		Card{HEART, "A"},
		Card{HEART, "9"},
		Card{HEART, "K"},
	}
	trick3 := []Card{
		Card{CLUBS, "J"},
		Card{CARO, "J"},
		Card{CARO, "7"},
	}
	sst.cardsPlayed = append(sst.cardsPlayed, trick1...)
	sst.cardsPlayed = append(sst.cardsPlayed, trick2...)
	sst.cardsPlayed = append(sst.cardsPlayed, trick3...)
	debugTacticsLog("Played: %v\n", sst.cardsPlayed)

	sst.trumpsInGame = makeTrumpDeck(sst.trump)
	sst.trumpsInGame = remove(sst.trumpsInGame, trick1...)
	sst.trumpsInGame = remove(sst.trumpsInGame, trick2...)
	sst.trumpsInGame = remove(sst.trumpsInGame, trick3...)
	debugTacticsLog("trumpsInGame: %v\n", sst.trumpsInGame)

	sst.declarerVoidSuit[SPADE] = true
	sst.opp2VoidCards = []Card{Card{HEART, "D"}, Card{HEART, "8"}, Card{HEART, "7"}}

	// skatState := SkatState{
	// 	sst,
	// 	playersP,
	// }

	// debugTacticsLog("Card %v\n", h1[0])
	// var skatStateP game.State
	// skatStateP = skatState.FindNextState(SkatAction{h1[0]})

	//p1.depth = 6
	card := p1.playerTactic(&sst, p1.hand)
	debugTacticsLog("Playing card %v\n", card)
	// minimax.DEBUG = true
	// minimax.MAXDEPTH = 3

	// skatStateP = &skatState
	// var a game.Action
	// var v float64

	// if true {
	// 	a, v = minimax.AlphaBeta(skatStateP)
	// }

	// debugTacticsLog("Action: %v, Value: %.4f\n", a, v)
	// if false {
	// 	t.Errorf("TEST")
	// }
}

// func TestAlphaBetaTactics2CDef(t *testing.T) {

// 	h1 := []Card{
// 		Card{CARO, "J"},
// 		Card{SPADE, "A"},
// 		}
// 	h2 := []Card{
// 		Card{CARO, "10"},
// 		Card{SPADE, "10"},
// 		}
// 	h3 := []Card{
// 		Card{HEART, "D"},
// 		Card{CLUBS, "A"},
// 		}

// 	dist := [][]Card{h1, h2, h3}

// 	p1 := makePlayer(dist[0])
// 	p2 := makePlayer(dist[1])
// 	p3 := makePlayer(dist[2])
// 	p1.name = "Decl"
// 	p2.name = "Opp1"
// 	p3.name = "Opp2"

// 	players = []PlayerI{&p1, &p2, &p3}

// 	sst := makeSuitState()
// 	sst.trump = CARO
// 	sst.declarer = &p2
// 	sst.opp1 = &p3
// 	sst.opp2 = &p1
// 	sst.leader = &p1

// 	skatState := SkatState{
// 		sst,
// 		CARO, // trump
// 		dist, 			
// 		[]Card{}, // trick 
// 		1, // declarer 
// 		0, // who's turn is it
// 		0, 
// 		0,
// 		true,
// 	}

// 	minimax.DEBUG = true

// 	var skatStateP game.State
// 	skatStateP = &skatState
// 	var a game.Action
// 	var v float64

// 	if (false) {
//  		a, v = minimax.AlphaBetaTactics(skatStateP)
// 	}
// 	debugTacticsLog("Action: %v, Value: %.4f\n", a, v)
// 	if false {
// 		t.Errorf("TEST")
// 	}
// }

// func TestAlphaBetaTactics10C(t *testing.T) {
// 	dist := make([][]Card, 3)
// 	dist[0] = []Card {
// 		Card{SPADE, "J"},
// 		Card{HEART, "J"},

// 		Card{CARO, "K"},
// 		Card{CARO, "10"},
// 		Card{CARO, "9"},
// 		Card{CARO, "8"},

// 		Card{SPADE, "A"},
// 		Card{SPADE, "K"},
// 		Card{SPADE, "9"},

// 		Card{CLUBS, "K"},
// 	}
// 	// dist[0] = Shuffle(dist[0])
// 	dist[1] = []Card {
// 		Card{CARO, "J"},

// 		Card{CARO, "A"},
// 		Card{CARO, "7"},

// 		Card{CLUBS, "A"},
// 		Card{CLUBS, "8"},
// 		Card{CLUBS, "D"},

// 		Card{SPADE, "8"},
// 		Card{SPADE, "D"},

// 		Card{HEART, "8"},
// 		Card{HEART, "10"},
// 	}
// 	// dist[1] = Shuffle(dist[1])
// 	dist[2] = []Card {
// 		Card{CLUBS, "J"},

// 		Card{CARO, "D"},

// 		Card{SPADE, "10"},
// 		Card{SPADE, "7"},

// 		Card{HEART, "K"},
// 		Card{HEART, "A"},

// 		Card{CLUBS, "D"},
// 		Card{CLUBS, "8"},
// 		Card{CLUBS, "9"},
// 		Card{CLUBS, "7"},
// 	}
// 	// dist[2] = Shuffle(dist[2])

// 	p1 := makePlayer(dist[0])
// 	p2 := makePlayer(dist[1])
// 	p3 := makePlayer(dist[2])
// 	p1.name = "Decl"
// 	p2.name = "Opp1"
// 	p3.name = "Opp2"

// 	players = []PlayerI{&p1, &p2, &p3}

// 	sst := makeSuitState()
// 	sst.trump = CARO
// 	sst.declarer = &p1
// 	sst.opp1 = &p2
// 	sst.opp2 = &p3
// 	sst.leader = &p1


// 	skatState := SkatState{
// 		sst,
// 		CARO, // trump
// 		dist, 			
// 		[]Card{}, // trick 
// 		0, // declarer 
// 		0, // who's turn is it
// 		10, // score because of skat  
// 		0,
// 		true,
// 	}

// 	var skatStateP game.State
// 	skatStateP = &skatState
// 	var a game.Action
// 	var v float64

// 	minimax.MAXDEPTH = 9


// 	startWhole := time.Now()

// 	if (false) {
//  		a, v = minimax.AlphaBetaTactics(skatStateP)
// 	}
	
// 	ti := time.Now()
// 	elapsed := ti.Sub(startWhole)		
// 	debugTacticsLog("TOTAL %v\n", elapsed)

// 	debugTacticsLog("Action: %v, Value: %.4f\n", a, v)
// 	if false {
// 		t.Errorf("TEST")
// 	}
// }
