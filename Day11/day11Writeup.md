# Day 11 Writeup
## Language Used: `Go`
### Part 1:
This was much simpler than the problems we had the last two days, which were kinda insane.

For this part, the main problem was parsing the input to get a proper graph. Once we get the graph, it was a simple DFS problem, and that was an extremely simple problem to solve:
```Go
func dfs(g Graph, cur, tar string, vis map[string]bool) int {
	if cur == tar {
		return 1
	}
	if vis[cur] {
		return 0
	}
	nb, ex := g[cur]
	if !ex {
		return 0
	}
	vis[cur] = true
	res := 0
	for _, neighbor := range nb {
		res += dfs(g, neighbor, tar, vis)
	}
	vis[cur] = false
	return res
}
```
Just putting the `start` as *you* and putting the `end` as *out*, it was easy to get the solution to this part.

### Part 2:
This part was not that bad, I thought just needed some requirement changes, so I made this program:
```Go
func dfs2(g Graph, cur, tar string, vis map[string]bool, req map[string]bool) int {
	if _, r := req[cur]; r {
		req[cur] = true
	}
	if cur == tar {
		t1 := true
		for _, been := range req {
			if !been {
				t1 = false
				break
			}
		}
		if _, r := req[cur]; r {
			req[cur] = false
		}
		if t1 {
			return 1
		}
		return 0
	}
	if vis[cur] {
		if _, r := req[cur]; r {
			req[cur] = false
		}
		return 0
	}
	nb, ex := g[cur]
	if !ex {
		if _, r := req[cur]; r {
			req[cur] = false
		}
		return 0
	}
	vis[cur] = true
	res := 0
	for _, neighbor := range nb {
		res += dfs2(g, neighbor, tar, vis, req)
	}
	vis[cur] = false
	if _, r := req[cur]; r {
		req[cur] = false
	}
	return res
}
```
Where `req` was this:
```Go
req := map[string]bool{"dac": false, "fft": false}
```
After running the program for sometime, I realized that I would need to heavily optimize the code, or else I will take days to get the answer. Since the brute force algorithm was finding paths exponentially, it hurt the efficiency, making me realize I will need to try and help the algorithm out.

I tried messing around with memoization a little bit, since I was feeling that was the main cause of bottleneck in the program, and tried looking for further solutions.

I started reading the [megathread](https://www.reddit.com/r/adventofcode/comments/1pjp1rm/2025_day_11_solutions/) and noticed that a lot of people had used a little bit of [topological sorting](https://cp-algorithms.com/graph/topological-sort.html) for further increasing the efficiency of their algorithms, as well as using [DP](https://www.geeksforgeeks.org/competitive-programming/dynamic-programming-dp-and-directed-acyclic-graphs-dag/) for saving states.

This lead me to the rabbit hole of SCCs([Strongly Connected Components](https://cp-algorithms.com/graph/strongly-connected-components.html)), which further lead me to [Trajan's Algorithm for SCCs](https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm). Now this was completely new territory for me. I have no clue of how to implement it, so as a newbie I consulted some help from [Claude's Latest Sonnet](https://www.anthropic.com/claude/sonnet) model and got this method from it:
```Go
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
```
```Go
func (t *TarjanSCC) findSccs() {
	for _, node := range t.nodes {
		if t.ids[node] == -1 {
			t.dfs(node)
		}
	}
}
```
```Go
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
```
```Go
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
```
This part of the code is completely unfamiliar to me, so I got some context as to what each method does:
- `TarjanSCC`:  
It is a struct that holds all the state needed to run Tarjan’s algorithm on the graph. It stores the adjacency list, the list of nodes to process, per-node discovery ids and low-link values, a stack and on-stack flags, and the resulting SCC id for each node. It also tracks a running counter for the DFS order and the total number of SCCs found.
- `findSccs`:  
It is the top-level method that kicks off Tarjan’s algorithm over the whole (relevant) graph. It loops over all nodes and, for any node that has not yet been assigned an id, calls dfs to explore its SCC. This guarantees that every node in the chosen node set is assigned to exactly one strongly connected component.
- `dfs`:  
It is the core recursive procedure implementing Tarjan’s algorithm for a single start node. It assigns a discovery id and low-link value, pushes the node onto the stack, and explores all outgoing edges to update low-link values based on reachable nodes and back edges. When it discovers that a node is the root of an SCC (its id equals its low-link), it pops nodes from the stack until it reaches that root and assigns them all to a new SCC id.
- `findRelevantSubgraph`:  
It trims the original graph down 
to only the part that can actually participate in paths from start to end. It first does a forward BFS/DFS from start to mark everything reachable from start, then a reverse BFS/DFS from end on the reversed graph to mark everything that can reach end. It intersects these sets to get “relevant” nodes and then builds a subgraph containing only those nodes and edges between them.

Using this, I got my answer, and realized that was the reason why I was not able to get my answer with the original DFS method. It just kept exponentially increasing the count.

However, at this point, I was still extremely curious as though why my Memoization results were not working. I went back to the [Advent of Code Reddit Megathread](https://www.reddit.com/r/adventofcode/comments/1pjp1rm/2025_day_11_solutions/) and came across [this comment](https://www.reddit.com/r/adventofcode/comments/1pjp1rm/comment/ntfli95/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button).

This comment shared this [code](https://github.com/ayoubzulfiqar/advent-of-code/blob/main/2025/Go/Day11/part_2.go) for part 2.  
Now this completely baffled me, since I was so confused why my memoization was not working, especially when apart from the bitmasking, the algorithm was pretty much the same. I tried to analyze by comparing the two codes (I save the codes in different files when I try a different approach so I don't lose progress on any one approach). 

I spent about an hour trying to fix my memoization code, since a lot of people were able to solve the question using that. After a lot of time debugging, I don't know what I was going wrong with. I just decided to scrap the entire thing, and restart the entire process from scratch with a new file, which lead me to this memoization:
```Go
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
	for _, neighbor := range nb {
		res += dfsMemo(g, neighbor, tar, vis, reqMask, reqIndex, fullMask, memo)
	}
	vis[cur] = false
	memo[key] = res
	return res
}
```
I later realized that the initial reason why the memoization was not working was because I was storing the entire string path, instead of masking it and storing only a part of it, actually harming the entire lookup time. 

Well, in the end, I got two versions of the code for solving part 2.  
I have kept both version of it in the same file for reference. (The solution code uses Trajan's SCC Algorithm since it was something new that I had to learn, but has the DFS memoization result present in the bottom)

And thus, I have completed both parts for the day!

Here are my solution codes:
- [Part 1](./graphOut.go)
- [Part 2](./graphReq.go)