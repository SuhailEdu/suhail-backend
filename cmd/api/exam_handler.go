package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/SuhailEdu/suhail-backend/internal/types"
	"github.com/SuhailEdu/suhail-backend/internal/validations"
	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
	_ "github.com/thedevsaddam/govalidator"
	_ "golang.org/x/crypto/bcrypt"
	"net/http"
)

func (config *Config) getExamsList(c echo.Context) error {

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID
	exams, err := config.db.GetUserExams(c.Request().Context(), userId)
	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeExams(exams))

}

func (config *Config) getParticipatedExams(c echo.Context) error {

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID
	exams, err := config.db.GetParticipatedExams(c.Request().Context(), userId)
	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeParticipatedExams(exams))

}

func (config *Config) getSingleExam(c echo.Context) error {

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID
	//examId := pgtype.UUID{Bytes: [16]byte([]byte(c.Param("id"))), Valid: true}
	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return serverError(c, err)
	}
	fmt.Println(c.Param("id"))

	exam, err := config.db.FindMyExam(c.Request().Context(), schema.FindMyExamParams{ID: examId, UserID: userId})
	fmt.Println(exam, err)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{})
		participatedExam, pErr := config.db.FindMyParticipatedExam(c.Request().Context(), schema.FindMyParticipatedExamParams{ID: examId, UserID: userId})
		if pErr != nil {
			return c.JSON(http.StatusNotFound, map[string]string{})
		}
		return dataResponse(c, participatedExam)

	}

	return dataResponse(c, types.SerializeSingleExam(exam))

}
func (config *Config) createExam(c echo.Context) error {

	var examSchema types.ExamInput

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

	isCorrect, questionErrors := validations.ValidateQuestions(examSchema.Questions...)
	if !isCorrect {
		return validationError(c, map[string]interface{}{
			"questions": questionErrors,
		})
	}

	examParams := schema.CreateExamParams{
		UserID:           authenticatedUser.ID,
		Title:            examSchema.ExamTitle,
		Slug:             pgtype.Text{String: examSchema.ExamTitle},
		IsAccessable:     pgtype.Bool{Bool: true},
		VisibilityStatus: examSchema.Status,
	}

	createdExam, err := config.db.CreateExam(c.Request().Context(), examParams)
	if err != nil {
		return serverError(c, err)
	}

	var questionsParams []schema.CreateExamQuestionsParams

	for _, question := range examSchema.Questions {
		answers, err := json.Marshal(question.Options)
		if err != nil {
			return serverError(c, err)
		}
		questionsParams = append(questionsParams, schema.CreateExamQuestionsParams{
			ExamID:   createdExam.ID,
			Question: question.Title,
			Type:     "options",
			Answers:  answers,
		})
	}

	_, err = config.db.CreateExamQuestions(c.Request().Context(), questionsParams)

	if err != nil {
		return serverError(c, err)
	}

	examId := createdExam.ID

	createdQuestions, _ := config.db.GetExamQuestions(c.Request().Context(), examId)

	return dataResponse(c, types.SerializeExamResource(createdExam, createdQuestions))

}

func (config *Config) updateExam(c echo.Context) error {

	examId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		return badRequestError(c, err)
	}

	canUpdateExam, exam := isExamAuthor(c, config, examId)
	if !canUpdateExam {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	var examSchema types.UpdateExamInput

	rules := govalidator.MapData{
		"exam_title": []string{"required", "min:4", "max:30"},
		"status":     []string{"required", "in:public,private"},
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

	updateParams := schema.UpdateExamParams{
		ID:               examId,
		Title:            examSchema.ExamTitle,
		VisibilityStatus: examSchema.Status,
	}

	exam.Title = examSchema.ExamTitle
	exam.VisibilityStatus = examSchema.Status

	err = config.db.UpdateExam(c.Request().Context(), updateParams)

	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeUpdateExam(exam))

}
func (config *Config) updateQuestion(c echo.Context) error {

	questionId, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		return badRequestError(c, err)
	}

	examId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		return badRequestError(c, err)
	}

	canUpdateExam, _ := isExamAuthor(c, config, examId)
	if !canUpdateExam {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	_, err = config.db.GetQuestionById(c.Request().Context(), questionId)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{})
	}

	//authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)

	var questionSchema types.UpdateQuestionInput

	rules := govalidator.MapData{
		"title":   []string{"required"},
		"options": []string{"required"},
	}

	opts := govalidator.Options{
		Request: c.Request(), // request object
		Rules:   rules,       // rules map
		Data:    &questionSchema,
	}
	// Create a new validator instance
	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		return validationError(c, e)
	}

	var options []types.OptionInput
	for _, option := range questionSchema.Options {
		options = append(options, types.OptionInput{
			Option:    option.Option,
			IsCorrect: option.IsCorrect,
		})
	}

	question := types.QuestionInput{
		Title:   questionSchema.Title,
		Options: options,
	}

	isCorrect, vErrors := validations.ValidateQuestions(question)

	if !isCorrect {
		return validationError(c, vErrors)
	}

	jsonOptions, _ := json.Marshal(question.Options)

	updateParams := schema.UpdateQuestionParams{
		ID:       questionId,
		Question: questionSchema.Title,
		Answers:  jsonOptions,
	}

	err = config.db.UpdateQuestion(c.Request().Context(), updateParams)

	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeUpdateQuestion(question, questionId, examId))

}

