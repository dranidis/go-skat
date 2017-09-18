package main

import (
	"encoding/json"
	_ "errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	_ "html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var r = rand.New(rand.NewSource(1))
var _ = rand.New(rand.NewSource(time.Now().Unix()))

var delayMs = 1000
var totalGames = 21
var oficialScoring = false

var verbose = false
/*
channels for html and ISS comm
*/
var playChannel chan CardPlayed
var trickChannel chan Card
var skatChannel chan Card
var winnerChannel chan string
var bidChannel chan string
var trumpChannel chan string
var scoreChannel chan Score
var pickUpChannel chan string
var declareChannel chan string
var discardChannel chan Card
var skatPositionChannel chan int

var gameNr int

var issConnect = false
var issOpp1 = "xskat"
var issOpp2 = "xskat"
var issSentDelay = 0
var issUsername string
var issDEBUG = false

var	minMaxPlayerFlag = false
var maxHandSizeFlag = 4
var MINIMAX_ALG = "abt"
var minmax2Flag = false
var minmax3Flag = false
var	minimaxSearching = false

var gameIndex = 1
var player1 PlayerI
var player2 PlayerI
var player3 PlayerI
var html = false

var players []PlayerI
var gamePlayers []PlayerI

type SuitState struct {
	declarer PlayerI
	opp1     PlayerI
	opp2     PlayerI
	trump    string
	leader   PlayerI
	follow   string
	trick    []Card
	skat             []Card
	// not necessary for game but for tactics
	trumpsInGame     []Card
	cardsPlayed      []Card
	declarerVoidSuit map[string]bool
	opp1VoidSuit  map[string]bool
	opp2VoidSuit  map[string]bool
}

func (s SuitState) String() string {
	voidString := func(m map[string]bool) string {
		str := "("
		for k, v := range m {
			if v {
				str += k
				str += ","
			}
		}
		str += ")"
		return str
	}
	return fmt.Sprintf("D:%s, O1:%s, O2:%s, T:%s, Leads:%s, Fol: %s, TRICK:%v, SKAT:%v, InGame: %v, Played: %v, D_VOID:%s,O1_VOID:%s,O2_VOID:%s\n", 
		s.declarer.getName(), s.opp1.getName(), s.opp2.getName(), s.trump, s.leader.getName(), s.follow, s.trick, 
		s.skat, s.trumpsInGame, s.cardsPlayed, 
		voidString(s.declarerVoidSuit), voidString(s.opp1VoidSuit), voidString(s.opp2VoidSuit))
}

func (s *SuitState) cloneSuitStateNotPlayers() SuitState {
	newSS := makeSuitState()
	newSS.declarer = s.declarer //is it necessary?
	newSS.opp1 = s.opp1 // swallow copy, will be replaced anyway
	newSS.opp2 = s.opp2 // swallow copy, will be replaced anyway 
	newSS.trump = s.trump
	newSS.leader = s.leader // swallow copy, will be replaced anyway 
	newSS.follow = s.follow
	newTrick := make([]Card, len(s.trick))
	copy(newTrick, s.trick)

	newSS.trick = newTrick

	newSS.skat = s.skat // Does not change. No need to clone

	newtrumpsInGame := make([]Card, len(s.trumpsInGame))
	copy(newtrumpsInGame, s.trumpsInGame)
	
	newcardsPlayed := make([]Card, len(s.cardsPlayed))
	copy(newcardsPlayed, s.cardsPlayed)
		
	newSS.declarerVoidSuit = map[string]bool{
			CLUBS: s.declarerVoidSuit[CLUBS],
			SPADE: s.declarerVoidSuit[SPADE],
			HEART: s.declarerVoidSuit[HEART],
			CARO:  s.declarerVoidSuit[CARO],
		}
	newSS.opp1VoidSuit = map[string]bool{
			CLUBS: s.opp1VoidSuit[CLUBS],
			SPADE: s.opp1VoidSuit[SPADE],
			HEART: s.opp1VoidSuit[HEART],
			CARO:  s.opp1VoidSuit[CARO],
		}
	newSS.opp2VoidSuit = map[string]bool{
			CLUBS: s.opp2VoidSuit[CLUBS],
			SPADE: s.opp2VoidSuit[SPADE],
			HEART: s.opp2VoidSuit[HEART],
			CARO:  s.opp2VoidSuit[CARO],
		}

		return newSS
}

func makeSuitState() SuitState {
	skat := make([]Card, 2)
	return SuitState{nil, nil, nil, "", nil, "", []Card{}, skat, []Card{}, []Card{},
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
	}
}

func trickWinner(s *SuitState) int {
	if s.greater(s.trick[0], s.trick[1]) && s.greater(s.trick[0], s.trick[2]) {
		return 0
	} else if s.greater(s.trick[1], s.trick[2]) {
		return 1
	} else {
		return 2
	}
}

func setNextTrickOrder(s *SuitState, players []PlayerI) []PlayerI {
	index := trickWinner(s)
	var winner PlayerI

	debugTacticsLog("PLAYERS: %v %d\n", players, index)
	winner = players[index]
	for i := 0; i < index; i++ {
		players = rotatePlayers(players)
	}

	winner.setScore(winner.getScore() + sum(s.trick))

	if s.declarer != nil && s.opp1 != nil && s.opp2 != nil {
		gameLog("TRICK %v\n", s.trick)
		debugTacticsLog("%d points: %d - %d\n", sum(s.trick), s.declarer.getScore(), s.opp1.getScore()+s.opp2.getScore())
	}
	if html && !minimaxSearching {
		winnerChannel <- winner.getName()
	}

	winner.setSchwarz(false)
	s.trick = []Card{}
	s.leader = players[0]

	if s.trump == NULL {
		if winner.getName() == s.declarer.getName() {
			// declarer lost
			return nil
		}
	}

	return players
}

type CardPlayed struct {
	C  Card
	Pi int
}

func round(s *SuitState, players []PlayerI) []PlayerI {
	card1 := play(s, players[0])
	if html {
		htmlLog("Sending to channel...%v \n", card1)
		playChannel <- CardPlayed{card1, 0}
	}
	time.Sleep(time.Duration(delayMs) * time.Millisecond)
	s.follow = getSuit(s.trump, s.trick[0])
	card2 := play(s, players[1])
	if html {
		htmlLog("Sending to channel...%v \n", card2)

		playChannel <- CardPlayed{card2, 1}
	}
	time.Sleep(time.Duration(delayMs) * time.Millisecond)
	card3 := play(s, players[2])
	if html {
		htmlLog("Sending to channel...%v \n", card3)
		playChannel <- CardPlayed{card3, 2}
	}
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	players = setNextTrickOrder(s, players)
	s.follow = ""
	return players
}

func analysePlay(s *SuitState, p PlayerI, card Card) {
	// debugTacticsLog("PLAY ANALYSIS: %v, %s, Card: %v\n", s.trick, p.getName(), card)
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
		// debugTacticsLog("PLAY ANALYSIS: OPP\n")
		if s.follow == s.trump {
			// debugTacticsLog("PLAY ANALYSIS: TRUMP\n")
			if getSuit(s.trump, card) == s.trump && (card.Rank == "A" || card.Rank == "10") {
				// debugTacticsLog("PLAY ANALYSIS: A/10\n")
				if isLosingTrick(s, p, card) {
					// TODO:
					debugTacticsLog("INFERENCE: **************************************\n")
					debugTacticsLog("INFERENCE: Playing a full Trump on a losing trick\n")
					debugTacticsLog("INFERENCE: Is the last of the player             \n")
					debugTacticsLog("INFERENCE: **************************************\n")

					if p.getName() == s.opp1.getName() {
					//	s.opp1VoidSuit[s.trump] = true
					}
					if p.getName() == s.opp2.getName() {
					//	s.opp2VoidSuit[s.trump] = true
					}
				}
			}
		}
	}
}

