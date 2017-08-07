package main

import (
	"fmt"
)

// Sudoku contains the (un)finished sudoku
type Sudoku struct {
	matrix [9][9]int // 0 means not filled in, 10 means marked, 1-9 is a digit
}

// Print prints the sudoku in pretty form
func (s *Sudoku) Print() {
	for i, row := range s.matrix {
		if i > 0 && i%3 == 0 { // We print a row divider every 3 rows
			fmt.Println("---+---+---")
		}

		for j, digit := range row {
			if j > 0 && j%3 == 0 { // We print a column divider every 3 digits
				fmt.Print("|")
			}

			switch digit {
			case 0: // not filled in yet
				fmt.Print("#")
			case 10: // marking
				fmt.Print("*")
			default:
				fmt.Print(digit)
			}
		}
		fmt.Println() // Newline
	}
}

// IsPlacementValid tests whether a placement interferes with other digits in the sudoku or not. It's not smart, it just looks at other digits in the row, column and square.
func (s *Sudoku) IsPlacementValid(rowIndex, columnIndex, targetDigit int) bool {
	currentValue := s.matrix[rowIndex][columnIndex]
	if currentValue >= 1 && currentValue <= 9 { // Already taken
		return false
	}
	for _, valueDigit := range s.matrix[rowIndex] {
		if valueDigit == targetDigit { // Digit already used in row
			return false
		}
	}

	for _, row := range s.matrix {
		if row[columnIndex] == targetDigit { // Digit already used in column
			return false
		}
	}

	// This is the upper left corner of the square the target is in
	squareRow := rowIndex - (rowIndex % 3)
	squareColumn := columnIndex - (columnIndex % 3)

	for testRow := squareRow; testRow < squareRow+3; testRow++ {
		for testColumn := squareColumn; testColumn < squareColumn+3; testColumn++ {
			if s.matrix[testRow][testColumn] == targetDigit { // Target already in square
				return false
			}
		}
	}

	return true
}

// Mark marks a place on the sudoku so you can visualize a point on the sudoku when you call sudoku.Print()
func (s *Sudoku) Mark(row, column int) {
	s.matrix[row][column] = 10
}
