# Day 2 Writeup
## Language Used: `Go`
### Part 1:
This time, we try and find for repeating substrings in the numbers, so first thing I do after is create a method that detects if that number is a repeating sequence or not:
```Go
func test(num int) bool {
	s := strconv.Itoa(num)
	n := len(s)
	if n%2 != 0 {
		return false
	}
	m := n / 2
	return s[:m] == s[m:]
}
```
After checking for this, we can easily get the answer.

Now time for part 2!

### Part 2:
As expected, I have to find multiple repeating sequences this time, so time to edit the testing method to account for multiple accounts:
```Go
func test(num int) bool {
	s := strconv.Itoa(num)
	n := len(s)
	for i := 1; i <= n/2; i++ {
		if n%i == 0 {
			p := s[:i]
			rep := n / i
			r := strings.Repeat(p, rep)
			if r == s {
				return true
			}
		}
	}
	return false
}
```
This time we check for each pattern possible until we get the right answer. I can try finding for something more efficient, but this was the simple method that I used currently.

Thus, I got the solution to both parts!

Here are the solutions code:
- [Part 1](./invalidSum.go)
- [Part 2](./moreThanTwice.go)