func isLosingTrick(s *SuitState, p PlayerI, card Card) bool {
	var c Card
	for _, c = range s.trick {
		if s.greater(c, card) {
			return true
		}
	}
	// if noHigherCard(s, false, []Card{}, c) {
	// 	return true
	// }
	return false // TODO!!!!
}

func play(s *SuitState, p PlayerI) Card {
	red := color.New(color.Bold, color.FgRed).SprintFunc()
	valid := sortSuit(s.trump, validCards(*s, p.getHand()))

	p.setHand(sortSuit(s.trump, p.getHand()))
	gameLog("Trick: %v\n", s.trick)
	pName := p.getName()
	if s.declarer.getName() == p.getName() {
		pName = red(pName)
	}
	debugTacticsLog("(%v) HAND %v ", pName, p.getHand())
	debugTacticsLog("valid: %v\n", valid)
	if s.declarer.getName() != p.getName() {
		debugTacticsLog("\tPrevious suit: %s:%v, %s:%v\n",
			s.opp1.getName(), s.opp1.getPreviousSuit(),
			s.opp2.getName(), s.opp2.getPreviousSuit())
	}
	if s.declarer.getName() == p.getName() {
		gameLog("(%s) ", red(p.getName()))
	} else {
		gameLog("(%v) ", p.getName())
	}

	card := p.playerTactic(s, valid)

	if issConnect && p.getName() == issUsername {
		playCard(card)
	}	

	analysePlay(s, p, card)

	p.setHand(remove(p.getHand(), card))
	s.trick = append(s.trick, card)
	if getSuit(s.trump, card) == s.trump {
		s.trumpsInGame = remove(s.trumpsInGame, card)
	}
	s.cardsPlayed = append(s.cardsPlayed, card)
	return card
}

// Returns a list of all cards that are playeable from the player's hand.
func validCards(s SuitState, playerHand []Card) []Card {
	return filter(playerHand, func(c Card) bool {
		return valid(s.follow, s.trump, playerHand, c)
	})
}

func valid(follow, trump string, playerHand []Card, card Card) bool {
	// debugTacticsLog(".. PLAYER HAND: %v\n", playerHand)
	for _, c := range playerHand {
		// if there is at least one card in your hand matching the followed suit
		// your played card should follow
		if follow == getSuit(trump, c) {
			// if s.follow != getSuit(s.trump, card) {
			// 	debugTacticsLog(".. INVALID CARD: %v. Valid in hand: %v\n", card, c)
			// }
			return follow == getSuit(trump, card)
		}
	}
	// otherwise any card is playable
	return true
}

var bids = []int{
	18, 20, 22, 23, 24,
	27, 30, 33, 35, 36,
	40, 44, 45, 46, 48, 50,
	54, 55, 59, 60,
	63, 66, 70, 72, 77,
	80, 81, 84, 88, 90, 96, 99, 100, 108, 110, 117,
	121, 126, 130, 132, 135, 140, 143, 144,
}

func bidding(listener, speaker PlayerI, bidIndex int) (int, PlayerI) {
	for speaker.accepts(bidIndex, false) {
		bidLog("\t(%v) %d\n", speaker.getName(), bids[bidIndex])
		if listener.accepts(bidIndex, true) {
			//bidLog("Yes %d\n", bids[bidIndex])
			bidLog("\t(%v) Yes\n", listener.getName())
			bidIndex++
		} else {
			//	bidLog("Listener (%v) Pass %d\n", listener.getName(), bids[bidIndex])
			bidLog("\t(%v) Pass\n", listener.getName())
			return bidIndex, speaker
		}
	}
	//bidLog("(%v) Pass %d\n", speaker.getName(), bids[bidIndex])
	bidLog("\t(%v) Pass \n", speaker.getName())
	bidIndex--
	return bidIndex, listener
}

func bid(players []PlayerI) (int, PlayerI) {
	bidLog("(%v) vs (%v)\n", players[0].getName(), players[1].getName())
	bidIndex, p := bidding(players[0], players[1], 0)
	bidIndex++
	bidLog("(%v) vs (%v)\n", p.getName(), players[2].getName())
	bidIndex, p = bidding(p, players[2], bidIndex)
	if bidIndex == -1 {
		if players[0].accepts(0, false) { // speaks at end
			bidLog("\t(%s) Yes %d\n", players[0].getName(), bids[0])
			return bids[0], players[0]
		} else {
			bidLog("\t(%s) Pass\n", players[0].getName())
			return 0, nil
		}
	}
	//	p.isDeclarer = true
	if bidIndex < 0 {
		return 0, p
	}
	return bids[bidIndex], p
}

