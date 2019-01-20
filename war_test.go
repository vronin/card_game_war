package main

import (
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Fatal("Expected ", b, "result ", a)
}

func Test_splitDeck_InHalf_SplitsDeckCorrectly(t *testing.T) {
	// Arrange
	deck := []Card{{"2C", 2}, {"3C", 3}, {"4C", 4}, {"5C", 5}}

	// Act
	splitDecks := splitDeck(deck, 2)

	// Assert
	AssertEqual(t, len(splitDecks[0]), 2)
	AssertEqual(t, len(splitDecks[1]), 2)
	AssertEqual(t, splitDecks[0][0].name, "2C")
	AssertEqual(t, splitDecks[0][1].name, "3C")
	AssertEqual(t, splitDecks[1][0].name, "4C")
	AssertEqual(t, splitDecks[1][1].name, "5C")
}

func Test_splitDeck_InDifferentSize_SplitsDeckCorrectly(t *testing.T) {
	// Arrange
	deck := []Card{{"2C", 2}, {"3C", 3}, {"4C", 4}, {"5C", 5}}

	// Act
	splitDecks := splitDeck(deck, 1)

	// Assert
	AssertEqual(t, len(splitDecks[0]), 1)
	AssertEqual(t, len(splitDecks[1]), 3)
	AssertEqual(t, splitDecks[0][0].name, "2C")
	AssertEqual(t, splitDecks[1][0].name, "3C")
	AssertEqual(t, splitDecks[1][1].name, "4C")
	AssertEqual(t, splitDecks[1][2].name, "5C")
}

func Test_pickTopCards_ifBothHaveCards_ReturnsThem(t *testing.T) {
	// Arrange
	decks := [][]Card{
		// First deck
		{{"2C", 2}},
		// Second deck
		{{"3C", 3}},
	}

	// Act
	card1, card2, result := pickTopCards(decks)

	// Assert
	AssertEqual(t, card1.name, "2C")
	AssertEqual(t, card2.name, "3C")
	AssertEqual(t, result, Ok)
}

func Test_pickTopCards_FirstPlayerOutOfCards_ReturnsSecondPlayerWon(t *testing.T) {
	// Arrange
	decks := [][]Card{
		// First deck
		{},
		// Second deck
		{{"2C", 2}},
	}

	// Act
	_, _, result := pickTopCards(decks)

	// Assert
	AssertEqual(t, result, GameCompletedSecondPlayerWon)
}

func Test_pickTopCards_SecondPlayerOutOfCards_ReturnsFirstPlayerWon(t *testing.T) {
	// Arrange
	decks := [][]Card{
		// First deck
		{{"2C", 2}},
		// Second deck
		{},
	}

	// Act
	_, _, result := pickTopCards(decks)

	// Assert
	AssertEqual(t, result, GameCompletedFirstPlayerWon)
}

func Test_pickTopCards_BothPlayerOutOfCards_ReturnsDraw(t *testing.T) {
	// Arrange
	decks := [][]Card{
		// First deck
		{},
		// Second deck
		{},
	}

	// Act
	_, _, result := pickTopCards(decks)

	// Assert
	AssertEqual(t, result, GameCompletedDraw)
}

func Test_winnerTakesAll_ForFirstPlayer_FirstPlayerDeckAdded(t *testing.T) {
	// Arrange
	game := Game{state: Ok}
	game.playersDecks = [][]Card{
		// First deck
		{{"2C", 2}},
		// Second deck
		{},
	}
	game.tableDeck = []Card{{"3C", 3}}

	// Act
	winnerTakesAll(&game, 0)

	// Assert
	AssertEqual(t, len(game.playersDecks[0]), 2)
	AssertEqual(t, game.playersDecks[0][0].name, "2C")
	AssertEqual(t, game.playersDecks[0][1].name, "3C")
	AssertEqual(t, len(game.tableDeck), 0)
}

