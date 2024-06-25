package types

import (
	"encoding/json"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type OptionInput struct {
	Option    string `json:"option"`
	IsCorrect bool   `json:"is_correct"`
}

type QuestionInput struct {
	Title   string        `json:"title"`
	Options []OptionInput `json:"options"`
}

type ExamInput struct {
	ExamTitle string          `json:"exam_title"`
	Status    string          `json:"status"`
	Questions []QuestionInput `json:"questions"`
}

type OptionResource struct {
	Option    string `json:"option"`
	IsCorrect bool   `json:"is_correct"`
}

type QuestionResource struct {
	Title   string           `json:"title"`
	Options []OptionResource `json:"options"`
}

type ExamResourceWithQuestions struct {
	ExamTitle string             `json:"exam_title"`
	Status    string             `json:"status"`
	Questions []QuestionResource `json:"questions"`
}

type ExamResource struct {
	Id             uuid.UUID `json:"id"`
	UserId         uuid.UUID `json:"user_id"`
	ExamTitle      string    `json:"exam_title"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	QuestionsCount int64     `json:"questions_count"`
}

func SerializeExamResource(exam schema.Exam, questions []schema.ExamQuestion) ExamResourceWithQuestions {

	qs := make([]QuestionResource, len(questions))

	var options []OptionResource

	for i, q := range questions {

		err := json.Unmarshal(q.Answers, &options)
		if err != nil {
			return ExamResourceWithQuestions{}
		}

		qs[i] = QuestionResource{
			Title:   q.Question,
			Options: options,
		}

	}

	return ExamResourceWithQuestions{
		ExamTitle: exam.Title,
		Status:    exam.VisibilityStatus,
		Questions: qs,
	}

}

func examSchemaToExamResource(exam schema.Exam, questions []schema.ExamQuestion) (ExamResourceWithQuestions, error) {

	//var examResource []ExamResource

	var questionsResource []QuestionResource

	for _, q := range questions {

		var options []OptionResource

		err := json.Unmarshal(q.Answers, &options)
		if err != nil {
			return ExamResourceWithQuestions{}, err

		}
		questionsResource = append(questionsResource, QuestionResource{
			Title:   q.Question,
			Options: options,
		})

	}

	return ExamResourceWithQuestions{
		ExamTitle: exam.Title,
		Status:    exam.VisibilityStatus,
		Questions: questionsResource,
	}, nil

}

func SerializeExams(exams []schema.GetUserExamsRow) []ExamResource {

	if len(exams) == 0 {

		return []ExamResource{}

	}

	var examResource []ExamResource

	for _, exam := range exams {
		examResource = append(examResource, ExamResource{
			Id:             pgUUIDtoGoogleUUID(exam.ID),
			UserId:         pgUUIDtoGoogleUUID(exam.UserID),
			ExamTitle:      exam.Title,
			Status:         exam.VisibilityStatus,
			QuestionsCount: exam.QuestionsCount,
		})
	}

	return examResource

}

func SerializeParticipatedExams(exams []schema.GetParticipatedExamsRow) []ExamResource {

	if len(exams) == 0 {

		return []ExamResource{}

	}

	var examResource []ExamResource

	for _, exam := range exams {
		examResource = append(examResource, ExamResource{
			Id:             pgUUIDtoGoogleUUID(exam.ID),
			UserId:         pgUUIDtoGoogleUUID(exam.UserID),
			ExamTitle:      exam.Title,
			Status:         exam.VisibilityStatus,
			QuestionsCount: exam.QuestionsCount,
		})
	}

	return examResource

}

//func SerializeExams(exams []schema.GetUserExamsRow  string) []ExamResource {

func pgUUIDtoGoogleUUID(inputId pgtype.UUID) uuid.UUID {
	idUUId, _ := inputId.UUIDValue()
	byteId := idUUId.Bytes[:]
	id, _ := uuid.FromBytes(byteId)
	return id

}
