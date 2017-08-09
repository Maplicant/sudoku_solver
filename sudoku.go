package main

import (
	"errors"
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

// Deduct tries to find values of squares purely by deducting from other squares
func (s *Sudoku) Deduct() error {
	for rowIndex := 0; rowIndex < 9; rowIndex++ {
		for columnIndex := 0; columnIndex < 9; columnIndex++ {
			square := s.matrix[rowIndex][columnIndex]
			n, err := square.Value()
			if err != nil {
				return err
			}
			if n != 0 {
				continue
			}
		possibilities:
			for _, possibility := range square.possibilities {
				// Check the row
			outerrow:
				for {
					for otherColumnIndex, otherSquare := range s.matrix[rowIndex] {
						if otherColumnIndex == columnIndex {
							continue
						}

						if Contains(otherSquare.possibilities, possibility) {
							break outerrow
						}
					}
					placement := Placement{
						row:         rowIndex,
						column:      columnIndex,
						targetDigit: possibility,
					}
					(*s).Apply(placement)
					break possibilities
				}

				// Check the column
			outercolumn:
				for {
					for otherRowIndex, otherRow := range s.matrix {
						if otherRowIndex == rowIndex {
							continue
						}
						if Contains(otherRow[columnIndex].possibilities, possibility) {
							break outercolumn
						}
					}
					placement := Placement{
						row:         rowIndex,
						column:      columnIndex,
						targetDigit: possibility,
					}
					(*s).Apply(placement)
					break possibilities
				}

				// Check the 3x3 square
			outersquare:
				for {
					squareRow := rowIndex - (rowIndex % 3)
					squareColumn := columnIndex - (columnIndex % 3)

					for otherRowIndex := squareRow; otherRowIndex < squareRow+3; otherRowIndex++ {
						for otherColumnIndex := squareColumn; otherColumnIndex < squareColumn+3; otherColumnIndex++ {
							if otherRowIndex == rowIndex && otherColumnIndex == columnIndex {
								continue
							}
							if Contains(s.matrix[otherRowIndex][otherColumnIndex].possibilities, possibility) { // Target already in square
								break outersquare
							}
						}
					}

					placement := Placement{
						row:         rowIndex,
						column:      columnIndex,
						targetDigit: possibility,
					}
					(*s).Apply(placement)
					break possibilities
				}
			}
		}
	}

	return nil
}

// Solve solves a sudoku using deducting as much as possible and trying different placements when needed
func (s *Sudoku) Solve() {
	for {
		// Deduct as much as we can
		for {
			oldsudoku := *s
			err := s.Deduct()
			if err != nil {
				fmt.Println("what is going on while deducting in Solve")
				s.Print()
				panic("")
			}
			if oldsudoku == *s {
				break
			}
		}
		fmt.Println("Done deducting")
		s.Print()

		if s.IsFinished() {
			break
		}

		rowToTry, columnToTry := s.FindSquareToTry()
		sudokus := make(chan Sudoku, 2)
		for _, possibility := range s.matrix[rowToTry][columnToTry].possibilities {
			if possibility < 1 {
				continue
			}
			placement := Placement{
				rowToTry, columnToTry, possibility,
			}
			go s.Try(placement, sudokus)
		}

		*s = <-sudokus
	}
}

// Try tries a placement to see whether it collides when you try to deduct it further.
func (ins *Sudoku) Try(placement Placement, solutions chan Sudoku) {
	// TODO: Think forward more
	s := *ins
	s.Apply(placement)
	for {
		oldsudoku := s
		err := s.Deduct()
		if err != nil {
			return
		}
		if s.IsFinished() {
			solutions <- s
			return
		}
		if oldsudoku == s {
			break
		}
	}
	solutions <- s
}

// FindSquareToTry finds the coordinates of the square with the least amount of possibilities to try.
func (s *Sudoku) FindSquareToTry() (int, int) {
	rowIndex, columnIndex := 0, 0
	lowestPossibilities := 10

	for potentialRowIndex, row := range s.matrix {
		for potentialColumnIndex, square := range row {
			possibilities := square.AmountOfPossibilities()
			if possibilities != 1 && possibilities < lowestPossibilities {
				lowestPossibilities = possibilities
				rowIndex, columnIndex = potentialRowIndex, potentialColumnIndex
			}
		}
	}

	return rowIndex, columnIndex
}

// Contains checks whether an array contains a certain integer.
// I'm sure there's something in the standard library, but I don't have any documentation at the moment
func Contains(arr [9]int, n int) bool {
	for _, i := range arr {
		if i == n {
			return true
		}
	}
	return false
}

// IsFinished checks whether a sudoku is finished
func (s *Sudoku) IsFinished() bool {
	for _, row := range s.matrix {
		for _, square := range row {
			n, err := square.Value()
			if err != nil {
				panic(err)
			}
			if n == 0 {
				return false
			}
		}
	}
	return true
}

// Amount of possibilities checks how many possibilities there are left for a certain square
func (s Square) AmountOfPossibilities() int {
	possibilities := 0
	for _, possibility := range s.possibilities {
		if possibility != 0 {
			possibilities++
		}
	}
	return possibilities
}

// Value gives you the value of a square. 0 means not taken, otherwise you get a digit
func (s Square) Value() (int, error) {
	amountOfDigits := 0
	lastDigit := 0
	for _, digit := range s.possibilities {
		if digit != 0 {
			amountOfDigits++
			lastDigit = digit
		}
	}

	if amountOfDigits == 1 {
		return lastDigit, nil
	} else if amountOfDigits > 1 {
		return 0, nil
	}

	return 0, errors.New(fmt.Sprintf("no digits left when trying to get value", s))
}

// ToString converts a square to a string (# for not taken, 1-9 for a digit, * for a marked square)
func (s Square) ToString() string {
	if s.marked {
		return "*"
	}

	digit, err := s.Value()
	if err != nil {
		panic(err)
	}

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
		sudoku.PrintPossibilities()
		for ri, row := range sudoku.matrix {
			for ci, square := range row {
				fmt.Println(ri, ci, square.possibilities)
			}
		}
		fmt.Println()
	}

	return &sudoku
}

