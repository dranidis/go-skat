package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"io"
	"log"
	"strconv"
	"errors"
	"time"
)

var username = "goskat"
var opp1name = ""
var opp2name = ""
var tableNr = int64(-1)
var playerNr = -1
var issTrump = ""
var issSkat = []Card{}

var	waitServer chan string

var connR io.Reader
var connW io.Writer
var serverLogFile *os.File



func Connect(usr, pwd string) error {
	waitServer = make(chan string)

	// connect to server socket
	if ! issDEBUG {
		var conniss io.ReadWriter
		var err error
		if conniss, err = net.Dial("tcp", "skatgame.net:7000"); err != nil {
			fmt.Printf("Error %v\n", err)
			return err
		}
		fmt.Println("Connected to server")
		connR = conniss
		connW = conniss

		fmt.Println("Sending username:", usr)
		fmt.Fprintf(conniss, usr+"\n")
		fmt.Printf("SENT: %v", usr)

		// listen for reply
		message, err := bufio.NewReader(conniss).ReadString('\n')
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return err
		} else {
			fmt.Printf("RCVD: %v\n", message)
		}

		if strings.Index(message, "password") == -1 {
			return errors.New("Error. Password not requested:" + message)
		}

		fmt.Fprintf(conniss, pwd+"\n")
		message, err = bufio.NewReader(conniss).ReadString('\n')
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return err
		} else {
			fmt.Printf("RCVD: %v\n", message)
		}

		if strings.Index(message, "Welcome") == -1 {
			return errors.New("Not logged in:" + message)
		}

		serverLogFile = createFile("server.log")
	} else {
		connR = os.Stdin
		connW = os.Stdout
	}

	username = usr

	go func() {
		createTable()
		invite(issOpp1, issOpp2)
		for i :=0 ; i < totalGames; i++ {
			ready()
			<- waitServer // wait for game end
			fmt.Println("GAME ENDED")			
		}

		leaveTable()

		// for i :=0 ; i < 10; i++ {
		// 	time.Sleep(time.Duration(delayMs) * time.Millisecond)	
		// }	
		fmt.Println("Exiting...")			

		os.Exit(1)		
	}()

	readFromServer() // BLOCKS
	return nil // no error
}

func readFromServer() {
	scanner := bufio.NewScanner(connR)
	text := ""
	for {
		// fmt.Println("..Waiting msg from server..")
	  	for scanner.Scan() {
	  		text = scanner.Text()
			fmt.Fprintf(serverLogFile, "%s\n", text)
		  	parseServer(text)
			// fmt.Println("..Waiting msg from server..")
	  	}
	  	if err := scanner.Err(); err != nil {
	  		fmt.Fprintln(os.Stderr, "ERROR:", err)
	  	}
	}
}

