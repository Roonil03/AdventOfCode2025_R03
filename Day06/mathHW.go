package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	filename := "./Day06/test06.txt"
	filename = "./Day06/input06.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	var grid [][]int
	var operators []string
	for sc.Scan() {
		line := sc.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0] == "*" || fields[0] == "+" {
			operators = fields
		} else {
			var row []int
			for _, field := range fields {
				num, err := strconv.Atoi(field)
				if err == nil {
					row = append(row, num)
				}
			}
			if len(row) > 0 {
				grid = append(grid, row)
			}
		}
	}
	rows := len(grid)
	cols := len(grid[0])
	var res int64 = 0
	for c := range cols {
		op := operators[c]
		switch op {
		case "*":
			var temp int64 = 1
			for r := range rows {
				temp *= int64(grid[r][c])
			}
			fmt.Printf("Multiplication %d: %d\n", c+1, temp)
			res += temp
		case "+":
			var temp int64 = 0
			for r := range rows {
				temp += int64(grid[r][c])
			}
			fmt.Printf("Add %d: %d\n", c+1, temp)
			res += temp
		}
	}
	fmt.Println("The answer is:", res)
}
