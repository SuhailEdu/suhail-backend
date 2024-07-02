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
	Type    string        `json:"type"`
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
	Type    string           `json:"type"`
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
	Id        uuid.UUID        `json:"id"`
	ExamId    uuid.UUID        `json:"exam_id"`
	Title     string           `json:"title"`
	Type      string           `json:"type"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Options   []OptionResource `json:"options"`
}
type LiveQuestionResource struct {
	Id        uuid.UUID `json:"id"`
	ExamId    uuid.UUID `json:"exam_id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Options   []string  `json:"options"`
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
	Id                uuid.UUID `json:"id"`
	UserId            uuid.UUID `json:"user_id"`
	ExamTitle         string    `json:"exam_title"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	QuestionsCount    int64     `json:"questions_count"`
	ParticipantsCount int64     `json:"participants_count"`
}
type ExamParticipant struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Status    string    `json:"status"`
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
			Id:      q.ID,
			ExamId:  q.ExamID,
			Title:   q.Question,
			Options: options,
			Type:    q.Type,
		}

	}

	return ExamResourceWithQuestions{
		Id:        exam.ID,
		UserId:    exam.UserID,
		ExamTitle: exam.Title,
		Status:    exam.VisibilityStatus,
		Questions: qs,
		CreatedAt: exam.CreatedAt.Time,
		UpdatedAt: exam.UpdatedAt.Time,
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
			Type:    q.Type,
		})

	}

	return ExamResourceWithQuestions{
		ExamTitle: exam.Title,
		Status:    exam.VisibilityStatus,
		Questions: questionsResource,
		CreatedAt: exam.CreatedAt.Time,
		UpdatedAt: exam.UpdatedAt.Time,
	}, nil

}

func SerializeExams(exams []schema.GetUserExamsRow) []ExamResource {

	if len(exams) == 0 {

		return []ExamResource{}

	}

	var examResource []ExamResource

	for _, exam := range exams {
		examResource = append(examResource, ExamResource{
			Id:                exam.ID,
			UserId:            exam.UserID,
			ExamTitle:         exam.Title,
			Status:            exam.VisibilityStatus,
			QuestionsCount:    exam.QuestionsCount,
			ParticipantsCount: exam.ParticpantsCount,
			CreatedAt:         exam.CreatedAt.Time,
			UpdatedAt:         exam.UpdatedAt.Time,
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
			CreatedAt:      exam.CreatedAt.Time,
			UpdatedAt:      exam.UpdatedAt.Time,
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
		Type:    question.Type,
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
		Type:    question.Type,
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
			Type:    ex.ExamQuestion.Type,
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

func SerializeSingleParicipatedExam(exam []schema.FindMyParticipatedExamRow) ExamResourceWithQuestions {

	qs := make([]QuestionResource, len(exam))
	for i, ex := range exam {
		var answers []OptionResource
		_ = json.Unmarshal(ex.ExamQuestion.Answers, &answers)

		qs[i] = QuestionResource{
			Id:      ex.ExamQuestion.ID,
			ExamId:  ex.ExamQuestion.ExamID,
			Title:   ex.ExamQuestion.Question,
			Options: answers,
			Type:    ex.ExamQuestion.Type,
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

func SerializeGetExamQuestions(questions []schema.ExamQuestion) []QuestionResource {
	var questionResource []QuestionResource
	for _, question := range questions {
		var answers []OptionResource
		_ = json.Unmarshal(question.Answers, &answers)
		questionResource = append(questionResource, QuestionResource{
			Id:        question.ID,
			ExamId:    question.ExamID,
			Title:     question.Question,
			Type:      question.Type,
			CreatedAt: question.CreatedAt.Time,
			UpdatedAt: question.UpdatedAt.Time,
			Options:   answers,
		})
	}

	return questionResource

}

func SerializeGetLiveExamQuestions(questions []schema.ExamQuestion) interface{} {

	var questionResource []LiveQuestionResource

	for _, question := range questions {
		var answers []OptionResource
		_ = json.Unmarshal(question.Answers, &answers)
		var fixedAnswers []string

		for _, answer := range answers {
			fixedAnswers = append(fixedAnswers, answer.Option)
		}
		questionResource = append(questionResource, LiveQuestionResource{
			Id:        question.ID,
			ExamId:    question.ExamID,
			Title:     question.Question,
			Type:      question.Type,
			CreatedAt: question.CreatedAt.Time,
			UpdatedAt: question.UpdatedAt.Time,
			Options:   fixedAnswers,
		})
	}

	return questionResource

}
func SerializeGetExamParticipants(participants []schema.GetExamParticipantsRow) []ExamParticipant {

	ps := make([]ExamParticipant, len(participants))
	for i, participant := range participants {
		uuidValue, _ := uuid.FromBytes(participant.ID.Bytes[:])
		ps[i] = ExamParticipant{
			ID:        uuidValue,
			Email:     participant.Email,
			FirstName: participant.FirstName.String,
			LastName:  participant.LastName.String,
			Status:    participant.Status,
		}
	}

	return ps

}

//func SerializeExams(exams []schema.GetUserExamsRow  string) []ExamResource {

func pgUUIDtoGoogleUUID(inputId pgtype.UUID) uuid.UUID {
	idUUId, _ := inputId.UUIDValue()
	byteId := idUUId.Bytes[:]
	id, _ := uuid.FromBytes(byteId)
	return id

}
