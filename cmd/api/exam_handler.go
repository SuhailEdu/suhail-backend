package main

import (
	"fmt"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
	_ "github.com/thedevsaddam/govalidator"
	_ "golang.org/x/crypto/bcrypt"
	"net/http"
	"reflect"
)

type Option struct {
	Option    string `json:"option"`
	IsCorrect bool   `json:"is_correct"`
}

type Question struct {
	Title   string   `json:"title"`
	Options []Option `json:"options"`
}

type QuestionValidationResponse struct {
	QuestionIndex   int    `json:"question_index"`
	IsQuestionError bool   `json:"is_question_error"`
	OptionIndex     int    `json:"option_index"`
	Message         string `json:"message"`
}

type ExamInput struct {
	ExamTitle string     `json:"exam_title"`
	Status    string     `json:"status"`
	Questions []Question `json:"questions"`
}

func (config *Config) createExam(c echo.Context) error {

	var examSchema ExamInput

	rules := govalidator.MapData{
		"exam_title": []string{"required", "min:4", "max:30"},
		"status":     []string{"required", "in:public,private"},
		"questions":  []string{"required"},
	}

	opts := govalidator.Options{
		Request: c.Request(), // request object
		Rules:   rules,       // rules map
		Data:    &examSchema,
	}
	// Create a new validator instance
	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {

		return validationError(c, e)
	}

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)

	examTitleExists, _ := config.db.CheckExamTitleExists(c.Request().Context(), schema.CheckExamTitleExistsParams{
		Title:  examSchema.ExamTitle,
		UserID: authenticatedUser.ID,
	})
	if examTitleExists {
		return validationError(c, map[string]interface{}{
			"exam_title": "You already have an exam with this title.",
		})

	}

	isCorrect, questionErrors := validateQuestions(examSchema.Questions)
	if !isCorrect {
		return validationError(c, map[string]interface{}{
			"questions": questionErrors,
		})
	}
	return c.JSON(http.StatusCreated, examSchema)

}

func validateQuestions(questionsInput []Question) (bool, interface{}) {

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
