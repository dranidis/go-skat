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

	if !s.getOpp1VoidSuit()[CARO] {
		t.Errorf("Error not follow void")
	}
	if s.getOpp2VoidSuit()[CARO] {
		t.Errorf("Error not follow void")
	}
	if s.getDeclarerVoidSuit()[CARO] {
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
	if !s.getOpp1VoidSuit()[s.trump] {
		t.Errorf("MH Player played A on a losing trick. Won by the declarer. It is his last card.")
	}
}

func TestInferenceSmearTrickTrumpWhenParnerWins1(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = SPADE
	s.leader = &d
	s.declarer = &d

	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = SPADE

	s.trick = []Card{
		Card{CARO, "J"}, // d
		Card{HEART, "J"}, // o1
	}
	players = []PlayerI{&d, &o1, &o2}

	card := Card{SPADE, "10"} // o2

	analysePlay(&s, s.opp2, card)
	if s.opp2VoidSuit[s.trump] {
		t.Errorf("MH Player played 10 on a trump won by partner. It is NOT his last card.")
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
	if s.getOpp1VoidSuit()[s.trump] {
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
	if !s.getOpp2VoidSuit()[s.trump] {
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
	if s.getOpp2VoidSuit()[s.trump] {
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

func TestInferencePlayerPlays_ValueCard_onTrickWonByDeclarer(t *testing.T) {
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

	players = []PlayerI{&d, &o1, &o2}

	s.trick = []Card{
		Card{HEART, "A"}, // d
		Card{HEART, "7"}, // o1
	}

	card := Card{HEART, "K"} // o2 does not have Q, 9, 8 , 7

	analysePlay(&s, s.opp2, card)
	if ! in(s.opp2VoidCards, Card{HEART, "D"}, Card{HEART, "9"},Card{HEART, "8"}, Card{HEART, "7"}) {
		t.Errorf("BH Player does not have cards lower than K.")
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
	if ! s.getDeclarerVoidSuit()[HEART] {
		t.Errorf("Declarer is void on suit.")
	}

}

func TestInference_Declarer_A10_at_Declarer(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.cardsPlayed = []Card{}
	s.leader = &o1

	s.declarer = &d
	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = HEART

	players = []PlayerI{&o1, &o2, &d}

	s.trick = []Card{
		Card{HEART, "K"}, 
		Card{HEART, "8"}, 
	}

	card := Card{HEART, "9"} // goes under, he has A but no 10

	analysePlay(&s, s.declarer, card)
	// if ! in(s.opp1VoidCards, Card{HEART, "A"}) {
	// 	t.Errorf("Opp1 does not have A")
	// }
	// if ! in(s.opp2VoidCards, Card{HEART, "A"}) {
	// 	t.Errorf("Opp2 does not have A")
	// }
	if ! in(s.declarerVoidCards, Card{HEART, "10"}) {
		t.Errorf("Declarer does not have 10 but A")
	}

}


// DECLARER CAN ONLY KNOW THAT
func TestInference_Opponent_Plays_K_Decl_Played_A_and_Has_10_or_In_Skat(t *testing.T) {
	d := makePlayer([]Card{
		Card{HEART, "10"},
		})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.cardsPlayed = []Card{}
	s.leader = &d

	s.declarer = &d
	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = HEART

	players = []PlayerI{&d, &o1, &o2}

	s.trick = []Card{
		Card{HEART, "A"}, 
	}

	card := Card{HEART, "K"} // opponent plays his lower card. He does not have any more HEART

	analysePlay(&s, s.opp1, card)
	// if ! in(s.opp1VoidCards, Card{HEART, "A"}) {
	// 	t.Errorf("Opp1 does not have A")
	// }
	// if ! in(s.opp2VoidCards, Card{HEART, "A"}) {
	// 	t.Errorf("Opp2 does not have A")
	// }
	if !d.getInference().opp1VoidSuitB[s.follow] {
		t.Errorf("MH Player played K on a losing trick (10 in declarer). It is his last card.")
	}
}

// DECLARER CAN ONLY KNOW THAT
func TestInference_Opponent_Plays_K_Decl_Played_A_and_Has_10_or_In_Skat_Opp1(t *testing.T) {
	d := makePlayer([]Card{
		Card{HEART, "10"},
		})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.leader = &o2

	s.declarer = &d
	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = HEART

	players = []PlayerI{&o2, &d, &o1}

	s.trick = []Card{
		Card{HEART, "7"}, 
		Card{HEART, "A"}, 
	}
	s.cardsPlayed = s.trick

	card := Card{HEART, "K"} // opponent plays his lower card. He does not have any more HEART

	analysePlay(&s, s.opp1, card)

	if !d.getInference().opp1VoidSuitB[s.follow] {
		t.Errorf("MH Player played K on a losing trick (10 in declarer). It is his last card.")
	}
}

// DECLARER CAN ONLY KNOW THAT
func TestInference_Opponent_Plays_D_Decl_Played_A_and_Has_10_or_In_Skat(t *testing.T) {
	d := makePlayer([]Card{
		Card{HEART, "10"},
		})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.cardsPlayed = []Card{}
	s.leader = &d

	s.declarer = &d
	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = HEART

	players = []PlayerI{&d, &o1, &o2}

	s.trick = []Card{
		Card{HEART, "A"}, 
	}

	card := Card{HEART, "D"} // opponent might still have HEART

	analysePlay(&s, s.opp1, card)
	if d.getInference().opp1VoidSuitB[s.follow] {
		t.Errorf("MH Player played D on a losing trick (10 in declarer). It is NOT his last card.")
	}
}

func TestInference_Opponent_Plays_D_Decl_Played_A_and_Has_10_or_In_Skat_2(t *testing.T) {
	d := makePlayer([]Card{
		Card{HEART, "10"},
		})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.cardsPlayed = []Card{}
	s.leader = &d

	s.declarer = &d
	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = HEART

	players = []PlayerI{&d, &o1, &o2}

	s.trick = []Card{
		Card{HEART, "7"}, 
	}

	card := Card{HEART, "D"} // opponent might still have HEART

	analysePlay(&s, s.opp1, card)
	if d.getInference().opp1VoidSuitB[s.follow] {
		t.Errorf("MH Player played D on a losing trick (10 in declarer). It is NOT his last card.")
	}
}

// DECLARER CAN ONLY KNOW THAT
func TestInference_Opponent_Plays_D_Decl_Played_A_and_Has_10_or_In_Skat_1(t *testing.T) {
	d := makePlayer([]Card{
		Card{HEART, "10"},
		})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = CLUBS
	s.cardsPlayed = []Card{}
	s.leader = &o1

	s.declarer = &d
	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = HEART

	players = []PlayerI{&o1, &o2, &d}

	s.trick = []Card{}

	card := Card{HEART, "A"} // opponent might still have HEART

	analysePlay(&s, s.opp1, card)
	if d.getInference().opp1VoidSuitB[s.follow] {
		t.Errorf("First card")
	}
}

func TestInference_Opponent_Plays_D_Decl_Played_A_and_Has_10_or_In_Skat_3(t *testing.T) {
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

	players = []PlayerI{&d, &o1, &o2}

	s.trick = []Card{Card{HEART, "9"}}
	s.cardsPlayed = s.trick

	card := Card{HEART, "10"} // opponent might still have HEART

	analysePlay(&s, s.opp1, card)
	if d.getInference().opp1VoidSuitB[s.follow] {
		t.Errorf("First card")
	}
}

func TestInference_Opponent_Plays_D_Decl_Played_A_and_Has_10_or_In_Skat_4(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = HEART
	s.leader = &o2

	s.declarer = &d
	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = HEART

	players = []PlayerI{&o2, &d, &o1}

	s.trick = []Card{
		Card{CARO, "J"},
		Card{HEART, "D"},
	}
	s.cardsPlayed = s.trick

	card := Card{SPADE, "D"} // opponent might still have HEART

	analysePlay(&s, s.opp1, card)
	if d.getInference().opp1VoidSuitB[HEART] {
		t.Errorf("Opp 2 wins: %v", s.trick)
	}
}

func TestInference_Opponent_Plays_D_Decl_Played_A_and_Has_10_or_In_Skat_5(t *testing.T) {
	d := makePlayer([]Card{})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})

	s := makeSuitState()
	s.trump = HEART
	s.leader = &d

	s.declarer = &d
	s.opp1 = &o1
	s.opp2 = &o2
	s.follow = HEART

	players = []PlayerI{&d, &o1, &o2}

	s.trick = []Card{
		Card{HEART, "A"},
	}
	s.cardsPlayed = s.trick

	card := Card{CARO, "10"} 

	analysePlay(&s, s.opp1, card)
	if d.getInference().opp1VoidSuitB[CARO] {
		t.Errorf("Opp 2 wins: %v", s.trick)
	}
}