func (config *Config) addQuestionsToExam(c echo.Context) error {

	examId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		return badRequestError(c, err)
	}

	_, err = config.db.GetExamQuestions(c.Request().Context(), examId)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{})
	}

	canUpdateExam, _ := isExamAuthor(c, config, examId)
	if !canUpdateExam {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	type input struct {
		Question types.QuestionInput `json:"question"`
	}

	var questionInput input

	err = json.NewDecoder(c.Request().Body).Decode(&questionInput)
	if err != nil {
		return badRequestError(c, err)
	}

	alreadyExists, err := config.db.CheckQuestionTitleExits(c.Request().Context(), schema.CheckQuestionTitleExitsParams{
		Question: questionInput.Question.Title,
		ExamID:   examId,
	})

	if err != nil {
		return serverError(c, err)
	}

	if alreadyExists {
		return c.JSON(http.StatusConflict, map[string]string{
			"title": "You already have a question with this title.",
		})
	}

	fmt.Println(questionInput.Question.Title)

	isCorrect, vError := validations.ValidateQuestions([]types.QuestionInput{questionInput.Question}...)
	if !isCorrect {
		return validationError(c, vError)
	}

	jsonOptions, err := json.Marshal(questionInput.Question.Options)
	if err != nil {
		return serverError(c, err)
	}

	questionParams := schema.CreateQuestionParams{

		ExamID:   examId,
		Question: questionInput.Question.Title,
		Type:     "options",
		Answers:  jsonOptions,
	}

	question, err := config.db.CreateQuestion(c.Request().Context(), questionParams)

	return c.JSON(http.StatusOK, types.SerializeCreateQuestion(question))

	//var options []types.OptionInput
	//for _, option := range questionSchema.Options {
	//	options = append(options, types.OptionInput{
	//		Option:    option.Option,
	//		IsCorrect: option.IsCorrect,
	//	})
	//}
	//
	//question := types.QuestionInput{
	//	Title:   questionSchema.Title,
	//	Options: options,
	//}
	//
	//isCorrect, vErrors := validations.ValidateQuestions(question)
	//
	//if !isCorrect {
	//	return validationError(c, vErrors)
	//}
	//
	//jsonOptions, _ := json.Marshal(question.Options)
	//
	//updateParams := schema.UpdateQuestionParams{
	//	ID:       questionId,
	//	Question: questionSchema.Title,
	//	Answers:  jsonOptions,
	//}
	//
	//err = config.db.UpdateQuestion(c.Request().Context(), updateParams)
	//
	//if err != nil {
	//	return serverError(c, err)
	//}
	//
	//return dataResponse(c, types.SerializeUpdateQuestion(question, questionId, examId))

}
func (config *Config) deleteQuestion(c echo.Context) error {

	questionId, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		return badRequestError(c, err)
	}

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	canUpdateExam, _ := isExamAuthor(c, config, examId)
	if !canUpdateExam {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID

	_ = config.db.DeleteQuestion(c.Request().Context(), schema.DeleteQuestionParams{
		ID:     questionId,
		UserID: userId,
		ID_2:   examId,
	})

	return c.JSON(http.StatusNoContent, nil)

}

func isExamAuthor(c echo.Context, config *Config, examId uuid.UUID) (bool, schema.Exam) {

	exam, err := config.db.GetExamById(c.Request().Context(), examId)

	if err != nil {
		return false, schema.Exam{}
	}

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID

	if exam.UserID != userId {
		return false, schema.Exam{}

	}
	return true, schema.Exam{}

}
