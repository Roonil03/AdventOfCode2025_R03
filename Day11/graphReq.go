package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Graph map[string][]string

func parseGraph(filename string) (Graph, Graph) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil
	}
	defer file.Close()
	g := make(Graph)
	rev := make(Graph)
	nodes := make(map[string]bool)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		node := strings.TrimSpace(parts[0])
		neighborsStr := strings.TrimSpace(parts[1])
		nodes[node] = true

		var neighbors []string
		if neighborsStr != "" {
			neighbors = strings.Fields(neighborsStr)
		}
		g[node] = neighbors

		for _, neighbor := range neighbors {
			nodes[neighbor] = true
		}
	}
	for node := range nodes {
		if _, ok := g[node]; !ok {
			g[node] = []string{}
		}
		if _, ok := rev[node]; !ok {
			rev[node] = []string{}
		}
	}
	for u, neighbors := range g {
		for _, v := range neighbors {
			rev[v] = append(rev[v], u)
		}
	}
	fmt.Println("Number of nodes:", len(g)-1)
	return g, rev
}

func main() {
	filename := "./Day11/test11_2.txt"
	filename = "./Day11/input11.txt"
	g, rev := parseGraph(filename)
	if g == nil {
		fmt.Println("Error reading file:", filename)
		return
	}
	start, end := "svr", "out"
	req := []string{"dac", "fft"}
	subgraph, relMap := findRelevantSubgraph(g, rev, start, end)
	if len(subgraph) == 0 {
		fmt.Println("No relevant subgraph: either 'svr' cannot reach 'out' or 'out' cannot be reached.")
		fmt.Println("Number of paths: 0")
		return
	}
	var rel []string
	for node := range relMap {
		rel = append(rel, node)
	}
	fmt.Printf("Relevant nodes: %d\n", len(rel))
	s := &TarjanSCC{
		graph:    subgraph,
		nodes:    rel,
		ids:      make(map[string]int),
		low:      make(map[string]int),
		onStack:  make(map[string]bool),
		at:       0,
		sccCount: 0,
		sccId:    make(map[string]int),
	}
	for _, node := range rel {
		s.ids[node] = -1
	}
	s.findSccs()
	sccGraph := make(map[int][]int)
	inDeg := make(map[int]int)
	sccMask := make(map[int]int)
	for i := 0; i < s.sccCount; i++ {
		sccGraph[i] = []int{}
		inDeg[i] = 0
	}
	for u, nei := range subgraph {
		sccU := s.sccId[u]
		for idx, reqNode := range req {
			if u == reqNode {
				sccMask[sccU] |= (1 << idx)
			}
		}
		for _, v := range nei {
			sccV := s.sccId[v]
			if sccU != sccV {
				sccGraph[sccU] = append(sccGraph[sccU], sccV)
				inDeg[sccV]++
			}
		}
	}
	var q []int
	for i := 0; i < s.sccCount; i++ {
		if inDeg[i] == 0 {
			q = append(q, i)
		}
	}
	topo := []int{}
	head := 0
	for head < len(q) {
		u := q[head]
		head++
		topo = append(topo, u)
		for _, v := range sccGraph[u] {
			inDeg[v]--
			if inDeg[v] == 0 {
				q = append(q, v)
			}
		}
	}
	dp := make([][4]int64, s.sccCount)
	startScc, okStart := s.sccId[start]
	endScc, okEnd := s.sccId[end]
	if !okStart || !okEnd {
		fmt.Println("No Paths")
		fmt.Println("Number of paths: 0")
		return
	}
	dp[startScc][sccMask[startScc]] = 1
	for _, u := range topo {
		for mask := 0; mask < 4; mask++ {
			if dp[u][mask] > 0 {
				for _, v := range sccGraph[u] {
					newMask := mask | sccMask[v]
					dp[v][newMask] += dp[u][mask]
				}
			}
		}
	}
	totalPaths := dp[endScc][3]
	fmt.Printf("Number of paths: %d\n", totalPaths)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type TarjanSCC struct {
	graph    Graph
	nodes    []string
	ids      map[string]int
	low      map[string]int
	onStack  map[string]bool
	stack    []string
	at       int
	sccCount int
	sccId    map[string]int
}

func (t *TarjanSCC) findSccs() {
	for _, node := range t.nodes {
		if t.ids[node] == -1 {
			t.dfs(node)
		}
	}
}

func (t *TarjanSCC) dfs(atNode string) {
	t.stack = append(t.stack, atNode)
	t.onStack[atNode] = true
	t.ids[atNode] = t.at
	t.low[atNode] = t.at
	t.at++
	for _, to := range t.graph[atNode] {
		if t.ids[to] == -1 {
			t.dfs(to)
			if t.low[to] < t.low[atNode] {
				t.low[atNode] = t.low[to]
			}
		} else if t.onStack[to] {
			if t.ids[to] < t.low[atNode] {
				t.low[atNode] = t.ids[to]
			}
		}
	}
	if t.ids[atNode] == t.low[atNode] {
		for len(t.stack) > 0 {
			node := t.stack[len(t.stack)-1]
			t.stack = t.stack[:len(t.stack)-1]
			t.onStack[node] = false
			t.sccId[node] = t.sccCount
			if node == atNode {
				break
			}
		}
		t.sccCount++
	}
}

func findRelevantSubgraph(g, rev Graph, start, end string) (Graph, map[string]bool) {
	if _, ok := g[start]; !ok {
		return make(Graph), make(map[string]bool)
	}
	if _, ok := g[end]; !ok {
		return make(Graph), make(map[string]bool)
	}
	reachable := make(map[string]bool)
	q := []string{start}
	reachable[start] = true
	head := 0
	for head < len(q) {
		u := q[head]
		head++
		for _, v := range g[u] {
			if !reachable[v] {
				reachable[v] = true
				q = append(q, v)
			}
		}
	}
	canReachEnd := make(map[string]bool)
	q = []string{end}
	canReachEnd[end] = true
	head = 0
	for head < len(q) {
		u := q[head]
		head++
		for _, v := range rev[u] {
			if !canReachEnd[v] {
				canReachEnd[v] = true
				q = append(q, v)
			}
		}
	}
	rel := make(map[string]bool)
	for node := range reachable {
		if canReachEnd[node] {
			rel[node] = true
		}
	}
	subgraph := make(Graph)
	for u := range rel {
		var relevantNeighbors []string
		for _, v := range g[u] {
			if rel[v] {
				relevantNeighbors = append(relevantNeighbors, v)
			}
		}
		subgraph[u] = relevantNeighbors
	}
	return subgraph, rel
}

/*
Solving the code using DFS memo:

func dfsMemo(g Graph, cur, tar string, vis map[string]bool, reqMask int, reqIndex map[string]int, fullMask int, memo map[[2]string]int) int {
	if i, ok := reqIndex[cur]; ok {
		reqMask |= (1 << i)
	}
	if cur == tar {
		if reqMask == fullMask {
			return 1
		}
		return 0
	}
	if vis[cur] {
		return 0
	}
	key := [2]string{cur, fmt.Sprintf("%d", reqMask)}
	if val, ok := memo[key]; ok {
		return val
	}
	nb, ex := g[cur]
	if !ex {
		memo[key] = 0
		return 0
	}
	vis[cur] = true
	res := 0
	for _, n := range nb {
		res += dfsMemo(g, n, tar, vis, reqMask, reqIndex, fullMask, memo)
	}
	vis[cur] = false
	memo[key] = res
	return res
}

*/
