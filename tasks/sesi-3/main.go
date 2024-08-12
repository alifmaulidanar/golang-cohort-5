package main

import (
	"fmt"
	"os"
	"sesi-3/data" // sertakan package data
	"strings"
)

func main() {
	if len(os.Args) < 2 { // pengecekan input pengguna pertama kali
		fmt.Println("Tolong masukan nama atau nomor absen")
		fmt.Println("Contoh: 'go run main.go Alif' atau 'go run main.go 1'")
		return
	}

	input := os.Args[1] // get input pengguna
	found := false      // untuk error handling

	for i, friend := range data.Friends { // looping data teman yang didapatkan dari package data
		if strings.EqualFold(friend.Nama, input) || fmt.Sprintf("%d", i) == input { // mencari data teman berdasarkan nama atau indeks
			fmt.Printf("ID: %d\n", i)
			fmt.Printf("Nama: %s\n", friend.Nama)
			fmt.Printf("Alamat: %s\n", friend.Alamat)
			fmt.Printf("Pekerjaan: %s\n", friend.Pekerjaan)
			fmt.Printf("Alasan: %s\n", friend.Alasan)
			found = true
			break
		}
	}

	if !found { // error handling di luar looping data teman agar output error handling tidak redundan
		fmt.Println("Data dengan nama/absen tsb tidak tersedia")
	}
}
