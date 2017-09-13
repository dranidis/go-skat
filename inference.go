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
}
