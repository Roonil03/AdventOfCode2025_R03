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
	filename := "./Day01/test01.txt"
	// filename := "./Day01/input01.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	res := 0
	pos := 50
	i := 1
	mod := 100
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
		// for a := 0; a < val; a++ {
		// 	if dir == 'L' {
		// 		pos--
		// 		if pos < 0 {
		// 			pos = 99
		// 			res++
		// 		}
		// 	} else {
		// 		pos++
		// 		if pos == 100 {
		// 			pos = 0
		// 			res++
		// 		}
		// 	}
		// 	// fmt.Printf("  Iteration %d.%d: Position %d\tDirection: %c\tValue: %d\n", i, a, pos, dir, val)
		// }
		if dir == 'L' {
			check := pos
			if pos == 0 {
				check = mod
			}
			if val >= check {
				res++
				res += (val - check) / mod
			}
			pos = (pos - val) % mod
			if pos < 0 {
				pos += mod
			}
		} else {
			tot := pos + val
			temp := tot / mod
			res += temp
			pos = tot % mod
		}
		fmt.Printf("Iteration %d: Position %d\tDirection: %c\tValue: %d\n", i, pos, dir, val)
		i++
	}
	fmt.Println("The answer is:", res)
	if err := sc.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}
}
