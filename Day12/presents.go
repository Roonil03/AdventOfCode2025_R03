package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Present struct {
	id   int
	size int
}

type RegReq struct {
	id    int
	count int
}

func hashesInString(s string) int {
	count := 0
	for _, ch := range s {
		if ch == '#' {
			count++
		}
	}
	return count
}

func parseInput(filename string) ([]Present, []struct {
	w    int
	h    int
	reqs []RegReq
}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	var p []Present
	var regions []struct {
		w    int
		h    int
		reqs []RegReq
	}
	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			continue
		}
		if strings.Contains(line, ":") && !strings.Contains(line, "x") {
			parts := strings.Split(line, ":")
			id, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			i++
			var shapeStr string
			for i < len(lines) {
				l := strings.TrimSpace(lines[i])
				if l == "" {
					i++
					break
				}
				if strings.Contains(l, ":") {
					break
				}
				shapeStr += l
				i++
			}
			size := hashesInString(shapeStr)
			p = append(p, Present{id: id, size: size})
		} else if strings.Contains(line, "x") && strings.Contains(line, ":") {
			re := regexp.MustCompile(`(\d+)x(\d+): (.*)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 4 {
				w, _ := strconv.Atoi(matches[1])
				h, _ := strconv.Atoi(matches[2])
				countsStr := matches[3]
				var reqs []RegReq
				countParts := strings.Fields(countsStr)
				for i, countStr := range countParts {
					count, _ := strconv.Atoi(countStr)
					reqs = append(reqs, RegReq{id: i, count: count})
				}
				regions = append(regions, struct {
					w    int
					h    int
					reqs []RegReq
				}{w: w, h: h, reqs: reqs})
			}
			i++
		} else {
			i++
		}
	}
	return p, regions, nil
}

func main() {
	filename := "./Day12/test12.txt"
	filename = "./Day12/input12.txt"
	p, regions, err := parseInput(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	res := 0
	for i, r := range regions {
		if fit(p, r) {
			fmt.Printf("Region %d (%dx%d): YES\n", i, r.w, r.h)
			res++
		} else {
			fmt.Printf("Region %d (%dx%d): NO\n", i, r.w, r.h)
		}
	}
	fmt.Println("The answer is:", res)
}

func fit(p []Present, r struct {
	w    int
	h    int
	reqs []RegReq
}) bool {
	a1 := r.w * r.h
	tot := 0
	a2 := 0
	for _, req := range r.reqs {
		tot += req.count
		for _, ps := range p {
			if ps.id == req.id {
				a2 += req.count * ps.size
				break
			}
		}
	}
	if tot*3*3 <= a1 {
		return true
	}
	if a2 > a1 {
		return false
	}
	return true
}
