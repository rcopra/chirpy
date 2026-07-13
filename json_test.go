package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithJSON(t *testing.T) {
	rec := httptest.NewRecorder()

	payload := map[string]string{"message": "hello"}
	if err := respondWithJSON(rec, http.StatusCreated, payload); err != nil {
		t.Fatalf("respondWithJSON returned error: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusCreated)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}

	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if got["message"] != "hello" {
		t.Errorf("body message = %q, want %q", got["message"], "hello")
	}
}

func TestRespondWithJSONUnmarshalablePayload(t *testing.T) {
	rec := httptest.NewRecorder()

	if err := respondWithJSON(rec, http.StatusOK, make(chan int)); err == nil {
		t.Error("respondWithJSON should return an error for an unmarshalable payload")
	}
}

func TestRespondWithError(t *testing.T) {
	rec := httptest.NewRecorder()

	if err := respondWithError(rec, http.StatusBadRequest, "something went wrong"); err != nil {
		t.Fatalf("respondWithError returned error: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if got["error"] != "something went wrong" {
		t.Errorf("body error = %q, want %q", got["error"], "something went wrong")
	}
}
