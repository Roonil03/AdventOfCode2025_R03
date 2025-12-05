# Day 4 Writeup
## Language Used: `Go`
### Part 1:
We are back to mazes again `T-T`...

Well, this one requires us to figure out based on current position. Therefore, time to check the entire grid for the accessible ones:
```Go
for r := range rows {
    for c := range cols {
        if grid[r][c] != '@' {
            continue
        }
        adj := 0
        for _, d := range dir {
            nr, nc := r+d[0], c+d[1]
            if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
                if grid[nr][nc] == '@' {
                    adj++
                }
            }
        }
        if adj < 4 {
            fmt.Println(res+1, ": [x,y]:", r+1, ",", c+1)
            res++
        }
    }
}
```

This is pretty standard to figure out which coordinates work, to cross verify and check if we are getting the right answer.

With mazes, I know that part 2 is going to be tough as well, so time to do part 2 now...

### Part 2:
This one seems to be computational hell. If we edit the previous code, and apply iterations to it, the brute force approach will take a lot of time. 

For now, this is the approach I went with, to attain the right answer:
```Go
    i := 0
	for {
		i++
		var temp [][2]int
		for r := range rows {
			for c := range cols {
				if grid[r][c] != '@' {
					continue
				}
				adj := 0
				for _, d := range dir {
					nr, nc := r+d[0], c+d[1]
					if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
						if grid[nr][nc] == '@' {
							adj++
						}
					}
				}
				if adj < 4 {
					temp = append(temp, [2]int{r, c})
				}
			}
		}
		if len(temp) == 0 {
			break
		}
		fmt.Printf("Iteration %d: Removing %d rolls\n", i, len(temp))
		for _, p := range temp {
			grid[p[0]][p[1]] = '.'
			res++
		}
	}
```

Well, that was a nightmare... 

For optimizations, I really can't think of anything that works out better and allows us to do the first part as well. If we want to just focus on the second part seperately, by making it how many iterations it will take, then we can convert the program to try another approach of brute force algorith, but instead of using it as one-pass, we are runing it as a queue approach, similar to a Breadth First Search Approach to figure out whether this code is optimal or not.

Here is the implementation of it, that I cooked up a little later:
```Go
h1 := func(r, c int) int {
    adj := 0
    for _, d := range dir {
        nr, nc := r+d[0], c+d[1]
        if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
            if grid[nr][nc] == '@' {
                adj++
            }
        }
    }
    return adj
}
chk := make(map[[2]int]bool)
for r := range rows {
    for c := range cols {
        if grid[r][c] == '@' {
            chk[[2]int{r, c}] = true
        }
    }
}
it := 0
for len(chk) > 0 {
    it++
    var temp [][2]int    
    for pos := range chk {
        r, c := pos[0], pos[1]
        if grid[r][c] != '@' {
            continue
        }
        if h1(r, c) < 4 {
            temp = append(temp, [2]int{r, c})
        }
    }    
    if len(temp) == 0 {
        break
    }    
    fmt.Printf("it %d: Removing %d rolls\n", it, len(temp))    
    chk = make(map[[2]int]bool)
    for _, p := range temp {
        grid[p[0]][p[1]] = '.'
        res++        
        for _, d := range dir {
            nr, nc := p[0]+d[0], p[1]+d[1]
            if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
                if grid[nr][nc] == '@' {
                    chk[[2]int{nr, nc}] = true
                }
            }
        }
    }
}
```
This does consideribly improve the efficiency of the code...

Thus, with that, both parts are done for the day!

Here are the solution codes:
- [Part 1](./rollAccess.go)
- [Part 2](./againAndAgain.go)