package main

import "fmt"

func main() {
	for i := 1; i <= 15; i++ { // looping nilai i <= 15
		switch {
		case i%3 == 0 && i%5 == 0: // kondisi kelipatan 3 dan 5
			fmt.Println("FizzBuzz")
		case i%3 == 0: // kondisi kelipatan 3
			fmt.Println("Fizz")
		case i%5 == 0: // kondisi kelipatan 5
			fmt.Println("Buzz")
		default: // kondisi bukan kelipatan 3 ataupun 5
			fmt.Println(i)
		}
	}
}
