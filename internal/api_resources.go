package internal

import (
	"database/sql"
	"encoding/json"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
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

type OptionResource struct {
	Option    string `json:"option"`
	IsCorrect bool   `json:"is_correct"`
}

type QuestionResource struct {
	Title   string           `json:"title"`
	Options []OptionResource `json:"options"`
}

type ExamResource struct {
	ExamTitle string             `json:"exam_title"`
	Status    string             `json:"status"`
	Questions []QuestionResource `json:"questions"`
}

func SerializeUserResource(user schema.User, token string) UserResource {

	return UserResource{
		//ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     token,
	}
}

func SerializeExamResource(exam schema.Exam, questions []schema.ExamQuestion) ExamResource {

	qs := make([]QuestionResource, len(questions))

	var options []OptionResource

	for i, q := range questions {
		//var option OptionResource
		err := json.Unmarshal(q.Answers, &options)
		if err != nil {
			return ExamResource{}
		}
		//options = append(options, option)

		qs[i] = QuestionResource{
			Title:   q.Question,
			Options: options,
		}

	}

	return ExamResource{
		ExamTitle: exam.Title,
		Status:    exam.VisibilityStatus,
		Questions: qs,
	}

}
