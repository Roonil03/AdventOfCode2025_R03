package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

func main() {
	filename := "./Day09/test09.txt"
	filename = "./Day09/input09.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	var pts []Point
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		pts = append(pts, Point{x, y})
	}
	n := len(pts)
	res := 0
	var t1, t2 Point
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			w := int(math.Abs(float64(pts[i].x - pts[j].x + 1)))
			h := int(math.Abs(float64(pts[i].y - pts[j].y + 1)))
			temp := w * h
			if temp > res {
				res = temp
				t1 = pts[i]
				t2 = pts[j]
				fmt.Printf("\tPoints [%d, %d] and [%d, %d]\t Area: %d\n", pts[i].x, pts[i].y, pts[j].x, pts[j].y, temp)
			}
			if temp < 0 {
				fmt.Println("Error, overflow detected!")
				break
			}
		}
	}
	fmt.Printf("\nPoints: [%d,%d] and [%d,%d]\n", t1.x, t1.y, t2.x, t2.y)
	fmt.Println("The answer is:", res)
}
