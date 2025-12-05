# Day 5 Writeup
## Language Used: `Go`
### Part 1:
We got segments, which got me instantly thinking of segment trees. I went to my Templates repository and pulled the code for segment trees and here's the segment tree I went with:
```Go
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
```
And with the segment tree, I solved the first part.

### Part 2:
Well, screw this segment tree approach. I changed the entire code to use a merged intervals approach, and then just enumerated through the entire thing.

HOW WAS THIS SECTION EASIER THAN THE FIRST ONE??!?

WHY DID I DO IT VIA SEGMENT TREE??


AHHHHH


Well... solved these two parts...

Here are the solution codes:
- [Part 1](./fresh.go)
- [Part 2](./mergedIntervals.go)