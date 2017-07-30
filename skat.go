package main

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"time"
)

var logFile io.Writer = nil
var debugTacticsLogFlag = false
var gameLogFlag = false
var delayMs = 0
var totalGames = 21

func logToFile(format string, a ...interface{}) {
	if logFile != nil {
		fmt.Fprintf(logFile, format, a...)
	}
}
func bidLog(format string, a ...interface{}) {
	if gameLogFlag {
		fmt.Printf(format, a...)
	}
	logToFile(format, a...)
}

func gameLog(format string, a ...interface{}) {
	if gameLogFlag {
		fmt.Printf(format, a...)
	}
	logToFile(format, a...)
}

func debugTacticsLog(format string, a ...interface{}) {
	if debugTacticsLogFlag {
		fmt.Printf(format, a...)
	}
	logToFile(format, a...)
}

type SuitState struct {
	declarer PlayerI
	opp1     PlayerI
	opp2     PlayerI
	trump    string
	leader   PlayerI
	follow   string
	trick    []Card
	// not necessary for game but for tactics
	trumpsInGame []Card
	cardsPlayed  []Card
}

func makeSuitState() SuitState {
	return SuitState{nil, nil, nil, "", nil, "", []Card{}, []Card{}, []Card{}}
}

func setNextTrickOrder(s *SuitState, players []PlayerI) []PlayerI {
	var newPlayers []PlayerI
	var winner PlayerI
	if s.greater(s.trick[0], s.trick[1]) && s.greater(s.trick[0], s.trick[2]) {
		winner = players[0]
		newPlayers = players
	} else if s.greater(s.trick[1], s.trick[2]) {
		winner = players[1]
		newPlayers = []PlayerI{players[1], players[2], players[0]}
	} else {
		winner = players[2]
		newPlayers = []PlayerI{players[2], players[0], players[1]}
	}

	winner.setScore(winner.getScore() + sum(s.trick))

	if s.declarer != nil && s.opp1 != nil && s.opp2 != nil {
		gameLog("TRICK %v\n", s.trick)
		debugTacticsLog("%d points: %d - %d\n", sum(s.trick), s.declarer.getScore(), s.opp1.getScore()+s.opp2.getScore())
	}

	winner.setSchwarz(false)
	s.trick = []Card{}
	s.leader = newPlayers[0]

	return newPlayers
}

func round(s *SuitState, players []PlayerI) []PlayerI {
	play(s, players[0])
	time.Sleep(time.Duration(delayMs) * time.Millisecond)
	s.follow = getSuite(s.trump, s.trick[0])
	play(s, players[1])
	time.Sleep(time.Duration(delayMs) * time.Millisecond)
	play(s, players[2])
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	players = setNextTrickOrder(s, players)
	s.follow = ""
	return players
}

func play(s *SuitState, p PlayerI) Card {
	valid := sortSuit(s.trump, validCards(*s, p.getHand()))

	p.setHand(sortSuit(s.trump, p.getHand()))
	gameLog("Trick: %v\n", s.trick)
	debugTacticsLog("(%v) HAND %v ", p.getName(), p.getHand())
	debugTacticsLog("valid: %v\n", valid)
	if s.opp1 != nil && s.opp2 != nil {
		debugTacticsLog("Previous suit: %s:%v, %s:%v\n",
			s.opp1.getName(), s.opp1.getPreviousSuit(),
			s.opp2.getName(), s.opp2.getPreviousSuit())
	}
	if s.declarer == p {
		red := color.New(color.Bold, color.FgRed).SprintFunc()
		gameLog("(%s) ", red(p.getName()))
	} else {
		gameLog("(%v) ", p.getName())
	}
	card := p.playerTactic(s, valid)
	p.setHand(remove(p.getHand(), card))
	s.trick = append(s.trick, card)
	if getSuite(s.trump, card) == s.trump {
		s.trumpsInGame = remove(s.trumpsInGame, card)
	}
	s.cardsPlayed = append(s.cardsPlayed, card)
	return card
}

