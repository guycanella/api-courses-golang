package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/guycanella/api-courses-golang/internal/domain"
	mysqlrepo "github.com/guycanella/api-courses-golang/internal/repository/mysql"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
)

var (
	flagTest = flag.Bool("test", false, "run migrate in test environment")
)

func main() {
	flag.Parse()

	isTest := *flagTest
	env := "development"

	if isTest {
		env = "test"
	}

	db, err := mysqlrepo.OpenDatabase(&mysqlrepo.Env{E: &env})
	if err != nil {
		log.Fatal(err)
	}

	gofakeit.Seed(time.Now().UnixNano())

	users := seedUsers(db)
	courses := seedCourses(db)
	enrolls := seedEnrollments(db, users, courses)

	if env == "test" {
		log.Printf("✅ Test seed ok: %d users, %d courses, %d enrollments\n", len(users), len(courses), len(enrolls))
		return
	}

	log.Printf("✅ seed ok: %d users, %d courses, %d enrollments\n", len(users), len(courses), len(enrolls))
}

func seedUsers(db *gorm.DB) []domain.User {
	var users []domain.User
	seenEmails := map[string]struct{}{}

	for len(users) < 3 {
		email := strings.ToLower(gofakeit.Email())
		if _, ok := seenEmails[email]; ok {
			continue
		}

		seenEmails[email] = struct{}{}
		users = append(users, domain.User{
			Name:  gofakeit.Name(),
			Email: email,
		})
	}

	if err := db.Create(&users).Error; err != nil {
		log.Fatal("creating users: ", err)
	}

	return users
}

func seedCourses(db *gorm.DB) []domain.Course {
	var courses []domain.Course

	for i := 0; i < 2; i++ {
		title := gofakeit.Sentence(4)
		desc := gofakeit.Paragraph(1, 3, 20, " ")
		courses = append(courses, domain.Course{
			Title:       title,
			Description: desc,
		})
	}

	if err := db.Create(&courses).Error; err != nil {
		log.Fatal("creating courses: ", err)
	}

	return courses
}

func seedEnrollments(db *gorm.DB, users []domain.User, courses []domain.Course) []domain.Enrollment {
	enrolls := []domain.Enrollment{
		{UserID: users[0].ID, CourseID: courses[0].ID},
		{UserID: users[1].ID, CourseID: courses[0].ID},
		{UserID: users[2].ID, CourseID: courses[1].ID},
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&enrolls).Error
	}); err != nil {
		log.Fatal("creating enrollments: ", err)
	}

	return enrolls
}
