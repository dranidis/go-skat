
# THrowing a Clubs J, when a CARO J would do!? See clubsJ... txt
	Trick: []
	(bernie) HAND [] valid: []
		Previous suit: bernie:, ddmits:
	(bernie) Trick: [A]
	(ddmits) HAND [] valid: []
		Previous suit: bernie:CLUBS, ddmits:
	(ddmits) Trick: [A 7]
	(goskat) HAND [J J A 10 8 A 7 A D 9] valid: [J J A 10 8 A 7 A D 9]
	(goskat) Valid: [J J A 10 8 A 7 A D 9]
	Equivalent: [A 10 A A D J J 9 7 8]
	Valid: [A 10 A A D J J 9 7 8]Similar: [J J A 8 A 7 A D 9]
	AB: ALL REMAINING CARDS (18): [10 K 9 8 J K 9 7 J 10 K D 9 8 10 K 8 7]
	AB: REMAINING after void: 18 cards: [10 K 9 8 J K 9 7 J 10 K D 9 8 10 K 8 7] [] []
	AB: (goskat) 20 Worlds, ab
	AB: Opp1: [9 K 7 10 9 10 K 10 8] Opp2: [D 8 K K J 8 7 J 9], SKAT:[D D]	J   90
	AB: Opp1: [8 K 9 D K J 8 9 10] Opp2: [K 10 7 K 8 J 9 10 7], SKAT:[D D]	J   90
	AB: Opp1: [8 7 K 10 10 9 10 8 J] Opp2: [J 7 K D 9 K 9 8 K], SKAT:[D D]	J   76
	AB: Opp1: [7 7 8 9 J 9 10 K J] Opp2: [10 K 8 K 9 8 10 D K], SKAT:[D D]	J   76
	AB: Opp1: [7 J 8 9 10 7 9 K 10] Opp2: [K J 9 K 8 10 8 K D], SKAT:[D D]	J   69
	AB: Opp1: [10 9 D K 8 K J J 10] Opp2: [K K 7 9 8 7 10 9 8], SKAT:[D D]	J   62
	AB: Opp1: [8 7 8 9 K 9 K K 10] Opp2: [K 7 D 10 10 9 J 8 J], SKAT:[D D]	J   58
	AB: Opp1: [K 8 10 7 9 K 9 K 8] Opp2: [10 8 10 J 9 7 J D K], SKAT:[D D]	J   66
	AB: Opp1: [8 K K 10 J 9 8 K 9] Opp2: [10 7 D K 9 8 J 10 7], SKAT:[D D]	J   91
	AB: TIMEOUT
	AB: Time: 5.197365995s
	AB: 9 Worlds: map[A:40.111111111111114 J:72.33333333333333 J:70.44444444444444 A:47.111111111111114 D:58.111111111111114 9:58 8:53.111111111111114 A:37.77777777777778 7:56.333333333333336]
	AB: (goskat) Hand: [J J A 10 8 A 7 A D 9], Playing card: J with value  72.3)


		goskat: void:	[]
		bernie: void:	[]
		bernie: (D) void:
		ddmits: void:	[]
			Played cards  : [A 7]
			Trumps in game: [J J J J]

	DECLARER Cards of suit Grand still in play: [J J       ]Cards of suit Grand still in play: [J J       ]..sure winners: [J]
	BACKHAND must not follow...CHECKING the A... in trick [A 7], valid [A A A D J J 7 8 9], played [A 7]..Winners: [J J]
	LAST RESORT: returning highWinnerLowLoser: J
	AB: 
	TACTICS suggest: J



# Optimize Legal actions in MinMax
	Evaluate only better actions.
	For example when a player surely loses the trick and he needs to choose between 10, K, D, 9 only the 9 should be evaluated as a legal action.

# NULL Evaluation has to be more strict

# MM (not ab) is taking too long!
	TIMEOUT

# A very interesting game: go build; go install; go-skat -minmax2 -minmax3 -html -r 999 -g 8
	

