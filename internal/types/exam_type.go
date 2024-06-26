package types

import (
	"encoding/json"
	"github.com/SuhailEdu/suhail-backend/models"
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

func SerializeExamResource(exam models.Exam, questions models.ExamQuestionSlice) ExamResourceWithQuestions {

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

func examSchemaToExamResource(exam models.Exam, questions models.ExamQuestionSlice) (ExamResourceWithQuestions, error) {

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

func SerializeExams(exams models.ExamSlice) []ExamResource {

	if len(exams) == 0 {

		return []ExamResource{}

	}

	var examResource []ExamResource

	for _, exam := range exams {
		examId, _ := uuid.FromBytes([]byte(exam.ID))
		userId, _ := uuid.FromBytes([]byte(exam.UserID))

		examResource = append(examResource, ExamResource{
			Id:             examId,
			UserId:         userId,
			ExamTitle:      exam.Title,
			Status:         exam.VisibilityStatus,
			QuestionsCount: 55,
		})
	}

	return examResource

}

func SerializeParticipatedExams(exams models.ExamSlice) []ExamResource {

	if len(exams) == 0 {

		return []ExamResource{}

	}

	var examResource []ExamResource

	for _, exam := range exams {
		examId, _ := uuid.FromBytes([]byte(exam.ID))
		userId, _ := uuid.FromBytes([]byte(exam.UserID))

		examResource = append(examResource, ExamResource{
			Id:             examId,
			UserId:         userId,
			ExamTitle:      exam.Title,
			Status:         exam.VisibilityStatus,
			QuestionsCount: 55,
		})
	}

	return examResource

}

func SerializeSingleExam(exam models.Exam) ExamResource {

	examId, _ := uuid.FromBytes([]byte(exam.ID))
	userId, _ := uuid.FromBytes([]byte(exam.UserID))

	return ExamResource{
		Id:             examId,
		UserId:         userId,
		ExamTitle:      exam.Title,
		Status:         exam.VisibilityStatus,
		QuestionsCount: 55,
	}

}

//func SerializeExams(exams []schema.GetUserExamsRow  string) []ExamResource {

func pgUUIDtoGoogleUUID(inputId pgtype.UUID) uuid.UUID {
	idUUId, _ := inputId.UUIDValue()
	byteId := idUUId.Bytes[:]
	id, _ := uuid.FromBytes(byteId)
	return id

}
