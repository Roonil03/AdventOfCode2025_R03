# Day 8 Writeup
## Language Used: `Go`
### Part 1:
Well, that really picked my brain. Initially I thought of brute forcing the entire solution via a matrix and then moving forward with that, but then I thought of using a subset method and minimal spanning trees if that was possible to get the answer through that method.

After thinking for a little bit, I realized that I was confused as to how the number of joint pairs were even found, since all I could read up was that they were based on shortest distances that were not based on anything. It took me some time to figure out that the distances only played a role here for deciding the connections and number of vertices in the set. Thus, I decided to go for a [Union Set](https://cp-algorithms.com/data_structures/disjoint_set_union.html) approach and then using Kruskal's Algorithm for finding the connections.

Thus with that, I landed with this algorithm:
```Go
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
    }
}
```

The `cap` was seperate since in the testcase, it was 10, but in the final version, I needed to keep it 1000. Thus, I changed it with the initial part of the code.

### Part 2:
This part was just some slight modification such that I was saving the two furthest away sets that could be found using the Kruskal's Algorithm.
Thus, for the `cap`, I changed it to be a maximum of `n-1` connections made, to ensure that I had two distinct furthest sets that were formed in this manner.

Which lead to make this snippet:
```Go
sort.Slice(edges, func(i, j int) bool {
    return edges[i].dist < edges[j].dist
})
uf := NewUnionFind(n)
var res_i, res_j int
conn := 0
for _, edge := range edges {
    if uf.Union(edge.i, edge.j) {
        conn++
        res_i = edge.i
        res_j = edge.j
        fmt.Printf("%d and %d\t(Dist: %.2f)\n", edge.i, edge.j, edge.dist)
        if conn == n-1 {
            break
        }
    }
}
res := pts[res_i].x * pts[res_j].x
fmt.Println("The answer is:", res)
```

Though I solved this using Kruskal's, I am pretty sure one can try solving this using Prim's Algorithm as well.
Basic implementation will also require a Priority Queue as one of the data structures to be used.   
I'd suggest someone try it out using Prim's as well, to be completely honest.

With that, I had the solution to both parts of this problem!

Here are the solution codes:
- [Part 1](./junctionBox.go)
- [Part 2](./lengthProblems.go)