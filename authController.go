package main

import (
	"errors"
	"fmt"
	"forum/database"
	"forum/models"
	"log"
	"net/http"
	"regexp"
	"github.com/gorilla/sessions"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`a-z0-9._%+\-]+@[a-z0-9._%+\-]+\. [a-z0-9._+\-]`)
	return Re.MatchString(email)
}

func (p *Program) Register(w http.ResponseWriter, r *http.Request) error {
	var data map[string]interface{}
	//var userData models.User
	if err := r.BodyParser(&data); err != nil {
	}
	fmt.Println("Unable to parse body")
	//Check if password is less than 6 characters
	if len(data["password"].(string)) <= 6 {
		http.Error(w, "Name is too short", http.StatusBadRequest)
	}

	//Check if email is valid
	if !validateEmail(data["email"].(string)) {
		http.Error(w, "Name is too short", http.StatusBadRequest)

		//Check if email already exists
		if DB.IsEmailTaken(data["email"].(string)) {
			http.Error(w, "Email already exist", http.StatusBadRequest)

		user := models.User{
			Name:     data["name"].(string),
			Email:    data["email"].(string),
			Password: data["password"].([]byte),
		}
		user.SetPassword(data["password"].(string))
		err := DB.CreateUser(&user)
		if err != nil {
			log.Println(err)
		}
		http.Error(w, "Name is too short", http.StatusBadRequest)
	}
}
return nil
}


var (
    store = sessions.NewCookieStore([]byte("my-secret-key"))
)

func (p *Program) Login(w http.ResponseWriter, r *http.Request) error {
    var data map[string]interface{}
    if err := r.BodyParser(&data); err != nil {
        log.Println(err)
        http.Error(w, "Unable to parse body", http.StatusBadRequest)
        return err
    }

    if err := validateInput(data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return err
    }

    user, err := createUser(data)
    if err != nil {
        log.Println(err)
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return err
    }

    // Create a new session for the user
    session, err := store.New(r, "my-session")
    if err != nil {
        log.Println(err)
        http.Error(w, "Failed to create session", http.StatusInternalServerError)
        return err
    }

    // Set the user ID in the session
    session.Values["user_id"] = user.ID

    // Save the session
    err = session.Save(r, w)
    if err != nil {
        log.Println(err)
        http.Error(w, "Failed to save session", http.StatusInternalServerError)
        return err
    }

    http.Error(w, "User created successfully", http.StatusOK)
    return nil
}



/* 
func (p *Program) Login(w http.ResponseWriter, r *http.Request) error {
	var data map[string]interface{}
	if err := r.BodyParser(&data); err != nil {
		log.Println(err)
		http.Error(w, "Unable to parse body", http.StatusBadRequest)
		return err
	}
	

	if err := validateInput(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	_, err := createUser(data)
if err != nil {
    log.Println(err)
    http.Error(w, "Failed to create user", http.StatusInternalServerError)
    return err
}


	http.Error(w, "User created successfully", http.StatusOK)
	return nil
}
 */
func validateInput(data map[string]interface{}) error {
	if len(data["password"].(string)) < 6 {
		return errors.New("Password is too short")
	}

	if !validateEmail(data["email"].(string)) {
		return errors.New("Invalid email format")
	}

	if DB.IsEmailTaken(data["email"].(string)) {
		return errors.New("Email already exists")
	}

	return nil
}

func createUser(data map[string]interface{}) (*models.User, error) {
	user := models.User{
		Name:     data["name"].(string),
		Email:    data["email"].(string),
		Password: data["password"].([]byte),
	}
	user.SetPassword(data["password"].(string))
	if err := database.DB.CreateUser(&user); err != nil {
		return nil, err
	}
	return &user, nil
}





/* func (p *Program) Login(w http.ResponseWriter, r *http.Request) error {
	var data map[string]interface{}
	if err := r.BodyParser(&data); err != nil {
		http.Error(w, "Unable to parse body", http.StatusBadRequest)
		return err
	}

	// Check if password is less than 6 characters
	if len(data["password"].(string)) < 6 {
		http.Error(w, "Password is too short", http.StatusBadRequest)
		return nil
	}

	// Check if email is valid
	if !validateEmail(data["email"].(string)) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return nil
	}

	// Check if email already exists
	if database.DB.IsEmailTaken(data["email"].(string)) {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return nil
	}

	user := models.User{
		Name:     data["name"].(string),
		Email:    data["email"].(string),
		Password: data["password"].([]byte),
	}
	user.SetPassword(data["password"].(string))
	err := database.DB.CreateUser(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return err
	}

	http.Error(w, "User created successfully", http.StatusOK)
	return nil
}
} */