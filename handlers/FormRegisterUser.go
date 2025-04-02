package handlers

import (
	"fmt"
	"net/http"
)
func RegisterHandler(w http.ResponseWriter, r *http.Request){

	fmt.Println("Entrou aqui")

	//Validações Formulario
	if r.Method != http.MethodPost{
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	//Dados Formularios
	if err := r.ParseForm(); err != nil{
		http.Error(w, "Erro ao processar o formulario", http.StatusBadGateway)
		return
	}


	//Ober dados Formularios
	nome := r.FormValue("name")
	nickname := r.FormValue("nickname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Exibe os valores no console (apenas para teste)
	fmt.Println("Nome:", nome)
	fmt.Println("Sobrenome:", nickname)
	fmt.Println("Email:", email)
	fmt.Println("Senha:", password) 
}