func parseServer(t string) {
	fmt.Printf("RECV: %s\n", green(t))


	// create .3 goskat 3 -1 3 ? ? ? 0
	if strings.HasPrefix(t, "create")  {
		s := strings.Split(t, " ")
		if len(s) > 2 && s[1][0] == '.' && s[2] == username {
			n := strings.Split(s[1], ".")
			// fmt.Printf("number: %v\n", n)
			tableNr, _ = strconv.ParseInt(n[1], 10, 64)
			// fmt.Printf("number: %v\n", number)
			fmt.Printf("Creating table: %d\n", tableNr)
			waitServer <- "OK"
			return
		}
	}

	// table .3 goskat ...................
			// s pick up skat
			// w ??.?? 179.9 239.9 240.0   skat cards
			// number bidding
			// p pass
			// D.??.??     game?
			// XX  card

	if tableNr >=0 && strings.HasPrefix(t, fmt.Sprintf("table .%d %s", tableNr, username)) {
		s := strings.Split(t, " ")

		// table .3 goskat start 1 goskat 240.0 xskat 240.0 xskat:2 240.0
		if s[3] == "start" {
			if s[5] == username {
				playerNr = 0
				opp1name = s[7]
				opp2name = s[9]
			}
			if s[7] == username {
				playerNr = 1
				opp1name = s[9]
				opp2name = s[5]
			}
			if s[9] == username {
				playerNr = 2
				opp1name = s[5]
				opp2name = s[7]
			}
			gameNr := s[4]
			fmt.Printf("You are player %d in game: %s with %s and %s \n", playerNr, gameNr, opp1name, opp2name)
			// waitServer <- "OK"
			return
		}

		if s[3] == "end" {
			s := strings.Split(t, "(")
			processGameFormat("(" + s[1])
				//
			// TODO: get cards info to calculate score based on declarer cards and skat
			//

			waitServer <- "OK"
			return
		}
		// table .3 goskat play ..................
		if len(s) > 5 && s[3] == "play" {
			player := s[4]
			action := s[5]
			// fmt.Printf("Player: %s, Action: %s in [%s]\n", player, action, t)

			if player != "w" && len(action) > 2 {
				ss := strings.Split(action, ".")
				if len(ss) > 0 && ss[0] == "SC" {
					// SHORTCUT SC: many cards played at once
					//table .2 goskat play 2 SC.D9.H8.HK 227.1 230.9 235.9
					fmt.Printf("Player: %s, PLAYED SHORTCARD: %s\n", player, action)
					rank := ss[1]
					for i := 2; i < len(ss); i++ {
						rank += "." + ss[i]
					}
					scCard := Card{"SC", rank}
					_ = scCard
					// ignore it
					// trickChannel <- scCard
					return
				}
			}

			// XX
			if bidNr, err := strconv.ParseInt(action, 10, 64); err == nil {
				fmt.Printf("Player: %s, Bidding: %d\n", player, bidNr)
				if player != fmt.Sprintf("%d", playerNr) { // only sent to ISSPLAYER
					bidChannel <- action
				}
				return
			} else if len(action) == 2 { // CARDS MOSTLY
				if player != "w" {
					if action == "RE" {
						fmt.Printf("Player: %s, RESIGNS, Ignored!\n", player)
						return
					}
					if card, err := parseCard(action); err == nil {
						fmt.Printf("Player: %s, PLAYED: %v\n", player, card)
						if player != fmt.Sprintf("%d", playerNr) { // only sent to ISSPLAYER
							trickChannel <- card
						}
						return					
					}
				} 
			}

			// TODO
			// Hand game declarations
			// example: Grand Hand
			// table .6 goskat play 1 GH 238.0 240.0 240.0

			// D.??.?? 179.9 239.9 240.0   declare game and skat

			if len(action) > 1 && (action[1] == 'H' || action[1] == 'O') || (len(action) > 6 && player != "w") {
				s := strings.Split(action, ".")
				switch s[0][0] {
				case 'C':
					issTrump = CLUBS
				case 'S':
					issTrump = SPADE
				case 'H':
					issTrump = HEART
				case 'D':
					issTrump = CARO
				case 'G':
					issTrump = GRAND
				case 'N':
					issTrump = NULL
				default:
					log.Fatal("Unrecognized game declared ", s[0], " in: ", action)
				}
				fmt.Printf("Player: %s declares %s\n", player, issTrump)
				if len(s[0]) > 1 {
					handIndex := 1
					if s[0][1] == 'O' {
						fmt.Printf("Player: %s plays an OUVERT game %s\n", player)
						handIndex++
					}
					if len(s[0]) > handIndex && s[0][handIndex] == 'H' {
						fmt.Printf("Player: %s plays a HAND game %s\n", player)
						if len(s[0]) > handIndex + 1 {
							if s[0][handIndex+1] == 'S' {
								fmt.Printf("Player: %s announces a SCHNEIDER game %s\n", player)
							}
							if s[0][handIndex+1] == 'Z' {
								fmt.Printf("Player: %s announces a SCHWARZ game %s\n", player)
							}
						}
					}
				}


				if player != fmt.Sprintf("%d", playerNr) { // only sent to ISSPLAYER
					declareChannel <- issTrump
				}

				if len(s) == 3 { // SKAT // DEAD CODE?????
					sk1 := s[1]
					if sk1 != "??" {
						issSkat = make([]Card, 2)
						issSkat[0], _ = parseCard(sk1)
						sk2 := s[2]
						issSkat[1], _ = parseCard(sk2)
						fmt.Printf("SKAT: %v\n", issSkat)
						state.skat = issSkat
						fmt.Printf("THE SKAT IS: %v\n", state.skat)
					}
				}
				return
			}



			// p
			if action == "p" {
				fmt.Printf("Player: %s, PASS\n", player)
				if player != fmt.Sprintf("%d", playerNr) { // only sent to ISSPLAYER
					bidChannel <- action
				}
				return
			}
			// y
			if action == "y" {
				fmt.Printf("Player: %s, ACCEPTS\n", player)
				if player != fmt.Sprintf("%d", playerNr) { // only sent to ISSPLAYER
					bidChannel <- action
				}
				return
			}			

			// s
			if action == "s" {
				fmt.Printf("Player: %s, PICK UP SKAT\n", player)
				// pickUpChannel <- "SKAT"
				return
				// TODO: what happens in a Hand game???
			}

			// w
			if player == "w" {
				if len(action) == 5 && action[2] == '.' {
					ss := strings.Split(action, ".")
					if ss[0] != "??" {
						card1, _ := parseCard(ss[0])
						card2, _ := parseCard(ss[1])
						skatChannel <- card1
						skatChannel <- card2
						fmt.Printf("Player: %s, SKAT: %v %v\n", player, card1, card2)
					}
					return
				} else if len(action) > 1 && action[0] == 'T' && action[1] == 'I' {
				// TI.0	   timeout ?
					ss := strings.Split(action, ".")
					pNr, _ := strconv.ParseInt(ss[1], 10, 64)
					fmt.Printf("Player: %d, TIMEOUT\n", pNr)
					return
				} else {
					// [table .3 goskat play w HK.D9.HA.DK.HJ.DJ.C8.SK.D8.ST|??.??.??.??.??.??.??.??.??.??|??.??.??.??.??.??.??.??.??.??|??.?? 240.0 240.0 240.0]
					hands := strings.Split(action, "|")
					if playerNr < 0 {
						log.Fatal("No player nr yet")
					}
					cards := parseCards(strings.Split(hands[playerNr], "."))
					fmt.Printf("Your hand: %v\n", cards)

					makeChannels()
					gamePlayers[0].setName(username)
					gamePlayers[0].setHand(sortSuit("", cards))
					gamePlayers[1].setName(opp1name)
					gamePlayers[2].setName(opp2name)
					players = []PlayerI{gamePlayers[0], gamePlayers[1], gamePlayers[2]}
					if playerNr == 1 {
						players = rotatePlayers(players)
						players = rotatePlayers(players)
					}
					if playerNr == 2 {
						players = rotatePlayers(players)
					}
					initState()
					
					//DealCards()
					initGame()
					go (func () {
						if bidPhase() == 0 {
							fmt.Println("ISS: All passed")
							// ??????????????
							return 
						}
						gs := declareAndPlay()
						fmt.Println("ISS: gs:", gs)
						})()
					return
				}
			}

		}
	}
	fmt.Println("Unhandled server message: ", t)
}

