package validations

import (
	"fmt"
	"github.com/SuhailEdu/suhail-backend/internal/types"
	"reflect"
)

type QuestionValidationResponse struct {
	QuestionIndex   int    `json:"question_index"`
	IsQuestionError bool   `json:"is_question_error"`
	OptionIndex     int    `json:"option_index"`
	Message         string `json:"message"`
}

func ValidateQuestions(questionsInput ...types.QuestionInput) (bool, interface{}) {

	for i, question := range questionsInput {

		//validate question title
		if question.Title == "" {
			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: true,
				Message:         "question has no title",
			}
		}

		if len(question.Title) < 5 {

			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: true,
				Message:         "question title should be at least 5 characters",
			}
		}

		if len(question.Title) > 60 {
			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: true,
				Message:         "question title should be less than 60 characters",
			}
		}

		//validate options length
		if len(question.Options) < 2 || len(question.Options) > 4 {
			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: true,
				Message:         "question should have 2 to 4 options.",
			}

		}

		//validate single correct option
		correctOptionFound := false
		for optionIndex, option := range question.Options {
			if option.IsCorrect {
				if correctOptionFound {
					return false, QuestionValidationResponse{
						QuestionIndex:   i,
						IsQuestionError: false,
						OptionIndex:     optionIndex,
						Message:         "question should have only one correct option",
					}
				}

				correctOptionFound = true
			}

		}

		var optionsTitles []string

		for o, option := range question.Options {
			if option.Option == "" {
				return false, QuestionValidationResponse{
					QuestionIndex:   i,
					IsQuestionError: false,
					OptionIndex:     o,
					Message:         "Question option has no title",
				}
			}

			if len(question.Title) < 5 {

				return false, QuestionValidationResponse{
					QuestionIndex:   i,
					IsQuestionError: false,
					OptionIndex:     o,
					Message:         "Question option title should be at least 5 characters",
				}
			}

			if len(question.Title) > 60 {
				return false, QuestionValidationResponse{
					QuestionIndex:   i,
					IsQuestionError: false,
					OptionIndex:     o,
					Message:         "Question option title should be less than 60 characters",
				}
			}
			optionsTitles = append(optionsTitles, option.Option)
		}

		if !isSliceUnique(optionsTitles) {
			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: true,
				Message:         "Question option titles must be unique",
			}
		}
	}

	//validate questions titles uniqueness
	var titles []string

	for _, q := range questionsInput {
		titles = append(titles, q.Title)
	}

	if !isSliceUnique(titles) {
		return false, "question titles must be unique"
	}

	return true, nil
}
func isSliceUnique(input []string) bool {

	set := make(map[string]interface{})
	for _, element := range input {
		set[element] = struct {
		}{}
	}

	uniqueTitles := reflect.ValueOf(set).MapKeys()

	fmt.Println(len(uniqueTitles), len(input))

	return len(uniqueTitles) == len(input)
}
