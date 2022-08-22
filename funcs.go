package H

import (
	"sync"
)

func Filter[T interface{}](v []T, f func(x T) bool) []T {
	var wg sync.WaitGroup
	chunkSize := 100
	numOfChunks := len(v) / chunkSize
	if len(v)%chunkSize != 0 {
		numOfChunks++
	}
	chunks := make([][]T, numOfChunks)
	for i := 0; i < numOfChunks; i++ {
		if i == numOfChunks-1 {
			chunks[i] = v[i*chunkSize:]
		} else {
			chunks[i] = v[i*chunkSize : (i+1)*chunkSize]
		}
	}
	var res []T
	resChan := make(chan interface{})
	nilChan := make(chan interface{}, numOfChunks)
	wg.Add(len(chunks))
	for i := 0; i < len(chunks); i++ {
		go func(i int) {
			for _, x := range chunks[i] {
				if f(x) {
					resChan <- x
				}
			}
			nilChan <- nil
			if len(nilChan) == numOfChunks {
				close(nilChan)
				resChan <- nil
				close(resChan)
			}
			wg.Done()
		}(i)
	}
	go func(x *[]T) {
		var res []T
		for v := range resChan {
			if v != nil {
				res = append(res, v.(T))
			}
			*x = res
		}
	}(&res)
	wg.Wait()
	if len(res) == 0 {
		return make([]T, 0)
	}
	return res
}

func Find[T interface{}](v []T, f func(x T) bool) interface{} {
	chunkSize := 100
	numOfChunks := len(v) / chunkSize
	if len(v)%chunkSize != 0 {
		numOfChunks++
	}
	chunks := make([][]T, numOfChunks)
	for i := 0; i < numOfChunks; i++ {
		if i == numOfChunks-1 {
			chunks[i] = v[i*chunkSize:]
		} else {
			chunks[i] = v[i*chunkSize : (i+1)*chunkSize]
		}
	}
	resChan := make(chan interface{})
	nilChan := make(chan interface{}, numOfChunks)
	for i := 0; i < len(chunks); i++ {
		go func(i int) {
			for _, v := range chunks[i] {
				if f(v) {
					resChan <- v
					close(resChan)
					return
				}
			}
			nilChan <- nil
			if len(nilChan) == numOfChunks {
				close(nilChan)
				resChan <- nil
				close(resChan)
			}
		}(i)
	}
	res := <-resChan
	if res == nil {
		return nil
	}
	return res
}
func FindItem[T comparable](v []T, x T) interface{} {
	chunkSize := 100
	numOfChunks := len(v) / chunkSize
	if len(v)%chunkSize != 0 {
		numOfChunks++
	}
	chunks := make([][]T, numOfChunks)
	for i := 0; i < numOfChunks; i++ {
		if i == numOfChunks-1 {
			chunks[i] = v[i*chunkSize:]
		} else {
			chunks[i] = v[i*chunkSize : (i+1)*chunkSize]
		}
	}
	nilChan := make(chan interface{}, numOfChunks)
	resChan := make(chan interface{})
	for i := 0; i < len(chunks); i++ {
		go func(i int) {
			for _, v := range chunks[i] {
				if v == x {
					resChan <- v
					close(resChan)
					return
				}
			}
			nilChan <- nil
			if len(nilChan) == numOfChunks {
				close(nilChan)
				resChan <- nil
				close(resChan)
			}
		}(i)

	}
	res := <-resChan
	if res == nil {
		return nil
	}
	return res
}
func Reduce[T interface{}](v []T, f func(j T, k T) T) T {
	var wg sync.WaitGroup
	chunkSize := 100
	numOfChunks := len(v) / chunkSize
	if len(v)%chunkSize != 0 {
		numOfChunks++
	}
	chunks := make([][]T, numOfChunks)
	for i := 0; i < numOfChunks; i++ {
		if i == numOfChunks-1 {
			chunks[i] = v[i*chunkSize:]
		} else {
			chunks[i] = v[i*chunkSize : (i+1)*chunkSize]
		}
	}
	var res T
	resChan := make(chan interface{})
	nilChan := make(chan interface{}, numOfChunks)
	wg.Add(len(chunks))
	for i := 0; i < len(chunks); i++ {
		go func(i int) {
			var acc T
			for _, x := range chunks[i] {
				acc = f(acc, x)
			}
			resChan <- acc
			nilChan <- nil

			if len(nilChan) == numOfChunks {
				close(nilChan)
				resChan <- nil
				close(resChan)
			}
			wg.Done()
		}(i)
	}
	go func(x *T) {
		var acc T
		for v := range resChan {
			if v != nil {
				acc = f(acc, v.(T))
			}
			*x = acc
		}
	}(&res)
	wg.Wait()
	return res
}
