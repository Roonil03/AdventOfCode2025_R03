package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z int
}

func distance(p1, p2 Point) float64 {
	dx := float64(p1.x - p2.x)
	dy := float64(p1.y - p2.y)
	dz := float64(p1.z - p2.z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

type Edge struct {
	i, j int
	dist float64
}

func main() {
	filename := "./Day08/test08.txt"
	cap := 10
	filename = "./Day08/input08.txt"
	cap = 1000
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
		z, _ := strconv.Atoi(parts[2])
		pts = append(pts, Point{x, y, z})
	}
	n := len(pts)
	mat := make([][]float64, n)
	for i := range mat {
		mat[i] = make([]float64, n)
	}
	for i := range n {
		for j := range n {
			if i == j {
				mat[i][j] = 0
			} else {
				mat[i][j] = distance(pts[i], pts[j])
			}
		}
	}
	var edges []Edge
	for i := range n {
		for j := i + 1; j < n; j++ {
			edges = append(edges, Edge{i, j, mat[i][j]})
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].dist < edges[j].dist
	})
	uf := NewUnionFind(n)
	conn := 0
	for i, edge := range edges {
		if i >= cap {
			break
		}
		if uf.Union(edge.i, edge.j) {
			conn++
			fmt.Printf("%d and %d\t(Dist: %.2f)\n", edge.i, edge.j, edge.dist)
		} else {
			fmt.Printf("%d and %d\t(Dist: %.2f)\tConnected\n", edge.i, edge.j, edge.dist)
		}
	}
	s := uf.GetSizes()
	sort.Slice(s, func(i, j int) bool {
		return s[i] > s[j]
	})
	fmt.Printf("\nTotal edges processed: %d\n", cap)
	fmt.Printf("Total successful connections: %d\n", conn)
	fmt.Printf("Number of circuits: %d\n", len(s))
	fmt.Printf("Circuit sizes: %v\n", s)
	if len(s) >= 3 {
		res := s[0] * s[1] * s[2]
		fmt.Println("The answer is:", res)
	} else if len(s) == 2 {
		res := s[0] * s[1] * 1
		fmt.Println("The answer is:", res)
	} else if len(s) == 1 {
		res := s[0] * 1 * 1
		fmt.Println("The answer is:", res)
	}
}

type UnionFind struct {
	parent []int
	size   []int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	size := make([]int, n)
	for i := range parent {
		parent[i] = i
		size[i] = 1
	}
	return &UnionFind{parent, size}
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)
	if rootX == rootY {
		return false
	}
	if uf.size[rootX] < uf.size[rootY] {
		rootX, rootY = rootY, rootX
	}
	uf.parent[rootY] = rootX
	uf.size[rootX] += uf.size[rootY]
	return true
}

func (uf *UnionFind) GetSizes() []int {
	sizeMap := make(map[int]int)
	for i := range uf.parent {
		root := uf.Find(i)
		sizeMap[root] = uf.size[root]
	}
	var sizes []int
	for _, size := range sizeMap {
		sizes = append(sizes, size)
	}
	return sizes
}
