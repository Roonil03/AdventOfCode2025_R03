package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	MAX_BUTTONS  = 20
	MAX_COUNTERS = 16
)

type Rational struct {
	num *big.Int
	den *big.Int
}

func NewRational(n, d int64) Rational {
	r := Rational{
		num: big.NewInt(n),
		den: big.NewInt(d),
	}
	if r.den.Sign() == 0 {
		r.den.SetInt64(1)
		return r
	}
	if r.den.Sign() < 0 {
		r.num.Neg(r.num)
		r.den.Neg(r.den)
	}
	g := new(big.Int).GCD(nil, nil, new(big.Int).Abs(r.num), new(big.Int).Abs(r.den))
	r.num.Div(r.num, g)
	r.den.Div(r.den, g)
	return r
}

func (r Rational) Add(other Rational) Rational {
	// a/b + c/d = (ad + bc) / bd
	num := new(big.Int).Mul(r.num, other.den)
	num.Add(num, new(big.Int).Mul(other.num, r.den))
	den := new(big.Int).Mul(r.den, other.den)
	return NewRational(num.Int64(), den.Int64())
}

func (r Rational) Sub(other Rational) Rational {
	// a/b - c/d = (ad - bc) / bd
	num := new(big.Int).Mul(r.num, other.den)
	num.Sub(num, new(big.Int).Mul(other.num, r.den))
	den := new(big.Int).Mul(r.den, other.den)
	return NewRational(num.Int64(), den.Int64())
}

func (r Rational) Mul(other Rational) Rational {
	num := new(big.Int).Mul(r.num, other.num)
	den := new(big.Int).Mul(r.den, other.den)
	return NewRational(num.Int64(), den.Int64())
}

func (r Rational) Div(other Rational) Rational {
	num := new(big.Int).Mul(r.num, other.den)
	den := new(big.Int).Mul(r.den, other.num)
	return NewRational(num.Int64(), den.Int64())
}

func (r Rational) IsZero() bool {
	return r.num.Sign() == 0
}

func (r Rational) ToInt() (int64, bool) {
	if r.den.Sign() == 0 {
		return 0, false
	}
	mod := new(big.Int).Mod(r.num, r.den)
	if mod.Sign() != 0 {
		return 0, false
	}
	return new(big.Int).Div(r.num, r.den).Int64(), true
}

func RationalFromInt(n int64) Rational {
	return Rational{
		num: big.NewInt(n),
		den: big.NewInt(1),
	}
}

type Job struct {
	line  string
	index int
}

type Result struct {
	cost  uint64
	index int
	err   error
}

func searchFreeVariables(
	matrix *[MAX_COUNTERS][MAX_BUTTONS + 1]Rational, numButtons int, numPivots int, pivotCol *[MAX_COUNTERS]int, freeVars *[MAX_BUTTONS]int, numFree int, bound uint64, freeVals *[MAX_BUTTONS]uint64, depth int, currentFreeCost uint64, minCost *uint64, mu *sync.Mutex) {
	mu.Lock()
	currentMin := *minCost
	mu.Unlock()
	if currentFreeCost >= currentMin {
		return
	}
	if depth == numFree {
		solution := make([]int64, MAX_BUTTONS)
		for f := 0; f < numFree; f++ {
			solution[freeVars[f]] = int64(freeVals[f])
		}
		totalCost := currentFreeCost
		valid := true
		for rowIdx := numPivots - 1; rowIdx >= 0; rowIdx-- {
			col := pivotCol[rowIdx]
			val := matrix[rowIdx][numButtons]
			for c := col + 1; c < numButtons; c++ {
				val = val.Sub(matrix[rowIdx][c].Mul(RationalFromInt(solution[c])))
			}
			if v, ok := val.ToInt(); ok {
				if v < 0 {
					valid = false
					break
				}
				solution[col] = v
				totalCost += uint64(v)
				mu.Lock()
				if totalCost >= *minCost {
					valid = false
					mu.Unlock()
					break
				}
				mu.Unlock()
			} else {
				valid = false
				break
			}
		}
		if valid {
			mu.Lock()
			if totalCost < *minCost {
				*minCost = totalCost
			}
			mu.Unlock()
		}
		return
	}
	mu.Lock()
	remainingBudget := uint64(0)
	if *minCost > currentFreeCost {
		remainingBudget = *minCost - currentFreeCost
	}
	mu.Unlock()
	thisBound := bound
	if remainingBudget < thisBound {
		thisBound = remainingBudget
	}
	for v := uint64(0); v < thisBound; v++ {
		freeVals[depth] = v
		searchFreeVariables(matrix, numButtons, numPivots, pivotCol, freeVars, numFree, bound, freeVals, depth+1, currentFreeCost+v, minCost, mu)
		mu.Lock()
		if *minCost <= currentFreeCost+v {
			mu.Unlock()
			break
		}
		mu.Unlock()
	}
}

