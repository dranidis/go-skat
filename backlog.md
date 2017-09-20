# USE MinMax to evaluate bidding

# Inference: add beliefs
	In a game xskat played 10 trump in a losing trick. He had A trump and J spade kept. 
	A belief should be added and retracted if impossible hands occur in distributions.



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

# GRAND vs SUIT after Bidding
	Prefer a solid SUIT against a mmm GRAND if bid allows it


-------------------------------------------------------------------------
# DONE:

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



