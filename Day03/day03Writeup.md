# Day 3 Writeup
## Language Used: `Go`
### Part 1:
This time we have to find the maximum two digit number that we can make preserving the order, so the obvious approach was trying to something with a suffix array of some kind. I settled for using a suffix maximum array, letting me figure out the 2 digit number easily by instantly looking up the second digit after the first digit and selecting the best numbers.

Thus, we get this method that we work with:
```Go
func test1(s string) int {
	n := len(s)
	suf := make([]byte, n)
	suf[n-1] = s[n-1]
	for i := n - 2; i >= 0; i-- {
		suf[i] = max(s[i], suf[i+1])
	}
	res := 0
	for i := 0; i < n-1; i++ {
		a := int(s[i] - '0')
		b := int(suf[i+1] - '0')
		jolts := a*10 + b
		if jolts > res {
			res = jolts
		}
	}
	return res
}
```

And we get the answer to the problem.

### Part 2:
Now we jump from 2 digits to 12 digits. Therefore, I changed my approach from a 2 digit number suffix maximum array to a greedy approach for finding the 12 digits, and landed with this approach:
```Go
func test2(s string) int64 {
	n := len(s)
	var res int64 = 0
	for st := 0; st <= n-12; st++ {
		var j int64 = 0
		pos := st
		for d := 0; d < 12; d++ {
			rem := 12 - d - 1
			searchEnd := n - rem
			t := byte('0')
			id := pos
			for i := pos; i < searchEnd; i++ {
				if s[i] > t {
					t = s[i]
					id = i
				}
			}
			j = j*10 + int64(t-'0')
			pos = id + 1
		}
		if j > res {
			res = j
		}
	}
	return res
}
```

With this, I calculated the 12 digit numbers, calculated their sum and got the output!

And with that, I completed both problems for the day!

Here are the solution codes:
- [Part 1](./joltage.go)
- [Part 2](./12jolts.go)