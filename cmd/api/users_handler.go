package main

import (
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
	_ "github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type validationSchema struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (config *Config) registerUser(c echo.Context) error {

	var userInput validationSchema

	rules := govalidator.MapData{
		"first_name": []string{"required", "between:3,12"},
		"last_name":  []string{"required", "between:3,12"},
		"email":      []string{"required", "min:4", "max:20", "email"},
		"password":   []string{"required", "min:8", "max:20"},
	}

	opts := govalidator.Options{
		Request: c.Request(), // request object
		Rules:   rules,       // rules map
		Data:    &userInput,
	}
	// Create a new validator instance
	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {

		err := map[string]interface{}{"validationError": e}
		// Validate the User struct
		// Validation failed, handle the error
		return validationError(c, err)
	}

	if checkEmailIsUnique(c, *config, userInput.Email) {

		emailError := map[string]interface{}{"email": []string{"Email is already taken"}}

		err := map[string]interface{}{"validationError": emailError}
		return validationError(c, err)
	}

	log.Println(userInput)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 12)

	if err != nil {
		return serverError(c, err)
	}

	userSchema := schema.CreateUserParams{
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Email:     userInput.Email,
		Password:  passwordHash,
	}

	createdUser, err := config.db.CreateUser(c.Request().Context(), userSchema)

	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(http.StatusOK, serializeUserResource(createdUser))
}

func checkEmailIsUnique(c echo.Context, config Config, email string) bool {

	isUsed, _ := config.db.CheckEmailUniqueness(c.Request().Context(), email)

	return isUsed

}
