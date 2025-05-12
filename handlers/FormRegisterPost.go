package handlers

import (
	"io"
	"net/http"
	"os"
	"techTalk/database"
	"time"
)

func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20)

	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	link := r.FormValue("link")

	userID := cookie.Value

	var imageURL string

	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		err = os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			http.Error(w, "Erro ao criar diretório", http.StatusInternalServerError)
			return
		}

		dst, err := os.Create("./uploads/" + handler.Filename)
		if err != nil {
			http.Error(w, "Erro ao salvar o arquivo", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Erro ao salvar o arquivo", http.StatusInternalServerError)
			return
		}

		imageURL = "/uploads/" + handler.Filename
	} else {
		// Se não tiver imagem, apenas continue. imageURL ficará vazio
		imageURL = ""
	}

	_, err = database.DB.Exec(`
		INSERT INTO posts (title, description, user_id, created_at, image_url, link)
		VALUES (?, ?, ?, ?, ?, ?)
	`, title, description, userID, time.Now(), imageURL, link)


	if err != nil {
		http.Error(w, "Erro ao salvar no banco de dados", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

