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
	// filename := "./Day01/test01.txt"
	filename := "./Day01/input01.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	res := 0
	pos := 50
	mod := 100
	i := 1
	defer file.Close()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if len(line) < 2 {
			continue
		}
		dir := line[0]
		numStr := line[1:]
		val, err := strconv.Atoi(numStr)
		if err != nil {
			log.Printf("Skipping invalid line %q: %v", line, err)
			continue
		}
		if dir == 'L' {
			pos = (pos - val)
			for pos < 0 {
				pos += 100
			}
		} else {
			pos = (pos + val) % mod
		}
		if pos == 0 {
			res++
		}
		fmt.Printf("Iteration %d: Position %d\tDirection: %c\tValue: %d\n", i, pos, dir, val)
		i++
	}
	fmt.Println("The answer is:", res)
	if err := sc.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
}
