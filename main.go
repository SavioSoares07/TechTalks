package main

import (
	"fmt"
	"net/http"
	"techTalks/handlers"
)

func main() {
	// Servindo arquivos estáticos
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Rotas
	http.HandleFunc("/signup", handlers.FormRegisterHandler)
	http.HandleFunc("/signin", handlers.FormLoginHandler)


	fmt.Println("Servidor rodando na porta 8000")
	http.ListenAndServe(":8000", nil)
}
