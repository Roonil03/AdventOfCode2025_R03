package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Point struct {
	x, y int
}

type RectangleResult struct {
	area int64
	p1   Point
	p2   Point
}

type Segment struct {
	x1, y1, x2, y2 int
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
	fmt.Printf("Points loaded: %d\n", n)
	minX, maxX := pts[0].x, pts[0].x
	minY, maxY := pts[0].y, pts[0].y
	for i := 1; i < n; i++ {
		if pts[i].x < minX {
			minX = pts[i].x
		}
		if pts[i].x > maxX {
			maxX = pts[i].x
		}
		if pts[i].y < minY {
			minY = pts[i].y
		}
		if pts[i].y > maxY {
			maxY = pts[i].y
		}
	}
	fmt.Printf("Bounds: x[%d,%d] y[%d,%d]\n", minX, maxX, minY, maxY)
	fmt.Printf("Grid span: %d x %d\n", maxX-minX, maxY-minY)
	segs := make([]Segment, n)
	for i := 0; i < n; i++ {
		next := (i + 1) % n
		segs[i] = Segment{
			x1: pts[i].x,
			y1: pts[i].y,
			x2: pts[next].x,
			y2: pts[next].y,
		}
	}
	red := make(map[[2]int]bool)
	for _, seg := range segs {
		if seg.x1 == seg.x2 {
			s, e := seg.y1, seg.y2
			if s > e {
				s, e = e, s
			}
			for y := s; y <= e; y++ {
				red[[2]int{seg.x1, y}] = true
			}
		} else if seg.y1 == seg.y2 {
			s, e := seg.x1, seg.x2
			if s > e {
				s, e = e, s
			}
			for x := s; x <= e; x++ {
				red[[2]int{x, seg.y1}] = true
			}
		}
	}
	fmt.Printf("Red tiles: %d\n", len(red))
	workers := runtime.NumCPU()
	if workers > 16 {
		workers = 16
	}
	fmt.Printf("Starting rectangle checks with %d workers...\n", workers)
	res := checkRectanglesOptimized(pts, segs, red, n, workers, minX, maxX, minY, maxY)
	fmt.Printf("Rect vertices: [%d,%d] and [%d,%d]\n", res.p1.x, res.p1.y, res.p2.x, res.p2.y)
	fmt.Printf("The answer is: %d\n", res.area)
}

func checkRectanglesOptimized(pts []Point, segs []Segment, red map[[2]int]bool,
	n, workers, minX, maxX, minY, maxY int) RectangleResult {
	type job struct {
		i, j int
	}
	tot := n * (n - 1) / 2
	job_chan := make(chan job, 10000)
	res_chan := make(chan RectangleResult, 1000)
	var wg sync.WaitGroup
	var proc int64
	var last int64
	go func() {
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				job_chan <- job{i, j}
			}
		}
		close(job_chan)
		fmt.Printf("Job generation complete: %d jobs\n", tot)
	}()
	stop := make(chan bool)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cur := atomic.LoadInt64(&proc)
				if cur > last {
					fmt.Printf("Progress: %d/%d (%.1f%%) - Rate: %.0f jobs/sec\n",
						cur, tot,
						float64(cur)*100/float64(tot),
						float64(cur-last)/5.0)
					last = cur
				}
			case <-stop:
				return
			}
		}
	}()
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range job_chan {
				r := checkRectangleFast(pts, segs, red, j.i, j.j,
					minX, maxX, minY, maxY)
				if r.area > 0 {
					res_chan <- r
				}
				atomic.AddInt64(&proc, 1)
			}
		}(w)
	}
	go func() {
		wg.Wait()
		close(res_chan)
		stop <- true
	}()
	max := RectangleResult{area: 0}
	count := 0
	for r := range res_chan {
		count++
		if r.area > max.area {
			max = r
			fmt.Printf("New max found: %d at [%d,%d]-[%d,%d]\n",
				r.area, r.p1.x, r.p1.y, r.p2.x, r.p2.y)
		}
	}
	fmt.Printf("\nTotal valid rectangles found: %d\n", count)
	return max
}

func checkRectangleFast(pts []Point, segs []Segment, red map[[2]int]bool,
	i, j, minX, maxX, minY, maxY int) RectangleResult {
	x1, y1 := pts[i].x, pts[i].y
	x2, y2 := pts[j].x, pts[j].y
	minRX, maxRX := x1, x2
	if x1 > x2 {
		minRX, maxRX = x2, x1
	}
	minRY, maxRY := y1, y2
	if y1 > y2 {
		minRY, maxRY = y2, y1
	}
	if minRX < minX || maxRX > maxX || minRY < minY || maxRY > maxY {
		return RectangleResult{area: 0}
	}
	w := int64(maxRX - minRX + 1)
	h := int64(maxRY - minRY + 1)
	maxSamp := 10000
	cells := w * h
	var samps [][2]int
	if cells <= int64(maxSamp) {
		for y := minRY; y <= maxRY; y++ {
			for x := minRX; x <= maxRX; x++ {
				samps = append(samps, [2]int{x, y})
			}
		}
	} else {
		samps = append(samps,
			[2]int{minRX, minRY},
			[2]int{maxRX, minRY},
			[2]int{minRX, maxRY},
			[2]int{maxRX, maxRY})
		edgeSamp := 100
		for i := 0; i < edgeSamp; i++ {
			t := float64(i) / float64(edgeSamp-1)
			x := minRX + int(float64(maxRX-minRX)*t)
			samps = append(samps, [2]int{x, minRY}, [2]int{x, maxRY})
			y := minRY + int(float64(maxRY-minRY)*t)
			samps = append(samps, [2]int{minRX, y}, [2]int{maxRX, y})
		}
		grid := 50
		for gy := 0; gy < grid; gy++ {
			for gx := 0; gx < grid; gx++ {
				x := minRX + int(float64(maxRX-minRX)*float64(gx)/float64(grid-1))
				y := minRY + int(float64(maxRY-minRY)*float64(gy)/float64(grid-1))
				samps = append(samps, [2]int{x, y})
			}
		}
	}
	for _, pt := range samps {
		x, y := pt[0], pt[1]
		if red[[2]int{x, y}] {
			continue
		}
		if !isPointInPolygon(x, y, segs) {
			return RectangleResult{area: 0}
		}
	}
	area := w * h
	return RectangleResult{
		area: area,
		p1:   pts[i],
		p2:   pts[j],
	}
}

func isPointInPolygon(x, y int, segs []Segment) bool {
	inside := false
	for _, seg := range segs {
		x1, y1 := seg.x1, seg.y1
		x2, y2 := seg.x2, seg.y2
		if ((y1 > y) != (y2 > y)) &&
			(x < (x2-x1)*(y-y1)/(y2-y1)+x1) {
			inside = !inside
		}
	}
	return inside
}
