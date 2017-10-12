package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/term"
)

func main() {
	g := NewGame()
	g.Run()
}

type Game struct {
	Heartbeat int     // time between downshifts in ms
	Board     [][]int // row * col starting from top left

	Piece    *Piece
	PieceRow int // combine X
	PieceCol int // combine Y

	endLowerRoutine chan bool
}

func NewGame() *Game {
	screen := newBoard()

	return &Game{
		Heartbeat:       1000,
		Board:           screen,
		endLowerRoutine: make(chan bool, 1),
	}
}

func (g *Game) Run() {
	moves := make(chan Move)
	go notifyKeystrokes(moves)
	ticker := time.NewTicker(1 * 1000 * 1000 * 1000)

	rowsDeleted := uint(0)
	for !g.IsDone() {
		g.Piece = RandomPiece()
		g.PieceCol = 4
		g.PieceRow = 3

		ticker.Stop()
		ticker := time.NewTicker(1 * 1000 * 1000 * 1000 / (1 << (rowsDeleted / 5)))

		for !g.PieceCollided() {
			g.Print()
			select {
			case <-ticker.C:
				g.PieceRow += 1
			case move := <-moves:
				switch move {
				case Left:
					g.PieceCol -= 1
					if g.PieceCollided() {
						g.PieceCol += 1
					}
				case Right:
					g.PieceCol += 1
					if g.PieceCollided() {
						g.PieceCol -= 1
					}
				case Rotate:
					g.Piece.Rotate()
					if g.PieceCollided() {
						// put it back hack
						g.Piece.Rotate()
						g.Piece.Rotate()
						g.Piece.Rotate()
					}
				case Down:
					g.PieceRow += 1
				}
			}
		}
		g.PieceRow -= 1
		var deletes int
		g.Board, deletes = SumBoards(g.Board, g.PieceBoard(), 1, 1).CleanBoard()
		rowsDeleted += uint(deletes)
	}
	fmt.Println("You died!")
}

func (g *Game) PieceCollided() bool {
	height, _ := g.Piece.Dims()
	if g.PieceRow+height-1 == len(g.Board) {
		return true
	}

	board := SumBoards(g.Board, g.PieceBoard(), 1, 1)
	for _, row := range board {
		for _, i := range row {
			if i == 2 {
				return true
			}
		}
	}
	return false
}

func (g *Game) IsDone() bool {
	for _, i := range g.Board[4] {
		if i != 0 {
			spew.Println(g.Board)
			return true
		}
	}
	return false
}

func (g *Game) CheckPieceAndShift() {
	height, width := g.Piece.Dims()
	if g.PieceCol+width > len(g.Board[0]) {
		g.PieceCol = len(g.Board[0]) - width
	}

	if g.PieceCol < 0 {
		g.PieceCol = 0
	}

	if g.PieceRow+height > len(g.Board) {
		g.PieceRow = len(g.Board) - height
	}
}

type Board [][]int

func newBoard() Board {
	screen := make([][]int, 24)
	for idx, _ := range screen {
		screen[idx] = make([]int, 10)
	}
	return screen
}

func (g *Game) PieceBoard() Board {
	g.CheckPieceAndShift()
	screen := newBoard()

	height, width := g.Piece.Dims()
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen[g.PieceRow+r][g.PieceCol+c] = g.Piece.Shape[r][c]
		}
	}

	return screen
}

func (g *Game) Print() {
	outboard := g.Board
	if g.Piece != nil {
		pieceboard := g.PieceBoard()
		outboard = SumBoards(outboard, pieceboard, 1, 2)
	}
	PrintBoard(outboard[4:])
}

func (g *Game) isMoveValid(m Move) {
	switch m {
	case Left:
		g.PieceCol -= 1
		defer func() { g.PieceCol += 1 }()
	}
}

// returns cleaned board and number of rows deleted
func (in Board) CleanBoard() (Board, int) {
	out := newBoard()

	allOnes := func(row []int) bool {
		for _, i := range row {
			if i != 1 {
				return false
			}
		}
		return true
	}

	outIdx := len(out) - 1
	var inIdx int
	for inIdx = len(in) - 1; inIdx >= 0; inIdx-- {
		if allOnes(in[inIdx]) {
			continue
		}
		copy(out[outIdx], in[inIdx])
		outIdx--
	}
	return out, outIdx - inIdx
}

func SumBoards(a, b Board, aMult, bMult int) Board {
	screen := newBoard()

	for r, row := range screen {
		for c, _ := range row {
			screen[r][c] = a[r][c]*aMult + b[r][c]*bMult
		}
	}

	return screen
}

func PrintBoard(b Board) {
	term, err := term.Open("/dev/tty")
	if err != nil {
		panic(err)
	}
	term.Write([]byte("\033[2J"))
	for _, row := range b {
		fmt.Print("|")
		for _, e := range row {
			if e == 0 {
				fmt.Print(" ")
			} else if e == 1 {
				fmt.Print("#")
			} else {
				fmt.Print("X")
			}
		}
		fmt.Println("|")
	}
	fmt.Println("------------")
	fmt.Println("ctrl-c to quit. Arrows to move/rotate.")
}

type Move int

const (
	Left Move = iota
	Right
	Rotate
	Down
)

func notifyKeystrokes(c chan Move) {
	term, err := term.Open("/dev/tty")
	if err != nil {
		panic(err)
	}
	err = term.SetCbreak()
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 4)
	for true {
		n, err := term.Read(buffer)
		if err != nil {
			panic(err)
		}

		if n == 3 && buffer[0] == 27 && buffer[1] == 91 {
			switch buffer[2] {
			case 68:
				c <- Left
			case 65:
				c <- Rotate
			case 67:
				c <- Right
			case 66:
				c <- Down
			}
		}
		if n == 1 && buffer[0] == 3 {
			os.Exit(0)
		}
	}
}
