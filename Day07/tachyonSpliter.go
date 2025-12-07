package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "./Day07/test07.txt"
	filename = "./Day07/input07.txt"
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
	for i := 0; i < c; i++ {
		if grid[0][i] == 'S' {
			grid[1][i] = '|'
		}
	}
	res := 0
	for i := 2; i < r; i++ {
		for j := 0; j < c; j++ {
			if grid[i-1][j] == '|' && grid[i][j] == '^' && j > 0 && j < c-1 {
				grid[i][j-1] = '|'
				grid[i][j+1] = '|'
				res++
			}
			if grid[i-1][j] == '|' && grid[i][j] == '.' {
				grid[i][j] = '|'
			}
		}
		fmt.Println("Line:", i, "\tCurrent:", res)
		// printMaze(grid)
	}
	fmt.Println("The answer is:", res)
}

func printMaze(grid [][]byte) {
	for _, i := range grid {
		for _, j := range i {
			fmt.Printf("%c ", j)
		}
		fmt.Println()
	}
	fmt.Printf("\n\n")
}
