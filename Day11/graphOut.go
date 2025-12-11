package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Graph map[string][]string

func parseGraph(filename string) (Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	g := make(Graph)
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
		str := strings.TrimSpace(parts[1])
		var nb []string
		if str != "" {
			nb = strings.Fields(str)
		}
		g[node] = nb
	}
	return g, sc.Err()
}

func main() {
	filename := "./Day11/test11.txt"
	filename = "./Day11/input11.txt"
	g, err := parseGraph(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	fmt.Println("Graph loaded successfully")
	fmt.Printf("Nodes: %d\n", len(g))
	start := "you"
	end := "out"
	vis := make(map[string]bool)
	res := dfs(g, start, end, vis)
	fmt.Printf("Paths from '%s' to '%s': %d\n", start, end, res)
}

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
