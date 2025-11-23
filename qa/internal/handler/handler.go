package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"qa/internal/model"
	"qa/internal/repo"
)

type Handler struct {
	Repo *repo.Repo
}

func NewHandler(r *repo.Repo) *Handler {
	return &Handler{Repo: r}
}

func (h *Handler) ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/questions/", h.questionsHandler)
	mux.HandleFunc("/answers/", h.answersHandler)
	return mux
}

func (h *Handler) questionsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/questions")

	if path == "" || path == "/" {
		switch r.Method {
		case http.MethodGet:
			h.listQuestions(w, r)
		case http.MethodPost:
			h.createQuestion(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if len(parts) > 1 && parts[1] == "answers" {
		if r.Method == http.MethodPost {
			h.createAnswer(w, r, id)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getQuestion(w, r, id)
	case http.MethodDelete:
		h.deleteQuestion(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listQuestions(w http.ResponseWriter, _ *http.Request) {
	qs, err := h.Repo.ListQuestions()
	if err != nil {
		log.Println("listQuestions:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, qs)
}

func (h *Handler) createQuestion(w http.ResponseWriter, r *http.Request) {
	var q model.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(q.Text) == "" {
		http.Error(w, "text required", http.StatusBadRequest)
		return
	}

	if err := h.Repo.CreateQuestion(&q); err != nil {
		log.Println("createQuestion:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, q)
}

func (h *Handler) getQuestion(w http.ResponseWriter, _ *http.Request, id int) {
	q, err := h.Repo.GetQuestionWithAnswers(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	respondJSON(w, q)
}

func (h *Handler) deleteQuestion(w http.ResponseWriter, _ *http.Request, id int) {
	if err := h.Repo.DeleteQuestion(id); err != nil {
		log.Println("deleteQuestion:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) answersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/answers/")
	id, err := strconv.Atoi(strings.Trim(path, "/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getAnswer(w, r, id)
	case http.MethodDelete:
		h.deleteAnswer(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createAnswer(w http.ResponseWriter, r *http.Request, questionID int) {
	var a model.Answer
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	a.QuestionID = questionID

	if strings.TrimSpace(a.UserID) == "" || strings.TrimSpace(a.Text) == "" {
		http.Error(w, "user_id and text required", http.StatusBadRequest)
		return
	}

	if err := h.Repo.CreateAnswer(&a); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "question not found", http.StatusNotFound)
			return
		}
		log.Println("createAnswer:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, a)
}

func (h *Handler) getAnswer(w http.ResponseWriter, _ *http.Request, id int) {
	a, err := h.Repo.GetAnswer(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	respondJSON(w, a)
}

func (h *Handler) deleteAnswer(w http.ResponseWriter, _ *http.Request, id int) {
	if err := h.Repo.DeleteAnswer(id); err != nil {
		log.Println("deleteAnswer:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("respondJSON:", err)
	}
}
