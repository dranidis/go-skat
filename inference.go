package main


func analysePlay(s *SuitState, p PlayerI, card Card) {
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
				if isLosingTrick(s, p, card) {
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
			if !card.equals(Card{s.follow, "10"}) {
				debugTacticsLog("INFERENCE: **************************************\n")
				debugTacticsLog("INFERENCE: Player does not have the %v           \n", Card{s.follow, "10"})
				debugTacticsLog("INFERENCE: **************************************\n")
				if p.getName() == s.opp1.getName() {
					s.opp1VoidCards = append(s.opp1VoidCards, Card{s.follow, "10"})
				}
				if p.getName() == s.opp2.getName() {
					s.opp2VoidCards = append(s.opp2VoidCards, Card{s.follow, "10"})
				}
			} 
		}

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

func isLosingTrick(s *SuitState, p PlayerI, card Card) bool {
	for i, c := range s.trick {
		if s.greater(c, card) && players[i].getName() == s.declarer.getName() {
			if len(s.trick) == 1 && noHigherCard(s, false, p.getHand(), c) {
				return true
			}
			if len(s.trick) == 2 {
				return true
			}
		}
	}
	return false 
}

