package main

type Inference struct {
//facts
	trumpsInGame     []Card
	cardsPlayed      []Card

	declarerVoidSuit map[string]bool
	opp1VoidSuit  map[string]bool
	opp2VoidSuit  map[string]bool

// beliefs
	declarerVoidSuitB map[string]bool
	opp1VoidSuitB  map[string]bool
	opp2VoidSuitB  map[string]bool

	declarerVoidCards  []Card
	opp1VoidCards []Card
	opp2VoidCards []Card
}

func makeInference() Inference {
	return Inference{
		[]Card{}, 
		[]Card{},
		map[string]bool{
			CLUBS: false,
			SPADE: false,
			HEART: false,
			CARO:  false,
		},
		map[string]bool{
			CLUBS: false,
			SPADE: false,
			HEART: false,
			CARO:  false,
		},
		map[string]bool{
			CLUBS: false,
			SPADE: false,
			HEART: false,
			CARO:  false,
		},
		map[string]bool{
			CLUBS: false,
			SPADE: false,
			HEART: false,
			CARO:  false,
		},
		map[string]bool{
			CLUBS: false,
			SPADE: false,
			HEART: false,
			CARO:  false,
		},
		map[string]bool{
			CLUBS: false,
			SPADE: false,
			HEART: false,
			CARO:  false,
		},
		[]Card{}, 
		[]Card{}, 
		[]Card{},
	}
}

func (s *Inference) cloneInference() Inference {
	newI := makeInference()

	newI.trumpsInGame = make([]Card, len(s.trumpsInGame))
	copy(newI.trumpsInGame, s.trumpsInGame)
	
	newI.cardsPlayed = make([]Card, len(s.cardsPlayed))
	copy(newI.cardsPlayed, s.cardsPlayed)
		
	newI.declarerVoidSuit = map[string]bool{
			CLUBS: s.declarerVoidSuit[CLUBS],
			SPADE: s.declarerVoidSuit[SPADE],
			HEART: s.declarerVoidSuit[HEART],
			CARO:  s.declarerVoidSuit[CARO],
		}
	newI.opp1VoidSuit = map[string]bool{
			CLUBS: s.opp1VoidSuit[CLUBS],
			SPADE: s.opp1VoidSuit[SPADE],
			HEART: s.opp1VoidSuit[HEART],
			CARO:  s.opp1VoidSuit[CARO],
		}
	newI.opp2VoidSuit = map[string]bool{
			CLUBS: s.opp2VoidSuit[CLUBS],
			SPADE: s.opp2VoidSuit[SPADE],
			HEART: s.opp2VoidSuit[HEART],
			CARO:  s.opp2VoidSuit[CARO],
		}

	// beliefs

	newI.declarerVoidSuitB = map[string]bool{
			CLUBS: s.declarerVoidSuitB[CLUBS],
			SPADE: s.declarerVoidSuitB[SPADE],
			HEART: s.declarerVoidSuitB[HEART],
			CARO:  s.declarerVoidSuitB[CARO],
		}
	newI.opp1VoidSuitB = map[string]bool{
			CLUBS: s.opp1VoidSuitB[CLUBS],
			SPADE: s.opp1VoidSuitB[SPADE],
			HEART: s.opp1VoidSuitB[HEART],
			CARO:  s.opp1VoidSuitB[CARO],
		}
	newI.opp2VoidSuitB = map[string]bool{
			CLUBS: s.opp2VoidSuitB[CLUBS],
			SPADE: s.opp2VoidSuitB[SPADE],
			HEART: s.opp2VoidSuitB[HEART],
			CARO:  s.opp2VoidSuitB[CARO],
		}

	newI.declarerVoidCards = make([]Card, len(s.declarerVoidCards))
	copy(newI.declarerVoidCards, s.declarerVoidCards)

	newI.opp1VoidCards = make([]Card, len(s.opp1VoidCards))
	copy(newI.opp1VoidCards, s.opp1VoidCards)

	newI.opp2VoidCards = make([]Card, len(s.opp2VoidCards))
	copy(newI.opp2VoidCards, s.opp2VoidCards)

	return newI
}


