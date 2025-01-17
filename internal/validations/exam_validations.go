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
				Message:         "لا يحتوي هذا السؤال على عنوان",
			}
		}

		if len(question.Title) < 3 {

			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: true,
				Message:         "يجب أن يحتوي عنوان السؤال على 5 أحرف كحد أدني",
			}
		}

		if len(question.Title) > 255 {
			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: true,
				Message:         "يجب أن لا يتجاوز عنوان السؤال 255 حرف كحد أقصى",
			}
		}

		if question.Type != "options" && question.Type != "yesOrNo" {

			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: true,
				Message:         "يحب أن يكون نوع السؤال  اختيارات او صح و خطأ",
			}

		}

		switch question.Type {
		case "options":
			isCorrect, vErr := validateOptionsQuestions(question)
			if !isCorrect {
				return false, vErr
			}
		case "yesOrNo":
			isCorrect, vErr := validateYesOrNoQuestions(question)
			if !isCorrect {
				return false, vErr
			}
		}

		var optionTitles []string
		for _, option := range question.Options {
			optionTitles = append(optionTitles, option.Option)
		}

		if !isSliceUnique(optionTitles) {
			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: true,
				Message:         "يحب أن تكون لا تتكرر الاختيارات",
			}
		}
	}

	//validate questions titles uniqueness
	var titles []string

	for _, q := range questionsInput {
		titles = append(titles, q.Title)
	}

	if !isSliceUnique(titles) {
		return false, map[string]string{
			"exam_title": "لديك اختبار بالفعل بهذا العنوان",
		}
	}

	return true, nil
}

func validateOptionsQuestions(questionInputs ...types.QuestionInput) (bool, interface{}) {

	//validate options length
	for i, question := range questionInputs {

		if len(question.Options) < 2 || len(question.Options) > 4 {
			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: true,
				Message:         "question should have 2 to 4 options.",
			}

		}

		//validate question has only one correct option
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

			if optionIndex == len(question.Options)-1 && !correctOptionFound {

				return false, QuestionValidationResponse{
					QuestionIndex:   i,
					IsQuestionError: false,
					OptionIndex:     optionIndex,
					Message:         "question should have one correct option",
				}
			}

		}

		for o, option := range question.Options {
			if option.Option == "" {
				return false, QuestionValidationResponse{
					QuestionIndex:   i,
					IsQuestionError: false,
					OptionIndex:     o,
					Message:         "Question option has no title",
				}
			}

			if len(option.Option) < 1 {

				return false, QuestionValidationResponse{
					QuestionIndex:   i,
					IsQuestionError: false,
					OptionIndex:     o,
					Message:         "Question option title should be at least 3 characters",
				}
			}

			if len(option.Option) > 255 {
				return false, QuestionValidationResponse{
					QuestionIndex:   i,
					IsQuestionError: false,
					OptionIndex:     o,
					Message:         "Question option title should be less than 60 characters",
				}
			}
		}

	}
	return true, nil
}
func validateYesOrNoQuestions(questionInputs ...types.QuestionInput) (bool, interface{}) {

	//validate question has only one correct option

	for i, question := range questionInputs {

		if len(question.Options) != 2 {

			return false, QuestionValidationResponse{
				QuestionIndex:   i,
				IsQuestionError: false,
				Message:         "question should have 2 options",
			}
		}

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

			if optionIndex == len(question.Options)-1 && !correctOptionFound {

				fmt.Println(option.IsCorrect, optionIndex)
				return false, QuestionValidationResponse{
					QuestionIndex:   i,
					IsQuestionError: false,
					OptionIndex:     optionIndex,
					Message:         "question should have one correct option",
				}
			}

		}

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

	return len(uniqueTitles) == len(input)
}
