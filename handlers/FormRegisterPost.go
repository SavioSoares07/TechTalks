package handlers

import (
	"net/http"
	"techTalk/database"
	"time"
)

func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar o formulário", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	userID := cookie.Value

	

	_, err = database.DB.Exec(`
		INSERT INTO posts (title, description, user_id, created_at)
		VALUES (?, ?, ?, ?)
	`, title, description, userID, time.Now())

	if err != nil {
		http.Error(w, "Erro ao salvar no banco de dados", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
