package main

import (
	"fmt"
	"time"
)

func printData(data []interface{}, times int) { // tanpa mutex
	for i := 0; i < times; i++ {
		fmt.Println(data, i+1)
	}
}

func main() {
	// data
	data1 := []interface{}{"coba1", "coba2", "coba3"}
	data2 := []interface{}{"bisa1", "bisa2", "bisa3"}

	// seharusnya menampilkan output go routine secara acak
	go printData(data1, 4)
	go printData(data2, 4)

	time.Sleep(time.Second * 1) // tunggu go routine selesai 1 detik
}