func gameScore(state SuitState, cs []Card, handGame bool) Score {
	withMatadors := 99
	schneider := false
	schwarz := false
	overbid := false
	ouvert := false
	gs := 0

	scoringString := ""
	if state.trump == NULL {
		scoringString += fmt.Sprintf("\nGame: %d\tNULL ", gameIndex)
		// gameLog("\nSCORING\n\tNULL ")
		if handGame {
			scoringString += fmt.Sprintf("HAND \n")
			// gameLog("HAND \n")
			if state.declarer.isSchwarz() {
				gs = 35
			} else {
				gs = -70
			}
		} else {
			if state.declarer.isSchwarz() {
				gs = 23
			} else {
				gs = -46
			}
		}
		scoringString += fmt.Sprintf("SCORE %d\t", gs)
		// gameLog("SCORE %d\t", gs)
	} else {
		mat := matadors(state.trump, cs)
		withMatadors = mat
		if mat < 0 {
			mat = mat * -1
		}
		multiplier := mat + 1

		if withMatadors > 0 {
			scoringString += fmt.Sprintf("\nGame: %d\t%s, With %d ", gameIndex, state.trump, mat)
			// gameLog("\nSCORING\n\t%s, With %d ", state.trump, mat)
		} else {
			scoringString += fmt.Sprintf("\nGame: %d\t%s, Without %d ", gameIndex, state.trump, mat)
			// gameLog("\nSCORING\n\t%s, Without %d ", state.trump, mat)
		}

		base := trumpBaseValue(state.trump)

		if handGame {
			multiplier++
			scoringString += fmt.Sprintf("Hand ")
			// gameLog("Hand ")
		}
		// Schneider?
		if state.declarer.getScore() > 89 || state.declarer.getScore() < 31 {
			multiplier++
			schneider = true
			scoringString += fmt.Sprintf("Schneider ")
			// gameLog("Schneider ")
		}

		if state.declarer.isSchwarz() || (state.opp1.isSchwarz() && state.opp2.isSchwarz()) {
			multiplier++
			schwarz = true
			scoringString += fmt.Sprintf("Schwarz ")
			// gameLog("Schwarz ")
		}
		//gameLog("\n\n")
		gs = multiplier * base

		// OVERBID?
		if gs < state.declarer.getDeclaredBid() {
			fmt.Printf(" --OVERBID!!! Game Value: %d < Bid: %d-- ", gs, state.declarer.getDeclaredBid())
			scoringString += fmt.Sprintf(" --OVERBID!!! Game Value: %d < Bid: %d-- ", gs, state.declarer.getDeclaredBid())
			// gameLog(" --OVERBID!!! Game Value: %d < Bid: %d-- ", gs, state.declarer.getDeclaredBid())
			overbid = true
			leastMult := 0
			for leastMult*base < state.declarer.getDeclaredBid() {
				leastMult++
			}
			//score = -2 * leastMult * base
			gs = -2 * leastMult * base
		} else if state.declarer.getScore() > 60 {
			//score = gs
		} else {
			gs = -2 * gs
		}
		scoringString += fmt.Sprintf("SCORE %d\t", gs)
		// gameLog("SCORE %d\t", gs)
	}

	gameLog(scoringString)

	if verbose {
		fmt.Printf(scoringString)
	}

	scoreStruct := Score{
		state.declarer.getName(),
		state.declarer.getScore(),
		state.opp1.getScore() + state.opp2.getScore(),
		gs,
		state.trump,
		withMatadors, 
		handGame,
		schneider, 
		schwarz, 
		ouvert, 
		overbid,
		0,0,0,
	}
	return scoreStruct
}

var grandGames = 0
var nullGames = 0

var state SuitState

var oldCards []Card

func DealCards() {
	debugTacticsLog("Dealing\n")

	var cards []Card 
	for {
		cards = Shuffle(makeDeck())
		gamePlayers[0].setHand(sortSuit("", cards[:10]))
		if dealGame == "" {
			break
		}
		player := gamePlayers[0].(*Player)
		g := player.canWin(false)	
		if g == dealGame {
			break
		}
	}

	oldCards = make([]Card, len(cards))
	copy(oldCards, cards)

	gamePlayers[0].setHand(sortSuit("", cards[:10]))
	gamePlayers[1].setHand(sortSuit("", cards[10:20]))
	gamePlayers[2].setHand(sortSuit("", cards[20:30]))
	copy(state.skat, cards[30:32])
}

func SameCards() {
	gamePlayers[0].setHand(sortSuit("", oldCards[:10]))
	gamePlayers[1].setHand(sortSuit("", oldCards[10:20]))
	gamePlayers[2].setHand(sortSuit("", oldCards[20:30]))
	copy(state.skat, oldCards[30:32])
}

/*

initstate()
------------> skat.go
				creates Suitstate state (GLOBAL)
initGame()
------------> skat.go
				reset all players
				player.calculateHighestBid(false)
				    canWin(false)
				    	p.handGame = false
				    change p.handGame in a HandGame
bidPhase()
------------> skat.go
				bid(players)
declareAndPlay()
------------> skat.go
				state.declarer.pickUpSkat(state.skat)

*/
func initState() {
	state = makeSuitState()
}

func initGame() {
	for _, p := range players {
		p.ResetPlayer()
		h := p.calculateHighestBid(false)
		debugTacticsLog("(%v) hand: %v Bid up to: %d\n", p.getName(), p.getHand(), h)
	}

	gameLog("\nPLAYER ORDER: %s - %s - %s\n\n", players[0].getName(), players[1].getName(), players[2].getName())
}

func bidPhase() int {
	// BIDDING
	bidDecl, declarer := bid(players)
	if bidDecl == 0 {
		gameLog("ALL PASSED\n")
		return 0
	}
	debugTacticsLog("Declarer %v\n", declarer)
	declarer.setDeclaredBid(bidDecl)

	state.setDeclarerAndOpps(players, declarer)

	gameLog("(%s) won the bidding\n", state.declarer.getName())
	return bidDecl
}

