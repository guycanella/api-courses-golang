package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/guycanella/api-courses-golang/internal/domain"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func mustCreateCourse(t *testing.T, db *gorm.DB, title string) domain.Course {
	t.Helper()
	course := domain.Course{Title: title, Description: ""}

	if err := db.Create(&course).Error; err != nil {
		t.Fatalf("Failed to create course: %v", err)
	}

	return course
}

func TestGetCourseByID_200(t *testing.T) {
	app, db := setupAll(t)
	title := "test-" + gofakeit.Sentence(4)
	course := mustCreateCourse(t, db, title)

	req := httptest.NewRequest("GET", "/courses/"+course.ID, nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to TestGetCourseByID_200 test request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var out struct {
		Course domain.Course `json:"course"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if out.Course.ID != course.ID {
		t.Fatalf("Unexpected courseID: got=%s want=%s", out.Course.ID, course.ID)
	}
}

func TestGetCourseByID_404NotFound(t *testing.T) {
	app, _ := setupAll(t)
	unknown := uuid.NewString()

	req := httptest.NewRequest("GET", "/courses/"+unknown, nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to TestGetCourseByID_404NotFound test request: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", resp.StatusCode)
	}
}

func TestGetCourseByID_400_InvalidUUID(t *testing.T) {
	app, _ := setupAll(t)

	req := httptest.NewRequest("GET", "/courses/1233", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to TestGetCourseByID_400_InvalidUUID test request: %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", resp.StatusCode)
	}
}

func TestGetCourseByID_500_DBError(t *testing.T) {
	app, db := setupAll(t)

	mysqlDB, _ := db.DB()
	_ = mysqlDB.Close()

	req := httptest.NewRequest("GET", "/courses/"+uuid.NewString(), nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to TestGetCourseByID_500_DBError test request: %v", err)
	}

	if resp.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", resp.StatusCode)
	}
}
