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
	xy := make(map[int]bool)
	for _, r := range ranges {
		xy[r.start] = true
		xy[r.end] = true
	}
	for _, i := range id {
		xy[i] = true
	}
	var coords []int
	for c := range xy {
		coords = append(coords, c)
	}
	sort.Ints(coords)
	st := NewSegmentTree(coords)
	for _, r := range ranges {
		st.Update(r.start, r.end)
	}
	res := 0
	for _, i := range id {
		if st.Query(i) {
			fmt.Println("Present:", i)
			res++
		}
	}
	fmt.Println("The answer is:", res)
}

type SegmentTree struct {
	tree   []bool
	lazy   []bool
	coords []int
}

func NewSegmentTree(coords []int) *SegmentTree {
	n := len(coords)
	return &SegmentTree{
		tree:   make([]bool, 4*n),
		lazy:   make([]bool, 4*n),
		coords: coords,
	}
}

func (st *SegmentTree) push(node, left, right int) {
	if st.lazy[node] {
		st.tree[node] = true
		if left != right {
			st.lazy[2*node] = true
			st.lazy[2*node+1] = true
		}
		st.lazy[node] = false
	}
}

func (st *SegmentTree) updateRange(node, left, right, qLeft, qRight int) {
	st.push(node, left, right)

	if qLeft > right || qRight < left {
		return
	}

	if qLeft <= left && right <= qRight {
		st.lazy[node] = true
		st.push(node, left, right)
		return
	}

	mid := (left + right) / 2
	st.updateRange(2*node, left, mid, qLeft, qRight)
	st.updateRange(2*node+1, mid+1, right, qLeft, qRight)
	st.push(2*node, left, mid)
	st.push(2*node+1, mid+1, right)
	st.tree[node] = st.tree[2*node] || st.tree[2*node+1]
}

func (st *SegmentTree) query(node, left, right, pos int) bool {
	st.push(node, left, right)
	if left == right {
		return st.tree[node]
	}
	mid := (left + right) / 2
	if pos <= mid {
		return st.query(2*node, left, mid, pos)
	}
	return st.query(2*node+1, mid+1, right, pos)
}

func (st *SegmentTree) Update(start, end int) {
	left := sort.SearchInts(st.coords, start)
	right := sort.SearchInts(st.coords, end)
	if left < len(st.coords) && right < len(st.coords) {
		st.updateRange(1, 0, len(st.coords)-1, left, right)
	}
}

func (st *SegmentTree) Query(val int) bool {
	idx := sort.SearchInts(st.coords, val)
	if idx < len(st.coords) && st.coords[idx] == val {
		return st.query(1, 0, len(st.coords)-1, idx)
	}
	return false
}
