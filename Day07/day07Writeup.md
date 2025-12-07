# Day 7 Writeup
## Language Used: `Go`
### Part 1:
![img1](https://image2url.com/images/1765091024029-352b2ec9-140a-4f83-a9cb-a394c5863cd0.jpg)  

Well, it was time that we got a proper grid type problem. This part was not even that bad. All we needed to check if there was a split character, and if there was a possibility that it would work.

Thus, I landed with this main functionality:
```Go
    for i := 2; i < r; i++ {
        for j := 0; j < c; j++ {
            if grid[i-1][j] == '|' && grid[i][j] == '^' && j > 0 && j < c-1 {
                grid[i][j-1] = '|'
                grid[i][j+1] = '|'
                res++
            }
            if grid[i-1][j] == '|' && grid[i][j] == '.' {
                grid[i][j] = '|'
            }
        }
    }
```
With this code, I was able to get the solution of the first part, fairly easily.

Now time for the second part...

### Part 2:
And here is where the grid becomes a little problematic...

I initially tried with a mathematical approach, testing if it would be possible to caluclate the number of split timelines and moving forward with that, but that ultimately lead me down the memoization route, allowing me to store the saved states each time I calculated the paths for each layer, ensuring that I don't have double paths, and accidentally don't double count, along with a slightly modified recursive function, and using principles of dynamic programming for this problem.

With that, I modified this helper function (which I named `test`) to get this:
```Go
    test = func(row, col int) int64 {
		if row >= r {
			return 1
		}
		if col < 0 || col >= c {
			return 0
		}
		key := fmt.Sprintf("%d,%d", row, col)
		if val, exists := memo[key]; exists {
			return val
		}
		var temp int64 = 0
		if grid[row][col] == '.' || grid[row][col] == 'S' || grid[row][col] == '|' {
			temp = test(row+1, col)
			fmt.Printf("Current Paths: T:%d\n", temp)
		} else if grid[row][col] == '^' {
			l1 := test(row+1, col-1)
			r1 := test(row+1, col+1)
			fmt.Printf("Current Paths: L:%d\tR:%d\n", l1, r1)
			temp = l1 + r1
		}
		memo[key] = temp
		return temp
	}
```
and with these base cases, I was able to get a monstrous answer for my test input, which was `0_0`...

and with that, I have solved both the parts!

Here are the solution codes:
- [Part 1](./tachyonSpliter.go)
- [Part 2](./timelineTachyon.go)