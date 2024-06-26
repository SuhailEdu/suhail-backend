package main

import (
	"encoding/json"
	"fmt"
	"github.com/SuhailEdu/suhail-backend/internal/types"
	"github.com/SuhailEdu/suhail-backend/internal/validations"
	"github.com/SuhailEdu/suhail-backend/models"
	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/govalidator"
	_ "github.com/thedevsaddam/govalidator"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	_ "golang.org/x/crypto/bcrypt"
)

func (config *Config) getExamsList(c echo.Context) error {

	authenticatedUser := c.Get("user").(*models.Token)

	userId := authenticatedUser.UserID
	//exams, err := config.db.GetUserExams(c.Request().Context(), userId)

	exams, err := models.Exams(qm.Limit(1), qm.Where("user_id = ?", userId)).AllG(c.Request().Context())

	if err != nil {
		return serverError(c, err)
	}

	return dataResponse(c, types.SerializeExams(exams))

}

//func (config *Config) getParticipatedExams(c echo.Context) error {
//
//	authenticatedUser := c.Get("user").(models.User)
//	userId := authenticatedUser.ID
//	exams, err := config.db.GetParticipatedExams(c.Request().Context(), userId)
//	if err != nil {
//		return serverError(c, err)
//	}
//
//	return dataResponse(c, types.SerializeParticipatedExams(exams))
//
//}

func (config *Config) getSingleExam(c echo.Context) error {

	//authenticatedUser := c.Get("user").(models.User)
	//userId := authenticatedUser.ID

	exam, err := models.FindExamG(c.Request().Context(), c.Param("id"))
	if err != nil {
		fmt.Println(err)
	}

	//if err.Error() == "no rows in result set" {
	//	participatedExam, pErr := config.db.FindMyParticipatedExam(c.Request().Context(), schema.FindMyParticipatedExamParams{ID: examId, UserID: userId})
	//	if pErr.Error() == "no rows in result set" {
	//		return c.JSON(http.StatusNotFound, map[string]string{})
	//	}
	//	return dataResponse(c, participatedExam)
	//
	//}

	return dataResponse(c, types.SerializeSingleExam(*exam))

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

	authenticatedUser := c.Get("user").(*models.Token)

	examTitleExists, err := models.Exams(qm.Where("title = ?", examSchema.ExamTitle)).ExistsG(c.Request().Context())

	if err != nil {
		return serverError(c, err)
	}

	if examTitleExists {
		return validationError(c, map[string]interface{}{
			"exam_title": "You already have an exam with this title.",
		})

	}

	isCorrect, questionErrors := validations.ValidateQuestions(examSchema.Questions)
	if !isCorrect {
		return validationError(c, map[string]interface{}{
			"questions": questionErrors,
		})
	}

	var examParams models.Exam
	examParams.UserID = authenticatedUser.UserID

	examParams.Title = examSchema.ExamTitle
	examParams.Slug = null.String{String: examSchema.ExamTitle, Valid: true}
	examParams.IsAccessable = null.Bool{Bool: true, Valid: true}
	examParams.VisibilityStatus = examSchema.Status

	err = examParams.InsertG(c.Request().Context(), boil.Blacklist("id"))
	err = examParams.GetID()
	if err != nil {
		fmt.Println("Error inserting exam", examParams.ID, err)
		return serverError(c, err)
	}

	fmt.Println("exam", examParams.ID)
	if err != nil {
		return serverError(c, err)
	}
	var questionsParams models.ExamQuestionSlice

	//err = examParams.AddExamQuestionsG(c.Request().Context() , true , &questionsParams)

	for _, question := range examSchema.Questions {
		answers, err := json.Marshal(question.Options)
		if err != nil {
			return serverError(c, err)
		}
		questionsParams = append(questionsParams, &models.ExamQuestion{
			ExamID:   examParams.ID,
			Question: question.Title,
			Type:     "options",
			Answers:  answers,
		})
	}

	if len(questionsParams) > 0 {
		_, err = questionsParams.InsertAll(c.Request().Context(), config.db, boil.Blacklist("id"))
		if err != nil {
			return serverError(c, err)
		}
	}

	return dataResponse(c, types.SerializeExamResource(examParams, questionsParams))

}