func Test_winnerTakesAll_ForSecondPlayer_SecondPlayerDeckAdded(t *testing.T) {
	// Arrange
	game := Game{state: Ok}
	game.playersDecks = [][]Card{
		// First deck
		{{"2C", 2}},
		// Second deck
		{},
	}
	game.tableDeck = []Card{{"3C", 3}}

	// Act
	winnerTakesAll(&game, 1)

	// Assert
	AssertEqual(t, len(game.playersDecks[1]), 1)
	AssertEqual(t, game.playersDecks[1][0].name, "3C")
	AssertEqual(t, len(game.tableDeck), 0)
}

func Test_playMove_NormalRoundWhenFirstPlayerTakes_CorrectGameEndState(t *testing.T) {
	// Arrange
	game := Game{state: Ok}
	game.playersDecks = [][]Card{
		// First deck
		{{"3C", 3}},
		// Second deck
		{{"2C", 2}},
	}
	game.tableDeck = []Card{}

	// Act
	playMove(&game)

	// Assert
	AssertEqual(t, game.state, Ok)
	AssertEqual(t, len(game.playersDecks[0]), 2)
	AssertEqual(t, len(game.playersDecks[1]), 0)
	AssertEqual(t, len(game.tableDeck), 0)
	AssertEqual(t, game.dummyRound, false)
}

func Test_playMove_NormalRoundWhenSecondPlayerTakes_CorrectGameEndState(t *testing.T) {
	// Arrange
	game := Game{state: Ok}
	game.playersDecks = [][]Card{
		// First deck
		{{"2C", 2}},
		// Second deck
		{{"3C", 3}},
	}
	game.tableDeck = []Card{}

	// Act
	playMove(&game)

	// Assert
	AssertEqual(t, game.state, Ok)
	AssertEqual(t, len(game.playersDecks[0]), 0)
	AssertEqual(t, len(game.playersDecks[1]), 2)
	AssertEqual(t, len(game.tableDeck), 0)
	AssertEqual(t, game.dummyRound, false)
}

func Test_playMove_ForEqualValueCards_GoesIntoDummyRound(t *testing.T) {
	// Arrange
	game := Game{state: Ok}
	game.playersDecks = [][]Card{
		// First deck
		{{"2C", 2}},
		// Second deck
		{{"2S", 2}},
	}
	game.tableDeck = []Card{}

	// Act
	playMove(&game)

	// Assert
	AssertEqual(t, game.state, Ok)
	AssertEqual(t, len(game.playersDecks[0]), 0)
	AssertEqual(t, len(game.playersDecks[1]), 0)
	AssertEqual(t, len(game.tableDeck), 2)
	AssertEqual(t, game.dummyRound, true)
}

func Test_playMove_ForDummyRound_TakesCardsAndExitsDummyRound(t *testing.T) {
	// Arrange
	game := Game{state: Ok, dummyRound: true}
	game.playersDecks = [][]Card{
		// First deck
		{{"2C", 2}},
		// Second deck
		{{"3C", 3}},
	}
	game.tableDeck = []Card{}

	// Act
	playMove(&game)

	// Assert
	AssertEqual(t, game.state, Ok)
	AssertEqual(t, len(game.playersDecks[0]), 0)
	AssertEqual(t, len(game.playersDecks[1]), 0)
	AssertEqual(t, len(game.tableDeck), 2)
	AssertEqual(t, game.dummyRound, false)
}

func Test_playGame_IntegrationTest_RunsFullGame_ReturnsCorrectWinnerAndNumberOfMoves(t *testing.T) {
	// Arrange
	game := Game{state: Ok, dummyRound: false}
	game.playersDecks = [][]Card{
		// First deck
		{{"2C", 2}, {"10C", 5}, {"5C", 4}, {"7C", 5}},
		// Second deck
		{{"3H", 3}, {"10H", 5}, {"4H", 4}, {"6H", 2}},
	}
	game.tableDeck = []Card{}

	// Act
	playGame(&game)

	// Assert
	AssertEqual(t, game.state, GameCompletedFirstPlayerWon)
	AssertEqual(t, game.roundsCount, 6)
}

func Test_generateSortedDeck_CreatesFullDeck(t *testing.T) {
	// Arrange

	// Act
	deck := generateSortedDeck()

	// Assert
	AssertEqual(t, len(deck), 52)
	AssertEqual(t, deck[0].name, "2H")
	AssertEqual(t, deck[51].name, "AD")
}
