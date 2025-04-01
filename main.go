package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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


	//Conexão com banco de dados
	dsn := "root:savio2002@tcp(localhost:3306)/techtalks"
	
	db, err := sql.Open("mysql", dsn)
	if err != nil{
		log.Fatal("Erro ao conectar ao banco de dados", err)
	}
	defer db.Close()

	
	// Testar a conexão
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	fmt.Println("Conexão bem-sucedida!")


	fmt.Println("Servidor rodando na porta 8000")
	http.ListenAndServe(":8000", nil)
}
