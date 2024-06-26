package types

import (
	"database/sql"
	"github.com/SuhailEdu/suhail-backend/models"
	"github.com/google/uuid"
	"time"
)

type UserResource struct {
	ID              uuid.UUID    `json:"-"`
	FirstName       string       `json:"first_name"`
	LastName        string       `json:"last_name"`
	Email           string       `json:"email"`
	Password        []byte       `json:"-"`
	EmailVerifiedAt sql.NullTime `json:"-"`
	CreatedAt       time.Time    `json:"-"`
	UpdatedAt       time.Time    `json:"-"`
	Token           string       `json:"token"`
}

func SerializeUserResource(user models.User, token string) UserResource {

	return UserResource{
		//ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     token,
	}
}
