package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func TestCreateCourse201_Created(t *testing.T) {
	t.Helper()
	app, _ := setupAll(t)
	title := "test-" + uuid.NewString()
	desc := gofakeit.Paragraph(1, 3, 20, " ")

	in := map[string]string{
		"title":       title,
		"description": desc,
	}

	body, _ := json.Marshal(in)

	req := httptest.NewRequest("POST", "/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to TestCreateCourse201_Created request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to TestCreateCourse201_Created status=%d want=%d", resp.StatusCode, http.StatusCreated)
	}

	var out struct {
		CourseID string `json:"courseId"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if out.CourseID == "" {
		t.Fatal("Expected non-empty id")
	}

	loc := resp.Header.Get("Location")
	if !strings.HasPrefix(loc, "/courses/") {
		t.Fatalf("Missing/invalid Location header: %q", loc)
	}

}

func TestCreateCourse400_InvalidJSON(t *testing.T) {
	t.Helper()
	app, _ := setupAll(t)

	body := []byte(`{"title":`)

	req := httptest.NewRequest("POST", "/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to TestCreateCourse400_InvalidJSON request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Failed to TestCreateCourse400_InvalidJSON status=%d want=%d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestCreateCourse409_Conflict(t *testing.T) {
	t.Helper()
	app, _ := setupAll(t)
	title := "dup-" + gofakeit.Sentence(4)

	in1 := map[string]string{
		"title":       title,
		"description": "First",
	}

	in2 := map[string]string{
		"title":       title,
		"description": "Second",
	}

	body1, _ := json.Marshal(in1)
	body2, _ := json.Marshal(in2)

	req1 := httptest.NewRequest("POST", "/courses", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")

	resp1, err1 := app.Test(req1)
	if err1 != nil {
		t.Fatalf("Failed to TestCreateCourse409_Conflict request: %v", err1)
	}
	defer resp1.Body.Close()
	// Avoid false positive if the first request fails
	if resp1.StatusCode != http.StatusCreated {
		t.Fatalf("Want 201, got %d", resp1.StatusCode)
	}

	req2 := httptest.NewRequest("POST", "/courses", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")

	resp2, err2 := app.Test(req2)
	if err2 != nil {
		t.Fatalf("Failed to TestCreateCourse409_Conflict request: %v", err2)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusConflict {
		t.Fatalf("Failed to TestCreateCourse409_Conflict status=%d want=%d", resp2.StatusCode, http.StatusConflict)
	}
}

type validationErrResp struct {
	Errors map[string]string `json:"errors"`
}

func TestCreateCourse422_Validation(t *testing.T) {
	t.Helper()
	app, _ := setupAll(t)

	tests := []struct {
		name    string
		payload map[string]string
		want    map[string]string // campos esperados no "errors"
	}{
		{
			name:    "missing title",
			payload: map[string]string{"description": "ok desc"},
			want:    map[string]string{"title": "is required"},
		},
		{
			name:    "short title",
			payload: map[string]string{"title": "ab"},
			want:    map[string]string{"title": "too short"},
		},
		{
			name:    "short description",
			payload: map[string]string{"title": "Valid " + gofakeit.Sentence(4), "description": "ab"}, // < 3
			want:    map[string]string{"description": "too short"},
		},
		{
			name:    "multiple errors",
			payload: map[string]string{"title": "", "description": "ab"},
			want:    map[string]string{"title": "is required", "description": "too short"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest("POST", "/courses", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed TestCreateCourse422_Validation request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusUnprocessableEntity {
				t.Fatalf("Failed TestCreateCourse422_Validation status=%d want=%d", resp.StatusCode, http.StatusUnprocessableEntity)
			}

			var out validationErrResp
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				t.Fatalf("Decode: %v", err)
			}

			for field, wantMsg := range tc.want {
				got, ok := out.Errors[field]
				if !ok {
					t.Fatalf("Expected field %q in errors; got=%v", field, out.Errors)
				}
				if got != wantMsg {
					t.Fatalf("Errors[%q]=%q want=%q; full=%v", field, got, wantMsg, out.Errors)
				}
			}
		})
	}
}

func TestCreateCourse500_InternalServerError(t *testing.T) {
	t.Helper()
	app, db := setupAll(t)

	sqlDB, _ := db.DB()
	_ = sqlDB.Close()

	body := []byte(`{"title":"ok-` + uuid.NewString() + `"}`)
	req := httptest.NewRequest("POST", "/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed TestCreateCourse500_InternalServerError request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Failed Internal Server Error: status=%d want=%d", resp.StatusCode, http.StatusInternalServerError)
	}
}