func declareAndPlay() int {
	handGame := true

	gameLog("... SKAT\n")
	if state.declarer.pickUpSkat(state.skat) {
		gameLog("... PICKED UP SKAT and DISCARDED\n")
		handGame = false
	}

	gameLog("... DECLARE GAME\n")

	state.trump = state.declarer.declareTrump()

	if issConnect {
		if len(state.skat) == 2 && state.skat[0].Rank != "" {
			if !handGame {
				fmt.Printf("sending trump %s and skat %v %v to server" , state.trump, state.skat[0], state.skat[1])
				iss_declare(state.trump, false, state.skat)
			} 
		}  
	}

	if issConnect && handGame {
		fmt.Printf("sending trump %s and Hand Game to server" , state.trump)
		iss_declare(state.trump, true, state.skat)
	}

	if html {
		htmlLog("Sending trump %v", state.trump)
		trumpChannel <- state.trump
	}

	if state.trump == GRAND {
		grandGames++
	}
	if state.trump == NULL {
		nullGames++
	}
	state.trumpsInGame = filter(makeDeck(), func(c Card) bool {
		return getSuit(state.trump, c) == state.trump
	})

	state.leader = players[0]
	state.follow = ""

	players[0].setHand(sortSuit(state.trump, players[0].getHand()))
	players[1].setHand(sortSuit(state.trump, players[1].getHand()))
	players[2].setHand(sortSuit(state.trump, players[2].getHand()))

	handString := ""
	if handGame {
		handString = "(HAND)"
	}
	gameLog("\n(%s) TRUMP: %s %s\n", red(state.declarer.getName()), state.trump, handString)

	if issConnect && state.declarer.getName() == "goskat" {
		debugTacticsLog("SKAT: %v\n", state.skat)
	}


	declarerCards := make([]Card, len(state.declarer.getHand()))
	copy(declarerCards, state.declarer.getHand())
	declarerCards = append(declarerCards, state.skat...)

	// PLAY
	for i := 0; i < 10; i++ {
		debugTacticsLog("TRUMPS IN PLAY %v\n", sortRank(state.trumpsInGame))
		gameLog("\n")
		players = round(&state, players)
		if players == nil {
			break
		}
	}
	gameLog("\nSKAT: %v\n", state.skat)

	// gameLog("SKAT: %v, %d\n", skat, sum(skat))
	state.declarer.setScore(state.declarer.getScore() + sum(state.skat))

	gameSc := gameScore(state, declarerCards, handGame)

	gs := gameSc.GameScore

	state.declarer.incTotalScore(gs)

	scoreString := ""
	if gs > 0 {
		if oficialScoring {
			state.declarer.incTotalScore(50)
		}
		if state.trump != NULL {
			scoreString = fmt.Sprintf("%s VICTORY: %d - %d, SCORE: %d\n", state.declarer.getName(),
				state.declarer.getScore(), state.opp1.getScore()+state.opp2.getScore(), gs)
			// gameLog("%s VICTORY: %d - %d, SCORE: %d\n", state.declarer.getName(),
			// 	state.declarer.getScore(), state.opp1.getScore()+state.opp2.getScore(), gs)
		} else {
			scoreString = fmt.Sprintf("%s VICTORY: %d\n", state.declarer.getName(), gs)
			// gameLog("%s VICTORY: %d\n", state.declarer.getName(), gs)
		}
	} else {
		if oficialScoring {
			state.declarer.incTotalScore(-50)
			state.opp1.incTotalScore(40)
			state.opp2.incTotalScore(40)
		}
		state.opp1.wonAsDefenders()
		state.opp2.wonAsDefenders()
		if state.trump != NULL {
			scoreString = fmt.Sprintf("%s DEFEAT: %d - %d, SCORE: %d\n", state.declarer.getName(),
				state.declarer.getScore(), state.opp1.getScore()+state.opp2.getScore(), gs)
			// gameLog("%s DEFEAT: %d - %d, SCORE: %d\n", state.declarer.getName(),
			// 	state.declarer.getScore(), state.opp1.getScore()+state.opp2.getScore(), gs)
		} else {
			scoreString = fmt.Sprintf("%s DEFEAT: %d\n", state.declarer.getName(), gs)
			// gameLog("%s DEFEAT: %d\n", state.declarer.getName(), gs)
		}
	}
	gameLog(scoreString)
	if verbose {
		fmt.Printf(scoreString)
	}

	if player1 != nil && player2 != nil && player3 != nil {
		gameSc.Total1 = player1.getTotalScore()
		gameSc.Total2 = player2.getTotalScore()
		gameSc.Total3 = player3.getTotalScore()
	} else {
		debugTacticsLog("NIL players: %v %v %v\n", player1, player2, player3)
	}
	//fmt.Println("%v\n", gameSc)
	if html {
		scoreChannel <- gameSc
	}
	return gs
}

type Score struct {
	Declarer string
	DeclarerPoints int
	DefenderPoints int
	GameScore      int
	GameDeclared string
	With int
	Hand bool
	Schneider bool
	Schwarz bool
	Ouvert bool
	Overbid bool
	Total1 int
	Total2 int
	Total3 int
}

func (s *SuitState) setDeclarerAndOpps(players []PlayerI, declarer PlayerI) {
	state.declarer = declarer

	if declarer == players[0] {
		state.opp1, state.opp2 = players[1], players[2]
	}
	if declarer == players[1] {
		state.opp2, state.opp1 = players[0], players[2]
	}
	if declarer == players[2] {
		state.opp1, state.opp2 = players[0], players[1]
	}
	state.opp1.setPartner(state.opp2)	
}

func skatGame() int {
	gameLog("\n\nGAME %d/%d\n", gameIndex, totalGames)
	initState()
	DealCards()

	players = []PlayerI{gamePlayers[0], gamePlayers[1], gamePlayers[2]}
	initGame()
	if bidPhase() == 0 {
		return 0
	}
	gs := declareAndPlay()
	return gs
}

func repeatGame() int {
	gameLog("\n\nGAME %d/%d\n", gameIndex, totalGames)
	// players = rotatePlayers(players)
	initState()
	SameCards()

	players = []PlayerI{gamePlayers[0], gamePlayers[1], gamePlayers[2]}
	initGame()
	if bidPhase() == 0 {
		return 0
	}
	gs := declareAndPlay()
	return gs
}

func rotatePlayers(players []PlayerI) []PlayerI {
	newPlayers := []PlayerI{}
	newPlayers = append(newPlayers, players[1])
	newPlayers = append(newPlayers, players[2])
	newPlayers = append(newPlayers, players[0])
	return newPlayers
}

