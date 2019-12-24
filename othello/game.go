package othello

import (
	"fmt"

	"github.com/unathi-skosana/gothello/gomcts"
)

// Pieces
const (
	EMPTY = iota
	BLACK
	WHITE
	OUTER
)

// DIR
const DIR = 8

// sizes
const BOARD_SIZE = 100
const BOARD_WIDTH = 8
const PIECE_SLOTS = 64

// BLOCKS
const FIRST_BLOCK = 11
const LAST_BLOCK = 88

// NO MORE MOVES
const NO_MOVES = -1

// DIRECTIONS
var ALLDIRECTIONS = []int{-11, -10, -9, -1, 1, 9, 10, 11}

// OthelloBoardGameAction - action on a othello board
type OthelloBoardGameAction struct {
	move  int
	value int
}

// OthelloGameState - Othello game state
type OthelloGameState struct {
	nextToMove   int
	board        []int
	emptySquares uint16
	ended        bool
	result       gomcts.GameResult
}

// PrintBoard - prints out game board to console
func (s OthelloGameState) PrintBoard() {
	board := s.board
	fmt.Printf("   1 2 3 4 5 6 7 8 [%s=%d %s=%d]\n", nameof(BLACK), count(board, BLACK), nameof(WHITE), count(board, WHITE))
	for row := 1; row <= BOARD_WIDTH; row++ {
		fmt.Printf("%d  ", row)
		for col := 1; col <= BOARD_WIDTH; col++ {
			fmt.Printf("%s ", nameof(board[col+(10*row)]))
		}
		fmt.Printf("\n")
	}
}

/*************GameState interface methods **************/

// CreateOthelloInitialGameState - initializes a othello game state
func New() OthelloGameState {
	board := initializeOthelloBoard()
	state := OthelloGameState{nextToMove: BLACK, board: board, emptySquares: uint16(PIECE_SLOTS) - 4}
	return state
}

// IsGameEnded - OthelloGameState implementation of IsGameEnded method of GameState interface
func (s OthelloGameState) IsGameEnded() bool {
	_, ended := s.EvaluateGame()
	return ended
}

// EvaluateGame - OthelloGameState implementation of EvaluateGame method of GameState interface
func (s OthelloGameState) EvaluateGame() (result gomcts.GameResult, ended bool) {

	defer func() {
		s.result = result
		s.ended = ended
	}()

	if s.ended {
		return s.result, s.ended
	}

	whiteSum := 0
	blackSum := 0
	board := s.board
	nextToMove := s.nextToMove

	for i := FIRST_BLOCK; i <= LAST_BLOCK; i++ {
		if board[i] == BLACK {
			blackSum++
		} else if board[i] == WHITE {
			whiteSum++
		}
	}

	player_moves := numLegalActions(board, nextToMove)
	other_moves := numLegalActions(board, opponent(nextToMove))

	if s.emptySquares == 0 || player_moves == 0 && other_moves == 0 {
		if blackSum > whiteSum {
			return gomcts.GameResult(BLACK), true
		} else if blackSum < whiteSum {
			return gomcts.GameResult(WHITE), true
		}
		return gomcts.GameResult(0), true
	}

	return gomcts.GameResult(0), false
}

// ApplyTo - OthelloBoardGameAction implementation of ApplyTo method of Action interface
func (a OthelloBoardGameAction) ApplyTo(s gomcts.GameState) gomcts.GameState {
	OthelloGameState := s.(OthelloGameState)
	OthelloGameState.board = copyOthelloBoard(OthelloGameState.board)

	if a.move == NO_MOVES {
		OthelloGameState.nextToMove = opponent(OthelloGameState.nextToMove)
		return OthelloGameState
	}

	if OthelloGameState.nextToMove != a.value {
		panic("*hands slapped*,  not your turn")
	}

	makeMove(a.move, OthelloGameState)

	OthelloGameState.nextToMove = opponent(OthelloGameState.nextToMove)

	OthelloGameState.emptySquares--

	return OthelloGameState
}

// GetLegalActions - OthelloGameState implementation of GetLegalAction method of GameState interface
func (s OthelloGameState) GetLegalActions() []gomcts.Action {
	cnt := 0
	actions := make([]gomcts.Action, 0, 0)
	board := s.board
	nextToMove := s.nextToMove

	for i := FIRST_BLOCK; i <= LAST_BLOCK; i++ {
		if legalPlayer(board, i, nextToMove) == 1 {
			cnt++
			actions = append(actions, OthelloBoardGameAction{move: i, value: s.nextToMove})
		}
	}

	if cnt == 0 {
		actions = append(actions, OthelloBoardGameAction{move: NO_MOVES, value: s.nextToMove})
	}

	return actions
}

// NextToMove - OthelloGameState implementation of NextToMove method of GameState interface
func (s OthelloGameState) NextToMove() int {
	return s.nextToMove
}
