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
	s.follow = CLUBS

	s.trick = []Card{
		Card{CLUBS, "J"}, // d
	}
	players = []PlayerI{&d, &o1, &o2}

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
	s.follow = CLUBS

	s.trick = []Card{
		Card{CARO, "J"}, // d
	}
	players = []PlayerI{&d, &o1, &o2}

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
	s.follow = CLUBS

	s.trick = []Card{
		Card{CARO, "J"}, // d
		Card{CLUBS, "7"}, // o1
	}

	players = []PlayerI{&d, &o1, &o2}

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

func TestInferencePlayerDoesNOtPlay_10_onTrick_A_OpenedByPartner(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.leader = &d
	s.declarer = &d

	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = HEART

	players = []PlayerI{&o2, &d, &o1}

	s.trick = []Card{
		Card{HEART, "A"}, // o2
		Card{HEART, "7"}, // d
	}

	card := Card{HEART, "K"} // o1 does not have 10

	analysePlay(&s, s.opp1, card)
	if ! in(s.opp1VoidCards, Card{HEART, "10"}) {
		t.Errorf("BH Player does not have the 10.")
	}

}

func TestInferenceDeclarerPlays_10_onTrick_A_OpenedByOpponents(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.leader = &d
	s.declarer = &d

	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = HEART

	players = []PlayerI{&o2, &d, &o1}

	s.trick = []Card{
		Card{HEART, "A"}, // o2
	}

	card := Card{HEART, "10"} // d does not have any other HEART

	analysePlay(&s, s.declarer, card)
	if ! s.declarerVoidSuit[HEART] {
		t.Errorf("Declarer is void on suit.")
	}

}