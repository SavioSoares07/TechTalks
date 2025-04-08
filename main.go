package main

import (
	"fmt"
	"net/http"
	"techTalks/database"
	"techTalks/handlers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Servindo arquivos estáticos
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Rotas
	http.HandleFunc("/signup", handlers.FormRegisterHandler)
	http.HandleFunc("/signin", handlers.FormLoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/home", handlers.HomeHandler)
	http.HandleFunc("/profile", handlers.ProfileHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	http.HandleFunc("/registerPost", handlers.RegisterPostHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	database.ConnectionDB()
	defer database.GetDB().Close()

	fmt.Println("Servidor rodando na porta 8000")
	http.ListenAndServe(":8000", nil)
}
