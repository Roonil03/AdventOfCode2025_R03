package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Interval struct {
	start, end int
}

func main() {
	filename := "./Day05/test05.txt"
	filename = "./Day05/input05.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	type Range struct {
		start, end int
	}
	var ranges []Range
	var id []int
	sc := bufio.NewScanner(file)
	t1 := true
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			t1 = false
			continue
		}
		if t1 {
			parts := strings.Split(line, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			ranges = append(ranges, Range{start, end})
		} else {
			i, _ := strconv.Atoi(line)
			id = append(id, i)
		}
	}
	var it []Interval
	for _, r := range ranges {
		it = append(it, Interval{r.start, r.end})
	}
	m := merging(it)
	res := 0
	for _, i := range m {
		a := i.end - i.start + 1
		res += a
		fmt.Printf("[%d-%d]:  %d \n", i.start, i.end, a)
	}
	fmt.Println("The answer is:", res)
}

func merging(it []Interval) []Interval {
	if len(it) == 0 {
		return it
	}
	sort.Slice(it, func(i, j int) bool {
		return it[i].start < it[j].start
	})
	m := []Interval{it[0]}
	for i := 1; i < len(it); i++ {
		last := &m[len(m)-1]
		cur := it[i]
		if cur.start <= last.end+1 {
			if cur.end > last.end {
				last.end = cur.end
			}
		} else {
			m = append(m, cur)
		}
	}
	return m
}
