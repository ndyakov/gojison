package gojison

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
)

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	var result map[string]string
	Error(w, errors.New("test error"), 0)

	if w.Code != 400 {
		t.Error("Expected response code to be 400 when 0 passed, but was %d.", w.Code)
	}

	json.NewDecoder(w.Body).Decode(&result)
	got := result["error"]
	if got != "test error" {
		t.Errorf("Expected error to be \"test error\", but was %s!", got)
	}

	w = httptest.NewRecorder()
	Error(w, errors.New("not found"), 404)
	if w.Code != 404 {
		t.Error("Expected response code to be 404 but was %d.", w.Code)
	}
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	var result map[string]string
	Success(w, "okey", 0)

	if w.Code != 200 {
		t.Error("Expected response code to be 400 when 0 passed, but was %d.", w.Code)
	}

	json.NewDecoder(w.Body).Decode(&result)
	got := result["success"]
	if got != "okey" {
		t.Errorf("Expected error to be \"okey\", but was %s!", got)
	}

	w = httptest.NewRecorder()
	Success(w, "created", 201)
	if w.Code != 201 {
		t.Error("Expected response code to be 201 but was %d.", w.Code)
	}
}
