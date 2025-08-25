package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/guycanella/api-courses-golang/internal/domain"
)

func TestGetCourses200(t *testing.T) {
	t.Helper()
	app, _ := setupAll(t)

	req := httptest.NewRequest("GET", "/courses", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to TestGetCourses200 request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to TestGetCourses200 status=%d want=%d", resp.StatusCode, http.StatusOK)
	}

	var out struct {
		Data  []domain.Course `json:"data"`
		Page  int             `json:"page"`
		Limit int             `json:"limit"`
		Total int64           `json:"total"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("Decode: %v", err)
	}

	value := reflect.ValueOf(out)
	typ := reflect.TypeOf(out)

	for i := 0; i < value.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := value.Field(i)

		if fieldValue.IsZero() {
			t.Fatalf("Missing key %q in response: %#v", fieldType.Name, out)
		}
	}
}

func TestGetCourses400_InvalidPage(t *testing.T) {
	t.Helper()
	app, _ := setupAll(t)

	req := httptest.NewRequest("GET", "/courses?page=invalid", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to TestGetCourses400_InvalidPage request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Failed to TestGetCourses400_InvalidPage status=%d want=%d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestGetCourses400_InvalidLimit(t *testing.T) {
	t.Helper()
	app, _ := setupAll(t)

	req := httptest.NewRequest("GET", "/courses?limit=invalid", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to TestGetCourses400_InvalidLimit request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Failed to TestGetCourses400_InvalidLimit status=%d want=%d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestGetCourses500_InternalServerError(t *testing.T) {
	t.Helper()
	app, db := setupAll(t)

	sqlDB, _ := db.DB()
	_ = sqlDB.Close()

	req := httptest.NewRequest("GET", "/courses", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to TestGetCourses500_InternalServerError request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Failed to TestGetCourses500_InternalServerError status=%d want=%d", resp.StatusCode, http.StatusInternalServerError)
	}

	var out struct {
		Error string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if strings.TrimSpace(out.Error) == "" {
		t.Fatalf("expected non-empty error message")
	}
}
