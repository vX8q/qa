package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"qa/internal/handler"
	"qa/internal/model"
	"qa/internal/repo"
)

func setupServer(t *testing.T) http.Handler {
	t.Helper()

	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, gdb.AutoMigrate(&model.Question{}, &model.Answer{}))

	r := repo.New(gdb)
	h := handler.NewHandler(r)

	return h.ServeMux()
}

func TestCreateQuestionAndAnswer(t *testing.T) {
	srv := setupServer(t)

	log.Println("=== Создаем вопрос ===")

	q := map[string]string{"text": "Тестовый вопрос"}
	b, err := json.Marshal(q)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewReader(b))
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)

	var created model.Question
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &created))
	require.NotZero(t, created.ID)
	log.Printf("Создан вопрос ID=%d\n", created.ID)

	log.Println("=== Создаем ответ ===")

	ans := map[string]string{
		"user_id": "00000000-0000-0000-0000-000000000001",
		"text":    "Ответ",
	}
	b2, err := json.Marshal(ans)
	require.NoError(t, err)

	req2 := httptest.NewRequest(http.MethodPost, "/questions/"+strconv.Itoa(created.ID)+"/answers/", bytes.NewReader(b2))
	rr2 := httptest.NewRecorder()
	srv.ServeHTTP(rr2, req2)
	require.Equal(t, http.StatusCreated, rr2.Code)

	var createdAnswer model.Answer
	require.NoError(t, json.Unmarshal(rr2.Body.Bytes(), &createdAnswer))
	require.NotZero(t, createdAnswer.ID)
	log.Printf("Создан ответ ID=%d для вопроса ID=%d\n", createdAnswer.ID, created.ID)
}
