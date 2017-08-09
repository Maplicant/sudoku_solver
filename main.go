package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// easySudoku := NewSudoku(
	// 	[9][9]int{
	// 		[9]int{0, 0, 0, 0, 8, 0, 7, 0, 0},
	// 		[9]int{0, 7, 8, 3, 0, 0, 5, 2, 0},
	// 		[9]int{4, 1, 5, 0, 0, 0, 3, 8, 0},
	// 		[9]int{0, 0, 0, 6, 0, 5, 0, 1, 0},
	// 		[9]int{8, 0, 0, 0, 0, 0, 0, 0, 7},
	// 		[9]int{0, 3, 0, 8, 0, 2, 0, 0, 0},
	// 		[9]int{0, 9, 7, 0, 0, 0, 4, 3, 2},
	// 		[9]int{0, 2, 3, 0, 0, 4, 6, 9, 0},
	// 		[9]int{0, 0, 4, 0, 9, 0, 0, 0, 0},
	// 	},
	// )

	// hardSudoku := NewSudoku(
	// 	[9][9]int{
	// 		[9]int{3, 0, 5, 0, 0, 0, 8, 0, 6},
	// 		[9]int{0, 7, 0, 0, 0, 0, 0, 2, 0},
	// 		[9]int{0, 1, 0, 3, 0, 6, 0, 9, 0},
	// 		[9]int{0, 0, 9, 0, 7, 0, 3, 0, 0},
	// 		[9]int{0, 0, 0, 9, 8, 5, 0, 0, 0},
	// 		[9]int{0, 0, 0, 0, 4, 0, 0, 0, 0},
	// 		[9]int{7, 0, 0, 0, 0, 0, 0, 0, 4},
	// 		[9]int{0, 8, 2, 0, 0, 0, 7, 6, 0},
	// 		[9]int{9, 0, 3, 0, 0, 0, 2, 0, 8},
	// 	},
	// )

	// // I forgot how to write unit tests and don't have StackOverflow ready because I don't have any internet, this is my temp solution
	// fmt.Println(easySudoku.IsValid(Placement{3, 4, 6})) // Should be invalid
	// fmt.Println(easySudoku.IsValid(Placement{3, 4, 5})) // Should be invalid
	// fmt.Println(easySudoku.IsValid(Placement{3, 4, 3})) // Should be valid
	// fmt.Println(easySudoku.IsValid(Placement{3, 4, 9})) // Should be invalid

	// // Works
	// fmt.Println("Easy:")
	// easySudoku.Solve()
	// easySudoku.Print()

	// // Doesn't work yet
	// fmt.Println("Hard:")
	// hardSudoku.Solve()
	// hardSudoku.Print()
	sudoku := ReadSudoku(os.Stdin)
	sudoku.Print()
	fmt.Println("Solving...")
	sudoku.Solve()
	sudoku.Print()
}

func ReadSudoku(r io.Reader) *Sudoku {
	input := make([]byte, 90)
	offset := 0
	for i := 0; i < 9; i++ {
		amount_read, err := r.Read(input[offset:])
		if err != nil {
			panic(err)
		}
		offset += amount_read
	}

	matrix := [9][9]int{}

	row := 0
	column := 0
	for _, bt := range input {
		num := bt - 48
		if num >= 10 {
			continue
		}
		if num == 35 {
			num = 0
		}
		matrix[row][column] = int(num)
		column += 1
		if column > 8 {
			row += 1
			column = 0
		}
	}

	fmt.Println(matrix)
	return NewSudoku(matrix)
}
