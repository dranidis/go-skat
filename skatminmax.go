package main 

import (
	"fmt"
	"log"
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
	SuitState
	players    []PlayerI // YOU, 1, 2
}

type SkatAction struct {
	card Card
}

func (s SkatAction) String() string {
	return fmt.Sprintf("%s", s.card)
}

func (s SkatState) String() string {
	return fmt.Sprintf("%v, %s, %v %T, TRICK: %v, SCORE: %d - %d", 
		s.SuitState, s.trump, s.players, s.players[0], s.trick, 
		s.declarer.getScore(), s.opp1.getScore() + s.opp2.getScore())
}

func (m SkatState) Heuristic() float64 {
	if m.IsTerminal() {
		return m.FindRewardNum()
	} else {
		m.playToTheEndWithTactics()
		return m.FindRewardNum() 
	}
}

// type PlayerData struct {
// 	name         string
// 	hand         []Card
// 	highestBid   int
// 	score        int
// 	schwarz      bool
// 	totalScore   int
// 	previousSuit string
// 	won          int
// 	lost         int
// 	defWon       int
// 	declaredBid  int
// }
func copyPlayer(p PlayerI) Player {
	player := makePlayer(p.getHand())
	player.name = p.getName() 
	player.score = p.getScore() 	
	player.previousSuit = p.getPreviousSuit() 	
	player.schwarz = p.isSchwarz() 	
	// player.highestBid = p.getHighestBid() 	
	player.declaredBid = p.getDeclaredBid() 	
	return player
}

func copyPlayerMM(p PlayerI) MinMaxPlayer {
	player := makeMinMaxPlayer(p.getHand())
	player.name = p.getName() 
	player.score = p.getScore() 	
	player.previousSuit = p.getPreviousSuit() 	
	player.schwarz = p.isSchwarz() 	
	// player.highestBid = p.getHighestBid() 	
	player.declaredBid = p.getDeclaredBid() 	
	debugTacticsLog("copyPlayerMM: %v %v\n", p, player)
	return player
}

