# Day 06
## Language Used: `Go`
### Part 1:
Instead of getting a maze type of problem, we got a file reading problem...

And it was not fun. File reading in general is not fun if it is not in a particular format, so this just becomes a file handling problem itself.  
No way to optimize it either, apart from just reading through the entire file and then converting as we go.

I took a normal grid, used some switch cases and then solved the problem, this was an easier problem. I am expecting the second part to be annoyingly difficult unfortunately.

### Part 2:
Well, as expected, they completely flipped on us, now we have to perform column wise numerics, and find the proper alignments, to ensure that the problem doesn't irritate later.

The problems weren't inherently difficult, just the file reading part that made it extremely annoying to solve and gain information to perform the necessary calculations from.

Here are the solution codes for the two parts:
- [Part 1](./mathHW.go)
- [Part 2](./colMaths.go)