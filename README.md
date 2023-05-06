# gobench
Write and plot your benchmark just in Go.

## How to use
```go
package main

import (
	"math/rand"
	"time"

	"github.com/kevwan/gobench"
)

func main() {
	b := gobench.NewBench()
	b.Run(10000, func() {
		n := rand.Intn(100)
		time.Sleep(time.Millisecond * time.Duration(n))
	})
}
```

After running a period, you can `Ctrl+C` to stop the benchmark and it will automatically open your browser and show the benchmark result like below:

![image](https://user-images.githubusercontent.com/1918356/236614347-fd038716-170b-4ef5-bc1e-d51778c2dc98.png)

## Give a Star! ⭐

If you like or are using this project to learn or start your solution, please give it a star. Thanks!
