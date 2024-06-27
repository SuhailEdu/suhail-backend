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

type UpdateExamInput struct {
	ExamTitle string `json:"exam_title"`
	Status    string `json:"status"`
}

type UpdateQuestionInput struct {
	Title   string           `json:"title"`
	Options []OptionResource `json:"options"`
}

type AddQuestionsToExamInput struct {
	Title   string           `json:"title"`
	Options []OptionResource `json:"options"`
}

type OptionResource struct {
	Option    string `json:"option"`
	IsCorrect bool   `json:"is_correct"`
}

type QuestionResource struct {
	Id      uuid.UUID        `json:"id"`
	ExamId  uuid.UUID        `json:"exam_id"`
	Title   string           `json:"title"`
	Options []OptionResource `json:"options"`
}

type ExamResourceWithQuestions struct {
	Id        uuid.UUID          `json:"id"`
	ExamTitle string             `json:"exam_title"`
	UserId    uuid.UUID          `json:"user_id"`
	Status    string             `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
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
			Id:             exam.ID,
			UserId:         exam.UserID,
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
			Id:             exam.ID,
			UserId:         exam.UserID,
			ExamTitle:      exam.Title,
			Status:         exam.VisibilityStatus,
			QuestionsCount: exam.QuestionsCount,
		})
	}

	return examResource

}

func SerializeUpdateExam(exam schema.Exam) ExamResource {

	return ExamResource{
		Id:        exam.ID,
		UserId:    exam.UserID,
		ExamTitle: exam.Title,
		Status:    exam.VisibilityStatus,
		CreatedAt: exam.CreatedAt.Time,
		UpdatedAt: exam.UpdatedAt.Time,
	}

}

func SerializeCreateQuestion(question schema.ExamQuestion) QuestionResource {

	var options []OptionResource
	err := json.Unmarshal(question.Answers, &options)
	if err != nil {
		return QuestionResource{}
	}

	return QuestionResource{
		Id:      question.ID,
		ExamId:  question.ExamID,
		Title:   question.Question,
		Options: options,
	}

}

func SerializeUpdateQuestion(question QuestionInput, questionId uuid.UUID, examId uuid.UUID) QuestionResource {

	options := make([]OptionResource, len(question.Options))

	for i, option := range question.Options {
		options[i] = OptionResource{
			Option:    option.Option,
			IsCorrect: option.IsCorrect,
		}

	}

	return QuestionResource{
		Id:      questionId,
		ExamId:  examId,
		Title:   question.Title,
		Options: options,
	}

}

func SerializeSingleExam(exam []schema.FindMyExamRow) ExamResourceWithQuestions {

	qs := make([]QuestionResource, len(exam))
	for i, ex := range exam {
		var answers []OptionResource
		_ = json.Unmarshal(ex.ExamQuestion.Answers, &answers)

		qs[i] = QuestionResource{
			Id:      ex.ExamQuestion.ID,
			ExamId:  ex.ExamQuestion.ExamID,
			Title:   ex.ExamQuestion.Question,
			Options: answers,
		}
	}

	return ExamResourceWithQuestions{
		Id:        exam[0].Exam.ID,
		UserId:    exam[0].Exam.UserID,
		ExamTitle: exam[0].Exam.Title,
		Status:    exam[0].Exam.VisibilityStatus,
		CreatedAt: exam[0].Exam.CreatedAt.Time,
		UpdatedAt: exam[0].Exam.UpdatedAt.Time,
		Questions: qs,
	}

}

//func SerializeExams(exams []schema.GetUserExamsRow  string) []ExamResource {

func pgUUIDtoGoogleUUID(inputId pgtype.UUID) uuid.UUID {
	idUUId, _ := inputId.UUIDValue()
	byteId := idUUId.Bytes[:]
	id, _ := uuid.FromBytes(byteId)
	return id

}