// Returns a list of all cards that are playeable from the player's hand.
func validCards(s SuitState, playerHand []Card) []Card {
	return filter(playerHand, func(c Card) bool {
		return s.valid(playerHand, c)
	})
}

func (s SuitState) valid(playerHand []Card, card Card) bool {
	for _, c := range playerHand {
		// if there is at least one card in your hand matching the followed suit
		// your played card should follow
		if s.follow == getSuite(s.trump, c) {
			return s.follow == getSuite(s.trump, card)
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
	for speaker.accepts(bidIndex) {
		bidLog("\t(%v) %d\n", speaker.getName(), bids[bidIndex])
		if listener.accepts(bidIndex) {
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
		if players[0].accepts(0) {
			bidLog("\t(%s) Yes %d\n", players[0].getName(), bids[0])
			return 0, players[0]
		} else {
			bidLog("\t(%s) Pass\n", players[0].getName())
			return -1, nil
		}
	}
	//	p.isDeclarer = true
	return bidIndex, p
}

func gameScore(state SuitState, cs []Card, score, bid int,
	decSchwarz, oppSchwarz, handGame bool) int {
	mat := matadors(state.trump, cs)
	if mat < 0 {
		mat = mat * -1
	}
	multiplier := mat + 1

	gameLog("\nSCORING\n\tWith %d ", mat)

	base := trumpBaseValue(state.trump)

	if handGame {
		multiplier++
		gameLog("Hand ")
	}
	// Schneider?
	if score > 89 || score < 31 {
		multiplier++
		gameLog("Schneider ")
	}

	if decSchwarz || oppSchwarz {
		multiplier++
		gameLog("Schwarz ")
	}
	gameLog("\n\n")
	gs := multiplier * base

	// OVERBID?
	if gs < bid {
		gameLog("OVERBID!!! Game Value: %d < Bid: %d", gs, bid)
		leastMult := 0
		for leastMult*base < bid {
			leastMult++
		}
		return -2 * leastMult * base
	}

	if score > 60 {
		return gs
	} else {
		return -2 * gs
	}
}

func game(players []PlayerI) int {
	gameLog("\n\nGAME %d/%d\n", gameIndex, totalGames)
	// DEALING
	cards := Shuffle(makeDeck())
	players[0].setHand(sortSuit("", cards[:10]))
	players[1].setHand(sortSuit("", cards[10:20]))
	players[2].setHand(sortSuit("", cards[20:30]))
	skat := make([]Card, 2)
	copy(skat, cards[30:32])
	for _, p := range players {
		debugTacticsLog("(%v) hand: %v Bid up to: %d\n", p.getName(), p.getHand(), p.calculateHighestBid())
	}

	gameLog("\nPLAYER ORDER: %s - %s - %s\n\n", players[0].getName(), players[1].getName(), players[2].getName())

	// BIDDING
	bidIndex, declarer := bid(players)
	if bidIndex == -1 {
		gameLog("ALL PASSED\n")
		return 0
	}
	var opp1, opp2 PlayerI
	if declarer == players[0] {
		opp1, opp2 = players[1], players[2]
	}
	if declarer == players[1] {
		opp2, opp1 = players[0], players[2]
	}
	if declarer == players[2] {
		opp1, opp2 = players[0], players[1]
	}

	// HAND GAME?
	handGame := true
	// fmt.Printf("\nHAND bef: %v\n", sortSuit(declarer.getHand()))
	// fmt.Printf("SKAT bef: %v\n", skat)

	if declarer.pickUpSkat(skat) {
		// fmt.Printf("HAND aft: %v\n", sortSuit(declarer.getHand()))
		handGame = false
		// fmt.Printf("SKAT aft: %v\n", skat)
	}

	trump := declarer.declareTrump()
	allTrumps := filter(makeDeck(), func(c Card) bool {
		return getSuite(trump, c) == trump
	})
	// DECLARE
	state := SuitState{
		declarer, opp1, opp2,
		trump,
		players[0],
		"",
		[]Card{},
		allTrumps,
		[]Card{},
	}
	players[0].setHand(sortSuit(state.trump, players[0].getHand()))
	players[1].setHand(sortSuit(state.trump, players[1].getHand()))
	players[2].setHand(sortSuit(state.trump, players[2].getHand()))

	gameLog("\n(%s) TRUMP: %s\n", red(declarer.getName()), state.trump)
	declarerCards := make([]Card, len(declarer.getHand()))
	copy(declarerCards, declarer.getHand())
	declarerCards = append(declarerCards, skat...)

	// PLAY
	for i := 0; i < 10; i++ {
		debugTacticsLog("TRUMPS IN PLAY %v\n", state.trumpsInGame)
		gameLog("\n")
		players = round(&state, players)
	}
	gameLog("\nSKAT: %v\n", skat)

	// gameLog("SKAT: %v, %d\n", skat, sum(skat))
	declarer.setScore(declarer.getScore() + sum(skat))

	gs := gameScore(state, declarerCards, declarer.getScore(), bids[bidIndex],
		declarer.isSchwarz(), opp1.isSchwarz() && opp2.isSchwarz(), handGame)

	declarer.incTotalScore(gs)

	if declarer.getScore() > 60 && gs > 0 {
		gameLog("VICTORY: %d - %d, SCORE: %d\n",
			declarer.getScore(), opp1.getScore()+opp2.getScore(), gs)
	} else {
		gameLog("DEFEAT: %d - %d, SCORE: %d\n",
			declarer.getScore(), opp1.getScore()+opp2.getScore(), gs)
	}

	return gs

}

func rotatePlayers(players []PlayerI) []PlayerI {
	newPlayers := []PlayerI{}
	newPlayers = append(newPlayers, players[1])
	newPlayers = append(newPlayers, players[2])
	newPlayers = append(newPlayers, players[0])
	return newPlayers
}

var gameIndex = 1

func main() {
	file, err := os.Create("gameLog.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	logFile = file
	defer file.Close()

	// player1 := makePlayer([]Card{})
	player1 := makeHumanPlayer([]Card{})
	gameLogFlag = true
	player2 := makePlayer([]Card{})
	player3 := makePlayer([]Card{})
	player1.setName("You")
	player2.setName("Bob")
	player3.setName("Ana")
	// Try a player with a first card tactic
	//player3.firstCardPlay = true

	players := []PlayerI{&player1, &player2, &player3}
	rotateTimes := r.Intn(5)
	for i := 0; i < rotateTimes; i++ {
		players = rotatePlayers(players)
	}

	passed := 0
	won := 0
	lost := 0
	for ; gameIndex <= totalGames; gameIndex++ {
		for _, p := range players {
			p.setScore(0)
			p.setSchwarz(true)
			p.setPreviousSuit("")
		}
		score := game(players)
		if score == 0 {
			passed++
		}
		if score > 0 {
			won++
		} else if score < 0 {
			lost++
		}
		fmt.Printf("\n(%s) %5d     (%s) %5d     (%s) %5d\n", player1.getName(), player1.getTotalScore(), player2.getName(), player2.getTotalScore(), player3.getName(), player3.getTotalScore())
		//time.Sleep(1000 * time.Millisecond)
		players = rotatePlayers(players)
	}
	avg := float64(player1.getTotalScore()+player2.getTotalScore()+player3.getTotalScore()) / float64(totalGames-passed)

	money1 := float64(2.0*player1.getTotalScore()-player2.getTotalScore()-player3.getTotalScore()) / 100.0
	money2 := float64(2.0*player2.getTotalScore()-player1.getTotalScore()-player3.getTotalScore()) / 100.0
	money3 := float64(2.0*player3.getTotalScore()-player1.getTotalScore()-player2.getTotalScore()) / 100.0
	fmt.Printf("\n(%s) %5.2f     (%s) %5.2f     (%s) %5.2f\n",
		player1.getName(), money1,
		player2.getName(), money2,
		player3.getName(), money3)
	fmt.Printf("AVG %3.1f, passed %d, won %d, lost %d / %d games\n", avg, passed, won, lost, totalGames)
}
