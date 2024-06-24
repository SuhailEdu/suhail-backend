package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
	_ "github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type registrationSchema struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type loginSchema struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (config *Config) loginUser(c echo.Context) error {

	var userInput loginSchema

	rules := govalidator.MapData{
		"email":    []string{"required", "min:4", "max:20", "email"},
		"password": []string{"required", "min:8", "max:20"},
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

		return validationError(c, e)
	}

	user, err := config.db.GetUserByEmail(c.Request().Context(), userInput.Email)
	if err != nil {
		return validationError(c, map[string]any{"email": []string{"Incorrect credentials"}})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(userInput.Password)); err != nil {
		return validationError(c, map[string]any{"email": []string{"Incorrect credentials"}})
	}

	authToken, err := createUserToken(user, c, *config)

	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(http.StatusOK, serializeUserResource(user, authToken))
}
func (config *Config) registerUser(c echo.Context) error {

	var userInput registrationSchema

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

	authToken, err := createUserToken(createdUser, c, *config)

	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(http.StatusOK, serializeUserResource(createdUser, authToken))
}

func checkEmailIsUnique(c echo.Context, config Config, email string) bool {

	isUsed, _ := config.db.CheckEmailUniqueness(c.Request().Context(), email)

	return isUsed

}

func createUserToken(user schema.User, c echo.Context, config Config) (string, error) {

	randomBytes := make([]byte, 16)

	// Create empty byte string
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	plainText := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(plainText))

	_, err = config.db.CreateUserToken(c.Request().Context(), schema.CreateUserTokenParams{
		Hash:   hash[:],
		UserID: user.ID,
		Expiry: time.Now().Add(24 * time.Hour),
		Scope:  "login",
	})

	if err != nil {
		return "", err
	}
	return plainText, nil

}
