package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB é uma variável global que armazenará a conexão
var DB *sql.DB

// ConectarBanco inicializa a conexão com o banco de dados
func ConnectionDB() {
	dsn := "root:savio2002@tcp(localhost:3306)/techtalks"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	// Testar a conexão
	err = DB.Ping()
	if err != nil {
		log.Fatal("Erro ao verificar conexão com o banco:", err)
	}

	fmt.Println("Conexão com o banco de dados bem-sucedida!")
}

// GetDB retorna a instância do banco de dados
func GetDB() *sql.DB {
	return DB
}
