package main 

import (
	"fmt"
	"github.com/dranidis/go-skat/game"
)

// State representation:
// Players are represented with the numbers: 0 1 2
// 0 is the player who initiated the minimax algorithm
// If 0 is the declarer then 	1 is opp1 	and 	2 is opp2
// If 0 is the opp1 	then 	1 is opp2 	and 	2 is declarer
// If 0 is the opp2 	then 	1 is declarer and 	2 is opp1
// 
// In the state we store the hands of the players: 0 1 2
// Who is the declarer: 0, 1 or 2
// And whose turn it is: 0, 1, or 2
// 
type SkatState struct {
	trump         string
	playerHand    [][]Card // YOU, 1, 2
	trick         []Card
	declarer      int
	turn          int // 0 you, 1 player 1, 2 player2
	declScore     int
	oppScore      int
	schneiderGoal bool
}

type SkatAction struct {
	card Card
}


func (s SkatAction) String() string {
	return fmt.Sprintf("%s", s.card)
}

func (s SkatState) String() string {
	return fmt.Sprintf("%s, %v, TRICK: %v, %d, %d, SCORE: %d - %d, %v", s.trump, s.playerHand, s.trick, s.declarer, s.turn, s.declScore, s.oppScore, s.schneiderGoal)
}

func (m SkatState) Heuristic() float64 {
	if m.IsTerminal() {
		return m.FindRewardNum()
	} else {
		return m.playToTheEndWithTactics() 
	}
}

