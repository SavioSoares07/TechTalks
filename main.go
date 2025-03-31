package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello world!!!")
}

func main() {
	http.HandleFunc("/", homeHandler)

	fmt.Println("Servidor rodando na porta 8000")
	http.ListenAndServe(":8000", nil)
}
