package repo

import (
	"errors"

	"qa/internal/model"

	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Repo { return &Repo{DB: db} }

func (r *Repo) CreateQuestion(q *model.Question) error {
	return r.DB.Create(q).Error
}

func (r *Repo) ListQuestions() ([]model.Question, error) {
	var qs []model.Question
	err := r.DB.Order("created_at desc").Find(&qs).Error
	return qs, err
}

func (r *Repo) GetQuestionWithAnswers(id int) (*model.Question, error) {
	var q model.Question
	if err := r.DB.Preload("Answers").First(&q, id).Error; err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *Repo) DeleteQuestion(id int) error {
	return r.DB.Delete(&model.Question{}, id).Error
}

func (r *Repo) CreateAnswer(a *model.Answer) error {

	var q model.Question
	if err := r.DB.First(&q, a.QuestionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return r.DB.Create(a).Error
}

func (r *Repo) GetAnswer(id int) (*model.Answer, error) {
	var a model.Answer
	if err := r.DB.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repo) DeleteAnswer(id int) error {
	return r.DB.Delete(&model.Answer{}, id).Error
}
