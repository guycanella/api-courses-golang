package handlers_test

import (
	"testing"

	"github.com/guycanella/api-courses-golang/internal/handlers"
	mysqlrepo "github.com/guycanella/api-courses-golang/internal/repository/mysql"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupAll(t *testing.T) (*fiber.App, *gorm.DB) {
	t.Helper()

	mode := "test"
	db, err := mysqlrepo.OpenDatabase(&mysqlrepo.Env{E: &mode})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	db = db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	h := handlers.NewCoursesHandler(db)

	app.Post("/courses", h.CreateCourse)
	app.Get("/courses/:courseId", h.GetCourseByID)
	app.Get("/courses", h.ListCourses)

	return app, db
}