func (s *Inference) getDeclarerVoidSuit() map[string]bool {
	return map[string]bool{
			CLUBS: s.declarerVoidSuit[CLUBS] || s.declarerVoidSuitB[CLUBS],
			SPADE: s.declarerVoidSuit[SPADE] || s.declarerVoidSuitB[SPADE],
			HEART: s.declarerVoidSuit[HEART] || s.declarerVoidSuitB[HEART],
			CARO:  s.declarerVoidSuit[CARO] || s.declarerVoidSuitB[CARO],
		}
}

func (s *Inference) getOpp1VoidSuit() map[string]bool {
	return map[string]bool{
			CLUBS: s.opp1VoidSuit[CLUBS] || s.opp1VoidSuitB[CLUBS],
			SPADE: s.opp1VoidSuit[SPADE] || s.opp1VoidSuitB[SPADE],
			HEART: s.opp1VoidSuit[HEART] || s.opp1VoidSuitB[HEART],
			CARO:  s.opp1VoidSuit[CARO] || s.opp1VoidSuitB[CARO],
		}
}

func (s *Inference) getOpp2VoidSuit() map[string]bool {
	return map[string]bool{
			CLUBS: s.opp2VoidSuit[CLUBS] || s.opp2VoidSuitB[CLUBS],
			SPADE: s.opp2VoidSuit[SPADE] || s.opp2VoidSuitB[SPADE],
			HEART: s.opp2VoidSuit[HEART] || s.opp2VoidSuitB[HEART],
			CARO:  s.opp2VoidSuit[CARO] || s.opp2VoidSuitB[CARO],
		}
}

// Remove beliefs concerning opp player, about cards belonging the cards argument.
func (s *SuitState) detractBeliefs(opp string, cards []Card) {
	
	getSuit := func(suit string) int {
		return len(filter(cards, func (c Card) bool {
				return getSuit(s.trump, c) == suit
			}))
	}
	switch opp {
	case "decl":
		for _, suit := range suits {
			if getSuit(suit) > 0 {
				s.declarerVoidSuitB[suit] = false
			}
		}
		// s.declarerVoidSuitB = map[string]bool{
		// 		CLUBS: false,
		// 		SPADE: false,
		// 		HEART: false,
		// 		CARO:  false,
		// 	}
		s.declarerVoidCards = remove(s.declarerVoidCards, cards...)
		debugTacticsLog("Updated beliefs: %v, %v", voidString(s.declarerVoidSuitB), s.declarerVoidCards)
	case "opp1":
		for _, suit := range suits {
			if getSuit(suit) > 0 {
				s.opp1VoidSuitB[suit] = false
				s.declarer.getInference().opp1VoidSuitB[suit] = false
			}
		}
		// s.opp1VoidSuitB = map[string]bool{
		// 	CLUBS: false,
		// 	SPADE: false,
		// 	HEART: false,
		// 	CARO:  false,
		// }
		s.opp1VoidCards = remove(s.opp1VoidCards, cards...)
		debugTacticsLog("Updated beliefs: %v, %v", voidString(s.opp1VoidSuitB), s.opp1VoidCards)
	case "opp2":
		for _, suit := range suits {
			if getSuit(suit) > 0 {
				s.opp2VoidSuitB[suit] = false
			}
		}
		// s.opp2VoidSuitB = map[string]bool{
		// 	CLUBS: false,
		// 	SPADE: false,
		// 	HEART: false,
		// 	CARO:  false,
		// }
		s.opp2VoidCards = remove(s.opp2VoidCards, cards...)
		debugTacticsLog("Updated beliefs: %v, %v", voidString(s.opp2VoidSuitB), s.opp2VoidCards)
	}
}