func parseCard(action string) (Card, error) {
	rank, suit := "", ""
	switch action[0] {
	case 'C':
		rank = "CLUBS"
	case 'S':
		rank = "SPADE"
	case 'H':
		rank = "HEART"
	case 'D':
		rank = "CARO"
	default:
		return Card{"", ""}, errors.New("Not a card suit:" + action)
	}
	switch action[1] {
	case '7':
		suit = "7"
	case '8':
		suit = "8"
	case '9':
		suit = "9"
	case 'T':
		suit = "10"
	case 'Q':
		suit = "D"
	case 'K':
		suit = "K"
	case 'A':
		suit = "A"
	case 'J':
		suit = "J"
	default:
		return Card{"", ""}, errors.New("Not a card rank:" + action)
	}
	return Card{rank, suit}, nil	
}

func parseCards(ss []string) []Card {
	cards := []Card{}
	for _, s := range ss {
		c, _ := parseCard(s)
		cards = append(cards, c)
	}
	return cards
}

// COMMANDS SENT TO SERVER
func createTable() {
	sendToServer("create / 3")
	fmt.Println("...Waiting server response")
	<- waitServer
}

func invite(p1, p2 string) {
	sendToServer(fmt.Sprintf("table .%d %s invite %s %s", tableNr, username, p1, p2))
	opp1name = p1
	opp2name = p2
	fmt.Println("...Waiting server response")
	// <- waitServer
}

