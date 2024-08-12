package data // saya berikan nama package data

type Friend struct { // struct untuk menampung struktur data teman
	Nama      string
	Alamat    string
	Pekerjaan string
	Alasan    string
}

var Friends = []Friend{ // buat variable public agar bisa diakses dari package lain (main)
	{"Alif", "Jalan Kopi", "Software Engineer", "Ingin menambah wawasan dan memahami bahasa pemrograman Golang"},
	{"Budi", "Jalan Teh", "Game Developer", "Ingin mencari teman baru"},
	{"Caca", "Jalan Mangga", "System Analyst", "Ingin menghabiskan waktu luang"},
}
