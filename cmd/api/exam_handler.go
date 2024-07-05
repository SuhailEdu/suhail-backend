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
	_ "golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
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
	exams, err := config.db.GetParticipatedExams(c.Request().Context(), pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeParticipatedExams(exams))

}

func (config *Config) getSingleExam(c echo.Context) error {

	authenticatedUser := c.Get("user").(schema.GetUserByTokenRow)
	userId := authenticatedUser.ID

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	exam, err := config.db.FindMyExam(c.Request().Context(), schema.FindMyExamParams{ID: examId, UserID: userId})

	fmt.Println(err)

	if err != nil {

		participatedExam, pErr := config.db.FindMyParticipatedExam(c.Request().Context(), schema.FindMyParticipatedExamParams{
			UserID: pgtype.UUID{Bytes: userId, Valid: true},
			ID:     examId,
		})
		fmt.Println(pErr, participatedExam)
		if len(participatedExam) == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{})
		}

		return dataResponse(c, types.SerializeSingleParicipatedExam(participatedExam))

	}
	return dataResponse(c, types.SerializeSingleExam(exam, true))

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
		"exam_title":     []string{"required", "min:4", "max:30"},
		"status":         []string{"required", "in:public,private"},
		"ip_range_start": []string{"ip_v4"},
		"ip_range_end":   []string{"ip_v4"},
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
		IpRangeStart: pgtype.Text{
			String: examSchema.IpRangeStart,
			Valid:  true,
		},
		IpRangeEnd: pgtype.Text{
			String: examSchema.IpRangeEnd,
			Valid:  true,
		},
	}

	exam.Title = examSchema.ExamTitle
	exam.VisibilityStatus = examSchema.Status
	exam.IpRangeStart = pgtype.Text{
		String: examSchema.IpRangeStart,
		Valid:  true,
	}
	exam.IpRangeEnd = pgtype.Text{
		String: examSchema.IpRangeEnd,
		Valid:  true,
	}

	err = config.db.UpdateExam(c.Request().Context(), updateParams)

	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeUpdateExam(exam))

}

func (config *Config) deleteExam(c echo.Context) error {
	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	err = config.db.DeleteExam(c.Request().Context(), examId)
	if err != nil {
		return serverError(c, err)
	}

	return c.NoContent(http.StatusNoContent)

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
		"type":    []string{"required"},
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
		Type:    questionSchema.Type,
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

	alreadyExists, err := config.db.CheckQuestionTitleExists(c.Request().Context(), schema.CheckQuestionTitleExistsParams{
		Question: questionInput.Question.Title,
		ExamID:   examId,
	})

	if err != nil {
		return serverError(c, err)
	}

	if alreadyExists {
		return validationError(c, map[string]string{
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

func (config *Config) inviteUsersToExam(c echo.Context) error {

	var emails struct {
		Emails []string `json:"emails"`
	}

	err := json.NewDecoder(c.Request().Body).Decode(&emails)
	if err != nil {
		return serverError(c, err)
	}
	if len(emails.Emails) == 0 {
		return validationError(c, map[string]string{"emails": "no email address"})
	}

	for _, email := range emails.Emails {
		_, err = mail.ParseAddress(email)
		if err != nil {
			return validationError(c, map[string]string{"emails": fmt.Sprintf("invalid email address: %s", email)})
		}
		//emailExists, _ := config.db.CheckEmailUniqueness(c.Request().Context(), email)
		//
		//if !emailExists {
		//	return validationError(c, map[string]string{"emails": fmt.Sprintf("User with this email : %s does not exist", email)})
		//}
	}

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	for _, email := range emails.Emails {
		err = config.db.CreateExamParticipant(c.Request().Context(), schema.CreateExamParticipantParams{
			ExamID: examId,
			Email:  email,
			Status: "pending",
		})
		if err != nil {
			fmt.Println("here", email)
			return serverError(c, err)
		}

	}

	return c.JSON(http.StatusOK, emails.Emails)

}

func (config *Config) getExamQuestions(c echo.Context) error {
	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	questions, err := config.db.GetExamQuestions(c.Request().Context(), examId)
	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeGetExamQuestions(questions))
}

func (config *Config) getExamParticipants(c echo.Context) error {
	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	isAuthor, _ := isExamAuthor(c, config, examId)
	if !isAuthor {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	participants, err := config.db.GetExamParticipants(c.Request().Context(), examId)
	if err != nil {
		return serverError(c, err)
	}

	fmt.Println(participants)

	return dataResponse(c, types.SerializeGetExamParticipants(participants))
}

func (config *Config) removeParticipants(c echo.Context) error {

	examId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return badRequestError(c, err)
	}

	isAuthor, _ := isExamAuthor(c, config, examId)
	if !isAuthor {
		return unAuthorizedError(c, errors.New("unauthorized access"))
	}

	var emails struct {
		Emails []string `json:"emails"`
	}

	err = json.NewDecoder(c.Request().Body).Decode(&emails)
	if err != nil {
		return serverError(c, err)
	}
	if len(emails.Emails) == 0 {
		return validationError(c, map[string]string{"emails": "no email address"})
	}

	for _, email := range emails.Emails {
		_, err = mail.ParseAddress(email)
		if err != nil {
			return validationError(c, map[string]string{"emails": fmt.Sprintf("invalid email address: %s", email)})
		}
		//emailExists, _ := config.db.CheckEmailUniqueness(c.Request().Context(), email)
		//
		//if !emailExists {
		//	return validationError(c, map[string]string{"emails": fmt.Sprintf("User with this email : %s does not exist", email)})
		//}
	}

	//myUd, _ := uuid.Parse("2b2258ab-f848-405b-a95c-0901e682d9e7")
	myUd, _ := uuid.FromBytes([]byte("2b2258ab-f848-405b-a95c-0901e682d9e7"))

	fmt.Println(myUd.String())

	err = config.db.DeleteParticipants(c.Request().Context(), schema.DeleteParticipantsParams{
		ExamID: examId,
		Emails: emails.Emails,
	})

	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(http.StatusOK, emails.Emails)

}

func (config *Config) checkExamTitle(c echo.Context) error {
	var examSchema types.ExamInput

	rules := govalidator.MapData{
		"exam_title": []string{"required", "min:4", "max:30"},
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

	return c.NoContent(http.StatusOK)

}

func (config *Config) checkExamQuestions(c echo.Context) error {
	var examSchema types.ExamInput

	rules := govalidator.MapData{
		"questions": []string{"required"},
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

	isCorrect, questionErrors := validations.ValidateQuestions(examSchema.Questions...)
	if !isCorrect {
		return validationError(c, map[string]interface{}{
			"questions": questionErrors,
		})
	}

	return c.NoContent(http.StatusOK)

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
