package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "./Day07/test07.txt"
	// filename = "./Day07/input07.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	var grid [][]byte
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		grid = append(grid, []byte(line))
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	r, c := len(grid), len(grid[0])
	st := -1
	for i := 0; i < c; i++ {
		if grid[0][i] == 'S' {
			st = i
			break
		}
	}
	memo := make(map[string]int64)
	var test func(row, col int) int64
	test = func(row, col int) int64 {
		if row >= r {
			return 1
		}
		if col < 0 || col >= c {
			return 0
		}
		key := fmt.Sprintf("%d,%d", row, col)
		if val, exists := memo[key]; exists {
			return val
		}
		var temp int64 = 0
		if grid[row][col] == '.' || grid[row][col] == 'S' || grid[row][col] == '|' {
			temp = test(row+1, col)
			fmt.Printf("Current Paths: T:%d\n", temp)
		} else if grid[row][col] == '^' {
			l1 := test(row+1, col-1)
			r1 := test(row+1, col+1)
			fmt.Printf("Current Paths: L:%d\tR:%d\n", l1, r1)
			temp = l1 + r1
		}
		memo[key] = temp
		return temp
	}
	var res int64 = test(1, st)
	fmt.Println("The answer is:", res)
}
