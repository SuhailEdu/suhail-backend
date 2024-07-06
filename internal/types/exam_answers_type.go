package types

import (
	"encoding/json"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
)

func SerializeGetParticipantAnswers(questions []schema.GetParticipantAnswersRow) []map[string]interface{} {
	//return questions
	var result []map[string]interface{}

	for _, question := range questions {
		var answers []OptionResource
		_ = json.Unmarshal(question.Answers, &answers)

		result = append(result, map[string]interface{}{
			"exam_id":            question.ExamID,
			"question_id":        question.QuestionID,
			"user_id":            question.UserID,
			"question":           question.Question,
			"type":               question.Type,
			"options":            answers,
			"participant_answer": question.Answer,
		})

	}

	return result

}
