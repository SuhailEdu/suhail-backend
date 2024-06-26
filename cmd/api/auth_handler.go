package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"github.com/SuhailEdu/suhail-backend/internal/types"
	"github.com/SuhailEdu/suhail-backend/models"
	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
	_ "github.com/thedevsaddam/govalidator"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
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

	user, err := models.Users(qm.Where("email = ?", "client@gmail.com"), qm.Limit(1)).AllG(c.Request().Context())

	if len(user) == 0 {
		return validationError(c, map[string]any{"email": []string{"Incorrect credentials"}})
	}

	if err := bcrypt.CompareHashAndPassword(user[0].Password.Bytes, []byte(userInput.Password)); err != nil {
		return validationError(c, map[string]any{"email": []string{"Incorrect credentials"}})
	}

	authToken, err := createUserToken(*user[0], c)

	if err != nil {
		fmt.Println("error:", err)
		return serverError(c, err)
	}

	return c.JSON(http.StatusOK, types.SerializeUserResource(*user[0], authToken))
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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 12)

	if err != nil {
		return serverError(c, err)
	}

	userSchema := models.User{
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Email:     userInput.Email,
		Password:  null.Bytes{Bytes: passwordHash, Valid: true},
	}

	err = userSchema.InsertG(c.Request().Context(), boil.Infer())
	if err != nil {
		return serverError(c, err)
	}

	authToken, err := createUserToken(userSchema, c)

	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(http.StatusOK, types.SerializeUserResource(userSchema, authToken))
}

func checkEmailIsUnique(c echo.Context, config Config, email string) bool {

	exists, _ := models.Users(qm.Where("email = ?", email)).ExistsG(c.Request().Context())

	return exists

}

func createUserToken(user models.User, c echo.Context) (string, error) {

	randomBytes := make([]byte, 16)

	// Create empty byte string
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	plainText := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(plainText))

	tokenRecord := models.Token{
		Hash:   hash[:],
		UserID: user.ID,
		Expiry: time.Now().Add(100 * time.Hour),
		Scope:  "login",
	}

	err = user.AddTokensG(c.Request().Context(), true, &tokenRecord)

	//err = tokenRecord.InsertG(c.Request().Context(), boil.Infer())

	if err != nil {
		return "", err
	}
	return plainText, nil

}