func makeChannels() {
	playChannel = make(chan CardPlayed)
	trickChannel = make(chan Card)
	skatChannel = make(chan Card)
	bidChannel = make(chan string)
	winnerChannel = make(chan string)
	trumpChannel = make(chan string)
	scoreChannel = make(chan Score)
	pickUpChannel = make(chan string)
	declareChannel = make(chan string)
	discardChannel = make(chan Card)
	skatPositionChannel = make(chan int)
}

func makePlayers(auto, html, issConnect, analysis bool, analysisPl, analysisPlayerBid int) {
	if analysis {
		fmt.Printf("Creating players for analysis. Player: %d\n", analysisPl)
		if analysisPlayerBid != 0 {
			fmt.Printf(".. Bid:  %d\n", analysisPlayerBid)
		}
		switch analysisPl {
		case 1:
			lplayer1 := makeAPlayer([]Card{})
			lplayer1.forcedBid = analysisPlayerBid
			lplayer2 := makePlayer([]Card{})
			lplayer3 := makePlayer([]Card{})
			gamePlayers = []PlayerI{&lplayer1, &lplayer2, &lplayer3}
		case 2:
			lplayer1 := makePlayer([]Card{})
			lplayer2 := makeAPlayer([]Card{})
			lplayer2.forcedBid = analysisPlayerBid
			lplayer3 := makePlayer([]Card{})
			gamePlayers = []PlayerI{&lplayer1, &lplayer2, &lplayer3}
		case 3:
			lplayer1 := makePlayer([]Card{})
			lplayer2 := makePlayer([]Card{})
			lplayer3 := makeAPlayer([]Card{})
			lplayer3.forcedBid = analysisPlayerBid
			gamePlayers = []PlayerI{&lplayer1, &lplayer2, &lplayer3}
		}
		delayMs = 0

		gamePlayers[0].setName("You")
		player2.setName("Bob")
		gamePlayers[2].setName("Ana")
		fmt.Printf("Analysed player: %s\n", gamePlayers[analysisPl-1].getName())
		return
	}
	if auto {
		debugTacticsLog("Creating CPU players only\n")
		player := makePlayer([]Card{})
		player1 = &player
		delayMs = 0
	} else {
		if html {
			player := makeHtmlPlayer([]Card{})
			player1 = &player
		} else if issConnect {
			var lplayer1 PlayerI
			cpuplayer := makePlayer([]Card{})
			lplayer1 = &cpuplayer

			if minMaxPlayerFlag {
				debugTacticsLog("PLAYERS: Creating a MinMax player for ISS\n")
				mmplayer := makeMinMaxPlayer([]Card{})
				lplayer1 = &mmplayer
			}
			// player1 = &player
			issplayer2 := makeISSPlayer([]Card{})
			issplayer3 := makeISSPlayer([]Card{})
			lplayer1.setName("")
			issplayer2.setName("ISS1") // this will change by ISS
			issplayer3.setName("ISS2") // this will change by ISS
			gamePlayers = []PlayerI{lplayer1, &issplayer2, &issplayer3} // this will change by ISS
			delayMs = 0

			return			
		} else {
			player := makeHumanPlayer([]Card{})
			player1 = &player
		}
		delayMs = 500
	}

	if minmax2Flag {
		p2 := makeMinMaxPlayer([]Card{})
		player2 = &p2
	} else {
		p2 := makePlayer([]Card{})
		player2 = &p2
	}

	if minmax3Flag {
		p3 := makeMinMaxPlayer([]Card{})
		player3 = &p3
	} else {
		p3 := makePlayer([]Card{})
		player3 = &p3
	}

	player1.setName("You")
	player2.setName("Bob")
	player3.setName("Ana")
	gamePlayers = []PlayerI{player1, player2, player3}
}


var dealGame = ""

