package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	sudoku := ReadSudoku(os.Stdin)
	sudoku.Print()
	fmt.Println("")
	fmt.Println("Solving...")
	fmt.Println("")
	begin := time.Now()
	sudoku.Solve()
	timediff := time.Now().Sub(begin)
	sudoku.Print()
	fmt.Printf("\nSolved in %dms\n", timediff.Nanoseconds()/1000000)
}

// ReadSudoku reads a sudoku from `r`
func ReadSudoku(r io.Reader) *Sudoku {
	input := make([]byte, 200)
	nums := make([]int, 0)
	offset := 0
outer:
	for {
		amount_read, err := r.Read(input[offset:])
		if err == io.EOF {
			// Only continue if the user entered numbers in the line
			for _, bt := range input[offset:] {
				num := bt - 48
				if num >= 10 {
					continue
				}
				if num == 35 {
					num = 0
				}
				nums = append(nums, int(num))
			}
			offset += amount_read

			if len(nums) > 81 {
				break outer
			} else {
				panic(fmt.Sprintf("EOF but not enough numbers, got", len(nums)))
			}
		}
		if err != nil {
			panic(err)
		}

		// Only continue if the user entered numbers in the line
		for _, bt := range input[offset:] {
			num := bt - 48
			if bt == 35 {
				num = 0
			}
			if num >= 10 {
				continue
			}
			nums = append(nums, int(num))
		}
		offset += amount_read

		if len(nums) >= 81 {
			break outer
		}
	}

	matrix := [9][9]int{}

	row := 0
	column := 0
	for _, num := range nums {
		matrix[row][column] = num
		column += 1
		if column > 8 {
			row += 1
			column = 0
		}
	}

	return NewSudoku(matrix)
}
