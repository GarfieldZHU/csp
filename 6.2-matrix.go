package main

import (
	"fmt"
	"sync"
)

func center(i, j, a int, leftCh, upCh, rightCh, downCh chan int) {
	for x := range leftCh {
		rightCh <- x
		downCh <- (<-upCh) + a * x
	}
}

func east(leftCh chan int) {
	for range leftCh {}
}

func north(downCh chan int) {
	for {
		downCh <- 0
	}
}

func initalize(n int, a [][]int) (vCh, hCh [][]chan int){
	vCh = make([][]chan int, n + 1)
	hCh = make([][]chan int, n + 1)

	for i := 0; i <= n; i++ {
		vCh[i] = make([]chan int, n + 1)
		hCh[i] = make([]chan int, n + 1)
	}

	for i := 0; i < n; i++ {
		vCh[n][i] = make(chan int)
		hCh[i][n] = make(chan int)
		for j := 0; j < n; j++ {
			vCh[i][j] = make(chan int)
			hCh[i][j] = make(chan int)
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			go center(i + 1, j + 1, a[i][j], hCh[i][j], vCh[i][j], hCh[i][j + 1], vCh[i + 1][j])
		}
	}

	for i := 0; i < n; i++ {
		go east(hCh[i][n])
		go north(vCh[0][i])
	}

	return vCh, hCh
}

func west(in int, rightCh chan int) {
	rightCh <- in
}

func south(i int, ret []int, upCh chan int) {
	ret[i] = <-upCh
}

func main() {
	n := 3

	a := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	in := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	vCh, hCh := initalize(3, a)

	ret := make([][]int, n)

	for i := 0; i < n; i++ {
		// for each row in IN matrix
		for j := 0; j < n; j++ {
			go west(in[i][j], hCh[j][0]);
		}

		ret[i] = make([]int, n)

		var wg sync.WaitGroup
		wg.Add(n)
		for j := 0; j < n; j++ {
			go func(j int) {
				south(j, ret[i], vCh[n][j])
				wg.Done()
			}(j)
		}
		wg.Wait()
	}

	fmt.Println(ret)
}