func ready() {
	sendToServer(fmt.Sprintf("table .%d %s ready", tableNr, username))

}

func leaveTable() {
	sendToServer(fmt.Sprintf("table .%d %s leave", tableNr, username))
	serverLogFile.Close()
}

func playCard(card Card) {
	cardString := cardString(card)
	sendToServer(fmt.Sprintf("table .%d %s play %s", tableNr, username, cardString))
}

func cardString(card Card) string {
	cardString := ""
	switch card.Suit {
	case CLUBS:
		cardString = "C"
	case SPADE:
		cardString = "S"
	case HEART:
		cardString = "H"
	case CARO:
		cardString = "D"
	}
	switch card.Rank {
	case "10":
		cardString += "T"
	case "D":
		cardString += "Q"
	default:
		cardString += card.Rank
	}
	return cardString	
}

func playBid(bid string) {
	sendToServer(fmt.Sprintf("table .%d %s play %s", tableNr, username, bid))
}

// TODO
// Declare and discard
func sendPickUpSkat() {
	sendToServer(fmt.Sprintf("table .%d %s play s", tableNr, username))
}


func iss_declare(trump string, hand bool, skat []Card) {
	game := ""
	switch trump {
	case CLUBS:
		game = "C"
	case SPADE:
		game = "S"
	case HEART:
		game = "H"
	case CARO:
		game = "D"
	case GRAND:
		game = "G"
	case NULL:
		game = "N"
	}
	if hand {
		game += "H"
	}
	if !hand {
		card1 := cardString(skat[0])	
		card2 := cardString(skat[1])	
		sendToServer(fmt.Sprintf("table .%d %s play %s.%s.%s", tableNr, username, game, card1, card2))
	} else {
		sendToServer(fmt.Sprintf("table .%d %s play %s", tableNr, username, game))
	}
}

func sendToServer(s string) {
	if !issDEBUG {
		time.Sleep(time.Duration(issSentDelay) * time.Millisecond)
	} else {
		time.Sleep(time.Duration(1000) * time.Millisecond)		
	}
	fmt.Printf("SENT: %s\n", yellow(s))
	fmt.Fprintf(connW, "%s\n", s)
}

func processGameFormat(s string) {
	// table .8 goskat end 
			// (;GM[Skat]PC[International Skat Server]CO[]SE[252547]ID[4940313]DT[2017-09-20/12:43:07/UTC]P0[zoot]P1[goskat]P2[bernie]R0[]R1[0.0]R2[]MV[w HK.CA.H8.CK.CQ.HT.CJ.SK.DA.D7.HJ.ST.SJ.S8.C9.H9.S9.D8.H7.CT.HA.D9.DK.S7.SQ.SA.DJ.C8.C7.DQ.DT.HQ 1 p 2 p 0 18 0 s w DT.HQ 0 H.SK.D7 0 H8 1 H9 2 HA 2 SA 0 HK 1 S9 0 DA 1 D8 2 D9 0 CA 1 C9 2 C8 0 CQ 1 CT 2 C7 1 S8 2 SQ 0 HT 0 DT 1 H7 2 DK 1 SJ 2 DJ 0 CJ 0 CK 1 HJ 2 DQ 1 ST 2 S7 0 HQ ]R[d:0 win v:20 m:1 bidok p:73 t:6 s:0 z:0 p0:0 p1:0 p2:0 l:-1 to:-1 r:0] ;)

}
