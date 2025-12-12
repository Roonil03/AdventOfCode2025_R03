# Day 12 Writeup
## Language Used: `Go`
### Part 1:
The final day is here...

The problem initially seemed like a Brute Force problem, taking too much time. I am slowly realizing a pattern with Brute force problems and Advent of Code Problems `._.`

Well, so I thought of an early pruning condition and ensuring that they fit properly.

Thus, I made this simple method:
```Go
func fit(p []Present, r struct {
	w    int
	h    int
	reqs []RegReq
}) bool {
	a1 := r.w * r.h
	tot := 0
	a2 := 0
	for _, req := range r.reqs {
		tot += req.count
		for _, ps := range p {
			if ps.id == req.id {
				a2 += req.count * ps.size
				break
			}
		}
	}
	if tot*3*3 <= a1 {
		return true
	}
	if a2 > a1 {
		return false
	}
	return true
}
```
With that, I got the solution to the problem and got the star.

## Part 2:
Due to completing all the problems that are present in this year, I got the final star.

Thank you so much for this wonderful challenge. Though there was some ups and downs, it was truly fun... Irrespective of the fact that I was boiling my blood during [Day 10](../Day10/)

Well, with that I solved all the challenges!

Here is the solution code:
- [Part 1](./presents.go)