# World generation
	Generate 1000 valid worlds
	Choose 20 from these
	At every move if you find a world is no more valid, make it invalid so that it is not picked up


# Analysis of games played
	Type of game played
	Player
	Percentages

	e.g.
	GRAND/NULL/SUIT/GRAND H/NULL H/SUIT H
	You Bob Ana
	10% 12% 13% declared
	68% 72% 65% won
	XX% XX% XX%  lost

# INFERENCE
	For each suit count the cards played to conclude about the SKAT cards
	7 cards of each suit.
	1st trick 3 cards
	2nd trick 3 cards
	High probability that card in the skat (if you don't have it)

# REFACTORING
	Move players array into SuitState
	Maybe also gamePlayers array

# Find out why MINMAX Tactics player is losing:
	go build; go install; go-skat -auto -n 10  -v  -r 15 -minmax2

# BUG in BIDDING HTML
	IN HTML if both players have passed the offer to the 3rd is 20 instead of 18

# HTML
	Change score back at a replay (or do not cound new score)

# Game Lost?
	HAND: gs: CJ, HJ, CA, C10, C9, C7, 		HA, H7,		DK, DD
	x1 	gs 	x2 	x1 	gs
MH:	SA 	CA 	S9 			22-0
FH:		HA 	H9 	H8 		33-0  (An A is wasted. Opponents throw off. Play a trump instead!)
FH: 	HJ 	CD 	CJ 		33-7	(All evaluations are already negative!)
MH: S10 C10 DJ 			33-29	(lost a full trump)
BH: 		DA 	D8 	DD 	33-56
BH: 		D10 C8 	DK 	33-70  (LOST)

 table .5 goskat end (;GM[Skat]PC[International Skat Server]CO[]SE[252902]ID[4947663]DT[2017-09-28/12:14:02/UTC]P0[xskat]P1[goskat]P2[bernie]R0[]R1[0.0]R2[]MV[w HQ.CJ.ST.HK.D8.S7.SQ.H8.SA.C8.SJ.HJ.CA.S8.HA.DK.H7.CT.C7.C9.DA.D7.CK.HT.DJ.S9.CQ.DT.D9.H9.SK.DQ 1 18 0 p 2 p 1 s w SK.DQ 1 C.S8.SK 0 SA 1 CA 2 S9 1 HA 2 H9 0 H8 1 HJ 2 CQ 0 CJ 0 ST 1 CT 2 DJ 2 HT 0 HQ 1 H7 2 DA 0 D8 1 DQ 2 DT 0 C8 1 DK 0 HK 1 SJ 2 D9 1 C9 2 CK 0 SQ 2 D7 0 S7 1 C7 ]R[d:1 loss v:-48 m:-1 bidok p:43 t:4 s:0 z:0 p0:0 p1:0 p2:0 l:-1 to:-1 r:0] ;)

-------------------------------------------------------------------------
# DONE:

# USE MinMax to evaluate declaring after winning the bid
	Run a Grand game, and suit game below the bid
	Choose game which gives better chances of winning.

# 	Equivalent cards
	When player holds two equivalent cards in the hand like 78 or 89 they should always be evaluated together

# INFERENCE:
	DO not retract everything.!

# DISCARD
	It is best to discard a doubleton than a singleton and another card.

# FLAGS
	Allow customization of auto players:
	e.g -minmax1 Make 1st player minmax

# BUG:
	go build; go install; go-skat -n 12 -auto -cpuprofile=mm.prof -minmax2 -minmax3 -v -r 3668
	At the 8th game crashes

# TACTICS win this game with 106-14, MinMax 67-53!!!!
	go build; go install; go-skat -html  -r 1453 -g 22 -minmax2
	Opening and wasting the trumps on zero tricks

	I will try with MINMAX MC examining all actions in all worlds.

# Inference: add beliefs
	In a game xskat played 10 trump in a losing trick. He had A trump and J spade kept. 
	A belief should be added and retracted if impossible hands occur in distributions.


