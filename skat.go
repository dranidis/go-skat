package main

import (
	"encoding/json"
	_ "errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	_ "html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var r = rand.New(rand.NewSource(1))
var _ = rand.New(rand.NewSource(time.Now().Unix()))
var logFile io.Writer = nil
var debugTacticsLogFlag = false
var gameLogFlag = false
var fileLogFlag = true
var delayMs = 1000
var totalGames = 21
var oficialScoring = false

var playChannel chan CardPlayed
var trickChannel chan Card
var winnerChannel chan string
var trumpChannel chan string
var scoreChannel chan Score
var pickUpChannel chan string
var declareChannel chan string
var discardChannel chan Card
var skatPositionChannel chan int

func logToFile(format string, a ...interface{}) {
	if fileLogFlag && logFile != nil {
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


func htmlLog(format string, a ...interface{}) {
	red := color.New(color.Bold, color.FgYellow).SprintFunc()
	s := fmt.Sprintf(format, a...)
	fmt.Printf(red(s))
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
	skat             []Card
	trumpsInGame     []Card
	cardsPlayed      []Card
	declarerVoidSuit map[string]bool
}

func makeSuitState() SuitState {
	return SuitState{nil, nil, nil, "", nil, "", []Card{}, []Card{}, []Card{}, []Card{},
		map[string]bool{
			CLUBS: false,
			SPADE: false,
			HEART: false,
			CARO:  false,
		},
	}
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
	if html {
		winnerChannel <- winner.getName()
	}

	winner.setSchwarz(false)
	s.trick = []Card{}
	s.leader = newPlayers[0]

	if s.trump == NULL {
		if winner == s.declarer {
			// declarer lost
			return nil
		}
	}

	return newPlayers
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

func play(s *SuitState, p PlayerI) Card {
	if len(p.getHand()) == 0 {
		log.Fatal("EMPTY HAND")
	}
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

	// Player VOID on suit
	if p == s.declarer {
		if s.follow != "" && card.Suit != s.follow {
			s.declarerVoidSuit[s.follow] = true
		}
	}

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
		return s.valid(playerHand, c)
	})
}

func (s SuitState) valid(playerHand []Card, card Card) bool {
	for _, c := range playerHand {
		// if there is at least one card in your hand matching the followed suit
		// your played card should follow
		if s.follow == getSuit(s.trump, c) {
			return s.follow == getSuit(s.trump, card)
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

	if state.trump == NULL {
		gameLog("\nSCORING\n\tNULL ")
		if handGame {
			gameLog("HAND \n")
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
		gameLog("SCORE %d\n", gs)
	} else {
		mat := matadors(state.trump, cs)
		withMatadors = mat
		if mat < 0 {
			mat = mat * -1
		}
		multiplier := mat + 1

		if withMatadors > 0 {
			gameLog("\nSCORING\n\t%s, With %d ", state.trump, mat)
		} else {
			gameLog("\nSCORING\n\t%s, Without %d ", state.trump, mat)
		}

		base := trumpBaseValue(state.trump)

		if handGame {
			multiplier++
			gameLog("Hand ")
		}
		// Schneider?
		if state.declarer.getScore() > 89 || state.declarer.getScore() < 31 {
			multiplier++
			schneider = true
			gameLog("Schneider ")
		}

		if state.declarer.isSchwarz() || (state.opp1.isSchwarz() && state.opp2.isSchwarz()) {
			multiplier++
			schwarz = true
			gameLog("Schwarz ")
		}
		//gameLog("\n\n")
		gs = multiplier * base

		// OVERBID?
		if gs < state.declarer.getDeclaredBid() {
			gameLog(" --OVERBID!!! Game Value: %d < Bid: %d-- ", gs, state.declarer.getDeclaredBid())
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
		gameLog("SCORE %d\n", gs)
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
var state SuitState

func initGame(players []PlayerI) {
	for _, p := range players {
		p.setScore(0)
		p.setSchwarz(true)
		p.setPreviousSuit("")
		p.setDeclaredBid(0)
	}

	state = makeSuitState()
	state.skat = make([]Card, 2)
	// DEALING
	// for {
	// 		TRYING = false

	debugTacticsLog("SHUFFLING..")
	cards := Shuffle(makeDeck())
	players[0].setHand(sortSuit("", cards[:10]))
	players[1].setHand(sortSuit("", cards[10:20]))
	players[2].setHand(sortSuit("", cards[20:30]))
	copy(state.skat, cards[30:32])
	// // if player1.canWin() == "GRAND" || player2.canWin() == "GRAND" || player3.canWin() == "GRAND" {
	// // 	break
	// // }
	// if player2.canWin() == "GRAND" && TRYING{
	// 	break
	// }	

	// // 	if len(sevens(player1.getHand())) == 4 {
	// // 		debugTacticsLog("FOUR 7\n")
	// // 		break
	// // 	}
	// }
	for _, p := range players {
		h := p.calculateHighestBid()
		gameLog("(%v) hand: %v Bid up to: %d\n", p.getName(), p.getHand(), h)
	}

	gameLog("\nPLAYER ORDER: %s - %s - %s\n\n", players[0].getName(), players[1].getName(), players[2].getName())
}

func declareAndPlay(players []PlayerI) int {
	// HAND GAME?
	handGame := true
	// fmt.Printf("\nHAND bef: %v\n", sortSuit(declarer.getHand()))
	// fmt.Printf("SKAT bef: %v\n", skat)

	if state.declarer.pickUpSkat(state.skat) {
		// fmt.Printf("HAND aft: %v\n", sortSuit(declarer.getHand()))
		handGame = false
		// fmt.Printf("SKAT aft: %v\n", skat)
	}

	state.trump = state.declarer.declareTrump()

	if html {
		htmlLog("Sending trump %v", state.trump)
		trumpChannel <- state.trump
	}

	if state.trump == GRAND {
		//fmt.Println(GRAND)
		grandGames++
	}
	state.trumpsInGame = filter(makeDeck(), func(c Card) bool {
		return getSuit(state.trump, c) == state.trump
	})

	state.leader = players[0]
	state.follow = ""

	players[0].setHand(sortSuit(state.trump, players[0].getHand()))
	players[1].setHand(sortSuit(state.trump, players[1].getHand()))
	players[2].setHand(sortSuit(state.trump, players[2].getHand()))

	gameLog("\n(%s) TRUMP: %s\n", red(state.declarer.getName()), state.trump)
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

	if gs > 0 {
		if oficialScoring {
			state.declarer.incTotalScore(50)
		}
		if state.trump != NULL {
			gameLog("VICTORY: %d - %d, SCORE: %d\n",
				state.declarer.getScore(), state.opp1.getScore()+state.opp2.getScore(), gs)
		} else {
			gameLog("VICTORY: %d\n", gs)
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
			gameLog("DEFEAT: %d - %d, SCORE: %d\n",
				state.declarer.getScore(), state.opp1.getScore()+state.opp2.getScore(), gs)
		} else {
			gameLog("DEFEAT: %d\n", gs)
		}

	}
	gameSc.Total1 = player1.getTotalScore()
	gameSc.Total2 = player2.getTotalScore()
	gameSc.Total3 = player3.getTotalScore()
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

func game(players []PlayerI) int {
	gameLog("\n\nGAME %d/%d\n", gameIndex, totalGames)
	initGame(players)

	// BIDDING
	bidDecl, declarer := bid(players)
	if bidDecl == 0 {
		gameLog("ALL PASSED\n")
		return 0
	}
	debugTacticsLog("Declarer %v\n", declarer)
	declarer.setDeclaredBid(bidDecl)

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

	gs := declareAndPlay(players)
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
var player1 PlayerI
var player2 Player
var player3 Player
var html = false

type V struct {
}

var players []PlayerI

func makeChannels() {
	playChannel = make(chan CardPlayed)
	trickChannel = make(chan Card)
	winnerChannel = make(chan string)
	trumpChannel = make(chan string)
	scoreChannel = make(chan Score)
	pickUpChannel = make(chan string)
	declareChannel = make(chan string)
	discardChannel = make(chan Card)
	skatPositionChannel = make(chan int)
}

func makePlayers(auto, html bool) {
	if auto {
		gameLog("Creating CPU players only\n")
		player := makePlayer([]Card{})
		player1 = &player
		player.risky = true
		delayMs = 0
	} else {
		if html {
			player := makeHtmlPlayer([]Card{})
			player1 = &player
		} else {
			player := makeHumanPlayer([]Card{})
			player1 = &player
		}
		delayMs = 500
	}
	player2 = makePlayer([]Card{})
	player3 = makePlayer([]Card{})
	player1.setName("You")
	player2.setName("Bob")
	player3.setName("Ana")
	player2.risky = true
	player3.risky = false
	players = []PlayerI{player1, &player2, &player3}
}

func main() {

	// COMMAND LINE FLAGS
	auto := false
	var randSeed int
	flag.IntVar(&totalGames, "n", 21, "total number of games")
	flag.IntVar(&randSeed, "r", 0, "Seed for random number generator. A value of 0 (default) uses the UNIX time as a seed.")
	flag.BoolVar(&auto, "auto", false, "Runs with CPU players only")
	flag.BoolVar(&fileLogFlag, "log", true, "Saves log in a file")
	flag.BoolVar(&html, "html", false, "Starts an HTTP server at localhost:8080")
	flag.Parse()

	gameLog("AUTO: %v\n", auto)
	//fmt.Println(fileLogFlag)
	if randSeed == 0 {
		r = rand.New(rand.NewSource(time.Now().Unix()))
	} else {
		r = rand.New(rand.NewSource(int64(randSeed)))
	}

	makePlayers(auto, html)

	rotateTimes := r.Intn(5)
	for i := 0; i < rotateTimes; i++ {
		players = rotatePlayers(players)
	}

	if fileLogFlag {
		file, err := os.Create("gameLog.txt")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		logFile = file
		defer file.Close()
	}

	if html {
		rt := startServer()
		port := ":3000"
		htmlLog("Starting server at %s\n", port)
		http.ListenAndServe(port, rt)
	}

	passed := 0
	won := 0
	lost := 0
	for ; gameIndex <= totalGames; gameIndex++ {

		score := game(players)
		if score == 0 {
			passed++
		}
		if score > 0 {
			won++
		} else if score < 0 {
			lost++
		}
		if !auto {
			fmt.Printf("\nGAME: %6d (%s) %5d     (%s) %5d     (%s) %5d\n", gameIndex, player1.getName(), player1.getTotalScore(), player2.getName(), player2.getTotalScore(), player3.getName(), player3.getTotalScore())
		}
		//time.Sleep(1000 * time.Millisecond)
		players = rotatePlayers(players)
	}

	avg := float64(player1.getTotalScore()+player2.getTotalScore()+player3.getTotalScore()) / float64(totalGames-passed)

	money1 := float64(2.0*player1.getTotalScore()-player2.getTotalScore()-player3.getTotalScore()) / 100.0
	money2 := float64(2.0*player2.getTotalScore()-player1.getTotalScore()-player3.getTotalScore()) / 100.0
	money3 := float64(2.0*player3.getTotalScore()-player1.getTotalScore()-player2.getTotalScore()) / 100.0

	fmt.Printf("\t%s\t%s\t%s\n", player1.getName(), player2.getName(), player3.getName())
	fmt.Printf("EURO %5.2f\t%5.2f\t%5.2f\n", money1, money2, money3)
	fmt.Printf("WON  %5d\t%5d\t%5d\n", player1.getWon(), player2.getWon(), player3.getWon())
	fmt.Printf("LOST %5d\t%5d\t%5d\t\n", player1.getLost(), player2.getLost(), player3.getLost())
	fmt.Printf("bidp %5.0f\t%5.0f\t%5.0f\t\n",
		100*float64(player1.getLost()+player1.getWon())/float64(totalGames-passed),
		100*float64(player2.getLost()+player2.getWon())/float64(totalGames-passed),
		100*float64(player3.getLost()+player3.getWon())/float64(totalGames-passed))
	fmt.Printf("pcw  %5.0f\t%5.0f\t%5.0f\t\n",
		100*float64(player1.getWon())/float64(player1.getLost()+player1.getWon()),
		100*float64(player2.getWon())/float64(player2.getLost()+player2.getWon()),
		100*float64(player3.getWon())/float64(player3.getLost()+player3.getWon()))
	fmt.Printf("pcwd %5.0f\t%5.0f\t%5.0f\t\n",
		100*float64(player1.getWonAsDefenders())/float64(totalGames-passed-(player1.getLost()+player1.getWon())),
		100*float64(player2.getWonAsDefenders())/float64(totalGames-passed-(player2.getLost()+player2.getWon())),
		100*float64(player3.getWonAsDefenders())/float64(totalGames-passed-(player3.getLost()+player3.getWon())))
	fmt.Printf("AVG  %3.1f, passed %d, won %d, lost %d / %d games\n", avg, passed, won, lost, totalGames)
	fmt.Printf("Grand games %d, perc: %5.2f", grandGames, 100*float64(grandGames)/float64(totalGames))
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

	rt.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		currentGame++
		htmlLog("Starting Game: %d\n", currentGame)
		currentBidIndex = -1
		secondBidRound = false
		makeChannels()

		players = rotatePlayers(players)
		initGame(players)

		position := 0
		for i, p := range players {
			if player1 == p {
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
		// htmlLog("currentBidIndex %d\n", currentBidIndex)

		if pl == 2 {
			secondBidRound = true
		}
		var data BidData
		if (pl == 1 && !secondBidRound) || pl == 2 { // SPEAKER
			htmlLog("SPEAKER\n")
			if players[pl].accepts(currentBidIndex + 1) {
				currentBidIndex++
				debugTacticsLog("Player %s: %d\n", players[pl].getName(), bids[currentBidIndex])
				data = BidData{bids[currentBidIndex], true}
			} else {
				debugTacticsLog("Player %s: PASS \n", players[pl].getName())
				data = BidData{bids[currentBidIndex+1], false}
			}
		} else { // LISTENER
			htmlLog("LISTENER\n")

			if currentBidIndex == -1 {
				currentBidIndex++
			}
			if players[pl].accepts(currentBidIndex) {
				debugTacticsLog("Player %s: yes (%d)\n", players[pl].getName(), bids[currentBidIndex])
				data = BidData{bids[currentBidIndex], true}
			} else {
				debugTacticsLog("Player %s: no (%d)\n", players[pl].getName(), bids[currentBidIndex])
				data = BidData{bids[currentBidIndex], false}
			}
		}
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
		time.Sleep(time.Duration(delayMs) * time.Millisecond)

		sendJson(w, data)
	})

	rt.HandleFunc("/getbidvalue/{pl}", func(w http.ResponseWriter, r *http.Request) {
		pl, _ := strconv.ParseInt(mux.Vars(r)["pl"], 10, 64)

		if pl == 2 {
			secondBidRound = true
		}

		bidvalue := currentBidIndex
		if pl == 2 || (pl == 1 && !secondBidRound) {
			bidvalue = currentBidIndex + 1
			currentBidIndex++
		}
		if currentBidIndex == -1 {
			bidvalue = currentBidIndex + 1
			currentBidIndex++
		}
		htmlLog("BIDVALUE: %v\n", bids[bidvalue])
		data := BidData{bids[bidvalue], true} // boolean value is ignored
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

		go declareAndPlay(players) // end of goroutine
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
		htmlLog("Received /playCard/%d/ %v \n", pl, card)
		if state.valid(players[pl].getHand(), card) {
			htmlLog("Sending %v to trickChannel...", card)
			trickChannel <- card
			htmlLog("sent\n", card)
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
