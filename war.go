package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Card struct {
	// a.e. QH = Queen of Hearts
	name string
	// Numberic value to easily compare cards (ignoring suites and queens and kings etc. converted to integer)
	value int
}

type Game struct {
	playersDecks [][]Card
	// Cards which are on the table
	tableDeck []Card
	// Next round is dummy (when we take cards, but don't compare them)
	dummyRound bool

	// Meta-data used for the analysis
	roundsCount int
	state       int
}

var suites = []string{"H", "S", "C", "D"}
var ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

var randomSource = rand.New(rand.NewSource(time.Now().Unix()))

const (
	Ok                           = iota
	RoundCompletedDraw           = iota
	GameCompletedFirstPlayerWon  = iota
	GameCompletedSecondPlayerWon = iota
	GameCompletedDraw            = iota
)

func generateSortedDeck() []Card {
	deck := make([]Card, 0, len(ranks)*len(suites))
	for _, suite := range suites {
		for index, rank := range ranks {
			deck = append(deck, Card{rank + suite, index + 2})
		}
	}
	return deck
}

func shuffleDeck(deck []Card) []Card {
	for n := len(deck); n > 0; n-- {
		randIndex := randomSource.Intn(n)
		deck[n-1], deck[randIndex] = deck[randIndex], deck[n-1]
	}
	return deck
}

func splitDeck(deck []Card, cardsInFirstDeck int) [][]Card {
	decks := make([][]Card, 2)
	decks[0] = make([]Card, 0, 54)
	decks[1] = make([]Card, 0, 54)

	for i := 0; i < len(deck); i++ {
		if i < cardsInFirstDeck {
			decks[0] = append(decks[0], deck[i])
		} else {
			decks[1] = append(decks[1], deck[i])
		}
	}

	return decks
}

func createNewGame() Game {
	game := Game{state: Ok, dummyRound: false}
	deck := shuffleDeck(generateSortedDeck())
	game.playersDecks = splitDeck(deck, len(deck)/2)

	game.tableDeck = make([]Card, 0, len(deck))

	return game
}

func pickTopCards(playersDecks [][]Card) (player1Card Card, player2Card Card, gameResult int) {
	// Check whether decks are empty
	player1DeckIsEmpty := len(playersDecks[0]) == 0
	player2DeckIsEmpty := len(playersDecks[1]) == 0

	if player1DeckIsEmpty && player2DeckIsEmpty {
		return Card{}, Card{}, GameCompletedDraw
	}
	if player1DeckIsEmpty {
		return Card{}, Card{}, GameCompletedSecondPlayerWon
	}
	if player2DeckIsEmpty {
		return Card{}, Card{}, GameCompletedFirstPlayerWon
	}

	player1Card = playersDecks[0][0]
	playersDecks[0] = playersDecks[0][1:]

	player2Card = playersDecks[1][0]
	playersDecks[1] = playersDecks[1][1:]

	return player1Card, player2Card, Ok
}

func winnerTakesAll(game *Game, winnerIndex int) {
	// Shuffle cards which we take (so we don't end up with infinite game cycles)
	game.playersDecks[winnerIndex] = append(game.playersDecks[winnerIndex], shuffleDeck(game.tableDeck)...)
	// Clear table deck
	game.tableDeck = game.tableDeck[:0]
}

func isGameComplete(result int) bool {
	return result == GameCompletedDraw || result == GameCompletedFirstPlayerWon || result == GameCompletedSecondPlayerWon
}

func playMove(game *Game) {
	player1Card, player2Card, result := pickTopCards(game.playersDecks)
	if isGameComplete(result) {
		game.state = result
		return
	}
	game.tableDeck = append(game.tableDeck, player1Card, player2Card)

	if game.dummyRound {
		game.dummyRound = false
	} else if player1Card.value > player2Card.value {
		winnerTakesAll(game, 0)
	} else if player2Card.value > player1Card.value {
		winnerTakesAll(game, 1)
	} else {
		game.dummyRound = true
	}
}

func playGame(game *Game) {
	for true {
		playMove(game)
		if isGameComplete(game.state) {
			break
		}
		game.roundsCount = game.roundsCount + 1
	}
}

func main() {
	gamesCount := 0
	totalRounds := 0
	totalFirstPlayerWins := 0
	for true {
		game := createNewGame()
		playGame(&game)

		// Calculate statistics
		gamesCount = gamesCount + 1
		totalRounds = totalRounds + game.roundsCount
		if game.state == GameCompletedFirstPlayerWon {
			totalFirstPlayerWins = totalFirstPlayerWins + 1
		}
		if game.state == GameCompletedDraw {
			fmt.Println("Game ", gamesCount, " Whoaaaa. It's a draw. ")
		}

		// Print statistics
		if gamesCount%100000 == 0 {
			fmt.Println()
			fmt.Println("Game ", gamesCount, "is done")
			fmt.Println("Statistics:")
			fmt.Println("Average moves per game:", float64(totalRounds)/float64(gamesCount))
			fmt.Println("First player wins in ", float64(totalFirstPlayerWins)/float64(gamesCount))
		}
	}
}
