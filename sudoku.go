package main

import (
	"fmt"
	"strconv"
)

// Sudoku contains the (un)finished sudoku
type Sudoku struct {
	matrix [9][9]Square
}

// Placement is a struct that contains information about a placement (row, column, digit)
type Placement struct {
	row, column, targetDigit int
}

// Square is the thing a sudoku consists of. There are 81 of them in a sudoku
type Square struct {
	possibilities [9]int
	marked        bool
}

func (s Square) AmountOfPossibilities() int {
	possibilities := 0
	for _, possibility := range s.possibilities {
		if possibility != 0 {
			possibilities++
		}
	}
	return possibilities
}

// 0 means not taken, otherwise a digit
func (s Square) Value() int {
	amountOfDigits := 0
	lastDigit := 0
	for _, digit := range s.possibilities {
		if digit != 0 {
			amountOfDigits++
			lastDigit = digit
		}
	}

	if amountOfDigits == 1 {
		return lastDigit
	} else if amountOfDigits > 1 {
		return 0
	}

	panic("no digits left when trying to get value")
}

func (s Square) ToString() string {
	if s.marked {
		return "*"
	}

	digit := s.Value()

	if digit != 0 {
		return strconv.Itoa(digit)
	}
	return "#"
}

// NewSudoku creates a Sudoku from 2D integer array. 0 means not filled in, 1-9 means a digit.
func NewSudoku(matrix [9][9]int) *Sudoku {
	var sudoku Sudoku

	// These are the placements we're going to make when the sudoku is initialized.
	initialPlacements := make([]Placement, 0)

	// Initialize the sudoku.
	// With that I mean that we fill the sudoku.matrix with unmarked squared that still have all possibilities open
	// (so [1, 2, 3, 4, 5, 6, 7, 8, 9])
	for row := 0; row < 9; row++ {
		for column := 0; column < 9; column++ {
			square := Square{
				possibilities: [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				marked:        false,
			}
			sudoku.matrix[row][column] = square

			// If there's a value in the passed matrix, we add that to the initialPlacements.
			value := matrix[row][column]
			if value != 0 {
				initialPlacements = append(
					initialPlacements,
					Placement{
						row:         row,
						column:      column,
						targetDigit: value,
					},
				)
			}
		}
	}

	// Now we're going to apply the placements.
	for _, placement := range initialPlacements {
		if !sudoku.IsValid(placement) {
			panic("The supplies sudoku isn't valid (has a collision somewhere)")
		}

		sudoku.Apply(placement)
	}

	return &sudoku
}

// Print prints the sudoku in pretty form
func (s *Sudoku) Print() {
	for i, row := range s.matrix {
		if i > 0 && i%3 == 0 { // We print a row divider every 3 rows
			fmt.Println("---+---+---")
		}

		for j, square := range row {
			if j > 0 && j%3 == 0 { // We print a column divider every 3 digits
				fmt.Print("|")
			}
			fmt.Print(square.ToString())
		}
		fmt.Println() // Newline
	}
}

// PrintPossibilities prints the amount of possibilities left on each position
func (s *Sudoku) PrintPossibilities() {
	for i, row := range s.matrix {
		if i > 0 && i%3 == 0 { // We print a row divider every 3 rows
			fmt.Println("---+---+---")
		}

		for j, square := range row {
			if j > 0 && j%3 == 0 { // We print a column divider every 3 digits
				fmt.Print("|")
			}
			fmt.Print(square.AmountOfPossibilities())
		}
		fmt.Println() // Newline
	}
}

// IsPlacementValid tests whether a placement interferes with other digits in the sudoku or not. It's not smart, it just looks at other digits in the row, column and square.
func (s *Sudoku) IsValid(p Placement) bool {
	value := s.matrix[p.row][p.column].Value()
	if value >= 1 && value <= 9 { // Already taken
		return false
	}
	for _, square := range s.matrix[p.row] {
		if square.Value() == p.targetDigit { // Digit already used in row
			return false
		}
	}

	for _, row := range s.matrix {
		if row[p.column].Value() == p.targetDigit { // Digit already used in column
			return false
		}
	}

	// This is the upper left corner of the square the target is in
	// With square I mean the 3x3 square, not the struct with the name `Square`
	squareRow := p.row - (p.row % 3)
	squareColumn := p.column - (p.column % 3)

	for testRow := squareRow; testRow < squareRow+3; testRow++ {
		for testColumn := squareColumn; testColumn < squareColumn+3; testColumn++ {
			if s.matrix[testRow][testColumn].Value() == p.targetDigit { // Target already in square
				return false
			}
		}
	}

	return true
}

// Mark marks a place on the sudoku so you can visualize a point on the sudoku when you call sudoku.Print()
func (s *Sudoku) Mark(row, column int) {
	s.matrix[row][column].marked = true
}

// Unmark removes the marking from a square in the sudoku
func (s *Sudoku) Unmark(row, column int) {
	s.matrix[row][column].marked = false
}

// Apply applies a placement. It assumes that the supplied placement is valid
func (s *Sudoku) Apply(placement Placement) {
	s.matrix[placement.row][placement.column].possibilities = [9]int{placement.targetDigit}
	// Remove the possibility of the target number from all digits in the row
	for columnIndex := 0; columnIndex < 9; columnIndex++ {
		if columnIndex == placement.column {
			continue
		}

		possibilities := &s.matrix[placement.row][columnIndex].possibilities
		for i, possibility := range possibilities {
			if possibility == placement.targetDigit {
				possibilities[i] = 0
				break
			}
		}
	}

	// Same here, but now for all squares in the column
	for rowIndex := 0; rowIndex < 9; rowIndex++ {
		if rowIndex == placement.row {
			continue
		}

		possibilities := &s.matrix[rowIndex][placement.column].possibilities
		for i, possibility := range possibilities {
			if possibility == placement.targetDigit {
				possibilities[i] = 0
				break
			}
		}
	}

	// Same here, but now for all the squares in the current 3x3 square
	squareRow := placement.row - (placement.row % 3)
	squareColumn := placement.column - (placement.column % 3)

	for row := squareRow; row < squareRow+3; row++ {
		for column := squareColumn; column < squareColumn+3; column++ {
			if row == placement.row && column == placement.column {
				continue
			}

			possibilities := &s.matrix[row][column].possibilities
			for i, possibility := range possibilities {
				if possibility == placement.targetDigit {
					possibilities[i] = 0
					break
				}
			}
		}
	}
}
