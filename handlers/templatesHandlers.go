package handlers

import (
	"html/template"
	"net/http"
)

//Home page

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/home.index.html")
	tmpl.Execute(w, nil)
}

// Login Page

func FormLoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/signin/index.html")
	tmpl.Execute(w, nil)
}

// Register Page

func FormRegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/signup/index.html")
	tmpl.Execute(w, nil)
}
