package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"techTalk/database"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil{
		http.Error(w, "Erro ao processar o formulário", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	var storeHash string
	var userID string

	err := database.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &storeHash)

	if err == sql.ErrNoRows{
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	} else if err != nil{
		http.Error(w, "Erro ao conectar no banco de dados",http.StatusInternalServerError)
		return
	}

	

	if err := bcrypt.CompareHashAndPassword([]byte(storeHash), []byte(password)); err != nil {
		http.Error(w, "Senha invalida", http.StatusUnauthorized)
		return
	}
	

	http.SetCookie(w, &http.Cookie{
		Name: "user_id",
		Value: userID,
		Path: "/",
		HttpOnly: true,
	})

	fmt.Println("Usuário logado com sucesso")
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}