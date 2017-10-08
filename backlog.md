
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


