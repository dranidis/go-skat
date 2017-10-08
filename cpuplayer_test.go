package main

import (
	"testing"
)

func TestResetPlayer(t *testing.T) {
	p := makePlayer([]Card{})
	p.handGame = true
	p.trumpToDeclare = CLUBS

	p.ResetPlayer()

	if p.handGame {
		t.Errorf("Error in reset")
	}
	if p.trumpToDeclare == CLUBS {
		t.Errorf("Error in reset")
	}
}

func TestPickUpSkat0(t *testing.T) {
	player2 := makePlayer([]Card{})
	player3 := makePlayer([]Card{})
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

	players = []PlayerI{&player, &player2, &player3}

	skat := player.cardsToDiscard(CARO)

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
		Card{CARO, "10"},
		Card{HEART, "K"},
	})

	skat := player.cardsToDiscard(SPADE)
	card1 := Card{HEART, "K"}
	card2 := Card{HEART, "9"}

	if !in(skat, card1, card2) {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
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

	p2 := makePlayer([]Card{})
	p3 := makePlayer([]Card{})
	players = []PlayerI{&player, &p2, &p3}

	player.pickUpSkat(skat)
	card1 := Card{HEART, "D"}
	card2 := Card{HEART, "7"}

	if player.trumpToDeclare != SPADE {
		t.Errorf("Expected SPADE declaration, found: %s", player.trumpToDeclare)
	}

	if !in(skat, card1, card2) {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
	}

	if len(player.hand) != 10 {
		t.Errorf("Wrong hand size after skat change: %d", len(player.hand))
	}
}

func TestPickUpSkatAndDeclareNew(t *testing.T) {
	player := makePlayer([]Card{
		Card{CARO, "J"},
		Card{HEART, "J"},

		Card{CLUBS, "9"},
		Card{CLUBS, "7"},

		Card{SPADE, "A"},
		Card{SPADE, "10"},

		Card{HEART, "A"},
		Card{HEART, "7"},

		Card{CARO, "A"},
		Card{CARO, "K"},
	})

	skat := []Card{
		Card{SPADE, "7"},
		Card{CLUBS, "8"},
	}

	player.pickUpSkat(skat)

	if player.trumpToDeclare != SPADE {
		t.Errorf("Expected SPADE, found: ", player.trumpToDeclare)
	}
}

func TestPickUpSkatAndDeclare_10_D_9(t *testing.T) {
	player := makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},

		Card{SPADE, "10"},
		Card{SPADE, "D"},
		Card{SPADE, "9"},

		Card{HEART, "10"},
		Card{HEART, "D"},
		Card{HEART, "7"},

		Card{CLUBS, "7"},

		Card{CARO, "8"},
		Card{CARO, "9"},
	})

	skat := []Card{
		Card{CARO, "8"},
		Card{CARO, "9"},
	}

	player.pickUpSkat(skat)

	if player.trumpToDeclare != HEART {
		t.Errorf("Expected HEART, since SPADE is a bit stronger, found: ", player.trumpToDeclare)
	}
}

func TestPickUpSkatGrandWith4Aces(t *testing.T) {
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

		Card{CLUBS, "A"},
		Card{CLUBS, "J"},
	})

	p2 := makePlayer([]Card{})
	p3 := makePlayer([]Card{})
	players = []PlayerI{&player, &p2, &p3}

	skat := player.cardsToDiscard(GRAND)

	card1 := Card{SPADE, "D"}
	card2 := Card{CARO, "D"}
	if !in(skat, card1, card2) {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
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
	card1 := Card{HEART, "K"}
	card2 := Card{HEART, "7"}
	if !in(skat, card1, card2) {
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

func TestPickUpSkatGRandWITH4jS(t *testing.T) {
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

	p2 := makePlayer([]Card{})
	p3 := makePlayer([]Card{})
	players = []PlayerI{&player, &p2, &p3}

	player.trumpToDeclare = GRAND
	skat := player.cardsToDiscard(player.trumpToDeclare)
	card1 := Card{CARO, "10"}
	card2 := Card{CARO, "K"}
	if !in(skat, card1, card2) {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
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

		Card{HEART, "A"},
		Card{SPADE, "A"},
	})

	// skat := []Card{
	// }

	p2 := makePlayer([]Card{})
	p3 := makePlayer([]Card{})
	players = []PlayerI{&player, &p2, &p3}

	player.trumpToDeclare = CARO
	// fmt.Println("TestPickUpSkat7")
	// fmt.Println(player.hand)
	// fmt.Println(skat)
	skat := player.cardsToDiscard(player.trumpToDeclare)
	// fmt.Println(sortSuit("", player.hand))
	// fmt.Println(skat)
	cc1 := skat[1].Suit != SPADE || skat[1].Rank != "A"
	cc2 := skat[0].Suit != HEART || skat[0].Rank != "A"
	if cc1 || cc2 {
		t.Errorf("Found in skat: %v %v", skat[0], skat[1])
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
		Card{HEART, "D"}, // loser
		Card{HEART, "9"}, // loser
		Card{SPADE, "8"}, // loser
	})
	player.risky = true

	p2 := makePlayer([]Card{})
	p3 := makePlayer([]Card{})
	players = []PlayerI{&player, &p2, &p3}

	player.calculateHighestBid(false)

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

func TestOpponentTacticFOREFollowPreviousSuit3(t *testing.T) {

	// if you have a card with suit played in a previous trick
	// started from you or your partner continue with it.

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate
	s.trump = CLUBS
	s.trick = []Card{}
	s.trumpsInGame = makeTrumpDeck(s.trump)

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
	s.trumpsInGame = makeTrumpDeck(s.trump)
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

func TestOpponentTacticMID_PartnerLeads_Figure_or_Number(t *testing.T) {

	// When the partner leads a figure it means that he does not have the A
	// play a number
	// When the partner leads a number it means that he does have the A
	// so play the 10
	// unless you have 3 cards: 7 D 10
	// in that case play the D

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.trump = CLUBS
	s.follow = CARO

	validCards := []Card{Card{CARO, "10"}, Card{CARO, "8"}}

	s.trick = []Card{Card{CARO, "D"}}

	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

	s.cardsPlayed = s.trick
	s.trumpsInGame = makeTrumpDeck(s.trump)

	s.leader = &teamMate
	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by opponent, valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}

	s.trick = []Card{Card{CARO, "7"}}
	s.cardsPlayed = s.trick

	card = player.playerTactic(&s, validCards)

	exp = Card{CARO, "10"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by opponent, and valid cards: %v expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}

	validCards = []Card{Card{CARO, "10"}, Card{CARO, "D"}, Card{CARO, "8"}}

	s.trick = []Card{Card{CARO, "K"}}
	s.cardsPlayed = s.trick

	card = player.playerTactic(&s, validCards)

	exp = Card{CARO, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by opponent, and valid cards: %v expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}

	s.trick = []Card{Card{CARO, "7"}}
	s.cardsPlayed = s.trick

	card = player.playerTactic(&s, validCards)

	exp = Card{CARO, "D"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by opponent, and valid cards: %v expected to play %v, played %v",
			s.trick, validCards, exp, card)
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
	s.opp1 = &player
	s.opp2 = &teamMate

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
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate
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

func TestOpponentTacticMIDTrumpJ(t *testing.T) {
	// MIDDLEHAND

	// if declarer leads a high trump, and there are still higher trumps
	// DON"T smear it! High Risk!

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

	s.trump = CLUBS
	s.trick = []Card{Card{HEART, "J"}}
	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
	}

	validCards := []Card{
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "10"},
		Card{CLUBS, "A"},
	}
	//

	card := player.playerTactic(&s, validCards)
	exp := Card{CLUBS, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, with trumps in game: %v and valid %v, expected to play (not smear) %v, played %v",
			s.trick, s.trumpsInGame, validCards, exp, card)
	}
}

