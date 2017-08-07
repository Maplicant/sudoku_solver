package main

import (
	"fmt"
)

func main() {
	sudoku := NewSudoku(
		[9][9]int{
			[9]int{0, 0, 0, 0, 8, 0, 7, 0, 0},
			[9]int{0, 7, 8, 3, 0, 0, 5, 2, 0},
			[9]int{4, 1, 5, 0, 0, 0, 3, 8, 0},
			[9]int{0, 0, 0, 6, 0, 5, 0, 1, 0},
			[9]int{8, 0, 0, 0, 0, 0, 0, 0, 7},
			[9]int{0, 3, 0, 8, 0, 2, 0, 0, 0},
			[9]int{0, 9, 7, 0, 0, 0, 4, 3, 2},
			[9]int{0, 2, 3, 0, 0, 4, 6, 9, 0},
			[9]int{0, 0, 4, 0, 9, 0, 0, 0, 0},
		},
	)

	// I forgot how to write unit tests and don't have StackOverflow ready, this is my temp solution
	sudoku.Mark(3, 4)
	sudoku.PrintPossibilities()
	fmt.Println(sudoku.IsValid(Placement{3, 4, 6})) // Should be invalid
	fmt.Println(sudoku.IsValid(Placement{3, 4, 5})) // Should be invalid
	fmt.Println(sudoku.IsValid(Placement{3, 4, 4})) // Should be valid
	fmt.Println(sudoku.IsValid(Placement{3, 4, 9})) // Should be invalid
	sudoku.Print()
}