func main() {
	// COMMAND LINE FLAGS
	auto := false
	analysis := false
	analysisPlayer := 1
	analysisPlayerBid := 0
	winAnalysis := true

	var randSeed int
	flag.IntVar(&gameNr, "g", 1, "Deal cards # times before you start. You can use this option to move to a specific game of a series of games.")
	flag.IntVar(&totalGames, "n", 36, "total number of games, default 36")
	flag.IntVar(&randSeed, "r", 0, "Seed for random number generator. A value of 0 generates a random number to be used as a seed.")
	flag.BoolVar(&auto, "auto", false, "Runs with CPU players only")
	flag.BoolVar(&verbose, "v", false, "More info in auto mode.")
	flag.BoolVar(&analysis, "analysis", false, "Exhaustively tries out all the moves of a player in a repeated game")
	flag.BoolVar(&winAnalysis, "win", true, "Win or Lose target of the analysed player")
	flag.IntVar(&analysisPlayer, "player", 1, "The player whose moves are being analysed")
	flag.IntVar(&analysisPlayerBid, "bid", 0, "Force the bid of the analysed player")
	flag.BoolVar(&fileLogFlag, "log", true, "Saves log in a file")
	flag.BoolVar(&html, "html", false, "Starts an HTTP server at localhost:3000")
	flag.BoolVar(&issConnect, "iss", false, "Connects to ISS skat server")
	flag.BoolVar(&issDEBUG, "issDEBUG", false, "Fakes the responses of the skat server for debugging.")
	flag.BoolVar(&minMaxPlayerFlag, "minmax", false, "Uses a MinMax CPU player a AI player at ISS")
	flag.BoolVar(&minmax2Flag, "minmax2", false, "Uses a MinMax CPU player as 2nd player in auto and html.")
	flag.BoolVar(&minmax3Flag, "minmax3", false, "Uses a MinMax CPU player as 3rd player in auto and html.")
	flag.StringVar(&MINIMAX_ALG, "mmalg", "abt", "Algorithm used by minmax player: ab alphabeta, abt alphabeta with tactics for opponents")
	flag.StringVar(&issOpp1, "opp1", "xskat", "Opponent to play with at ISS skat server")
	flag.StringVar(&issOpp2, "opp2", "xskat", "Opponent to play with at ISS skat server")
	flag.IntVar(&issSentDelay, "issdelay", 0, "Delay (in ms) before sending an action to ISS server. Useful for debugging and for observing a game.")
	flag.IntVar(&maxHandSizeFlag, "mm-max", 4, "Max hand size for the minimax player. Below that normal tactics are used.")
	flag.StringVar(&dealGame, "deal", "", "Force a specific game deal: Grand, Null")
	flag.Parse()


	if auto {
		gameLogFlag = false
	}

	if randSeed == 0 {
		r = rand.New(rand.NewSource(time.Now().Unix()))
		randSeed = r.Intn(9999)
	}

	fmt.Printf("SEED: %d\n", randSeed)
	gameLog("SEED: %d\n", randSeed)
	fmt.Printf("Game: %d\n", gameNr)
	gameLog("Game: %d\n", gameNr)
	r = rand.New(rand.NewSource(int64(randSeed)))

	if fileLogFlag {
		file, err := os.Create(logFileName)
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		logFile = file
		defer file.Close()
	}

	makePlayers(auto, html, issConnect, analysis, analysisPlayer, analysisPlayerBid)

	if issConnect {
		issUsername = os.Getenv("ISS_USR")
		pwd := os.Getenv("ISS_PWD")
		if issUsername == "" || pwd == "" {
			log.Fatal("To connect to the ISS skat server, you have to set the ISS_USR and ISS_PWD environment variables.")
		}

		err := Connect(issUsername, pwd) // blocks
		if err != nil {
			log.Fatal("Error in server connection: ", err)
			return
		}
		return
	}

	rotateTimes := r.Intn(5) + gameNr - 1
	for i := 0; i < rotateTimes; i++ {
		gamePlayers = rotatePlayers(gamePlayers)
	}

	debugTacticsLog("SHUFFLING %d time(s)\n", gameNr)
	for i := 0; i < gameNr - 1; i++ {
		DealCards()
	}

	if html {
		gameLogFlag = false
		rt := startServer()
		port := ":3000"
		fmt.Printf("Starting server at %s\n", port)
		fmt.Printf("Open page :3000/html/\n")
		http.ListenAndServe(port, rt)
		// does not return
	}

	if analysis {
		if !winAnalysis {
			fmt.Println("Target: Declarer should lose.")
		} else {
			fmt.Println("Target: Declarer should win.")
		}		
		gameLogFlag = false
		gamePlayers = rotatePlayers(gamePlayers)
		score := skatGame()
		s := score
		// printScore(gamePlayers)
		i := 0
		condition := func(s int) bool {
			if winAnalysis {
				return s < 0
			} else {
				return s > 0 
			}
		}

		anim := animation()

		previousGameAnalysis = false
		for condition(s) && !analysisEnded {
			nextGameForAnalysis()
			s = repeatGame()		
			anim()	
			i++
			// printScore(gamePlayers)
		}
		if analysisEnded {
			fmt.Printf("No chance! %d repetitions\n", i)
		} else {
			if fileLogFlag {
				// logFile.Close() // close log file
			}
			file, err := os.Create("analysis.txt")
			if err != nil {
				log.Fatal("Cannot create file", err)
			}
			logFile = file
			fileLogFlag = true

			previousGameAnalysis = true
			nextGameForAnalysis()
			s = repeatGame()
			fmt.Printf("Won! %d repetitions\n", i)
		}
		return //exit
	}

	passed := 0
	won := 0
	lost := 0
	anim := animation()
	for ; gameIndex <= totalGames; gameIndex++ {
		gamePlayers = rotatePlayers(gamePlayers)
		score := skatGame()
		if score == 0 {
			passed++
		}
		if score > 0 {
			won++
		} else if score < 0 {
			lost++
		}
		if !auto {
			fmt.Printf("\nGAME: %6d (%s) %5d     (%s) %5d     (%s) %5d\n", gameIndex, gamePlayers[0].getName(), gamePlayers[0].getTotalScore(), gamePlayers[1].getName(), gamePlayers[1].getTotalScore(), gamePlayers[2].getName(), gamePlayers[2].getTotalScore())
		} else {
			gameLog("\nGAME: %6d (%s) %5d     (%s) %5d     (%s) %5d\n", gameIndex, gamePlayers[0].getName(), gamePlayers[0].getTotalScore(), gamePlayers[1].getName(), gamePlayers[1].getTotalScore(), gamePlayers[2].getName(), gamePlayers[2].getTotalScore())
			anim()
		}
	}

	avg := float64(gamePlayers[0].getTotalScore()+gamePlayers[1].getTotalScore()+gamePlayers[2].getTotalScore()) / float64(totalGames-passed)

	money1 := float64(2.0*gamePlayers[0].getTotalScore()-gamePlayers[1].getTotalScore()-gamePlayers[2].getTotalScore()) / 100.0
	money2 := float64(2.0*gamePlayers[1].getTotalScore()-gamePlayers[0].getTotalScore()-gamePlayers[2].getTotalScore()) / 100.0
	money3 := float64(2.0*gamePlayers[2].getTotalScore()-gamePlayers[0].getTotalScore()-gamePlayers[1].getTotalScore()) / 100.0

	fmt.Printf("\t%s\t%s\t%s\n", gamePlayers[0].getName(), gamePlayers[1].getName(), gamePlayers[2].getName())
	fmt.Printf("EURO %5.2f\t%5.2f\t%5.2f\n", money1, money2, money3)
	fmt.Printf("WON  %5d\t%5d\t%5d\n", gamePlayers[0].getWon(), gamePlayers[1].getWon(), gamePlayers[2].getWon())
	fmt.Printf("LOST %5d\t%5d\t%5d\t\n", gamePlayers[0].getLost(), gamePlayers[1].getLost(), gamePlayers[2].getLost())
	fmt.Printf("bidp %5.0f\t%5.0f\t%5.0f\t\n",
		100*float64(gamePlayers[0].getLost()+gamePlayers[0].getWon())/float64(totalGames-passed),
		100*float64(gamePlayers[1].getLost()+gamePlayers[1].getWon())/float64(totalGames-passed),
		100*float64(gamePlayers[2].getLost()+gamePlayers[2].getWon())/float64(totalGames-passed))
	fmt.Printf("pcw  %5.0f\t%5.0f\t%5.0f\t\n",
		100*float64(gamePlayers[0].getWon())/float64(gamePlayers[0].getLost()+gamePlayers[0].getWon()),
		100*float64(gamePlayers[1].getWon())/float64(gamePlayers[1].getLost()+gamePlayers[1].getWon()),
		100*float64(gamePlayers[2].getWon())/float64(gamePlayers[2].getLost()+gamePlayers[2].getWon()))
	fmt.Printf("pcwd %5.0f\t%5.0f\t%5.0f\t\n",
		100*float64(gamePlayers[0].getWonAsDefenders())/float64(totalGames-passed-(gamePlayers[0].getLost()+gamePlayers[0].getWon())),
		100*float64(gamePlayers[1].getWonAsDefenders())/float64(totalGames-passed-(gamePlayers[1].getLost()+gamePlayers[1].getWon())),
		100*float64(gamePlayers[2].getWonAsDefenders())/float64(totalGames-passed-(gamePlayers[2].getLost()+gamePlayers[2].getWon())))
	fmt.Printf("AVG  %3.1f, passed %d, won %d, lost %d / %d games\n", avg, passed, won, lost, totalGames)
	fmt.Printf("Grand games %d, perc: %5.2f\n", grandGames, 100*float64(grandGames)/float64(totalGames))
	fmt.Printf("Null games %d, perc: %5.2f\n", nullGames, 100*float64(nullGames)/float64(totalGames))
}

