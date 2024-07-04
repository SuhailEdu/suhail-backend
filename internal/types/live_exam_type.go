package types

import (
	"encoding/json"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/google/uuid"
	"time"
)

type LiveExamResource struct {
	Id         uuid.UUID `json:"id"`
	ExamTitle  string    `json:"exam_title"`
	UserId     uuid.UUID `json:"user_id"`
	Status     string    `json:"status"`
	LiveStatus string    `json:"live_status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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

type LiveExamForManager struct {
	Id         uuid.UUID `json:"id"`
	ExamTitle  string    `json:"exam_title"`
	UserId     uuid.UUID `json:"user_id"`
	Status     string    `json:"status"`
	LiveStatus string    `json:"live_status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type LiveExamQuestionForManagerResource struct {
	Id        uuid.UUID        `json:"id"`
	ExamId    uuid.UUID        `json:"exam_id"`
	Title     string           `json:"title"`
	Type      string           `json:"type"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Options   []OptionResource `json:"options"`
}
type LiveExamParticipant struct {
	ID               uuid.UUID `json:"id"`
	ExamId           uuid.UUID `json:"exam_id"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Status           string    `json:"status"`
	AttendanceStatus string    `json:"attendance_status"`
}

func SerializeGetLiveExamQuestionsForManager(questions []schema.GetLiveExamQuestionForManagerRow) interface{} {
	var questionResource []LiveExamQuestionForManagerResource
	for _, question := range questions {
		var answers []OptionResource
		_ = json.Unmarshal(question.ExamQuestion.Answers, &answers)
		questionResource = append(questionResource, LiveExamQuestionForManagerResource{
			Id:        question.ExamQuestion.ID,
			ExamId:    question.ExamQuestion.ExamID,
			Title:     question.ExamQuestion.Question,
			Type:      question.ExamQuestion.Type,
			CreatedAt: question.ExamQuestion.CreatedAt.Time,
			UpdatedAt: question.ExamQuestion.UpdatedAt.Time,
			Options:   answers,
		})
	}

	exam := LiveExamForManager{
		Id:         questions[0].Exam.ID,
		UserId:     questions[0].Exam.UserID,
		ExamTitle:  questions[0].Exam.Title,
		Status:     questions[0].Exam.VisibilityStatus,
		LiveStatus: questions[0].Exam.LiveStatus.String,
		CreatedAt:  questions[0].Exam.CreatedAt.Time,
		UpdatedAt:  questions[0].Exam.UpdatedAt.Time,
	}

	return map[string]interface{}{
		"questions": questionResource,
		"exam":      exam,
	}

}

func SerializeGetLiveExamParticipants(participants []schema.GetLiveExamParticipantsRow) []LiveExamParticipant {

	ps := make([]LiveExamParticipant, len(participants))
	for i, participant := range participants {
		ps[i] = LiveExamParticipant{
			ID:               participant.User.ID,
			ExamId:           participant.ExamParticipant.ExamID,
			Email:            participant.User.Email,
			FirstName:        participant.User.FirstName,
			LastName:         participant.User.LastName,
			Status:           participant.ExamParticipant.Status,
			AttendanceStatus: participant.ExamParticipant.AttendanceStatus.String,
		}
	}

	return ps

}
func SerializeGetLiveExam(exam schema.Exam, questions []schema.ExamQuestion) interface{} {

	var fixedQuestion []LiveQuestionResource

	for _, question := range questions {
		var answers []OptionResource
		_ = json.Unmarshal(question.Answers, &answers)
		var fixedAnswers []string
		for _, answer := range answers {
			fixedAnswers = append(fixedAnswers, answer.Option)
		}
		fixedQuestion = append(fixedQuestion, LiveQuestionResource{
			Id:        question.ID,
			ExamId:    question.ExamID,
			Title:     question.Title,
			Type:      question.Type,
			Options:   fixedAnswers,
			CreatedAt: question.CreatedAt.Time,
			UpdatedAt: question.UpdatedAt.Time,
		})

	}

	return map[string]interface{}{
		"questions": fixedQuestion,
		"exam": LiveExamResource{
			Id:         exam.ID,
			ExamTitle:  exam.Title,
			UserId:     exam.UserID,
			Status:     exam.VisibilityStatus,
			CreatedAt:  exam.CreatedAt.Time,
			UpdatedAt:  exam.UpdatedAt.Time,
			LiveStatus: exam.LiveStatus.String,
		},
	}

}
