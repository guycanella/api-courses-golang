// @title Go API Courses
// @version 1.0
// @description API to manage courses, users and enrollments
// @BasePath /
// @schemes http
// @host localhost:3333

// @contact.name Guilherme Arantes Canella
// @contact.url https://github.com/guycanella
// @contact.email guycanella@gmail.com

// @tag.name        courses
// @tag.description Operations on courses
package main

import (
	"log"
	"os"

	"github.com/guycanella/api-courses-golang/internal/handlers"
	"github.com/guycanella/api-courses-golang/internal/httpx"

	_ "github.com/guycanella/api-courses-golang/internal/docs"
	mysqlrepo "github.com/guycanella/api-courses-golang/internal/repository/mysql"

	swagger "github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	debug := os.Getenv("APP_DEBUG") == "true"
	httpx.SetDebug(debug)

	db, err := mysqlrepo.OpenDatabase(nil)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	h := handlers.NewCoursesHandler(db)

	app.Get("/courses", h.ListCourses)
	app.Get("/courses/:courseId", h.GetCourseByID)
	app.Post("/courses", h.CreateCourse)

	app.Get("/swagger/*", swagger.New(swagger.Config{
		Title: "Go API Courses",
	}))

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3333"
	}

	log.Fatal(app.Listen(":" + port))
}