// Print prints the sudoku in pretty form
func (s *Sudoku) Print() {
	// fmt.Println(s.matrix)
	for i, row := range s.matrix {
		if i > 0 && i%3 == 0 { // We print a row divider every 3 rows
			fmt.Println("---+---+---")
		}

		for j, square := range row {
			if j > 0 && j%3 == 0 { // We print a column divider every 3 digits
				fmt.Print("|")
			}
			// fmt.Println("\n", i, j)
			fmt.Print(square.ToString())
			// fmt.Println("done")
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

// IsValid tests whether a placement interferes with other digits in the sudoku or not. It's not smart, it just looks at other digits in the row, column and square.
func (s *Sudoku) IsValid(p Placement) bool {
	value, err := s.matrix[p.row][p.column].Value()
	if err != nil {
		panic(err)
	}
	if value >= 1 && value <= 9 && !(value == p.targetDigit) { // Already taken with a number that's not the target digit
		return false
	}
	for columnIndex, square := range s.matrix[p.row] {
		if columnIndex == p.column {
			continue
		}
		n, err := square.Value()
		if err != nil {
			panic(err)
		}
		if n == p.targetDigit { // Digit already used in row
			return false
		}
	}

	for rowIndex, row := range s.matrix {
		if rowIndex == p.row {
			continue
		}
		n, err := row[p.column].Value()
		if err != nil {
			panic(err)
		}
		if n == p.targetDigit { // Digit already used in column
			return false
		}
	}

	// This is the upper left corner of the square the target is in
	// With square I mean the 3x3 square, not the struct with the name `Square`
	squareRow := p.row - (p.row % 3)
	squareColumn := p.column - (p.column % 3)

	for testRow := squareRow; testRow < squareRow+3; testRow++ {
		for testColumn := squareColumn; testColumn < squareColumn+3; testColumn++ {
			if testRow == p.row && testColumn == p.column {
				continue
			}
			n, err := s.matrix[testRow][testColumn].Value()
			if err != nil {
				panic(err)
			}
			if n == p.targetDigit { // Target already in square
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
