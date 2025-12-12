package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Machine struct {
	numLights int
	tar       []int
	buttons   [][]int
}

func parseInput(line string) Machine {
	diagramRegex := regexp.MustCompile(`\[([^\]]+)\]`)
	diagramMatch := diagramRegex.FindString(line)
	diagram := diagramMatch[1 : len(diagramMatch)-1]
	tar := make([]int, len(diagram))
	for i, c := range diagram {
		if c == '#' {
			tar[i] = 1
		}
	}
	buttonRegex := regexp.MustCompile(`\(([^)]*)\)`)
	buttonMatches := buttonRegex.FindAllString(line, -1)
	var buttons [][]int
	for _, match := range buttonMatches {
		buttonStr := match[1 : len(match)-1]
		var button []int
		if strings.TrimSpace(buttonStr) != "" {
			parts := strings.Split(buttonStr, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part != "" {
					val, _ := strconv.Atoi(part)
					button = append(button, val)
				}
			}
		}
		buttons = append(buttons, button)
	}
	return Machine{
		numLights: len(tar),
		tar:       tar,
		buttons:   buttons,
	}
}

func main() {
	filename := "./Day10/test10.txt"
	filename = "./Day10/input10.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	var macs []Machine
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		machine := parseInput(line)
		macs = append(macs, machine)
	}
	fmt.Printf("Machines: %d\n", len(macs))
	res := 0
	for i, machine := range macs {
		temp := solve(machine.numLights, machine.tar, machine.buttons)
		res += temp
		fmt.Printf("\tMachine %d: %d\n", i+1, temp)
	}
	fmt.Println("The answer is:", res)
}

func solve(n_lig int, tar []int, buttons [][]int) int {
	n_but := len(buttons)
	mat := make([][]int, n_lig)
	for i := range mat {
		mat[i] = make([]int, n_but+1)
	}
	butSet := make([]map[int]bool, n_but)
	for i, btn := range buttons {
		butSet[i] = make(map[int]bool)
		for _, j := range btn {
			butSet[i][j] = true
		}
	}
	for i := range n_lig {
		for j := range n_but {
			if butSet[j][i] {
				mat[i][j] = 1
			}
		}
		mat[i][n_but] = tar[i]
	}
	pR := 0
	var pC []int
	for col := 0; col < n_but && pR < n_lig; col++ {
		t1 := false
		for row := pR; row < n_lig; row++ {
			if mat[row][col] == 1 {
				mat[pR], mat[row] = mat[row], mat[pR]
				t1 = true
				break
			}
		}
		if !t1 {
			continue
		}
		pC = append(pC, col)
		for row := range n_lig {
			if row != pR && mat[row][col] == 1 {
				for i := 0; i <= n_but; i++ {
					mat[row][i] ^= mat[pR][i]
				}
			}
		}
		pR++
	}
	for i := pR; i < n_lig; i++ {
		if mat[i][n_but] == 1 {
			return -1
		}
	}
	p := make(map[int]bool)
	for _, col := range pC {
		p[col] = true
	}
	var f_var []int
	for col := range n_but {
		if !p[col] {
			f_var = append(f_var, col)
		}
	}
	if len(f_var) == 0 {
		sol := make([]int, n_but)
		for i := len(pC) - 1; i >= 0; i-- {
			rid := i
			cid := pC[i]
			value := mat[rid][n_but]
			for j := cid + 1; j < n_but; j++ {
				value ^= mat[rid][j] * sol[j]
			}
			sol[cid] = value
		}
		sum := 0
		for _, val := range sol {
			sum += val
		}
		return sum
	}
	res := math.MaxInt
	comb := 1 << uint(len(f_var))
	for c := range comb {
		sol := make([]int, n_but)
		for i, freeVar := range f_var {
			if (c>>uint(i))&1 == 1 {
				sol[freeVar] = 1
			}
		}
		for i := len(pC) - 1; i >= 0; i-- {
			rid := i
			cid := pC[i]
			value := mat[rid][n_but]
			for j := cid + 1; j < n_but; j++ {
				value ^= mat[rid][j] * sol[j]
			}
			sol[cid] = value
		}
		cur := make([]int, n_lig)
		for i := range n_but {
			if sol[i] == 1 {
				for _, j := range buttons[i] {
					cur[j] ^= 1
				}
			}
		}
		t1 := true
		for i := range n_lig {
			if cur[i] != tar[i] {
				t1 = false
				break
			}
		}
		if t1 {
			temp := 0
			for _, val := range sol {
				temp += val
			}
			if temp < res {
				res = temp
			}
		}
	}
	if res == math.MaxInt {
		return -1
	}
	return res
}
