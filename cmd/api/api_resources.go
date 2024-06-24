package main

import (
	"database/sql"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/google/uuid"
	"time"
)

type UserResource struct {
	ID              uuid.UUID    `json:"id"`
	FirstName       string       `json:"first_name"`
	LastName        string       `json:"last_name"`
	Email           string       `json:"email"`
	Password        []byte       `json:"-"`
	EmailVerifiedAt sql.NullTime `json:"-"`
	CreatedAt       time.Time    `json:"-"`
	UpdatedAt       time.Time    `json:"-"`
	Token           string       `json:"token"`
}

func serializeUserResource(user schema.User, token string) UserResource {

	return UserResource{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     token,
	}

}
