package main

import (
	"math/rand"
	"time"
)

// piece generation

type Piece struct {
	Shape [][]int
}

func NewI() *Piece {
	return &Piece{
		Shape: [][]int{
			[]int{1},
			[]int{1},
			[]int{1},
			[]int{1},
		},
	}
}

func NewJ() *Piece {

	return &Piece{
		Shape: [][]int{
			[]int{0, 1},
			[]int{0, 1},
			[]int{1, 1},
		},
	}
}

func NewL() *Piece {

	return &Piece{
		Shape: [][]int{
			[]int{1, 0},
			[]int{1, 0},
			[]int{1, 1},
		},
	}
}

func NewO() *Piece {

	return &Piece{
		Shape: [][]int{
			[]int{1, 1},
			[]int{1, 1},
		},
	}
}

func NewS() *Piece {

	return &Piece{
		Shape: [][]int{
			[]int{1, 0},
			[]int{1, 1},
			[]int{0, 1},
		},
	}
}

func NewT() *Piece {

	return &Piece{
		Shape: [][]int{
			[]int{0, 1, 0},
			[]int{1, 1, 1},
		},
	}
}

func NewZ() *Piece {

	return &Piece{
		Shape: [][]int{
			[]int{0, 1},
			[]int{1, 1},
			[]int{1, 0},
		},
	}
}

func (p *Piece) Rotate() {
	rowCount := len(p.Shape[0])
	newShape := make([][]int, rowCount)
	for rowIdx, _ := range newShape {
		newShape[rowIdx] = make([]int, len(p.Shape))
		for colIdx, _ := range newShape[rowIdx] {
			newShape[rowIdx][colIdx] = p.Shape[colIdx][rowCount-1-rowIdx]
		}
	}
	p.Shape = newShape
}

// returns height, width
func (p *Piece) Dims() (int, int) {
	return len(p.Shape), len(p.Shape[0])
}

var PieceGenerators = []func() *Piece{
	NewI,
	NewJ,
	NewL,
	NewO,
	NewS,
	NewT,
	NewZ,
}

func RandomPiece() *Piece {
	pos := rand.Intn(len(PieceGenerators))
	// fmt.Println(pos)
	gen := PieceGenerators[pos]
	return gen()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
