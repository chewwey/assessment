//go:build unit
// +build unit

package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

//test require
// story1: create expenses DONE
// story2: get expenses by id DONE
// story3: update expenses by id DONE
// story4: get all expenses DONE

func TestUnCreateExpensesHandler(t *testing.T) {
	data := bytes.NewBufferString(`{
		"title": "test-title",
		"amount": 1234,
		"note": "test-note", 
		"tags": ["test-tag1", "test-tag2"]
	}`)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", data)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	x := float64(1234)
	mockRows := sqlmock.NewRows([]string{"id"}).AddRow("1")

	expected := "{\"id\":1,\"title\":\"test-title\",\"amount\":1234,\"note\":\"test-note\",\"tags\":[\"test-tag1\",\"test-tag2\"]}"

	db, mock, err := sqlmock.New()
	mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id")).WithArgs(
		"test-title", x, "test-note", pq.Array([]string{"test-tag1", "test-tag2"})).WillReturnRows(mockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}

	h := Handler{DB: db}
	ct := e.NewContext(req, rec)

	h.CreateExpensesHandler(ct)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUnGetExpensesIdHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses/:id", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow("1", "test-title", "1234", "test-note", pq.Array([]string{"test-tag1", "test-tag2"}))
	expected := "{\"id\":1,\"title\":\"test-title\",\"amount\":1234,\"note\":\"test-note\",\"tags\":[\"test-tag1\",\"test-tag2\"]}"

	db, mock, err := sqlmock.New()
	get := mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1"))

	get.ExpectQuery().WithArgs("1").WillReturnRows(mockRows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}

	h := Handler{DB: db}
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	err = h.GetExpensesByIdHandler(ctx)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUnUpdateExpensesHandler(t *testing.T) {
	data := bytes.NewBufferString(`{
        "title": "test-updated-title",
        "amount": 1234.00,
        "note": "test-updated-note",
        "tags": ["test-updated-tag1", "test-updated-tag2"]
    }`)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses/:id", data)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	expected := "{\"id\":1,\"title\":\"test-updated-title\",\"amount\":1234,\"note\":\"test-updated-note\",\"tags\":[\"test-updated-tag1\",\"test-updated-tag2\"]}"

	db, mock, err := sqlmock.New()
	get := mock.ExpectPrepare(regexp.QuoteMeta("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id=$5;")) //row id and rowAffected
	get.ExpectExec().WithArgs("test-updated-title", 1234.00, "test-updated-note", pq.Array([]string{"test-updated-tag1", "test-updated-tag2"}), "1").WillReturnResult(sqlmock.NewResult(1, 1))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}

	h := Handler{DB: db}
	ctx := e.NewContext(req, rec)

	ctx.SetParamNames("id")
	ctx.SetParamValues("1")
	err = h.UpdateByIdHandler(ctx)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUnGetAllExpensesHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "test-title", "1234", "test-note", pq.Array([]string{"test-tag1", "test-tag2"}))

	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT (.+) FROM expenses").WillReturnRows(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a database connection", err)
	}
	h := Handler{DB: db}
	ctx := e.NewContext(req, rec)

	expected := "[{\"id\":1,\"title\":\"test-title\",\"amount\":1234,\"note\":\"test-note\",\"tags\":[\"test-tag1\",\"test-tag2\"]}]"

	err = h.GetAllUserHandler(ctx)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