func analysePlay(s *SuitState, p PlayerI, card Card) {

	if p.getName() == s.opp1.getName() && opponentIsLosingTrick(s, p, card, true){
		if noHigherCard(s, true, s.declarer.getHand(), card) {
			debugTacticsLog("TRICK: %v, Card: %v\n", s.trick, card)
			debugTacticsLog("INFERENCE: ***************DECLARER***************\n")
			debugTacticsLog("INFERENCE: Opp1 does not have any lower cards    \n")
			debugTacticsLog("INFERENCE: Void on suit                          \n")
			debugTacticsLog("INFERENCE: **************************************\n")
			s.declarer.getInference().opp1VoidSuitB[s.follow] = true
			debugTacticsLog("SuitState: %v, Player: %v, Card %v\n", s, p, card)
		}
	}

	// Player VOID on suit
	if len(s.trick) > 0 && getSuit(s.trump, card) != getSuit(s.trump, s.trick[0]) {
		debugTacticsLog("TRICK: %v, Card: %v\n", s.trick, card)
		debugTacticsLog("INFERENCE: **************************************\n")
		debugTacticsLog("INFERENCE: Not following tricks                  \n")
		debugTacticsLog("INFERENCE: Void on suit                          \n")
		debugTacticsLog("INFERENCE: **************************************\n")
		if p.getName() == s.declarer.getName() {
			s.declarerVoidSuit[getSuit(s.trump, s.trick[0])] = true
		}
		if p.getName() == s.opp1.getName() {
			s.opp1VoidSuit[getSuit(s.trump, s.trick[0])] = true
		}
		if p.getName() == s.opp2.getName() {
			s.opp2VoidSuit[getSuit(s.trump, s.trick[0])] = true
		}
	}

	if p.getName() != s.declarer.getName() {
		if s.follow == s.trump {
			if getSuit(s.trump, card) == s.trump && (card.Rank == "A" || card.Rank == "10") {
				if opponentIsLosingTrick(s, p, card, false) {
					// TODO:
					debugTacticsLog("INFERENCE: **************************************\n")
					debugTacticsLog("INFERENCE: Playing a full Trump on a losing trick\n")
					debugTacticsLog("INFERENCE: Is the last of the player             \n")
					debugTacticsLog("INFERENCE: **************************************\n")

					if p.getName() == s.opp1.getName() {
						s.opp1VoidSuit[s.trump] = true
					}
					if p.getName() == s.opp2.getName() {
						s.opp2VoidSuit[s.trump] = true
					}
				}
			}
		}
	}

	if len(s.trick) == 2 {
		condition1 := s.trick[0].Rank == "A" && s.greater(s.trick[0], s.trick[1]) && players[0].getName() == partner(s, p).getName()
		condition2 := s.trick[1].equals(Card{s.follow, "A"}) && s.greater(s.trick[1], s.trick[0]) &&  players[1].getName() == partner(s, p).getName()
		if condition1 || condition2 {
			if getSuit(s.trump, card) == s.follow && cardValue(card) < 10 {
				voidCards := []Card{}
				voidCards = append(voidCards, Card{s.follow, "10"})
				if cardValue(card) < 4 {
					voidCards = append(voidCards, Card{s.follow, "K"})
				}
				if cardValue(card) < 3 {
					voidCards = append(voidCards, Card{s.follow, "D"})
				}
				debugTacticsLog("INFERENCE: **************************************\n")
				debugTacticsLog("INFERENCE: Player does not have the cards %v           \n", voidCards)
				debugTacticsLog("INFERENCE: **************************************\n")

				if p.getName() == s.opp1.getName() {
					s.opp1VoidCards = append(s.opp1VoidCards, voidCards...)
				}
				if p.getName() == s.opp2.getName() {
					s.opp2VoidCards = append(s.opp2VoidCards, voidCards...)
				}
			}
		}	

		// TODO
		// only opp1 can know that if he does not have the 10
		// inference from a point of view of the player is needed.
		// TODO
		// test: TestInference_Declarer_A10_at_Declarer DISABLED
		condition0 := getSuit(s.trump, s.trick[0]) != s.trump
		condition1 = s.trick[0].Rank == "K" && getSuit(s.trump, s.trick[1]) == s.follow && players[0].getName() == s.opp1.getName()
		condition2 =  getSuit(s.trump, card) == s.follow 
		condition3 :=  s.greater(s.trick[0], s.trick[1]) && s.greater(s.trick[0], card)
		if condition0 && condition1 && condition2 && condition3 {
			debugTacticsLog("INFERENCE: **************************************\n")
			debugTacticsLog("INFERENCE: Partner and Declarer go under         \n")
			debugTacticsLog("INFERENCE: Partner has 10 and declarer A         \n")
			debugTacticsLog("INFERENCE: **************************************\n")

			// s.opp1VoidCards = append(s.opp1VoidCards, Card{s.follow, "A"})  // This OR
			// s.opp2VoidCards = append(s.opp2VoidCards, Card{s.follow, "A"})  // this
			s.declarerVoidCards = append(s.declarerVoidCards, Card{s.follow, "10"})
		}



			// if !card.equals(Card{s.follow, "10"}) {
			// 	debugTacticsLog("INFERENCE: **************************************\n")
			// 	debugTacticsLog("INFERENCE: Player does not have the %v           \n", Card{s.follow, "10"})
			// 	debugTacticsLog("INFERENCE: **************************************\n")
			// 	if p.getName() == s.opp1.getName() {
			// 		s.opp1VoidCards = append(s.opp1VoidCards, Card{s.follow, "10"})
			// 	}
			// 	if p.getName() == s.opp2.getName() {
			// 		s.opp2VoidCards = append(s.opp2VoidCards, Card{s.follow, "10"})
			// 	}
			// } 
		

		condition1 = s.greater(s.trick[0], s.trick[1]) && players[0].getName() == s.declarer.getName()
		condition2 = s.greater(s.trick[1], s.trick[0]) && players[1].getName() == s.declarer.getName()
		if condition1 || condition2 {
			if !s.greater(card, s.trick...) && cardValue(card) > 0 && card.Rank != "J" {
				debugTacticsLog("INFERENCE: **************************************\n")
				debugTacticsLog("INFERENCE: Player does not have lower than %v    \n", card)
				debugTacticsLog("INFERENCE: **************************************\n")
				i := 0
				r := ""
				for i, r = range ranks {
					if card.Rank == r {
						break
					}
				}
				for j := i + 1 ; j < len(ranks); j++ {
					c := Card{card.Suit, ranks[j]}
					if !in(s.cardsPlayed, c) {
						if p.getName() == s.opp1.getName() {
							s.opp1VoidCards = append(s.opp1VoidCards, c)
						}
						if p.getName() == s.opp2.getName() {
							s.opp2VoidCards = append(s.opp2VoidCards, c)
						}					
					}
				}
			}
		}
	}


	if len(s.trick) > 0 {
		if p.getName() == s.declarer.getName() && card.equals(Card{s.follow, "10"}) && !s.greater(card, s.trick...) {
			debugTacticsLog("INFERENCE: **************************************\n")
			debugTacticsLog("INFERENCE: Declarer void on suit %v              \n", s.follow)
			debugTacticsLog("INFERENCE: **************************************\n")
			s.declarerVoidSuit[s.follow] = true
		}
	}
}

func opponentIsLosingTrick(s *SuitState, p PlayerI, card Card, lookSkat bool) bool {
	dIndex := 0
	for dIndex, _ = range players {
		if players[dIndex].getName() == s.declarer.getName() {
			break
		}
	}
	for i, c := range s.trick {
		if s.greater(c, card) && players[i].getName() == s.declarer.getName() {
			if len(s.trick) == 1 && noHigherCard(s, false, p.getHand(), c) {
				return true
			}
			if len(s.trick) == 2 {
				other := 0
				if dIndex == 0 {
					other = 1
				}
				if s.greater(s.trick[dIndex], s.trick[other]) {
					return true
				}
			}
		}
	}
	return false 
}

