package handlers

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"

	"techTalk/database"
)

func ValidateForm(name, nickname, email, password string)error{
	return validation.Errors{
		"name":     validation.Validate(name, validation.Required, validation.Length(3, 200)),
		"nickname": validation.Validate(nickname, validation.Required, validation.Length(3, 50)),
		"password": validation.Validate(password, validation.Required, validation.Length(6, 100)),
	}.Filter()
}


func HashPassword(password string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func RegisterHandler(w http.ResponseWriter, r *http.Request){

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


	//Obter dados Formularios
	name := r.FormValue("name")
	nickname := r.FormValue("nickname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	

	//Validação
	if err := ValidateForm(name, nickname, email, password); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	//Criação do Hash
	hashedPassword, err := HashPassword(password)
	if err != nil{
		http.Error(w, "Erro ao criptografar a senha", http.StatusInternalServerError)
		return
	}
	// Insere os dados no banco de dados
	_, err = database.DB.Exec("INSERT INTO users (name, nickname, email, password) VALUES (?,?,?,?)",
	name, nickname, email, hashedPassword)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Erro ao salvar no banco de dados", http.StatusInternalServerError)
		return
	}

	http.Redirect(w,r,"/signin", http.StatusSeeOther)
}