func (m *SkatState) playToTheEndWithTactics() float64 {

	var copyplayers []PlayerI
	// copy the current players in the copyplayers array in order to restore them after play
	// TODO
	copyplayers = []PlayerI{
		players[0].clone(),
		players[1].clone(),
		players[2].clone(),
	}

	// debugTacticsLog("Current state %v\n", m)


	s := makeSuitState();
	s.trump = m.trump

	if len(m.trick) == 3 {
		s.trick = []Card{}
	} else {
		s.trick = make([]Card, len(m.trick))
		copy(s.trick, m.trick)
	}
	
	p0 := makePlayer(m.playerHand[0]) 
	p1 := makePlayer(m.playerHand[1]) 
	p2 := makePlayer(m.playerHand[2]) 

	// can be refactored with an array [p0, p1, p2] and rotation m.declarer times.
	players = []PlayerI{&p0, &p1, &p2}

	if m.declarer == 0 {
		s.declarer = &p0
		s.opp1 = &p1
		s.opp2 = &p2
	}
	if m.declarer == 1 {
		s.declarer = &p1
		s.opp1 = &p2
		s.opp2 = &p0
	}
	if m.declarer == 2 {
		s.declarer = &p2
		s.opp1 = &p0
		s.opp2 = &p1
	}

	rotateTimes := m.turn - len(s.trick)
	if rotateTimes < 0 {
		rotateTimes += 3
	}
	for i := 0; i < rotateTimes; i++ {
		players = rotatePlayers(players)
	}
	// we have reached the turn order in the current trick

	// debugTacticsLog("Players %v, trick %v\n", players, s.trick)
	// playerNow := players[len(s.trick)]
	s.leader = players[0]



	f1 := debugTacticsLogFlag
	f2 := gameLogFlag
	f3 := fileLogFlag
	debugTacticsLogFlag = false
	gameLogFlag = false
	fileLogFlag = false

	for len(players[2].getHand()) > 0 {
		_, players = moveOne(&s, players)
	}


	// minimax chooses cards with high value at the beginning 
	// in order to maximize the gain.
	// For example it might prefer to play an A over a low trump (if both win) 
	// in the first trick since it 
	// adds 11 to the score. Instead if it plays a low trump 0 is added.
	//
	// create a test:
	// trick = DA
	// hand = J J K D 8 7     8     A 9 7
	// MINIMAX plays highest Jack. 7 could win and keep the Jack

// Another issue:
//	
// Trick: []
// (xskat) HAND [] valid: []
// 	Previous suit: xskat:SPADE, xskat:2:CARO
// (xskat) PLAY ANALYSIS: [], xskat, Card: 8
// PLAY ANALYSIS: OPP
// Trick: [8]
// (xskat:2) HAND [] valid: []
// 	Previous suit: xskat:CARO, xskat:2:CARO
// (xskat:2) PLAY ANALYSIS: [8], xskat:2, Card: 9
// PLAY ANALYSIS: OPP
// Trick: [8 9] both CARO
// (goskat) HAND [10 CLubs D 9 Heart ] valid: [10 D 9]
// (goskat) MM: ALL REMAINING CARDS (4): [K D 9 K]
// MM: ..Opponent 2 is VOID in CLUBS
// MM: REMAINING after void: 4 cards: [K D 9 SPADE K HEART ] [] []
// MM: (goskat) 6 Worlds
// MM: MinMaxPlayer
// MM: MINMAX: cards xskat: [K D], xskat:2: [9 K], SKAT:[10 10]
// MM: Decl: 0
// MM: Suggesting card: D with value 68.0000
// MM: MinMaxPlayer
// MM: MINMAX: cards xskat: [K 9], xskat:2: [D K], SKAT:[10 10]
// MM: Decl: 0
// MM: Suggesting card: D with value 68.0000
// MM: MinMaxPlayer
// MM: MINMAX: cards xskat: [K K], xskat:2: [D 9], SKAT:[10 10]
// MM: Decl: 0
// MM: Suggesting card: D with value 72.0000
// MM: MinMaxPlayer
// MM: MINMAX: cards xskat: [D 9], xskat:2: [K K], SKAT:[10 10]
// MM: Decl: 0
// MM: Suggesting card: D with value 69.0000
// MM: ..PRELIMINARY END!
// MM: Time: 2.35922ms
// MM: (goskat) Hand: [10 D 9], Playing card: D (at least 4 times) 

	// WHY GIVE 3 points away??
	// Although it does not change the outcome of the game
	// D would be lost anyway in next trick (taken by K)

// 	goskat: void:SPADE 
// 	xskat: void:
// 	xskat:2: void:CLUBS 
// DECLARER ..sure winners: [10]BACKHAND TRUMP OR No cards of suit played...ZERO valued trick. DO not trump!...LOSING Cards (last)[D 9]
// MM: 
// TACTICS suggest: 9



	score := s.declarer.getScore() + m.declScore
	debugTacticsLog("FINAL score: %d\n", score)
	// var skatStateP game.State
	// skatStateP = m
	// for !skatStateP.IsTerminal() {
	// // for i := 0; i < 1; i++ {
	// 	var a game.Action
	// 	skatAction := m.playWithTactics()
	// 	a = skatAction
	// 	debugTacticsLog("Action %v\n", skatAction)
	// 	skatStateP = skatStateP.FindNextState(a)	
	// 	debugTacticsLog("State %v\n", skatStateP)
	// }

	// restore game players using copyplayers
	players = []PlayerI{
		copyplayers[0],
		copyplayers[1],
		copyplayers[2],
	}

	debugTacticsLogFlag = f1
	gameLogFlag = f2
	fileLogFlag = f3


	return float64(score)

	// return m.FindRewardNum()
}

// func (m SkatState) playWithTactics() SkatAction {
// 	s := makeSuitState();
// 	s.trump = m.trump
	
// 	if len(m.trick) == 3 {
// 		s.trick = []Card{}
// 	} else {
// 		s.trick = make([]Card, len(m.trick))
// 		copy(s.trick, m.trick)
// 	}
	
// 	p0 := makePlayer(m.playerHand[0]) 
// 	p1 := makePlayer(m.playerHand[1]) 
// 	p2 := makePlayer(m.playerHand[2]) 

// 	// can be refactored with an array [p0, p1, p2] and rotation m.declarer times.
// 	players = []PlayerI{&p0, &p1, &p2}

// 	if m.declarer == 0 {
// 		s.declarer = &p0
// 		s.opp1 = &p1
// 		s.opp2 = &p2
// 	}
// 	if m.declarer == 1 {
// 		s.declarer = &p1
// 		s.opp1 = &p2
// 		s.opp2 = &p0
// 	}
// 	if m.declarer == 2 {
// 		s.declarer = &p2
// 		s.opp1 = &p0
// 		s.opp2 = &p1
// 	}

