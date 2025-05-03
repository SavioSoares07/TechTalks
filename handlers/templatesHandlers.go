package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"techTalk/database"
	"time"
)

type Response struct {
	ID        int
	Content   string
	CreatedAt time.Time
	AuthorName string
}

type Post struct {
	ID          int
	Title       string
	Description string
	AuthorName  string
	CreatedAt   time.Time
	DateStr     string
	ImageURL    string
	Responses   []Response  
}


//Home page

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
		return
	}

	// Buscar posts
	rows, err := database.DB.Query(`
		SELECT p.id, p.title, p.description, p.created_at, p.image_url, u.nickname
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
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
		var createdAtStr string
		var imageURL sql.NullString

		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&createdAtStr,
			&imageURL,
			&p.AuthorName,
		); err != nil {
			http.Error(w, "Erro ao ler os dados", http.StatusInternalServerError)
			log.Printf("Erro scan: %v", err)
			return
		}

		// Adicionando a URL da imagem
		if imageURL.Valid {
			p.ImageURL = imageURL.String
		} else {
			p.ImageURL = ""
		}

		// Parse do timestamp para a data
		parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Erro ao fazer parse da data: %v", err)
			p.CreatedAt = time.Now()
		} else {
			loc, _ := time.LoadLocation("America/Sao_Paulo")
			p.CreatedAt = parsedTime.In(loc)
		}
		p.DateStr = p.CreatedAt.Format("02/01/2006 15:04")

		// Buscar respostas para cada post
		var responses []Response
		responseRows, err := database.DB.Query(`
			SELECT r.id, r.content, r.created_at, u.nickname
			FROM responses r
			JOIN users u ON r.user_id = u.id
			WHERE r.post_id = ?
			ORDER BY r.created_at DESC
		`, p.ID)

		if err != nil {
			log.Printf("Erro ao buscar respostas para o post %d: %v", p.ID, err)
			continue
		}

		defer responseRows.Close()

		for responseRows.Next() {
			var r Response
			var createdAtStr string
			if err := responseRows.Scan(&r.ID, &r.Content, &createdAtStr, &r.AuthorName); err != nil {
				log.Printf("Erro ao ler resposta: %v", err)
				continue
			}

			parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
			if err != nil {
				r.CreatedAt = time.Now()
			} else {
				loc, _ := time.LoadLocation("America/Sao_Paulo")
				r.CreatedAt = parsedTime.In(loc)
			}
			responses = append(responses, r)
		}

		// Adicionar respostas ao post
		p.Responses = responses
		responseRows.Close() 
		posts = append(posts, p)
	}

	// Carregar o template
	tmpl, err := template.ParseFiles("templates/home/index.html")
	if err != nil {
		http.Error(w, "Erro ao carregar o template", http.StatusInternalServerError)
		log.Printf("Erro ao carregar template: %v", err)
		return
	}

	// Executar o template
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
		var (
			id int
			userID string
			p Post
			createdAtStr string
		)
	
		if err := rows.Scan(&id, &userID, &p.Title, &p.Description, &createdAtStr); err != nil {
			http.Error(w, "Erro ao ler os dados", http.StatusInternalServerError)
			log.Printf("Erro scan profile: %v", err)
			return
		}
	
		parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Erro ao fazer parse da data: %v", err)
			p.CreatedAt = time.Now()
		} else {
			loc, _ := time.LoadLocation("America/Sao_Paulo")
			p.CreatedAt = parsedTime.In(loc)
		}
		p.DateStr = p.CreatedAt.Format("02/01/2006 15:04")
	
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

//Post Page
func PostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/post/index.html")
	tmpl.Execute(w, nil)
}

//Response func

func ResponsePostHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao processar o formulario", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
		return
	}
	userId := cookie.Value

	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	_, err = database.DB.Exec(`
			INSERT INTO responses (post_id, user_id, content, created_at)
		VALUES (?, ?, ?, ?)
	`, postID, userId, content, time.Now())

	if err != nil {
		http.Error(w, "Erro ao salvar a resposta", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)

}