func printScore(players []PlayerI) {
	fmt.Printf("\nGAME: %6d (%s) %5d     (%s) %5d     (%s) %5d\n", gameIndex, players[0].getName(), players[0].getTotalScore(), players[1].getName(), players[1].getTotalScore(), players[2].getName(), players[2].getTotalScore())
}

func animation() func() {
	i := 0
	sym := 0
	symbol := []string{
		"\b/",
		"\b-",
		"\b\\",
		"\b|",
	}

	nextSymbol := func () string {
		sym++
		if sym == len(symbol) {
			sym = 0
		}
		return symbol[sym]
	}

	next := func() {
		i++
		if i % 50 == 0 {
			fmt.Print(nextSymbol())
		}
		if i % 1000 == 0 {
			if i % 10000 == 0 {
				fmt.Print("\b# ")
			} else {
				fmt.Print("\b. ")
			}
			if i % 50000 == 0 {
				fmt.Printf(" (%d)\n", i)
			}
		}	
	}
	return next
}

func startServer() *mux.Router {
	rt := mux.NewRouter()
	// Static route for CSS and JS files
	rt.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))
	fmt.Println("Starting a server")
	//templates := template.Must(template.ParseFiles("html/index.html"))
	// rt.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	if err := templates.ExecuteTemplate(w, "index.html", V{}); err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })
	var currentBidIndex = -1
	var secondBidRound = false
	var currentGame = 0
	var ForeHandAnswered = false

	rt.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		currentGame++
		htmlLog("Starting Game: %d\n", currentGame)
		currentBidIndex = -1
		secondBidRound = false
		makeChannels()

		gamePlayers = rotatePlayers(gamePlayers)
		initState()
		DealCards()

		players = []PlayerI{gamePlayers[0], gamePlayers[1], gamePlayers[2]}
		initGame()

		position := 0
		for i, p := range players {
			if player1.getName() == p.getName() {
				position = i
				break
			}
		}
		data := initData{
			player1.getHand(), 
			position,
			player1.getTotalScore(),
			player2.getTotalScore(),
			player3.getTotalScore(),
			currentGame,
		}
		sendJson(w, data)
	})

	rt.HandleFunc("/repeat", func(w http.ResponseWriter, r *http.Request) {
		// currentGame++
		htmlLog("Repeating Game: %d\n", currentGame)
		currentBidIndex = -1
		secondBidRound = false
		makeChannels()

		initState()
		SameCards()

		players = []PlayerI{gamePlayers[0], gamePlayers[1], gamePlayers[2]}
		initGame()

		position := 0
		for i, p := range players {
			if player1.getName() == p.getName() {
				position = i
				break
			}
		}
		data := initData{
			player1.getHand(), 
			position,
			player1.getTotalScore(),
			player2.getTotalScore(),
			player3.getTotalScore(),
			currentGame,
		}
		sendJson(w, data)
	})

	rt.HandleFunc("/getHand/{pl}", func(w http.ResponseWriter, r *http.Request) {
		pl, err := strconv.ParseInt(mux.Vars(r)["pl"], 10, 64)
		if err != nil {
			htmlLog("Invalid index")
			http.Error(w, "Invalid index", http.StatusInternalServerError)
		}
		pi := int(pl)
		if pi >= 0 && pi < len(players) {
			data := players[pi].getHand()
			sendJson(w, data)
		} else {
			htmlLog("Invalid index")
			http.Error(w, "Invalid index", http.StatusInternalServerError)
		}
	})

	rt.HandleFunc("/getHandAndSkat/{pl}", func(w http.ResponseWriter, r *http.Request) {
		pl, err := strconv.ParseInt(mux.Vars(r)["pl"], 10, 64)
		if err != nil {
			htmlLog("Invalid index")
			http.Error(w, "Invalid index", http.StatusInternalServerError)
		}
		pi := int(pl)

		
		// var player *HtmlPlayer
		// player, _ = players[pi].(* HtmlPlayer)

		skat1 := <-skatPositionChannel 
		skat2 := <-skatPositionChannel 

		if pi >= 0 && pi < len(players) {
			data := SkatData{
				players[pi].getHand(),
				skat1,
				skat2,
			}
			sendJson(w, data)
		} else {
			htmlLog("Invalid index")
			http.Error(w, "Invalid index", http.StatusInternalServerError)
		}
	})

	rt.HandleFunc("/bid/{pl}", func(w http.ResponseWriter, r *http.Request) {
		pl, _ := strconv.ParseInt(mux.Vars(r)["pl"], 10, 64)

		htmlLog("/bid/%d\n", pl)
		htmlLog("currentBidIndex %d\n", currentBidIndex)
		if pl == 0 {
			ForeHandAnswered = true
		}
		if pl == 2 {
			secondBidRound = true
		}
		htmlLog("secondBidRound:%v\n", secondBidRound)

		var data BidData

		if (pl == 1 && !secondBidRound) || (pl == 2 && ForeHandAnswered) {
			currentBidIndex++
		}
		if players[pl].accepts(currentBidIndex, false) { // false?? ONLY FOR ISS
			debugTacticsLog("Player %s: %d YES \n", players[pl].getName(), bids[currentBidIndex])
			data = BidData{bids[currentBidIndex], true}
		} else {
			debugTacticsLog("Player %s: PASS %d \n", players[pl].getName(), bids[currentBidIndex])
			data = BidData{bids[currentBidIndex], false}
		}

		time.Sleep(time.Duration(delayMs) * time.Millisecond)
		time.Sleep(time.Duration(delayMs) * time.Millisecond)

		sendJson(w, data)
	})

	rt.HandleFunc("/getbidvalue/{pl}", func(w http.ResponseWriter, r *http.Request) {
		pl, _ := strconv.ParseInt(mux.Vars(r)["pl"], 10, 64)

		htmlLog("/getbidvalue/%d\n", pl)
		htmlLog("currentBidIndex %d\n", currentBidIndex)

		if pl == 0 {
			ForeHandAnswered = true
		}
		if pl == 2 {
			secondBidRound = true
		}

		if (pl == 1 && !secondBidRound) || (pl == 2 && ForeHandAnswered) {
			currentBidIndex++
		}

		htmlLog("BIDVALUE: %v\n", bids[currentBidIndex])
		data := BidData{bids[currentBidIndex], true} // boolean value is ignored
		// time.Sleep(time.Duration(delayMs) * time.Millisecond)
		// time.Sleep(time.Duration(delayMs) * time.Millisecond)

		sendJson(w, data)
	})

	rt.HandleFunc("/declarer/{pl}/{bid}", func(w http.ResponseWriter, r *http.Request) {
		pl, err := strconv.ParseInt(mux.Vars(r)["pl"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bid, err1 := strconv.ParseInt(mux.Vars(r)["bid"], 10, 64)
		if err1 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if pl >= int64(len(players)) {
			http.Error(w, "Index error", http.StatusInternalServerError)
			return
		}
		state.declarer = players[pl]
		players[pl].setDeclaredBid(int(bid))
		htmlLog("Declarer: %s with Bid %d\n", players[pl].getName(), bid)

		if state.declarer == players[0] {
			state.opp1, state.opp2 = players[1], players[2]
		}
		if state.declarer == players[1] {
			state.opp2, state.opp1 = players[0], players[2]
		}
		if state.declarer == players[2] {
			state.opp1, state.opp2 = players[0], players[1]
		}
		state.opp1.setPartner(state.opp2)

		go declareAndPlay() // end of goroutine
	})

	rt.HandleFunc("/getCardPlayed", func(w http.ResponseWriter, r *http.Request) {
		htmlLog("Received /getCardPlayed, reading from playChannel...")
		cp := <-playChannel
		htmlLog("read %v\n", cp)

		//time.Sleep(time.Duration(delayMs) * time.Millisecond)
		//time.Sleep(time.Duration(delayMs) * time.Millisecond)

		sendJson(w, cp)
	})

	rt.HandleFunc("/playCard/{pl}/{suit}/{rank}", func(w http.ResponseWriter, r *http.Request) {
		pl, err := strconv.ParseInt(mux.Vars(r)["pl"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		suit := mux.Vars(r)["suit"]
		rank := mux.Vars(r)["rank"]
		card := Card{suit, rank}
		htmlLog("\nReceived /playCard/%d/ %v \n", pl, card)
		if valid(state.follow, state.trump, gamePlayers[pl].getHand(), card) {
			htmlLog("Sending %v to trickChannel...", card)
			trickChannel <- card
			htmlLog("sent %v\n", card)
		} else {
			htmlLog("Invalid card")
			http.Error(w, "Invalid Card", http.StatusInternalServerError)
		}
	})

	rt.HandleFunc("/discardCard/{pl}/{suit1}/{rank1}/{suit2}/{rank2}", func(w http.ResponseWriter, r *http.Request) {
		pl, err := strconv.ParseInt(mux.Vars(r)["pl"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		suit1 := mux.Vars(r)["suit1"]
		rank1 := mux.Vars(r)["rank1"]
		suit2 := mux.Vars(r)["suit2"]
		rank2 := mux.Vars(r)["rank2"]
		card1 := Card{suit1, rank1}
		card2 := Card{suit2, rank2}
		htmlLog("Received /discardCard from Player %d, %v %v\n", pl, card1, card2)
		if in(players[pl].getHand(), card1, card2) {
			htmlLog("Sending %v %v to discardChannel...", card1, card2)
			discardChannel <- card1
			discardChannel <- card2
			htmlLog("sent\n")
		} else {
			htmlLog("Cards not in hand")
			http.Error(w, "Cards not in hand", http.StatusInternalServerError)
		}
	})

	rt.HandleFunc("/getTrickWinner", func(w http.ResponseWriter, r *http.Request) {
		htmlLog("Wating for card...")

		winnerName := <-winnerChannel
		htmlLog("Sending winner %v\n", winnerName)

		// time.Sleep(time.Duration(delayMs) * time.Millisecond)
		// time.Sleep(time.Duration(delayMs) * time.Millisecond)
		// time.Sleep(time.Duration(delayMs) * time.Millisecond)

		sendJson(w, winnerName)
	})

	rt.HandleFunc("/getTrump", func(w http.ResponseWriter, r *http.Request) {
		htmlLog("Getting trump...")
		sendJson(w, <-trumpChannel)
	})

	rt.HandleFunc("/getScore", func(w http.ResponseWriter, r *http.Request) {
		htmlLog("Getting score...")
		sendJson(w, <-scoreChannel)
	})

	rt.HandleFunc("/pickUp/{b}", func(w http.ResponseWriter, r *http.Request) {
		b := mux.Vars(r)["b"]
		if b != "pick" && b != "hand" {
			http.Error(w, "Expected skat/hand", http.StatusInternalServerError)
			return
		}
		pickUpChannel <- b
	})

	rt.HandleFunc("/declare/{b}", func(w http.ResponseWriter, r *http.Request) {
		b := mux.Vars(r)["b"]
		declareChannel <- b
	})

	return rt

}

func sendJson(w http.ResponseWriter, data interface{}) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type initData struct {
	Hand     []Card
	Position int
	Score1 int
	Score2 int
	Score3 int
	Game int
}

type BidData struct {
	Bid      int
	Accepted bool
}

type SkatData struct {
	Hand     []Card
	SkatPos1	int
	SkatPos2	int
}
