package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"techTalk/database"
	"time"
)

type Post struct{
	Title string
	Description string
	CreatedAt time.Time
	DateStr string

}

//Home page

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query(`
		SELECT title, description, created_at
		FROM posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, "Erro ao buscar posts", http.StatusInternalServerError)
		log.Printf("Erro DB: %v", err)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var createdAtStr string // Recebe o valor bruto vindo do banco
	
		if err := rows.Scan(&p.Title, &p.Description, &createdAtStr); err != nil {
			http.Error(w, "Erro ao ler os dados", http.StatusInternalServerError)
			log.Printf("Erro scan: %v", err)
			return
		}
	
		// Parse: ajustar de acordo com o formato real no banco
		parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Erro ao fazer parse da data: %v", err)
			p.CreatedAt = time.Now() // fallback
		} else {
			loc, _ := time.LoadLocation("America/Sao_Paulo")
			p.CreatedAt = parsedTime.In(loc)
		}
	
		// Agora formata do jeito que quiser, exemplo: 15/04/2025 10:23
		p.DateStr = p.CreatedAt.Format("02/01/2006 15:04")
	
		posts = append(posts, p)
	}
	

	tmpl, err := template.ParseFiles("templates/home/index.html")
	if err != nil {
		http.Error(w, "Erro ao carregar o template", http.StatusInternalServerError)
		log.Printf("Erro ao carregar template: %v", err)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		http.Error(w, "Erro ao executar o template", http.StatusInternalServerError)
		log.Printf("Erro ao executar template: %v", err)
	}
}


//Login Page

func FormLoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/signin/index.html")
	tmpl.Execute(w, nil)
}

//Register Page

func FormRegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/signup/index.html")
	tmpl.Execute(w, nil)
}

//Profile Page

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Usuário não está logado", http.StatusUnauthorized)
		return
	}

	userID := cookie.Value // aqui já é "6", por exemplo
	fmt.Println("userID:", userID)

	rows, err := database.DB.Query(`
		SELECT id, user_id, title, description, created_at FROM posts WHERE user_id = ?
	`, userID)
	if err != nil {
		http.Error(w, "Erro ao buscar post do usuário", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan( &p.Title, &p.Description, &p.CreatedAt); err != nil {
			http.Error(w, "Erro ao ler os dados", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	tmpl, err := template.ParseFiles("templates/profile/index.html")
	if err != nil {
		http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		http.Error(w, "Erro ao renderizar template", http.StatusInternalServerError)
	}
}

//Profile Page

func PostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/post/index.html")
	tmpl.Execute(w, nil)
}