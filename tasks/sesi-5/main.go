package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct { // struct user
	Email    string
	Password string
	Address  string
	Job      string
	Reason   string
}

var users = map[string]User{ // user terdaftar
	"john.doe@example.com": {
		Email:    "john.doe@example.com",
		Password: "password123",
		Address:  "123 Main St, Cityville, USA",
		Job:      "Software Engineer",
		Reason:   "Sunshine and warmth.",
	},
	"jane.smith@example.com": {
		Email:    "jane.smith@example.com",
		Password: "mypassword",
		Address:  "456 Elm St, Cityville, USA",
		Job:      "Data Analyst",
		Reason:   "Coffee and data.",
	},
	"alif@mauli.danar": {
		Email:    "alif@maulid.danar",
		Password: "alifmaulidanar",
		Address:  "Jalan Kelapa No. 12, Jakarta Selatan, Indonesia",
		Job:      "Backend Engineer",
		Reason:   "Family, games, and coding.",
	},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // halaman utama (index.html / login page)
		tmpl := template.Must(template.ParseFiles("client/index.html"))
		data := struct { // data yang akan dikirim ke template
			Emails   []string
			AlertMsg string
		}{
			Emails:   getEmails(), // daftar email terdaftar
			AlertMsg: "",          // alert kosong (jika salah password akan muncul alert)
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { // endpoint login
		if r.Method == "POST" {
			email := r.FormValue("email")
			password := r.FormValue("password")
			user, exists := users[email] // pengecekan user dengan email terdaftar

			if !exists {
				tmpl := template.Must(template.ParseFiles("client/not_registered.html")) // redirect ke halaman not_registered.html jika email tidak terdaftar
				tmpl.Execute(w, nil)
				return
			}

			if user.Password != password {
				tmpl := template.Must(template.ParseFiles("client/index.html")) // kirim alert jika password salah
				data := struct {
					Emails   []string
					AlertMsg string
				}{
					Emails:   getEmails(),
					AlertMsg: "Password salah!",
				}
				tmpl.Execute(w, data)
				return
			}

			tmpl := template.Must(template.ParseFiles("client/profile.html")) // redirect ke halaman profile.html jika login berhasil
			tmpl.Execute(w, user)
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) { // endpoint logout
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func getEmails() []string { // mendapatkan daftar email untuk ditampilkan pada "List of available emails" di login page
	var emails []string
	for email := range users {
		emails = append(emails, email)
	}
	return emails
}
