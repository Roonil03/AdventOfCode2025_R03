package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "./Day04/test04.txt"
	filename = "./Day04/input04.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	var grid []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		grid = append(grid, sc.Text())
	}
	rows := len(grid)
	if rows == 0 {
		fmt.Println("The answer is: 0")
		return
	}
	cols := len(grid[0])
	res := 0
	dir := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	for r := range rows {
		for c := range cols {
			if grid[r][c] != '@' {
				continue
			}
			adj := 0
			for _, d := range dir {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
					if grid[nr][nc] == '@' {
						adj++
					}
				}
			}
			if adj < 4 {
				fmt.Println(res+1, ": [x,y]:", r+1, ",", c+1)
				res++
			}
		}
	}
	fmt.Println("The answer is:", res)
}
