package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Print("Enter a sentence: ")        // mendapatkan input pengguna
	reader := bufio.NewReader(os.Stdin)    // ternyata memerlukan library bufio untuk mendapatkan lebih dari 1 kata input pengguna
	sentence, _ := reader.ReadString('\n') // membaca input pengguna
	sentence = strings.TrimSpace(sentence) // menampung input pengguna sebagai sebuah string kalimat
	words := strings.Split(sentence, "")   // memecah kalimat per kata

	for _, word := range words { // looping untuk print kata per karakter secara 1 per 1 vertikal
		fmt.Println(word)
	}

	wordCount := make(map[string]int) // membuat var map untuk menampung jumlah masing2 karakter dari input pengguna
	for _, word := range words {      // dilanjutkan dengan menghitung jumlah setiap karakter yang muncul
		wordCount[word]++
	}
	fmt.Println(wordCount) // print penutup untuk menampilkan hasil perhitungan karakter
}
