package handlers

import (
	"fmt"
	"net/http"
)

func RegisterPostHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err!=nil{
		http.Error(w, "Erro ao processar o forumulario", http.StatusBadGateway)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	fmt.Println(title)
	fmt.Println(description)
}