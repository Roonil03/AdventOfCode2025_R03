package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "./Day03/test03.txt"
	filename = "./Day03/input03.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	var res int64 = 0
	for sc.Scan() {
		line := sc.Text()
		jolts := test2(line)
		res += jolts
		fmt.Println("  Jolts:", jolts)
	}
	fmt.Println("The answer is:", res)
}

func test2(s string) int64 {
	n := len(s)
	var res int64 = 0
	for st := 0; st <= n-12; st++ {
		var j int64 = 0
		pos := st
		for d := 0; d < 12; d++ {
			rem := 12 - d - 1
			searchEnd := n - rem
			t := byte('0')
			id := pos
			for i := pos; i < searchEnd; i++ {
				if s[i] > t {
					t = s[i]
					id = i
				}
			}
			j = j*10 + int64(t-'0')
			pos = id + 1
		}
		if j > res {
			res = j
		}
	}
	return res
}
