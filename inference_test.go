package main

import (
	"testing"
)

func TestInferenceNotFollowingSuit(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.leader = &d
	s.declarer = &d

	s.opp1 = &o1
	s.opp2 = &o2

	s.trick = []Card{Card{CARO, "7"}}
	card := Card{HEART, "7"}

	analysePlay(&s, s.opp1, card)

	if !s.opp1VoidSuit[CARO] {
		t.Errorf("Error not follow void")
	}
	if s.opp2VoidSuit[CARO] {
		t.Errorf("Error not follow void")
	}
	if s.declarerVoidSuit[CARO] {
		t.Errorf("Error not follow void")
	}	
}

func TestInferenceLastTrumpDeclarerWins1(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.leader = &d
	s.declarer = &d

	s.opp1 = &o1
	s.opp2 = &o2

	s.trick = []Card{
		Card{CLUBS, "J"}, // d
	}

	card := Card{CLUBS, "A"} // o1

	analysePlay(&s, s.opp1, card)
	if !s.opp1VoidSuit[s.trump] {
		t.Errorf("MH Player played A on a losing trick. Won by the declarer. It is his last card.")
	}
}

func TestInferenceLastTrumpDSmearing1(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.leader = &d
	s.declarer = &d

	s.opp1 = &o1
	s.opp2 = &o2

	s.trick = []Card{
		Card{CARO, "J"}, // d
	}

	card := Card{CLUBS, "A"} // o1

	analysePlay(&s, s.opp1, card)
	if s.opp1VoidSuit[s.trump] {
		t.Errorf("MH Player played A on a losing trick, BUT there are still higher Trumps in game, expecting partner to take the trick.")
	}
}

func TestInferenceLastTrumpDeclarerWins2(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.leader = &d
	s.declarer = &d

	s.opp1 = &o1
	s.opp2 = &o2

	s.trick = []Card{
		Card{CARO, "J"}, // d
		Card{CLUBS, "7"}, // o1
	}

	card := Card{CLUBS, "A"} // o2

	analysePlay(&s, s.opp2, card)
	if !s.opp2VoidSuit[s.trump] {
		t.Errorf("BH Player played A on a losing trick. Won by the declarer. It is his last card.")
	}
}

func TestInferenceLastTrumpPartnerWins(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.leader = &d
	s.declarer = &d

	s.opp1 = &o1
	s.opp2 = &o2

	s.trick = []Card{
		Card{CLUBS, "7"}, // d
		Card{CLUBS, "J"}, // o1
	}

	card := Card{CLUBS, "A"} // o2, smearing

	analysePlay(&s, s.opp2, card)
	if s.opp2VoidSuit[s.trump] {
		t.Errorf("BH Player SMEARING A on a losing trick. Won by the partner.")
	}

}