func TestOpponentTacticMIDTrumpJ_Opp_Cannot_Go_higher(t *testing.T) {
	// MIDDLEHAND

	// if declarer leads a trump, and partner cannot go higher
	// DON"T smear it!

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

	s.trump = CLUBS
	s.trick = []Card{Card{CLUBS, "A"}}
	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
	}

	validCards := []Card{
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "K"},
		Card{CLUBS, "D"},
	}
	//

	s.opp2VoidCards = []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
	}

	card := player.playerTactic(&s, validCards)
	exp := Card{CLUBS, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, with trumps in game: %v and valid %v, expected to play (not smear, opp void higher) %v, played %v",
			s.trick, s.trumpsInGame, validCards, exp, card)
	}
}

func TestOpponentTacticMIDTrump_NonJ_SMEAR(t *testing.T) {
	// MIDDLEHAND

	// if declarer leads a low trump, and there are still higher trumps
	// smear it!

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

	s.trump = CLUBS
	s.trick = []Card{Card{CLUBS, "K"}}
	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
	}

	validCards := []Card{
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "D"},
	}
	//

	card := player.playerTactic(&s, validCards)
	exp := Card{CLUBS, "D"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, with trumps in game: %v and valid %v, expected to play (smear) %v, played %v",
			s.trick, s.trumpsInGame, validCards, exp, card)
	}
}

func TestOpponentTacticMID_PartnerLeadsAVoidCard_Trump_only_if_still_points(t *testing.T) {
	validCards := []Card{
		Card{HEART, "J"},
		Card{CARO, "A"},
		Card{CARO, "K"},
		Card{CARO, "7"},
	}
	player := makePlayer([]Card{})
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer
	s.opp2 = &player
	s.opp1 = &teamMate
	players = []PlayerI{&teamMate, &player, &otherPlayer}

	s.trump = CLUBS
	s.trick = []Card{Card{HEART, "8"}}
	s.follow = HEART
	s.trumpsInGame = []Card{
		Card{HEART, "J"},
		Card{CLUBS, "D"},
		Card{CLUBS, "9"},
	}

	s.cardsPlayed = makeSuitDeck(HEART)
	s.cardsPlayed = remove(s.cardsPlayed, s.trick...)
	s.cardsPlayed = remove(s.cardsPlayed, Card{HEART, "9"})

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "7"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, with trumps in game: %v and valid %v, expected to play %v, played %v",
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

func TestOpponentTacticMIDPlayerLeadsLosingCard_Donot_Smear_if_partner_void(t *testing.T) {
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
	s.opp2VoidSuit[SPADE] = true

	//
	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "9"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, with cards played %v and valid %v, not expected to smear %v (partner void), played %v",
			s.trick, s.cardsPlayed, validCards, exp, card)
	}
}

