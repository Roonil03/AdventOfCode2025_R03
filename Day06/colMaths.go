package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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
	var grid [][]byte
	l := 0
	for sc.Scan() {
		line := sc.Text()
		row := []byte(line)
		if len(row) > l {
			l = len(row)
		}
		grid = append(grid, row)
	}
	for i := range grid {
		for len(grid[i]) < l {
			grid[i] = append(grid[i], ' ')
		}
	}
	rows := len(grid)
	cols := l
	var res int64 = 0
	vis := make([]bool, cols)
	for c := cols - 1; c >= 0; c-- {
		if vis[c] {
			continue
		}
		hasDigit := false
		for r := 0; r < rows-1; r++ {
			if c < len(grid[r]) && grid[r][c] >= '0' && grid[r][c] <= '9' {
				hasDigit = true
				break
			}
		}
		if !hasDigit {
			continue
		}
		start := c
		for start > 0 {
			t1 := true
			for r := 0; r < rows-1; r++ {
				if grid[r][start-1] != ' ' {
					t1 = false
					break
				}
			}
			if t1 {
				break
			}
			start--
		}
		var nums []int64
		var op string
		for col := start; col <= c; col++ {
			var numStr string
			for r := 0; r < rows-1; r++ {
				if grid[r][col] >= '0' && grid[r][col] <= '9' {
					numStr += string(grid[r][col])
				}
			}
			if numStr != "" {
				num, _ := strconv.ParseInt(numStr, 10, 64)
				nums = append(nums, num)
			}
			if grid[rows-1][col] == '*' || grid[rows-1][col] == '+' {
				op = string(grid[rows-1][col])
			}
			vis[col] = true
		}
		switch op {
		case "*":
			var temp int64 = 1
			for _, n := range nums {
				temp *= n
				fmt.Printf("%d\t", temp)
			}
			fmt.Printf("Multiplication: %d\n", temp)
			res += temp
		case "+":
			var temp int64 = 0
			for _, n := range nums {
				temp += n
				fmt.Printf("%d\t", temp)
			}
			fmt.Printf("Add: %d\n", temp)
			res += temp
		}
	}
	fmt.Println("The answer is:", res)
}

func extract(grid [][]byte, col int) ([]int64, string) {
	rows := len(grid)
	if col >= len(grid[0]) {
		return nil, ""
	}
	var dig []byte
	var op string
	for r := 0; r < rows; r++ {
		ch := grid[r][col]
		if ch == '*' || ch == '+' {
			op = string(ch)
		} else if ch >= '0' && ch <= '9' {
			dig = append(dig, ch)
		}
	}
	if len(dig) == 0 {
		return nil, op
	}
	var nums []int64
	cur := col
	for cur >= 0 {
		var d []byte
		for r := 0; r < rows-1; r++ {
			if cur < len(grid[r]) {
				ch := grid[r][cur]
				if ch >= '0' && ch <= '9' {
					d = append(d, ch)
				} else if ch != ' ' {
					break
				}
			}
		}
		if len(d) > 0 {
			numStr := string(d)
			num, _ := strconv.ParseInt(numStr, 10, 64)
			nums = append(nums, num)
		}
		if cur == 0 {
			break
		}
		t1 := true
		for r := 0; r < rows-1; r++ {
			if cur-1 < len(grid[r]) && grid[r][cur-1] != ' ' {
				t1 = false
				break
			}
		}
		if t1 {
			break
		}
		cur--
	}
	return nums, op
}
