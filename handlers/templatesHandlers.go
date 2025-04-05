package handlers

import (
	"html/template"
	"log"
	"net/http"
)

//Home page

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home/index.html")
	if err != nil {
		http.Error(w, "Erro ao carregar o template", http.StatusInternalServerError)
		log.Printf("Erro ao carregar template: %v", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Erro ao executar o template", http.StatusInternalServerError)
		log.Printf("Erro ao executar template: %v", err)
	}
}

//Login Page

func FormLoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/signin/index.html")
	tmpl.Execute(w, nil)
}

//Register Page

func FormRegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/signup/index.html")
	tmpl.Execute(w, nil)
}

//