func TestOpponentTacticMIDPlayerLeadsLosingTRUMP_Donot_Smear_if_partner_void(t *testing.T) {
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
	s.follow = CLUBS
	s.trick = []Card{Card{CLUBS, "7"}}
	s.opp2VoidSuit[CLUBS] = true
	s.trumpsInGame = makeTrumpDeck(CLUBS)

	//
	card := player.playerTactic(&s, validCards)
	debugTacticsLog("Played: %v\n", card)
	notexp := Card{CARO, "10"}
	if card.equals(notexp) {
		t.Errorf("In trick %v, with cards played %v and valid %v, not expected to smear %v (partner void), played %v",
			s.trick, s.cardsPlayed, validCards, notexp, card)
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

func TestOpponentTacticMID1(t *testing.T) {
	// MIDDLEHAND

	// if partner leads a very low card
	// don't SMEAR the trick

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

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
	s.trumpsInGame = []Card{Card{CLUBS, "J"}, Card{CLUBS, "8"}, Card{CLUBS, "7"}, Card{SPADE, "J"}}

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
	s.opp1 = &teamMate
	s.opp2 = &player

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
	s.opp1 = &teamMate
	s.opp2 = &player

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

func TestOpponentTacticMIDPartnerLeadsVoidSuitWinner_SmearNoTrump(t *testing.T) {
	// MIDDLEHAND

	// if partner leads a winner card in a suit (check other cards in game)
	// that you are void and smear the trick.

	validCards := []Card{
		Card{CARO, "J"},
		Card{CLUBS, "D"},
		Card{HEART, "10"},
		Card{HEART, "9"},
		Card{CARO, "9"},
	}

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer
	s.opp1 = &teamMate
	s.opp2 = &player

	s.trump = GRAND
	s.trick = []Card{Card{SPADE, "K"}}
	s.follow = SPADE
	s.cardsPlayed = []Card{
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "9"},
	}

	//

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "10"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by teammate (winner), and valid %v, expected: %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticMIDPartnerLeads_Trump_PlayLowTrump(t *testing.T) {
	// MIDDLEHAND

	validCards := []Card{
		Card{HEART, "J"},
		Card{CLUBS, "10"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
	}

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer
	s.opp1 = &teamMate
	s.opp2 = &player

	s.trump = CLUBS
	s.trick = []Card{Card{CARO, "J"}}
	s.follow = CLUBS
	s.cardsPlayed = []Card{
		Card{CLUBS, "D"},
		Card{CLUBS, "7"},
	}
	s.trumpsInGame = makeTrumpDeck(CLUBS)
	s.trumpsInGame = remove(s.trumpsInGame, s.cardsPlayed...)

	//

	card := player.playerTactic(&s, validCards)
	exp := Card{CLUBS, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by teammate, and valid %v, expected: %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticMIDPartnerLeads_VOID_SUIT(t *testing.T) {
	// MIDDLEHAND

	validCards := []Card{
		Card{HEART, "J"},
		Card{CLUBS, "10"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
	}

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer
	s.opp1 = &teamMate
	s.opp2 = &player

	s.trump = CLUBS
	s.trick = []Card{Card{SPADE, "10"}}
	s.follow = SPADE
	s.cardsPlayed = []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{CLUBS, "D"},
		Card{CLUBS, "7"},
	}
	s.cardsPlayed = append(s.cardsPlayed, makeSuitDeck(SPADE)...)
	s.trumpsInGame = makeTrumpDeck(CLUBS)
	s.trumpsInGame = remove(s.trumpsInGame, s.cardsPlayed...)
	s.declarerVoidSuit[SPADE] = true
	//

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "J"}
	if !card.equals(exp) {
		t.Errorf("In trick %v by teammate, and valid %v, expected: %v, played %v",
			s.trick, validCards, exp, card)
	}
}

// func TestOpponentTacticFORE_short_long(t *testing.T) {
// 	// FOREHAND

// 	// if declarer BACK short
// 	// if declarer MID long

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
// 	}
// 	// declarer MID
// 	s.opp2 = &player
// 	s.opp1 = &teamMate

// 	card := player.playerTactic(&s, validCards)
// 	exp := Card{CARO, "K"}
// 	if !card.equals(exp) {
// 		t.Errorf("FOREHAND, DECLARER MID, valid %v, expected: %v, played %v",
// 			validCards, exp, card)
// 	}
// 	// declarer BACK
// 	s.opp1 = &player
// 	s.opp2 = &teamMate
// 	teamMate.previousSuit = ""
// 	player.previousSuit = ""

// 	card = player.playerTactic(&s, validCards)
// 	exp = Card{HEART, "A"}
// 	if !card.equals(exp) {
// 		t.Errorf("FOREHAND, DECLARER BACK, valid %v, expected: %v, played %v",
// 			validCards, exp, card)
// 	}
// }

//TODO
func TestOpponentTacticFORE_long_Not_Full_if_trumps_in_play(t *testing.T) {
	// FOREHAND

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
		Card{CARO, "A"},
	}
	s.trumpsInGame = []Card{Card{CARO, "J"}}
	// declarer MID
	s.opp2 = &player
	s.opp1 = &teamMate

	card := player.playerTactic(&s, validCards)
	exp1 := Card{CARO, "K"}
	exp2 := Card{CARO, "D"}
	if !card.equals(exp1) && !card.equals(exp1) {
		t.Errorf("FOREHAND, DECLARER MID, valid %v, expected: %v or %v, played a full one %v",
			validCards, exp1, exp2, card)
	}
}

func TestOpponentTacticFORE_Protect_the_10_of_the_Partner(t *testing.T) {
	// FOREHAND
	// if declarer BACK short
	// never 2 numbers, or D number suit
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer

	s.trump = HEART
	s.trick = []Card{}
	_ = teamMate

	validCards := []Card{
		Card{CLUBS, "J"},
		Card{CARO, "J"},
		Card{HEART, "10"},

		Card{CLUBS, "10"},
		Card{CLUBS, "K"},
		Card{CLUBS, "D"},

		Card{SPADE, "K"},
		Card{SPADE, "8"},

		Card{CARO, "8"},
		Card{CARO, "9"},
	}

	// declarer BACK
	s.opp1 = &player
	s.opp2 = &teamMate
	teamMate.previousSuit = ""
	player.previousSuit = ""

	card := player.playerTactic(&s, validCards)
	exp := Card{SPADE, "K"}
	if !card.equals(exp) {
		t.Errorf("FOREHAND, DECLARER BACK, expected: %v, played %v",
			exp, card)
	}

	validCards = remove(validCards, Card{CARO, "9"})
	validCards = append(validCards, Card{CARO, "D"})
	teamMate.previousSuit = ""
	player.previousSuit = ""

	card = player.playerTactic(&s, validCards)
	exp = Card{SPADE, "K"}
	if !card.equals(exp) {
		t.Errorf("FOREHAND, DECLARER BACK, expected: %v, played %v",
			exp, card)
	}

	validCards = remove(validCards, Card{SPADE, "K"})
	validCards = append(validCards, Card{SPADE, "9"})
	teamMate.previousSuit = ""
	player.previousSuit = ""

	card = player.playerTactic(&s, validCards)
	exp = Card{CARO, "D"}
	if !card.equals(exp) {
		t.Errorf("FOREHAND, DECLARER BACK, expected: %v, played %v",
			exp, card)
	}
}

func TestOpponentTacticFORE_Protect_the_10_of_the_Partner_2(t *testing.T) {
	// FOREHAND
	// if declarer BACK short
	// never 2 numbers, or D number suit
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer

	s.trump = CLUBS
	s.trick = []Card{}
	_ = teamMate

	validCards := []Card{
		Card{CLUBS, "J"},
		Card{CLUBS, "9"},

		Card{SPADE, "A"},
		Card{SPADE, "D"}, // not neceassary here since we have the A
		Card{SPADE, "9"},
		Card{SPADE, "7"}, // best card so that parner goes over

		Card{HEART, "10"},
		Card{HEART, "8"}, // you don't want to discard this not to make 10 blank

		Card{CARO, "8"}, // bac choice
		Card{CARO, "9"},
	}

	// declarer BACK
	s.opp1 = &player
	s.opp2 = &teamMate
	s.trumpsInGame = makeTrumpDeck(CLUBS)
	teamMate.previousSuit = ""
	player.previousSuit = ""

	card := player.playerTactic(&s, validCards)
	exp := Card{SPADE, "7"}
	if !card.equals(exp) {
		t.Errorf("FOREHAND, DECLARER BACK, expected: %v, played %v",
			exp, card)
	}
}

func TestOpponentTacticFORE_Protect_the_10_of_the_Partner_3(t *testing.T) {
	// FOREHAND
	// if declarer BACK short
	// never 2 numbers, or D number suit
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer

	s.trump = SPADE
	s.trick = []Card{}
	_ = teamMate

	validCards := []Card{
		Card{SPADE, "D"},

		Card{CLUBS, "D"}, // not so good
		Card{CLUBS, "7"},

		Card{HEART, "K"}, // too long
		Card{HEART, "9"}, //
		Card{HEART, "8"},
		Card{HEART, "7"}, //

		Card{CARO, "K"}, // good choice
		Card{CARO, "8"},
	}

	// declarer BACK
	s.opp1 = &player
	s.opp2 = &teamMate
	s.trumpsInGame = makeTrumpDeck(SPADE)
	teamMate.previousSuit = ""
	player.previousSuit = ""

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "K"}
	debugTacticsLog("PLayed: %v\n", card)
	if !card.equals(exp) {
		t.Errorf("FOREHAND, DECLARER BACK, expected: %v, played %v",
			exp, card)
	}
}

// func TestOpponentTacticFORE_short_TOD_SUENDE_1_1(t *testing.T) {
// 	// FOREHAND

// 	// if declarer BACK short

// 	// never 2 numbers, or D number suit
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
// 		Card{CLUBS, "9"},
// 		Card{CLUBS, "8"},
// 		Card{CLUBS, "7"},

// 		Card{HEART, "9"},
// 		Card{HEART, "8"},

// 		Card{SPADE, "A"},
// 		Card{SPADE, "10"},

// 		Card{CARO, "D"},
// 		Card{CARO, "9"},

// 	}

// 	// declarer BACK
// 	s.opp1 = &player
// 	s.opp2 = &teamMate
// 	teamMate.previousSuit = ""
// 	player.previousSuit = ""

// 	card := player.playerTactic(&s, validCards)
// 	notExpSuit1 := HEART
// 	if card.Suit == notExpSuit1 {
// 		t.Errorf("FOREHAND, DECLARER BACK, not expected suit of 2 numbers: %v, played %v",
// 			validCards, card)
// 	}
// 	notExpSuit2 := CARO
// 	if card.Suit == notExpSuit2 {
// 		t.Errorf("FOREHAND, DECLARER BACK, not expected suit of D-number: %v, played %v",
// 			validCards, card)
// 	}
// }

// func TestOpponentTacticFORE_short_TOD_SUENDE_1_2(t *testing.T) {
// 	// FOREHAND

// 	// if declarer BACK short

// 	// never 2 numbers, or D number suit
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
// 		Card{CLUBS, "9"},
// 		Card{CLUBS, "8"},
// 		Card{CLUBS, "7"},

// 		Card{HEART, "9"},
// 		Card{HEART, "8"},

// 		Card{SPADE, "A"},
// 		Card{SPADE, "K"},
// 		Card{SPADE, "9"},
// 	}

// 	s.trumpsInGame = makeTrumpDeck(s.trump)

// 	// declarer BACK
// 	s.opp1 = &player
// 	s.opp2 = &teamMate
// 	teamMate.previousSuit = ""
// 	player.previousSuit = ""

// 	card := player.playerTactic(&s, validCards)
// 	notExpSuit := HEART
// 	if card.Suit == notExpSuit {
// 		t.Errorf("FOREHAND, DECLARER BACK, not expected suit of 2 numbers: %v, played %v",
// 			validCards, card)
// 	}
// }

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

	card := player.playerTactic(&s, validCards)
	if card.Suit == "" && card.Rank == "" {
		t.Errorf("Error validCards: %v, played %v",
			validCards, card)
	}
}

func TestOpponentTacticFORE_StrongTrumps(t *testing.T) {
	// FOREHAND
	// if declarer is short on Trumps and you are strong

	validCards := []Card{
		Card{CLUBS, "J"},
		Card{SPADE, "J"},
		Card{HEART, "8"},

		Card{CLUBS, "8"},
		Card{CARO, "9"},
		Card{CARO, "8"},
	}
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer

	teamMate.previousSuit = ""
	player.previousSuit = ""
	s.trump = HEART
	s.trick = []Card{}
	_ = teamMate

	// declarer MIDDLE
	s.opp2 = &player
	s.opp1 = &teamMate

	s.trumpsInGame = []Card{
		Card{HEART, "J"},
		Card{HEART, "10"},
	}

	card := player.playerTactic(&s, validCards)
	if !card.equals(Card{CLUBS, "J"}) && !card.equals(Card{SPADE, "J"}) {
		t.Errorf("Error validCards: %v, played %v",
			validCards, card)
	}
}

func TestOpponentTacticBACK_MateLeads(t *testing.T) {

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

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

func TestOpponentTacticBACK_MateLeads_PlayerWins_Dont_Play_A_trump_on_A_Zero_trick(t *testing.T) {

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

	s.trump = CARO
	s.trick = []Card{Card{CLUBS, "7"}, Card{CLUBS, "9"}}
	s.follow = CLUBS

	validCards := []Card{
		Card{CARO, "J"},
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

func TestOpponentTacticBACK_MateLeads_PlayerWins_Dont_Play_A_trump_on_A_Zero_trick_Unless_Saving_a_FULL(t *testing.T) {

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

	s.trump = CARO
	s.trick = []Card{Card{CLUBS, "7"}, Card{CLUBS, "9"}}
	s.follow = CLUBS

	validCards := []Card{
		Card{CARO, "10"},
		Card{SPADE, "D"},
		Card{HEART, "10"},
		Card{HEART, "D"},
		Card{HEART, "8"},
	}
	//

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "10"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, TRUMPS CARO, and valid %v, expected to save a trump: %v, played: %v",
			s.trick, validCards, exp, card)
	}
}

// Does not improve??
// func TestOpponentTacticBACK_PlayerLeads_(t *testing.T) {
// 	validCards := []Card{
// 		Card{HEART, "J"},
// 		Card{SPADE, "9"},
// 		Card{HEART, "8"},
// 	}
// 	player := makePlayer(validCards)

// 	otherPlayer := makePlayer([]Card{})
// 	// teamMate := makePlayer([]Card{})
// 	s := makeSuitState()
// 	s.leader = &otherPlayer
// 	s.declarer = &otherPlayer

// 	s.trump = CLUBS
// 	s.trick = []Card{Card{CARO, "9"}, Card{CARO, "K"}}
// 	s.follow = CARO

// 	card := player.playerTactic(&s, validCards)
// 	if getSuit(s.trump, card) == s.trump {
// 		t.Errorf("In trick led by declarer %v, and valid %v, NOT expected a trump: %v",
// 			s.trick, validCards, card)
// 	}
// }

// func TestOpponentTacticBACK_PlayerLeads_PutPlayerInMiddleHand(t *testing.T) {
// 	validCards := []Card{
// 		Card{CLUBS, "J"},
// 		Card{HEART, "J"},
// 		Card{CARO, "8"},
// 	}
// 	player := makePlayer(validCards)
// 	otherPlayer := makePlayer([]Card{})
// 	teamMate := makePlayer([]Card{})
// 	s := makeSuitState()
// 	s.leader = &otherPlayer
// 	s.declarer = &otherPlayer
// 	s.opp1 = &teamMate
// 	s.opp2 = &player

// 	s.trump = CARO
// 	s.trick = []Card{Card{CARO, "7"}, Card{CARO, "A"}}
// 	s.follow = CARO

// 	card := player.playerTactic(&s, validCards)
// 	exp := Card{HEART, "J"}
// 	if !card.equals(exp) {
// 		t.Errorf("In trick led by declarer %v, and valid %v, expected: %v to bring declarer at MIDDLEHAND, played: %v",
// 			s.trick, validCards, exp, card)
// 	}
// }

func TestOpponentTacticBACK_PlayerLeads_Trump(t *testing.T) {
	validCards := []Card{
		Card{CLUBS, "J"},
		Card{CARO, "8"},
	}
	player := makePlayer(validCards)
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &teamMate
	s.opp2 = &player

	s.trump = CARO
	s.trick = []Card{Card{CARO, "J"}, Card{HEART, "J"}}
	s.follow = CARO

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick led by declarer %v, and valid %v, expected: %v, played: %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticBACK2(t *testing.T) {

	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &teamMate
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

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

func TestOpponentTacticBACK_PlayHighTrumpWhenParterSMEARS(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	validCards := []Card{
		Card{SPADE, "J"},
		Card{SPADE, "9"},
		Card{CLUBS, "9"},
		Card{CARO, "D"},
	}
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &teamMate
	s.opp2 = &player

	s.trump = SPADE
	s.trick = []Card{Card{SPADE, "A"}, Card{SPADE, "10"}}
	s.follow = SPADE
	//
	card := player.playerTactic(&s, validCards)

	exp := Card{SPADE, "J"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, TRUMPS SPADE, and valid %v, expected : %v, got %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticBACK_DoNotWasteHighTrump(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	validCards := []Card{
		Card{SPADE, "J"},
		Card{SPADE, "9"},
		Card{CLUBS, "9"},
		Card{CARO, "D"},
	}
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &teamMate
	s.opp2 = &player

	s.trump = SPADE
	s.trick = []Card{Card{HEART, "K"}, Card{HEART, "D"}}
	s.follow = HEART
	//
	card := player.playerTactic(&s, validCards)

	exp := Card{SPADE, "9"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, TRUMPS SPADE, and valid %v, expected : %v, got %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticBACK_DoNotWasteTrumpOnZeroValueTrickIfYouCanWinAnotherTrump(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	validCards := []Card{
		Card{SPADE, "J"},
		Card{CARO, "D"},
	}
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &teamMate
	s.opp2 = &player

	s.trump = SPADE
	s.trick = []Card{Card{HEART, "8"}, Card{CARO, "7"}}
	s.follow = HEART

	s.trumpsInGame = []Card{
		Card{SPADE, "J"},
		Card{SPADE, "K"},
	}
	s.cardsPlayed = remove(makeTrumpDeck(SPADE), s.trumpsInGame...)
	//
	card := player.playerTactic(&s, validCards)

	exp := Card{CARO, "D"}
	if !card.equals(exp) {
		t.Errorf("In trick %v, TRUMPS SPADE, and valid %v, expected : %v, got %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticBACK_DoNotWasteA_FullOne_on_a_trick_to_save_Trumps(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	validCards := []Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{SPADE, "D"},

		Card{CLUBS, "10"},
	}
	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &teamMate
	s.opp2 = &player

	s.trump = SPADE
	s.trick = []Card{Card{HEART, "8"}, Card{CLUBS, "9"}}
	s.follow = HEART

	s.trumpsInGame = makeTrumpDeck(s.trump)

	// s.cardsPlayed = remove(makeTrumpDeck(SPADE), s.trumpsInGame...)
	//
	card := player.playerTactic(&s, validCards)

	notexp := Card{CLUBS, "10"}
	if card.equals(notexp) {
		t.Errorf("In trick %v, TRUMPS SPADE, and valid %v, NOT expected : %v, got %v",
			s.trick, validCards, notexp, card)
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
	s.cardsPlayed = remove(makeTrumpDeck(s.trump), s.trumpsInGame...)

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "J"}
	if !card.equals(exp) {
		t.Errorf("Trump: CLUBS, In trick %v and valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticFORE_DontPlayA10ifJsout(t *testing.T) {
	// don't play your A-10 trumps if Js still there
	validCards := []Card{
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CLUBS, "10"},
		Card{SPADE, "A"},
		Card{SPADE, "D"},
	}

	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &player
	s.declarer = &player

	s.trump = CARO
	s.trick = []Card{}

	s.trumpsInGame = []Card{
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "7"},
	}
	s.cardsPlayed = remove(makeTrumpDeck(s.trump), s.trumpsInGame...)

	card := player.playerTactic(&s, validCards)
	exp := Card{SPADE, "A"}
	if !card.equals(exp) {
		t.Errorf("Trump: CLUBS, In trick %v and valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}
}

// not sure about it
// func TestDeclarerTacticFORE_LowTrumps(t *testing.T) {
// 	// -r 2612 Bob
// 	// don't play trumps if you are low
// 	validCards := []Card{
// 		Card{CARO, "10"},
// 		Card{CARO, "K"},
// 		Card{CARO, "9"},
// 		Card{SPADE, "A"},
// 		Card{SPADE, "D"},
// 	}

// 	player := makePlayer(validCards)
// 	other := makePlayer(validCards)
// 	s := makeSuitState()
// 	s.leader = &player
// 	s.declarer = &player
// 	s.opp1 = &other

// 	s.trump = CARO
// 	s.trick = []Card{}
// 	s.opp1VoidSuit[s.trump] = true

// 	s.trumpsInGame = []Card{
// 		Card{HEART, "J"},
// 		Card{CARO, "J"},
// 		Card{CARO, "10"},
// 		Card{CARO, "K"},
// 		Card{CARO, "Q"},
// 		Card{CARO, "9"},
// 	}
// 	s.cardsPlayed = remove(makeTrumpDeck(s.trump), s.trumpsInGame...)

// 	card := player.playerTactic(&s, validCards)
// 	exp := Card{SPADE, "A"}
// 	if !card.equals(exp) {
// 		t.Errorf("Trump: CLUBS, opp1 has no more trumps. In trick %v and valid %v, expected to play %v, played %v",
// 			s.trick, validCards, exp, card)
// 	}
// }

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

func TestDeclarerTacticMID_LowTrump(t *testing.T) {
	// don't play a high trump if you are not strong

	validCards := []Card{
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{SPADE, "9"},
	}
	player := makePlayer(validCards)
	other := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = SPADE
	s.trick = []Card{Card{HEART, "9"}}

	card := player.playerTactic(&s, validCards)
	exp := Card{SPADE, "9"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v. Played %v",
			s.trump, s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticMID_ThrowOff(t *testing.T) {
	validCards := []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{SPADE, "D"},

		Card{CARO, "9"},
		Card{HEART, "A"},
		Card{HEART, "K"},
		Card{HEART, "8"},
	}
	player := makePlayer(validCards)
	other := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = SPADE
	s.trick = []Card{Card{CLUBS, "9"}}

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "9"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v. Played %v",
			s.trump, s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticMID_DontThrowOffTheProtectorOfa10(t *testing.T) {
	validCards := []Card{
		Card{HEART, "J"},
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "D"},

		Card{CLUBS, "10"},
		Card{CLUBS, "7"},
	}
	player := makePlayer(validCards)
	other := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = HEART
	s.trick = []Card{Card{SPADE, "9"}}
	s.trumpsInGame = []Card{
		Card{HEART, "J"},
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "K"}, //
		Card{HEART, "D"},
		Card{HEART, "9"},
		Card{HEART, "8"},
	}
	s.cardsPlayed = remove(makeTrumpDeck(s.trump), s.trumpsInGame...)

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v (sure winner) to keep the protector of 10. Played %v",
			s.trump, s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticMID_DontThrowOffAor10(t *testing.T) {
	validCards := []Card{
		Card{SPADE, "9"},
		Card{HEART, "A"},
		Card{HEART, "10"},
	}
	player := makePlayer(validCards)
	other := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = SPADE
	s.trick = []Card{Card{CARO, "7"}}
	s.trumpsInGame = []Card{
		Card{SPADE, "9"},
	}
	s.cardsPlayed = remove(makeTrumpDeck(s.trump), s.trumpsInGame...)

	card := player.playerTactic(&s, validCards)
	exp := Card{SPADE, "9"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v not an A or 10. Played %v",
			s.trump, s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticMID_ThrowOffToGoBack(t *testing.T) {
	// throw off on a null opener to go in backhand

	validCards := []Card{
		Card{CLUBS, "J"},
		Card{HEART, "A"},
		Card{HEART, "D"},
		Card{SPADE, "A"},
		Card{SPADE, "7"},
	}
	player := makePlayer(validCards)
	other := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = HEART
	s.trick = []Card{Card{CARO, "9"}}
	s.follow = CARO

	card := player.playerTactic(&s, validCards)
	exp := Card{SPADE, "7"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v. Played %v",
			s.trump, s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticBACK_TrumpWithLowJack(t *testing.T) {

	validCards := []Card{
		Card{SPADE, "J"},
		Card{CARO, "J"},
		Card{HEART, "D"},
		Card{SPADE, "A"},
		Card{SPADE, "7"},
	}
	player := makePlayer(validCards)
	other := makePlayer([]Card{})
	other1 := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player
	s.opp1 = &other
	s.opp2 = &other1

	s.trump = GRAND
	s.trick = []Card{Card{CLUBS, "A"}, Card{CLUBS, "8"}}
	s.follow = CLUBS

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "J"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v. Played %v",
			s.trump, s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticBACK_DontWasteYourAonaZeroTrick(t *testing.T) {

	validCards := []Card{
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "D"},
		Card{HEART, "A"},
		Card{HEART, "8"},
	}
	player := makePlayer(validCards)
	other := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = CARO
	s.trick = []Card{Card{HEART, "7"}, Card{HEART, "9"}}
	s.follow = HEART

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "8"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v. Played %v",
			s.trump, s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticBACK_DontWasteYourAonaZeroTrick_UnlessYouHavethe10(t *testing.T) {

	validCards := []Card{
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CARO, "D"},
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "K"},
		Card{HEART, "8"},
	}
	player := makePlayer(validCards)
	other := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player

	s.trump = CARO
	s.trick = []Card{Card{HEART, "7"}, Card{HEART, "9"}}
	s.follow = HEART

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v. Played %v",
			s.trump, s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticBACK_Dont_Throw_a_FullOne_on_a_zero_trick(t *testing.T) {

	validCards := []Card{
		Card{HEART, "J"},
		Card{HEART, "10"},
		Card{HEART, "K"},
		Card{HEART, "8"},
		Card{HEART, "7"},
		Card{CARO, "10"},
	}
	player := makePlayer(validCards)
	other1 := makePlayer([]Card{})
	other2 := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other1
	s.declarer = &player
	s.opp1 = &other1
	s.opp2 = &other2
	players = []PlayerI{&other1, &other2, &player}

	s.trump = HEART
	s.trick = []Card{Card{CLUBS, "7"}, Card{CLUBS, "9"}}
	s.follow = CLUBS

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "10"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v. Played %v",
			s.trump, s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticBACK_Throw_MAX_a_D_a_zero_trick(t *testing.T) {

	validCards := []Card{
		Card{HEART, "J"},
		Card{HEART, "10"},
		Card{HEART, "K"},
		Card{HEART, "8"},
		Card{HEART, "7"},
		Card{CARO, "10"},
		Card{CARO, "D"},
	}
	player := makePlayer(validCards)
	other1 := makePlayer([]Card{})
	other2 := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other1
	s.declarer = &player
	s.opp1 = &other1
	s.opp2 = &other2
	players = []PlayerI{&other1, &other2, &player}

	s.trump = HEART
	s.trick = []Card{Card{CLUBS, "7"}, Card{CLUBS, "9"}}
	s.follow = CLUBS

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "D"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v. Played %v",
			s.trump, s.trick, validCards, exp, card)
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

	s.cardsPlayed = makeDeck()
	s.cardsPlayed = remove(s.cardsPlayed, player.hand...)

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
	s.cardsPlayed = makeDeck()
	s.cardsPlayed = remove(s.cardsPlayed, player.hand...)
	s.cardsPlayed = remove(s.cardsPlayed, Card{HEART, "10"})

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

func TestDiscardInSkat1(t *testing.T) {
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
	// skat := []Card{Card{CARO, "7"}, Card{SPADE, "D"}}
	p := makePlayer(cards)
	skat := p.cardsToDiscard(CLUBS)

	if in(skat, Card{SPADE, "A"}) || in(skat, Card{HEART, "A"}) || in(skat, Card{CARO, "A"}) {
		t.Errorf("A discarded in SKAT: %v", skat)
	}
}

func TestDiscardInSkatKeepLongSuit(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},

		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "8"},
		Card{CARO, "7"},

		Card{SPADE, "K"},
		Card{SPADE, "D"},
		Card{SPADE, "9"},

		Card{CLUBS, "K"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
	}
	// skat := []Card{Card{CARO, "7"}, Card{SPADE, "D"}}
	p := makePlayer(cards)
	p.trumpToDeclare = SPADE
	p.declaredBid = 20

	skat := p.cardsToDiscard(p.trumpToDeclare)

	// fmt.Printf("SKAT: %v\n", skat)
	if in(skat, Card{CARO, "A"}) || in(skat, Card{CARO, "10"}) || in(skat, Card{CARO, "8"}) || in(skat, Card{CARO, "7"}) {
		t.Errorf("Cards from a long suit %v discarded in SKAT: %v", p.hand, skat)
	}
}

func TestDiscardInSkatDiscart10s(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},

		Card{HEART, "A"},
		Card{HEART, "D"},
		Card{HEART, "8"},
		Card{HEART, "7"},

		Card{SPADE, "A"},
		Card{SPADE, "7"},

		Card{CARO, "K"},
		Card{CARO, "7"},

		Card{CLUBS, "10"},
		Card{CLUBS, "8"},
	}
	skat := []Card{Card{CLUBS, "8"}, Card{CARO, "K"}}
	p := makePlayer(cards)

	p.discardInSkat(skat)

	// fmt.Printf("SKAT: %v\n", skat)
	if !in(skat, Card{CLUBS, "10"}, Card{CLUBS, "8"}) {
		t.Errorf("Cards from 10-X suit %v NOT discarded in SKAT: %v", p.hand, skat)
	}
}

func TestDiscardInSkatNULLNoRisk(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "7"},
		Card{CLUBS, "9"},
		Card{CLUBS, "10"},
		Card{CLUBS, "D"},

		Card{SPADE, "7"},
		Card{SPADE, "8"},
		Card{SPADE, "10"},
		Card{SPADE, "D"},

		Card{CARO, "7"},
		Card{CARO, "8"},
		Card{CARO, "J"},
		Card{CARO, "D"},
	}
	skat := []Card{Card{SPADE, "8"}, Card{SPADE, "D"}}
	p := makePlayer(cards)
	p.discardInSkat(skat)

	if len(p.hand) != 10 {
		t.Errorf("Not discarded (12 cards in hand): %v", p.hand)
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

// TODO::
// TODO::
// func TestDiscardInSkatBigCards(t *testing.T) {
// 	cards := []Card{
// 		Card{CLUBS, "J"},
// 		Card{HEART, "J"},
// 		Card{SPADE, "10"},
// 		Card{SPADE, "D"},
// 		Card{SPADE, "A"},
// 		Card{SPADE, "9"},
// 		Card{SPADE, "7"},

// 		Card{CLUBS, "K"},

// 		Card{HEART, "10"},
// 		Card{HEART, "7"},

// 		Card{CARO, "9"},
// 		Card{CARO, "8"},
// 	}
// 	skat := []Card{Card{SPADE, "D"}, Card{HEART, "7"}}
// 	p := makePlayer(cards)
// 	p.trumpToDeclare = SPADE
// 	p.risky = true
// 	p.discardInSkat(skat)

// 	cardsToDiscard := []Card{Card{HEART, "10"}, Card{CLUBS, "K"}}

// 	if !in(skat, cardsToDiscard...){
// 		t.Errorf("From hand: %v discarded in SKAT: %v instead of %v", p.hand, skat, cardsToDiscard)
// 	}
// }

func TestDiscardInSkatGRAND_10s(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},
		Card{CLUBS, "10"}, // DISCARD
		Card{CLUBS, "D"},
		Card{SPADE, "A"},
		Card{SPADE, "D"},
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{CARO, "10"}, // DISCARD
		Card{CARO, "D"},
		Card{CARO, "7"},
	}
	// skat := []Card{Card{SPADE, "9"}, Card{HEART, "9"}}
	p := makePlayer(cards)
	p.trumpToDeclare = GRAND
	p.risky = true
	skat := p.cardsToDiscard(p.trumpToDeclare)
	//
	if !in(skat, Card{CLUBS, "10"}) || !in(skat, Card{CARO, "10"}) {
		t.Errorf("From hand: %v discarded in SKAT: %v instead of 2 10s", p.hand, skat)
	}
}

func TestDiscardInSkatGRANDBlank(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "J"},
		Card{CARO, "J"}, //LOSER
		Card{CLUBS, "A"},
		Card{CLUBS, "8"}, //LOSER
		Card{CLUBS, "9"}, //LOSER
		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "D"}, //LOSER
		Card{SPADE, "9"}, //LOSER
		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{CARO, "9"}, //LOSER
	}
	// skat := []Card{Card{CLUBS, "9"}, Card{SPADE, "D"}}
	p := makePlayer(cards)
	p.trumpToDeclare = GRAND
	skat := p.cardsToDiscard(p.trumpToDeclare)

	if !in(skat, Card{SPADE, "D"}, Card{CARO, "9"}) {
		t.Errorf("hand: %v discarded in SKAT: %v", p.hand, skat)
	}
}

func TestDiscardInSkat_2CardsOfASuit(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"},
		Card{HEART, "J"},
		Card{CARO, "J"},

		Card{CLUBS, "10"},
		Card{CLUBS, "D"},
		Card{CLUBS, "8"},

		Card{HEART, "A"},
		Card{HEART, "D"},
		Card{HEART, "7"},

		Card{SPADE, "D"},
		Card{SPADE, "9"},

		Card{CARO, "7"},
	}
	skat := []Card{Card{CLUBS, "8"}, Card{SPADE, "D"}}
	p := makePlayer(cards)
	p.trumpToDeclare = CLUBS
	p.discardInSkat(skat)

	if !in(skat, Card{SPADE, "D"}, Card{SPADE, "9"}) {
		t.Errorf("Final hand: %v discarded in SKAT: %v", p.hand, skat)
	}
}

func TestDiscardInSkat_Blank10_2CardsOfASuit(t *testing.T) {
	cards := []Card{
		Card{SPADE, "J"},
		Card{CARO, "J"},

		Card{HEART, "A"},
		Card{HEART, "10"},
		Card{HEART, "K"},
		Card{HEART, "D"},
		Card{HEART, "8"},

		Card{CLUBS, "10"},
		Card{SPADE, "10"},
		Card{SPADE, "9"},

		Card{CARO, "8"},
		Card{CARO, "7"},
	}
	skat := []Card{Card{CLUBS, "10"}, Card{SPADE, "10"}}
	p := makePlayer(cards)
	p.trumpToDeclare = CLUBS
	p.discardInSkat(skat)

	if !in(skat, Card{CLUBS, "10"}, Card{SPADE, "10"}) {
		t.Errorf("Final hand: %v discarded in SKAT: %v", p.hand, skat)
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
	p := makePlayer(cards)
	p.trumpToDeclare = CLUBS

	skat := p.cardsToDiscard(p.trumpToDeclare)

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

// DECREASE PERCENTAGE

// func TestDiscardInSkat_CanWinNonAsuitWithLessCardsThanAsuit(t *testing.T) {
// 	cards := []Card{
// 		Card{CARO, "J"},
// 		Card{SPADE, "J"},

// 		Card{CARO, "K"}, // declare this
// 		Card{CARO, "9"},
// 		Card{CARO, "8"},

// 		Card{CLUBS, "D"},

// 		Card{SPADE, "A"}, // keep this as a strong by-suit
// 		Card{SPADE, "10"},
// 		Card{SPADE, "K"},
// 		Card{SPADE, "8"},

// 		Card{HEART, "10"},
// 		Card{HEART, "8"},
// 	}
// 	p := makePlayer(cards)

// 	canWin := p.canWin(true)
// 	if canWin != "SUIT" {
// 		t.Errorf("Hand %v can win CARO to keep the strong by-suit. Got: %v", p.hand, canWin)
// 	}

// 	skat := []Card{Card{HEART, "10"}, Card{HEART, "8"}}

// 	p.discardInSkat(skat)

// 	exp := CARO
// 	act := p.trumpToDeclare
// 	if exp != act {
// 		t.Errorf("Hand %v can win %v to keep the strong by-suit HEART. Got: %v", p.hand, exp, act)
// 	}
// }

func TestCanWinNULL1(t *testing.T) {
	cards := []Card{
		Card{HEART, "D"},
		Card{HEART, "K"},

		Card{CLUBS, "A"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
		Card{CARO, "10"},
		Card{CARO, "9"},
		Card{CARO, "K"},
		Card{CARO, "7"},
	}
	p := makePlayer(cards)

	canWin := p.canWin(false)
	if canWin == NULL {
		t.Errorf("Hand %v can not win NULL 2 cards to discard before SKAT pick up. Got: %v", p.hand, canWin)
	}
}

func TestCanWinNULL2(t *testing.T) {
	cards := []Card{
		Card{CLUBS, "D"},

		Card{HEART, "9"},

		Card{CLUBS, "A"},
		Card{CLUBS, "9"},
		Card{CLUBS, "8"},
		Card{CLUBS, "7"},
		Card{CARO, "10"},
		Card{CARO, "9"},
		Card{CARO, "K"},
		Card{CARO, "7"},
	}
	p := makePlayer(cards)

	canWin := p.canWin(false)
	if canWin != NULL {
		t.Errorf("Hand %v can win NULL. 1 card to discard before SKAT pick up. Got: %v", p.hand, canWin)
	}
}

func TestCanWinGRAND1(t *testing.T) {
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
	p2 := makePlayer(cards)
	p3 := makePlayer(cards)
	p.risky = true

	players = []PlayerI{&p, &p2, &p3}
	canWin := p.canWin(false)
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
	p2 := makePlayer(cards)
	p3 := makePlayer(cards)
	p.risky = true

	players = []PlayerI{&p, &p2, &p3}

	canWin := p.canWin(false)
	if canWin != "GRAND" {
		t.Errorf("Hand %v can win GRAND. Got: %v", p.hand, canWin)
	}
}

func TestOpponentTacticNULLBack1(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &teamMate
	s.opp2 = &player

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

func TestOpponentTacticNULLBack2(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &otherPlayer
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate

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
	s.opp1 = &teamMate
	s.opp2 = &player

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
	s.opp1 = &teamMate
	s.opp2 = &player

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

	s.opp1 = &teamMate
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
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()

	s.opp1 = &teamMate
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
	s.opp1 = &player
	s.opp2 = &teamMate
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

func TestOpponentTacticFORE_PlayerVoid(t *testing.T) {
	validCards := []Card{
		Card{CARO, "J"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "8"},
	}
	player := makePlayer(validCards)

	teamMate := makePlayer([]Card{})
	otherPlayer := makePlayer([]Card{})
	s := makeSuitState()
	s.declarer = &otherPlayer
	s.leader = &player
	s.opp1 = &teamMate
	s.opp2 = &player
	s.trump = CLUBS
	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
	}

	teamMate.previousSuit = "CARO"
	s.trick = []Card{}

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "A"}
	if !card.equals(exp) {
		t.Errorf("VOID: %v, In trick %v and valid %v, expected to follow previous trick and play %v, played %v",
			s.declarerVoidSuit, s.trick, validCards, exp, card)
	}

	s.declarerVoidSuit[CARO] = true
	card = player.playerTactic(&s, validCards)
	exp = Card{CARO, "8"}
	if !card.equals(exp) {
		t.Errorf("VOID: %v,In trick %v and valid %v, decl CARO void, expected to play Low %v, played %v",
			s.declarerVoidSuit, s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticFORE_PlayerVoid2(t *testing.T) {
	validCards := []Card{
		Card{CARO, "J"},
		Card{SPADE, "K"},
		Card{SPADE, "D"},
		Card{SPADE, "8"},
		Card{CARO, "10"},
		Card{HEART, "A"},
		Card{HEART, "D"},
		Card{HEART, "8"},
	}
	player := makePlayer(validCards)

	teamMate := makePlayer([]Card{})
	otherPlayer := makePlayer([]Card{})
	s := makeSuitState()
	s.declarer = &otherPlayer
	s.leader = &player
	s.opp1 = &teamMate
	s.opp2 = &player
	s.trump = CLUBS
	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
	}

	teamMate.previousSuit = "CARO"
	s.trick = []Card{}

	s.declarerVoidSuit[CARO] = true

	card := player.playerTactic(&s, validCards)
	unexp := Card{CARO, "10"}
	if card.equals(unexp) {
		t.Errorf("DEclarer void: %v. In trick %v and valid %v, not expected to play %v",
			s.declarerVoidSuit, s.trick, validCards, unexp)
	}
}

func TestOpponentTacticFORE_PlayerNoTrumps(t *testing.T) {
	validCards := []Card{
		Card{CARO, "J"},
		Card{CLUBS, "10"},
		Card{CLUBS, "7"},
		Card{CARO, "7"},
		Card{HEART, "7"},
	}
	player := makePlayer(validCards)

	teamMate := makePlayer([]Card{})
	otherPlayer := makePlayer([]Card{})
	s := makeSuitState()
	s.declarer = &otherPlayer
	s.leader = &player
	s.opp1 = &teamMate
	s.opp2 = &player
	s.trump = SPADE
	s.trumpsInGame = []Card{
		Card{CARO, "J"},
	}
	s.cardsPlayed = remove(makeTrumpDeck(s.trump), s.trumpsInGame...)

	player.previousSuit = "HEART"
	s.trick = []Card{}

	card := player.playerTactic(&s, validCards)
	unexp := Card{CARO, "J"}
	if card.equals(unexp) {
		t.Errorf("In trick %v and valid %v, all trumps exhausted it is not expected to play %v",
			s.trick, validCards, unexp)
	}
}

func TestOpponentTacticFORE_ToPartnerLongSuitWithAss(t *testing.T) {
	// declarer at backhand
	// play a low card to allow the partner take the trick
	validCards := []Card{
		Card{SPADE, "J"},
		Card{CARO, "J"},
		Card{HEART, "10"},
		Card{CARO, "A"},
		Card{CARO, "10"},
		Card{CARO, "K"},
		Card{CARO, "8"},
	}
	player := makePlayer(validCards)

	teamMate := makePlayer([]Card{})
	otherPlayer := makePlayer([]Card{})
	s := makeSuitState()
	s.declarer = &otherPlayer
	s.leader = &player
	s.opp2 = &teamMate
	s.opp1 = &player
	s.trump = CLUBS
	s.trumpsInGame = []Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},
	}

	s.trick = []Card{}

	card := player.playerTactic(&s, validCards)
	exp := Card{CARO, "8"}
	if !card.equals(exp) {
		t.Errorf("In trick %v and valid %v, expected to play %v, played %v",
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

func TestOverbid(t *testing.T) {
	p := makePlayer([]Card{
		Card{CLUBS, "J"},
		Card{HEART, "J"},

		Card{CARO, "8"},
		Card{CARO, "D"},

		Card{SPADE, "9"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},

		Card{HEART, "A"}, // loser
		Card{HEART, "D"}, // loser
		Card{HEART, "8"}, // loser
	})

	skat := []Card{
		Card{HEART, "9"},
		Card{HEART, "K"},
	}

	p.calculateHighestBid(false)
	p.declaredBid = p.highestBid
	p.pickUpSkat(skat)
	trump := p.declareTrump()
	if p.getGamevalue(trump) < p.declaredBid {
		t.Errorf("OVERBID")
	}
}

func TestChangeGrandToSuit(t *testing.T) {
	p := makePlayer([]Card{
		Card{SPADE, "J"},

		Card{CLUBS, "A"},
		// Card{CLUBS, "10"},
		Card{CLUBS, "9"},

		Card{SPADE, "A"},
		Card{SPADE, "10"},
		Card{SPADE, "K"},
		Card{SPADE, "D"},

		Card{HEART, "A"}, // loser
		// Card{HEART, "K"},// loser
		Card{HEART, "8"}, // loser
		Card{HEART, "7"}, // loser
	})
	o1 := makePlayer([]Card{})
	o2 := makePlayer([]Card{})
	players = []PlayerI{&o2, &p, &o1}

	s := makeSuitState()
	s.skat = []Card{
		Card{CLUBS, "10"},
		Card{HEART, "K"},
	}

	p.declaredBid = 20
	p.calculateHighestBid(true)
	p.pickUpSkat(s.skat)
	trump := p.declareTrump()
	if p.getGamevalue(trump) < p.declaredBid {
		t.Errorf("OVERBID")
	}

	debugTacticsLog("DECLARE: %v\n", p.trumpToDeclare)
	if trump == GRAND {
		t.Errorf("NO")
	}
}

func TestOpponentTacticFORE_No_trumps_in_Game1(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})
	player := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate
	s.trump = CLUBS
	s.trumpsInGame = []Card{}
	s.trick = []Card{}
	s.cardsPlayed = makeTrumpDeck(s.trump)

	player.previousSuit = CARO

	validCards := []Card{
		Card{SPADE, "9"},
		Card{HEART, "D"},
		Card{HEART, "A"},
	}

	player.hand = validCards
	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("No trumps in Game, In trick %v and valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestOpponentTacticFORE_No_trumps_in_Game2(t *testing.T) {
	otherPlayer := makePlayer([]Card{})
	teamMate := makePlayer([]Card{})

	validCards := []Card{
		Card{CARO, "J"},
		Card{HEART, "D"},
		Card{HEART, "A"},
	}

	player := makePlayer(validCards)
	s := makeSuitState()
	s.leader = &player
	s.declarer = &otherPlayer
	s.opp1 = &player
	s.opp2 = &teamMate
	s.trump = GRAND
	s.trumpsInGame = []Card{Card{CARO, "J"}}
	s.trick = []Card{}
	s.cardsPlayed = makeTrumpDeck(s.trump)
	s.cardsPlayed = remove(s.cardsPlayed, s.trumpsInGame...)

	debugTacticsLog("Suistate: %v\n", s)

	card := player.playerTactic(&s, validCards)
	exp := Card{HEART, "A"}
	if !card.equals(exp) {
		t.Errorf("No trumps in Game, In trick %v and valid %v, expected to play %v, played %v",
			s.trick, validCards, exp, card)
	}
}

func TestDeclarerTacticBACK_Just_Enough_To_Win(t *testing.T) {

	validCards := []Card{
		Card{HEART, "J"},
		Card{CARO, "10"},
		Card{CARO, "K"},
	}
	player := makePlayer(validCards)
	other := makePlayer([]Card{})
	other1 := makePlayer([]Card{})
	s := makeSuitState()
	s.leader = &other
	s.declarer = &player
	s.opp1 = &other
	s.opp2 = &other1
	s.skat = []Card{
		Card{HEART, "9"},
		Card{SPADE, "8"},
	}

	s.trump = CLUBS
	s.trick = []Card{Card{SPADE, "K"}, Card{SPADE, "7"}}
	s.follow = SPADE

	player.score = 38
	s.cardsPlayed = makeDeck()

	// remaining 11 points: 38 + 11 = 49
	// +10 (caro 10) + 2 (HEART J) = 61!
	remaining := []Card{
		Card{CLUBS, "7"},
		Card{SPADE, "9"},
		Card{SPADE, "A"},
		Card{CARO, "8"},
	}

	s.cardsPlayed = remove(s.cardsPlayed, remaining...)
	s.cardsPlayed = remove(s.cardsPlayed, s.skat...)
	s.cardsPlayed = remove(s.cardsPlayed, player.hand...)
	s.trumpsInGame = makeTrumpDeck(s.trump)
	s.trumpsInGame = remove(s.trumpsInGame, s.cardsPlayed...)

	card := player.playerTactic(&s, validCards)

	// the only chance to win
	exp := Card{CARO, "K"}
	if !card.equals(exp) {
		t.Errorf("Trump: %s, In trick %v and valid %v, was expected to play %v. Played %v",
			s.trump, s.trick, validCards, exp, card)
	}
}