// 	rotateTimes := m.turn - len(s.trick)
// 	if rotateTimes < 0 {
// 		rotateTimes += 3
// 	}
// 	for i := 0; i < rotateTimes; i++ {
// 		players = rotatePlayers(players)
// 	}
// 	// we have reached the turn order in the current trick

// 	debugTacticsLog("Players %v, trick %v\n", players, s.trick)
// 	playerNow := players[len(s.trick)]

// 	card := play(&s, playerNow)

// 	return SkatAction{card}
// }
	

// TODO:
// Replay all played tricks so that they get analysed
// to be used by tactics.	
func (m SkatState) GetTacticsMove() game.Action {

// COPY PASTE FROM ABOVE
	// NEEDS REFACTORING

	// var copyplayers []PlayerI
	// // copy the current players in the copyplayers array in order to restore them after play
	// // TODO
	// copyplayers = []PlayerI{
	// 	players[0].clone(),
	// 	players[1].clone(),
	// 	players[2].clone(),
	// }



	s := makeSuitState();
	s.trump = m.trump

	if len(m.trick) == 3 {
		s.trick = []Card{}
	} else {
		s.trick = make([]Card, len(m.trick))
		copy(s.trick, m.trick)
	}
	
	p0 := makePlayer(m.playerHand[0]) 
	p1 := makePlayer(m.playerHand[1]) 
	p2 := makePlayer(m.playerHand[2]) 

	// can be refactored with an array [p0, p1, p2] and rotation m.declarer times.
	players = []PlayerI{&p0, &p1, &p2}

	if m.declarer == 0 {
		s.declarer = &p0
		s.opp1 = &p1
		s.opp2 = &p2
	}
	if m.declarer == 1 {
		s.declarer = &p1
		s.opp1 = &p2
		s.opp2 = &p0
	}
	if m.declarer == 2 {
		s.declarer = &p2
		s.opp1 = &p0
		s.opp2 = &p1
	}

	rotateTimes := m.turn - len(s.trick)
	if rotateTimes < 0 {
		rotateTimes += 3
	}
	for i := 0; i < rotateTimes; i++ {
		players = rotatePlayers(players)
	}
	// we have reached the turn order in the current trick

	// debugTacticsLog("Players %v, trick %v\n", players, s.trick)
	// playerNow := players[len(s.trick)]
	s.leader = players[0]



	f1 := debugTacticsLogFlag
	f2 := gameLogFlag
	f3 := fileLogFlag
	debugTacticsLogFlag = false
	gameLogFlag = false
	fileLogFlag = false
// END OF COPY PASTE

	card, _ := moveOne(&s, players)

// COPY PASTE
	// players = []PlayerI{
	// 	copyplayers[0],
	// 	copyplayers[1],
	// 	copyplayers[2],
	// }

	debugTacticsLogFlag = f1
	gameLogFlag = f2
	fileLogFlag = f3

// END OF COPY PASTE

	return SkatAction{card}
}

func moveOne(s *SuitState, players []PlayerI) (Card, []PlayerI) {
	l := len(s.trick)
	var card Card
	if l == 0 {
		card = play(s, players[0])
		s.follow = getSuit(s.trump, s.trick[0])
	} 
	if l == 1 { 					// USING else if bevause play changes the s.trick
		debugTacticsLog("moveOne: %v\n", s.trick)
		card = play(s, players[1])
	} 
	if l == 2 {
		debugTacticsLog("moveOne: %v\n", s.trick)
		card = play(s, players[2])
		players = setNextTrickOrder(s, players)
		s.follow = ""
	}
	return card, players
}


// func (m SkatState) IsOpponentTurn() bool {
// 	return m.declarer != m.turn
// }

func (m SkatState) IsOpponentTurn() bool {
 	if m.declarer == 0 && m.turn == 0 {
 		return false
 	}
 	if m.declarer == 0 && m.turn != 0 {
 		return true
 	}
 	if m.declarer == m.turn {
 		return true
 	}
 	return false
 }

func (m *SkatState) IsTerminal() bool {
	if !m.schneiderGoal {
		if m.declScore > 60 {
			return true
		}
		if m.oppScore > 59 {
			return true
		}
	}
	return len(m.playerHand[0])+len(m.playerHand[1])+len(m.playerHand[2]) == 0
}