func (m *SkatState) playToTheEndWithTactics() {
	cpuPlayers := make([]Player, 3)
	for i := 0; i < 3; i++ {
		cpuPlayers[i] = copyPlayer(m.players[i])
		// cpuPlayers[i] = makePlayer(m.players[i].getHand())
		// cpuPlayers[i].name = m.players[i].getName() 
		// cpuPlayers[i].score = m.players[i].getScore() 
		// cpuPlayers[i].name += "-TAC" 
		if m.declarer.getName() == m.players[i].getName() {
			m.declarer = &cpuPlayers[i]
		}
		if m.opp1.getName() == m.players[i].getName() {
			m.opp1 = &cpuPlayers[i]
		}
		if m.opp2.getName() == m.players[i].getName() {
			m.opp2 = &cpuPlayers[i]
		}
		if m.leader.getName() == m.players[i].getName() {
			m.leader = &cpuPlayers[i]
		}
		// debugTacticsLog("CPU Player %v - %v\n", cpuPlayers[i], m.players[i])
		m.players[i] = &cpuPlayers[i]
		// debugTacticsLog("CPU Player %v - %v\n", cpuPlayers[i], m.players[i])
	}

	f1 := debugTacticsLogFlag
	f2 := gameLogFlag
	f3 := fileLogFlag
	debugTacticsLogFlag = false
	gameLogFlag = false
	fileLogFlag = false

	for len(m.players[2].getHand()) > 0 {
		_ = m.moveOne()
	}

	debugTacticsLogFlag = f1
	gameLogFlag = f2
	fileLogFlag = f3
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
	
// func (m SkatState) getGameSuitStateAndPlayers() (*SuitState, []PlayerI) {
// 	s := m.SuitState
// 	// s := makeSuitState();
// 	// s.trump = m.trump

// 	// if len(m.trick) == 3 {
// 	// 	s.trick = []Card{}
// 	// } else {
// 	// 	s.trick = make([]Card, len(m.trick))
// 	// 	copy(s.trick, m.trick)
// 	// }
	
// 	p0 := makePlayer(m.playerHand[0]) 
// 	p1 := makePlayer(m.playerHand[1]) 
// 	p2 := makePlayer(m.playerHand[2]) 

// 	// can be refactored with an array [p0, p1, p2] and rotation m.declarer times.
// 	players = []PlayerI{&p0, &p1, &p2}

// 	// if m.declarer == 0 {
// 	// 	s.declarer = &p0
// 	// 	s.opp1 = &p1
// 	// 	s.opp2 = &p2
// 	// }
// 	// if m.declarer == 1 {
// 	// 	s.declarer = &p1
// 	// 	s.opp1 = &p2
// 	// 	s.opp2 = &p0
// 	// }
// 	// if m.declarer == 2 {
// 	// 	s.declarer = &p2
// 	// 	s.opp1 = &p0
// 	// 	s.opp2 = &p1
// 	// }

// 	rotateTimes := m.turn - len(s.trick)
// 	if rotateTimes < 0 {
// 		rotateTimes += 3
// 	}
// 	for i := 0; i < rotateTimes; i++ {
// 		players = rotatePlayers(players)
// 	}
// 	// we have reached the turn order in the current trick

// 	// debugTacticsLog("Players %v, trick %v\n", players, s.trick)
// 	// playerNow := players[len(s.trick)]
// 	s.leader = players[0]

// 	return &s, players

// }


// TODO:
// Replay all played tricks so that they get analysed
// to be used by tactics.	
func (m SkatState) GetTacticsMove() game.Action {
	tmpState := m.copySkatState()
	cpuPlayers := make([]Player, 3)
	for i := 0; i < 3; i++ {
		cpuPlayers[i] = makePlayer(m.players[i].getHand())
		cpuPlayers[i].name = m.players[i].getName() 
		// cpuPlayers[i].name += "-TAC1" 
		if m.declarer.getName() == m.players[i].getName() {
			tmpState.declarer = &cpuPlayers[i]
		}
		if m.opp1.getName() == m.players[i].getName() {
			tmpState.opp1 = &cpuPlayers[i]
		}
		if m.opp2.getName() == m.players[i].getName() {
			tmpState.opp2 = &cpuPlayers[i]
		}
		if m.leader.getName() == m.players[i].getName() {
			tmpState.leader = &cpuPlayers[i]
		}
	}

	f1 := debugTacticsLogFlag
	f2 := gameLogFlag
	f3 := fileLogFlag
	debugTacticsLogFlag = false
	gameLogFlag = false
	fileLogFlag = false

	card := tmpState.moveOne()

	debugTacticsLogFlag = f1
	gameLogFlag = f2
	fileLogFlag = f3

	return SkatAction{card}
}

func (m *SkatState) moveOne() Card {
	// debugTacticsLog("MOVEONE %v\n", m)
	l := len(m.trick)
	var card Card
	if l == 0 {
		card = play(&m.SuitState, m.players[0])
		m.follow = getSuit(m.trump, m.trick[0])
	} 
	if l == 1 { 					// USING else if bevause play changes the s.trick
		card = play(&m.SuitState, m.players[1])
	} 
	if l == 2 {
		card = play(&m.SuitState, m.players[2])
		// var players = []PlayerI{&m.players[0], &m.players[1], &m.players[2]}
		m.players = setNextTrickOrder(&m.SuitState, m.players)
		// for i := 0; i < 3; i++ {
		// 	var p = players[i].(*MinMaxPlayer)
		// 	m.players[i] = *p
		// }
		m.follow = ""
	}
	// debugTacticsLog("MOVEONE card %v\n", card)

	return card
}


func (m SkatState) IsOpponentTurn() bool {
	return m.players[len(m.trick)].getName() != m.declarer.getName()
}

func (m *SkatState) IsTerminal() bool {
	for i := 0; i < 3; i++ {
		if len(m.players[i].getHand()) > 0 {
			// debugTacticsLog("Player %v\n", m.players[i])
			return false
		}
	}
	return true
}

func (m *SkatState) FindRewardNum() float64 {
	return float64(m.declarer.getScore() - m.opp1.getScore() - m.opp2.getScore())  
}

func (m *SkatState) FindReward() float64 {
	log.Fatal("Not used.")
	return float64(0.0) // TODO
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
	for _, card := range m.validCards(m.players[len(m.trick)].getHand()) {
		actions = append(actions, SkatAction{card})
	}
	return actions
}

func (m *SkatState) FindNextState(a game.Action) game.State {
	// debugTacticsLog("STATE: %v\n", m)
	ma := a.(SkatAction)

	// deep copy before you make any changes
	newState := m.copySkatState()

	currentP := newState.players[len(newState.trick)]
	// p := &currentP 
	analysePlay(&newState.SuitState, currentP, ma.card)

	currentP.setHand(remove(currentP.getHand(), ma.card))
	newState.players[len(newState.trick)] = currentP // STRANGE!!
	// fmt.Printf("Players: %v %v\n", newState.players, currentP)
	// fmt.Printf("END\n")
	// log.Fatal("END")

	// debugTacticsLog("TRICK: %v\n", newState.trick)
	newState.trick = append(newState.trick, ma.card)
	// debugTacticsLog("TRICK: %v\n", newState.trick)
	if getSuit(newState.trump, ma.card) == newState.trump {
		newState.trumpsInGame = remove(newState.trumpsInGame, ma.card)
	}
	newState.cardsPlayed = append(newState.cardsPlayed, ma.card)

	if len(newState.trick) == 1 {
		newState.follow = getSuit(newState.trump, newState.trick[0])
	}

	if len(newState.trick) == 3 {
		// var players = []PlayerI{&newState.players[0], &newState.players[1], &newState.players[2]}
		// debugTacticsLog("Players: %v\n", players)
		newState.players = setNextTrickOrder(&newState.SuitState, newState.players)
		// fmt.Printf("PLAYERS: %v\n", players)
		// fmt.Printf("PLAYER[0]: %v\n", players[0])
		// fmt.Printf("PLAYER[1]: %v\n", players[1])
		// fmt.Printf("PLAYER[2]: %v\n", players[2])


		// // p0 := players[0].(*MinMaxPlayer)
		// // fmt.Printf("p0: %v\n", p0)

		// // newState.players[0] = *p0
		// pp0 := players[0]
		// fmt.Printf("pp0: %v\n", pp0)
		// fmt.Printf("PLAYER[1]: %v\n", players[1])

		// p0 := pp0.(*MinMaxPlayer)
		// fmt.Printf("p0: %v\n", p0)
		// fmt.Printf("PLAYER[1]: %v\n", players[1])
		
		// newState.players[0] = *p0
		// fmt.Printf("newState.players[0]: %v\n", newState.players[0])
		// fmt.Printf("PLAYER[1]: %v\n", players[1])

		// pp1 := players[1]
		// fmt.Printf("pp1: %v\n", pp1)
		// p1 := pp1.(*MinMaxPlayer)
		// fmt.Printf("p1: %v\n", p1)
		// newState.players[1] = *p1
		// fmt.Printf("newState.players[1]: %v\n", newState.players[1])

		// p2 := players[2].(*MinMaxPlayer)
		// fmt.Printf("p0: %v\n", p2)
		// newState.players[2] = *p2

		// for i := 0; i < 3; i++ {
		// 	pp := players[i]
		// 	p := pp.(*MinMaxPlayer)
		// 	fmt.Printf("Player[%d]: %v\n", i, players[i])
		// 	fmt.Printf("Player: %v\n", p)
		// 	newState.players[i] = *p
		// }
		// fmt.Printf("PLAYERS: %v\n", newState.players)
		newState.follow = ""
	}
	var state game.State
	state = &newState
	// debugTacticsLog("NEWSTATE: %v\n", state)
	return state
}

// func (m *SkatState) FindNextState(a game.Action) game.State {
// 	ma := a.(SkatAction)

// 	// deep copy before you make any changes
// 	newState := m.copySkatState()

// 	if len(newState.trick) == 3 {
// 		newState.trick = []Card{}
// 	}
// 	// remove the card from the player and add it to the trick
// 	newState.playerHand[m.turn] = remove(newState.playerHand[m.turn], ma.card)
// 	newState.trick = append(newState.trick, ma.card)

// 	strump := newState.trump
// 	sfollow := getSuit(strump, newState.trick[0])

// 	if len(newState.trick) == 3 {
// 		winnerCard := -1
// 		// find winner
// 		if greater(strump, sfollow, newState.trick[0], newState.trick[1], newState.trick[2]) {
// 			winnerCard = 0
// 			newState.turn = m.turn + 1
// 		} else if greater(strump, sfollow, newState.trick[1], newState.trick[0], newState.trick[2]) {
// 			winnerCard = 1
// 			newState.turn = m.turn + 2
// 		} else {
// 			winnerCard = 2
// 			newState.turn = m.turn
// 		}
// 		if newState.turn > 2 {
// 			newState.turn -= 3
// 		}

// 		// set the scores, depending on who played when
// 		if m.turn == 0 { // YOU PLAYED THE 3RD CARD
// 			newState.setScores(2, winnerCard)
// 		}
// 		if m.turn == 1 { // YOU PLAYED THE 2nd CARD
// 			newState.setScores(1, winnerCard)
// 		}
// 		if m.turn == 2 { // YOU PLAYED THE 1st CARD
// 			newState.setScores(0, winnerCard)
// 		}

// 	} else {
// 		// set next player turn
// 		newState.turn++
// 		if newState.turn > 2 {
// 			newState.turn = 0
// 		}
// 	}

// 	var state game.State
// 	state = &newState
// 	return state
// }

// func (m *SkatState) setScores(w int, winner int) {
// 	if winner == w {
// 		// you won
// 		if m.declarer == 0 {
// 			m.declScore += sum(m.trick)
// 		} else {
// 			m.oppScore += sum(m.trick)
// 		}
// 	} else {
// 		// you lost
// 		if m.declarer == 0 {
// 			m.oppScore += sum(m.trick)
// 		} else {
// 			m.declScore += sum(m.trick)
// 		}
// 	}
// }

func (m SkatState) copySkatState() SkatState {
	suitState := m.SuitState.cloneSuitStateNotPlayers() 
	// copyPlayers := make([]MinMaxPlayer, 3)
	copyIPlayers := make([]PlayerI, 3)
	for i := 0; i < 3; i++ {
		p := m.players[i].clone()

		// var clone = p.(*MinMaxPlayer)
		// copyPlayers[i] = *clone
		copyIPlayers[i] = p


		if m.declarer.getName() == copyIPlayers[i].getName() {
			suitState.declarer = copyIPlayers[i]
		}
		if m.opp1.getName() == copyIPlayers[i].getName() {
			suitState.opp1 = copyIPlayers[i]
		}
		if m.opp2.getName() == copyIPlayers[i].getName() {
			suitState.opp2 = copyIPlayers[i]
		}
		if m.leader.getName() == copyIPlayers[i].getName() {
			suitState.leader = copyIPlayers[i]
		}
	}


	return SkatState{
		suitState,
		copyIPlayers,
	}
}