func solveLinearSystem(buttons []uint32, targets []uint32) (uint64, error) {
	numButtons := len(buttons)
	numCounters := len(targets)
	var matrix [MAX_COUNTERS][MAX_BUTTONS + 1]Rational
	for row := 0; row < numCounters; row++ {
		for col := 0; col < numButtons; col++ {
			bit := uint(row)
			if (buttons[col]>>bit)&1 == 1 {
				matrix[row][col] = RationalFromInt(1)
			} else {
				matrix[row][col] = RationalFromInt(0)
			}
		}
		matrix[row][numButtons] = RationalFromInt(int64(targets[row]))
	}
	var pivotCol [MAX_COUNTERS]int
	for i := range pivotCol {
		pivotCol[i] = -1
	}
	currentRow := 0
	for col := 0; col < numButtons; col++ {
		pivotRow := -1
		for row := currentRow; row < numCounters; row++ {
			if !matrix[row][col].IsZero() {
				pivotRow = row
				break
			}
		}
		if pivotRow == -1 {
			continue
		}
		if pivotRow != currentRow {
			for c := 0; c < numButtons+1; c++ {
				matrix[currentRow][c], matrix[pivotRow][c] = matrix[pivotRow][c], matrix[currentRow][c]
			}
		}
		pivotVal := matrix[currentRow][col]
		for c := 0; c < numButtons+1; c++ {
			matrix[currentRow][c] = matrix[currentRow][c].Div(pivotVal)
		}
		for row := 0; row < numCounters; row++ {
			if row != currentRow && !matrix[row][col].IsZero() {
				factor := matrix[row][col]
				for c := 0; c < numButtons+1; c++ {
					matrix[row][c] = matrix[row][c].Sub(factor.Mul(matrix[currentRow][c]))
				}
			}
		}
		pivotCol[currentRow] = col
		currentRow++
	}
	numPivots := currentRow
	for row := numPivots; row < numCounters; row++ {
		if !matrix[row][numButtons].IsZero() {
			return 0, fmt.Errorf("no solution")
		}
	}
	isPivot := make([]bool, MAX_BUTTONS)
	for row := 0; row < numPivots; row++ {
		if pivotCol[row] >= 0 {
			isPivot[pivotCol[row]] = true
		}
	}
	var freeVars [MAX_BUTTONS]int
	numFree := 0
	for col := 0; col < numButtons; col++ {
		if !isPivot[col] {
			freeVars[numFree] = col
			numFree++
		}
	}
	maxTarget := uint32(0)
	for _, t := range targets {
		if t > maxTarget {
			maxTarget = t
		}
	}
	searchBound := uint64(maxTarget + 1)
	minCost := uint64(^uint64(0))
	var freeVals [MAX_BUTTONS]uint64
	var mu sync.Mutex
	searchFreeVariables(&matrix, numButtons, numPivots, &pivotCol, &freeVars, numFree, searchBound, &freeVals, 0, 0, &minCost, &mu)
	if minCost == ^uint64(0) {
		return 0, fmt.Errorf("no solution")
	}
	return minCost, nil
}

func solvePartTwoMachine(line string) (uint64, error) {
	var buttons [MAX_BUTTONS]uint32
	numButtons := 0
	var targets [MAX_COUNTERS]uint32
	numCounters := 0
	i := 0
	for i < len(line) && line[i] != ']' {
		i++
	}
	i++
	for i < len(line) {
		if line[i] == '(' {
			i++
			mask := uint32(0)
			for line[i] != ')' {
				if line[i] >= '0' && line[i] <= '9' {
					num := uint32(0)
					for i < len(line) && line[i] >= '0' && line[i] <= '9' {
						num = num*10 + uint32(line[i]-'0')
						i++
					}
					mask |= 1 << num
				} else {
					i++
				}
			}
			buttons[numButtons] = mask
			numButtons++
			i++
		} else if line[i] == '{' {
			i++
			for line[i] != '}' {
				if line[i] >= '0' && line[i] <= '9' {
					num := uint32(0)
					for i < len(line) && line[i] >= '0' && line[i] <= '9' {
						num = num*10 + uint32(line[i]-'0')
						i++
					}
					targets[numCounters] = num
					numCounters++
				} else {
					i++
				}
			}
			break
		} else {
			i++
		}
	}
	return solveLinearSystem(buttons[:numButtons], targets[:numCounters])
}

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		cost, err := solvePartTwoMachine(job.line)
		results <- Result{cost: cost, index: job.index, err: err}
	}
}

func helper(filename string) (uint64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	var lines []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := sc.Err(); err != nil {
		return 0, err
	}
	numWorkers := runtime.NumCPU()
	jobs := make(chan Job, len(lines))
	results := make(chan Result, len(lines))
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}
	for i, line := range lines {
		jobs <- Job{line: line, index: i}
	}
	close(jobs)
	go func() {
		wg.Wait()
		close(results)
	}()
	total := uint64(0)
	i := 1
	for result := range results {
		if result.err != nil {
			return 0, fmt.Errorf("error solving machine %d: %v", result.index, result.err)
		}
		total += result.cost
		fmt.Println("\tMachine", i, ":", result.cost)
		i++
	}
	return total, nil
}

func main() {
	filename := "./Day10/test10.txt"
	filename = "./Day10/input10.txt"
	res, err := helper(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("The answer is: %d\n", res)
}