func (m *SkatState) FindRewardNum() float64 {
	// winsScore := 61
	// if m.schneiderGoal {
	// 	winsScore = 90
	// }

	// if m.declarer == 0 { //YOU
	// 	if m.declScore >= winsScore {
	// 		return float64(1.0)
	// 	} else {
	// 		return float64(0.0)
	// 	}
	// }
	// if m.declScore > winsScore - 1 {
	// 	return float64(0.0)
	// }
	// return float64(m.declScore) // TODO
	return float64(m.declScore - m.oppScore)  // TODO
}

func (m *SkatState) FindReward() float64 {
	winsScore := 61
	if m.schneiderGoal {
		winsScore = 90
	}

	if m.declarer == 0 { //YOU
		if m.declScore >= winsScore {
			return float64(1.0)
		} else {
			return float64(0.0)
		}
	}
	if m.declScore > winsScore-1 {
		return float64(0.0)
	}
	return float64(1.0) // TODO
}


func (m SkatState) validCards(cards []Card) []Card {
	if len(m.trick) == 0 {
		return cards
	}
	return filter(cards, func(c Card) bool {
		return valid(getSuit(m.trump, m.trick[0]), m.trump, cards, c)
	})
}

func (m *SkatState) FindLegals() []game.Action {
	actions := []game.Action{}
	for _, card := range m.validCards(m.playerHand[m.turn]) {
		actions = append(actions, SkatAction{card})
	}
	return actions
}

func (m *SkatState) FindNextState(a game.Action) game.State {
	ma := a.(SkatAction)

	// deep copy before you make any changes
	newState := copySkatState(*m)

	if len(newState.trick) == 3 {
		newState.trick = []Card{}
	}
	// remove the card from the player and add it to the trick
	newState.playerHand[m.turn] = remove(newState.playerHand[m.turn], ma.card)
	newState.trick = append(newState.trick, ma.card)

	strump := newState.trump
	sfollow := getSuit(strump, newState.trick[0])

	if len(newState.trick) == 3 {
		winnerCard := -1
		// find winner
		if greater(strump, sfollow, newState.trick[0], newState.trick[1], newState.trick[2]) {
			winnerCard = 0
			newState.turn = m.turn + 1
		} else if greater(strump, sfollow, newState.trick[1], newState.trick[0], newState.trick[2]) {
			winnerCard = 1
			newState.turn = m.turn + 2
		} else {
			winnerCard = 2
			newState.turn = m.turn
		}
		if newState.turn > 2 {
			newState.turn -= 3
		}

		// set the scores, depending on who played when
		if m.turn == 0 { // YOU PLAYED THE 3RD CARD
			newState.setScores(2, winnerCard)
		}
		if m.turn == 1 { // YOU PLAYED THE 2nd CARD
			newState.setScores(1, winnerCard)
		}
		if m.turn == 2 { // YOU PLAYED THE 1st CARD
			newState.setScores(0, winnerCard)
		}

	} else {
		// set next player turn
		newState.turn++
		if newState.turn > 2 {
			newState.turn = 0
		}
	}

	var state game.State
	state = &newState
	return state
}

func (m *SkatState) setScores(w int, winner int) {
	if winner == w {
		// you won
		if m.declarer == 0 {
			m.declScore += sum(m.trick)
		} else {
			m.oppScore += sum(m.trick)
		}
	} else {
		// you lost
		if m.declarer == 0 {
			m.oppScore += sum(m.trick)
		} else {
			m.declScore += sum(m.trick)
		}
	}
}

func copySkatState(m SkatState) SkatState {
	strick := make([]Card, len(m.trick))
	copy(strick, m.trick)

	p0hand := make([]Card, len(m.playerHand[0]))
	copy(p0hand, m.playerHand[0])
	p1hand := make([]Card, len(m.playerHand[1]))
	copy(p1hand, m.playerHand[1])
	p2hand := make([]Card, len(m.playerHand[2]))
	copy(p2hand, m.playerHand[2])

	return SkatState{
		m.trump,
		[][]Card{
			p0hand,
			p1hand,
			p2hand,
		},
		strick,
		m.declarer,
		m.turn,
		m.declScore,
		m.oppScore,
		m.schneiderGoal,
	}
}
