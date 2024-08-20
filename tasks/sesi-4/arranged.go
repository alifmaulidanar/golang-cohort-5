package main

import (
	"fmt"
	"sync"
	"time"
)

func printDataWithMutex(data []interface{}, times int, mutex *sync.Mutex) { // menggunakan mutex
	for i := 0; i < times; i++ {
		mutex.Lock() // lock
		fmt.Println(data, i+1)
		mutex.Unlock() // unlock
		time.Sleep(time.Millisecond * 10)
	}
}

func main() {
	// data
	data1 := []interface{}{"coba1", "coba2", "coba3"}
	data2 := []interface{}{"bisa1", "bisa2", "bisa3"}

	var mutex sync.Mutex

	// menampilkan output go routine secara berurutan
	go printDataWithMutex(data1, 4, &mutex)
	go printDataWithMutex(data2, 4, &mutex)

	time.Sleep(time.Second * 1) // tunggu go routine selesai 1 detik
}
