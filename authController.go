package main

import (
	"fmt"
	"forum/database"
	"forum/models"
	"log"
	"net/http"
	"regexp"

	"github.com/gofiber/fiber"
)

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`a-z0-9._%+\-]+@[a-z0-9._%+\-]+\. [a-z0-9._+\-]`)
	return Re.MatchString(email)
}

func (p *Program) Register(w http.ResponseWriter, r *http.Request) error {
	var data map[string]interface{}
	var userData models.User
	if err := c.BodyParser(&data); err != nil {
	}
	fmt.Println("Unable to parse body")
	//Check if password is less than 6 characters
	if len(data["password"].(string)) <= 6 {
		http.Error(w, "Name is too short", http.StatusBadRequest)
		return
	}

	//Check if email is valid
	if !validateEmail(data["email"].(string)) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid email",
		})
		//Check if email already exists
		if database.DB.IsEmailTaken(data["email"].(string)) {
			c.Status(400)
			return c.JSON(fiber.Map{
				"message": "Email already exists",
			})
		}
		user := models.User{
			Name:     data["name"].(string),
			Email:    data["email"].(string),
			Password: data["password"].(string),
		}
		user.SetPassword(data["password"].(string))
		err := database.DB.CreateUser(&user)
		if err != nil {
			log.Println(err)
		}
		c.Status(200)
		return c.JSON(fiber.Map{
			"user":    user,
			"message": "User created successfully",
		})
	}
}
