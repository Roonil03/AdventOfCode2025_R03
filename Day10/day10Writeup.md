# Day 10
## Language Used: `Go`
### Part 1:
After a long time I was actually baffled when originally reading the question. Similar to [Day 8](../Day08/README.md)'s problem, it took me sometime to actually decipher what the question actually meant, and then move forward with the explaination that was provided.  

First I made the method to parse the input in, since I will be using it for both parts:
```Go
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
```

As verbose as possible, I had to ensure that nothing could go wrong. This question at first seemed computationally hell to solve, which I did not like.

Since it's a system of linear equations, I opted for making a Gaussian Elimination Algorithm over GF(2):
```Go
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
```

### Part 2:
This...  
Took...  
Way...  
TOOOOO...  
MUCHHHH...  
TIMMMMMEEEE...  


Through so many iterations, I had actually given up at one point, and just scrolled through the massive [megathread](https://www.reddit.com/r/adventofcode/comments/1pity70/2025_day_10_solutions/) that contained so little solutions for part 2.

So many of them took hours to compute, and that was a similar issue that I was facing.

It actually felt like I was trying to solve an NP-Hard problem.

One of the tips that I got from one of my friends who had managed to solve it in 7 hours was trying to use Rational Numbers to solve the question.

Thus, for the first time, I had to actually get [Perplexity](https://www.perplexity.ai/) to generate the code for me. This... hurt me since I had absolutely no idea for this question:
```Go
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
```
This was the rational code that I had generated.

Apart from this, I was told to you RREF for this, and I had it generate the entire code. The only thing I am glad is that I was able to figure out the answer after much pain. Hopefully I never have to solve such questions again.

Here are the two solution codes:
- [Part 1](./machineLights.go)
- [Part 2](./machineActive.go)