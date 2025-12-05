# Day 01 Writeup:
## Language Used: `Go`
### Part 1:
We start advent of code again this year with a simple problem, requiring some simple file reading and arithmetic. This time, we have to calculate the number of times the dial reaches `0` at the final position, so we can calculate that using some simple `'%'` functions for the right side:
```Go
    else {
        pos = (pos + val) % mod
    }
```
and for the left side, we need to perform some calculations to ensure we don't go too much into the negative side:
```Go
    pos = (pos - val)
    for pos < 0 {
        pos += 100
    }
```
Updating `res` each time that we reach 0, we get the solution of this problem.

Pretty effective solution for a simple problem.

### Part 2:
Now we get to the annoying part. 

The first method I wanted to use was to just brute force through the entire thing, which led me to this:
```Go
    for a := 0; a < val; a++ {
        if dir == 'L' {
            pos--
            if pos < 0 {
                pos = 99
                res++
            }
        } else {
            pos++
            if pos == 100 {
                pos = 0
                res++
            }
        }
        fmt.Printf("  Iteration %d.%d: Position %d\tDirection: %c\tValue: %d\n", i, a, pos, dir, val)
    }
```
which was obviously pretty bad, and for some reason was not providing me the right answer.

I cross verified the question again this time to check and see if I had not misunderstood the question and...
Okay I was not accounting for the fact that at the end, I can still return back to the same position and have the multiple rounds of reaching zero, which I finally corrected with this mathematical approach:
```Go
    if dir == 'L' {
        check := pos
        if pos == 0 {
            check = mod
        }
        if val >= check {
            res++
            res += (val - check) / mod
        }
        pos = (pos - val) % mod
        if pos < 0 {
            pos += mod
        }
    } else {
        tot := pos + val
        temp := tot / mod
        res += temp
        pos = tot % mod
    }
```

Thus, giving me the solution to both parts in this manner.

Here are the solutions codes:
- [Part 1](./rotations.go)
- [Part 2](./backToZero.go)