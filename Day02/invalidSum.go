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
	// filename := "./Day02/test02.txt"
	filename := "./Day02/input02.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	var input string
	if sc.Scan() {
		input = sc.Text()
	}
	ranges := strings.Split(input, ",")
	sum := 0
	for _, r := range ranges {
		parts := strings.Split(r, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ := strconv.Atoi(parts[1])
		for num := start; num <= end; num++ {
			if test(num) {
				fmt.Print(num, " ")
				sum += num
			}
		}
		fmt.Println()
	}
	fmt.Println("The answer is:", sum)
}

func test(num int) bool {
	s := strconv.Itoa(num)
	n := len(s)
	if n%2 != 0 {
		return false
	}
	m := n / 2
	return s[:m] == s[m:]
}
