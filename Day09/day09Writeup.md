# Day 9 Writeup:
## Language Used: `Go`
### Part 1:
![img](https://image2url.com/images/1765356666655-6cb49f9f-d76a-4c3a-bb01-27a87ff03801.png)  
This part was quite simple, all I was doing is brute forcing and checking to see if this was making the largest rectangle or not, and then storing it and moving forward:
```Go

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			w := int(math.Abs(float64(pts[i].x - pts[j].x + 1)))
			h := int(math.Abs(float64(pts[i].y - pts[j].y + 1)))
			temp := w * h
			if temp > res {
				res = temp
				t1 = pts[i]
				t2 = pts[j]
				fmt.Printf("\tPoints [%d, %d] and [%d, %d]\t Area: %d\n", pts[i].x, pts[i].y, pts[j].x, pts[j].y, temp)
			}
			if temp < 0 {
				fmt.Println("Error, overflow detected!")
				break
			}
		}
	}
```
Thus, that gave me a simple answer...

and then came the part, that took me soo long...

### Part 2:
At first when I read the question I was a little bit confident that I could complete it again in sometime using brute force. I modified the code a bit to make a grid that stored the locations and then checked if the rectangle was present in the rectangle or not.

That ended up taking so much memory space, I had to force shut the execution multiple times. 

Then I tried for optimizing the code using histogram, which proved to be even more difficult, as managing the entire grid was becomng problematic, and I was not able to configure the entire code properly.

Thus, for the first time, I began my descent online to test for methods which would help me figure out the code faster, and help me learn something new.

Scrolling through the subreddit, I saw that there are multiple memes already made on the Part 2 of this problem like [this one](https://www.reddit.com/r/adventofcode/comments/1pilhch/2025_day_9_part_2_at_least_it_worked/) and [this one](https://www.reddit.com/r/adventofcode/comments/1pi0pek/2025_day_9_part_2/).

After spending a bit more time online, I finally found help via this thread [here](https://www.reddit.com/r/adventofcode/comments/1phywvn/2025_day_9_solutions/)

This got me thinking into Sparse Matrices. Since they don't store the entire matrix, and I could use geometry to check and see if they were present in the correct positions or not, I started researching a little more on techniques I could use to solve this problem.  
Which ended up leading me down this [rabbit hole](https://en.wikipedia.org/wiki/Sparse_matrix)...

With the help of [Claude AI's Sonnet Model](https://www.anthropic.com/claude/sonnet), I learnt how to implement the code into [Go](https://lodev.org/cgtutor/raycasting.html), and then moved forward with this implementation.

I went further to reduce the time by introducing [Go's goroutines](https://gobyexample.com/goroutines) to multithread and get the answer quicker.
```Go
go func() {
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            job_chan <- job{i, j}
        }
    }
    close(job_chan)
    fmt.Printf("Job generation complete: %d jobs\n", tot)
}()
//Job Generator
```
```Go
stop := make(chan bool)
go func() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            cur := atomic.LoadInt64(&proc)
            if cur > last {
                fmt.Printf("Progress: %d/%d (%.1f%%) - Rate: %.0f jobs/sec\n",
                    cur, tot,
                    float64(cur)*100/float64(tot),
                    float64(cur-last)/5.0)
                last = cur
            }
        case <-stop:
            return
        }
    }
}()
//Progress Report
```
```Go
for w := 0; w < workers; w++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        for j := range job_chan {
            r := checkRectangleFast(pts, segs, red, j.i, j.j,
                minX, maxX, minY, maxY)
            if r.area > 0 {
                res_chan <- r
            }
            atomic.AddInt64(&proc, 1)
        }
    }(w)
}
//Worker Pool
```
```Go
go func() {
    wg.Wait()
    close(res_chan)
    stop <- true
}()
//Cleanup
```
Here are some details about the methods used in this code:
- `RectangleResult`:  
A struct that stores the result of a valid rectangle check. It holds the rectangle's area (as int64 to handle large values), and the two Point corners (p1, p2) that define the rectangle. This is used to pass rectangle data between goroutines via channels and to track the maximum area rectangle found during parallel processing.
- `Segment`:  
Represents a line segment of the polygon boundary connecting two consecutive red tile points. Each segment stores the coordinates of its start point (x1, y1) and end point (x2, y2). These segments form the polygon edges and are used by the ray-casting algorithm to determine if interior points are inside or outside the polygon.
- `checkRectanglesOptimized`:  
Orchestrates the parallel rectangle checking process by spawning worker goroutines, distributing jobs via channels, and aggregating results. It creates 122,760 jobs (all point pairs), launches up to 16 workers to process them concurrently, tracks progress with atomic counters, and collects valid rectangles through a results channel, ultimately returning the one with maximum area.
- `checkRectangleFast`:  
Validates whether a rectangle between two points contains only red/green tiles. For small rectangles it checks all cells, but for large ones it strategically samples corners, edges, and an interior grid to avoid scanning millions of cells. Each sampled point is tested using ray-casting to verify it's inside the polygon, returning the rectangle area if valid or zero if invalid.

With that, I have officially completed both parts, taking me a lot of time ```._.```

Here are the solutions to both the parts:
- [Part 1](./largestRect.go)
- [Part 2](./withinRectangle.go)