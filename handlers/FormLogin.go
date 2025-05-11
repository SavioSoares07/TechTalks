package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"techTalk/database"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "templates/signin/index.html")
		return
	}

	if err := r.ParseForm(); err != nil {
		renderLoginWithError(w, "Erro ao processar o formulário")
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	var storeHash string
	var userID string

	err := database.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &storeHash)

	if err == sql.ErrNoRows {
		renderLoginWithError(w, "Usuário não encontrado")
		return
	} else if err != nil {
		renderLoginWithError(w, "Erro no servidor")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storeHash), []byte(password)); err != nil {
		renderLoginWithError(w, "Senha inválida")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    userID,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func renderLoginWithError(w http.ResponseWriter, errorMessage string) {
	tmpl, err := template.ParseFiles("templates/signin/index.html")
	if err != nil {
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]string{
		"Error": errorMessage,
	})
}
