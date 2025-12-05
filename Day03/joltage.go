package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "./Day03/test03.txt"
	// filename = "./Day03/input03.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	res := 0
	for sc.Scan() {
		line := sc.Text()
		jolts := test1(line)
		res += jolts
		fmt.Println("  Jolts:", jolts)
	}
	fmt.Println("The answer is:", res)
}

func test1(s string) int {
	n := len(s)
	suf := make([]byte, n)
	suf[n-1] = s[n-1]
	for i := n - 2; i >= 0; i-- {
		suf[i] = max(s[i], suf[i+1])
	}
	res := 0
	for i := 0; i < n-1; i++ {
		a := int(s[i] - '0')
		b := int(suf[i+1] - '0')
		jolts := a*10 + b
		if jolts > res {
			res = jolts
		}
	}
	return res
}
