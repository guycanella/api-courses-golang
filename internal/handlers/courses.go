package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guycanella/api-courses-golang/internal/domain"
	"github.com/guycanella/api-courses-golang/internal/httpx"
	"gorm.io/gorm"
)

type CoursesHandler struct {
	db *gorm.DB
}

func NewCoursesHandler(db *gorm.DB) *CoursesHandler {
	return &CoursesHandler{
		db: db,
	}
}

// ListCourses godoc
// @Summary      Get courses
// @Description  Returns paginated courses, with optional search by title.
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        page   query     int    false  "Page (>=1)"              minimum(1) default(1)
// @Param        limit  query     int    false  "Items per page [1..100]" minimum(1) maximum(100) default(10)
// @Param        q      query     string false  "Title fragment (2..100)" minlength(2) maxlength(100)
// @Success      200    {object}  handlers.CoursesResponse
// @Failure      400    {object}  handlers.ErrorResponse
// @Failure      500    {object}  handlers.ErrorResponse
// @Router       /courses [get]
func (handler *CoursesHandler) ListCourses(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page",
		})
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit",
		})
	}

	q := strings.TrimSpace(ctx.Query("q", ""))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var courses []domain.Course
	tx := handler.db.Model(&domain.Course{})

	if q != "" {
		tx = tx.Where("title LIKE ?", "%"+q+"%")
	}

	var total int64
	tx.Count(&total)
	if err := tx.
		Order("created_at desc").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&courses).Error; err != nil {
		return httpx.InternalServerError(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"data":  courses,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GetCourseByID godoc
// @Summary      Get course by ID
// @Tags         courses
// @Produce      json
// @Param        courseId  path      string  true  "Course ID"  format(uuid)
// @Success      200       {object}  handlers.CourseResponse
// @Failure      400       {object}  handlers.ErrorResponse
// @Failure      404       {object}  handlers.ErrorResponse
// @Failure      500       {object}  handlers.ErrorResponse
// @Router       /courses/{courseId} [get]
func (handler *CoursesHandler) GetCourseByID(ctx *fiber.Ctx) error {
	var course domain.Course
	courseId := ctx.Params("courseId")

	if courseId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "courseId is required",
		})
	}

	if _, fieldErr := uuid.Parse(courseId); fieldErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid courseId",
		})
	}

	if err := handler.db.First(&course, "id = ?", courseId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "course not found",
			})
		}

		return httpx.InternalServerError(ctx, err)
	}

	return ctx.JSON(fiber.Map{"course": course})
}

var validate = validator.New()

// CreateCourse godoc
// @Summary      Create course
// @Description  Creates a course and returns only the courseId and the Location header.
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        payload  body      handlers.CreateCourseDTO  true  "New course"
// @Success      201      {object}  handlers.CreatedIDResponse
// @Failure      400      {object}  handlers.ErrorResponse
// @Failure      409      {object}  handlers.ErrorResponse
// @Failure      422      {object}  handlers.ValidationErrorResponse
// @Failure      500      {object}  handlers.ErrorResponse
// @Router       /courses [post]
func (handler *CoursesHandler) CreateCourse(ctx *fiber.Ctx) error {
	var Body struct {
		Title       string `json:"title" validate:"required,min=3"`
		Description string `json:"description" validate:"omitempty,min=3"`
	}

	if err := ctx.BodyParser(&Body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON body",
		})
	}

	if err := validate.Struct(&Body); err != nil {
		errs := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errs[field] = "is required"
			case "min":
				errs[field] = "too short"
			default:
				errs[field] = "is invalid"
			}
		}

		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": errs,
		})
	}

	title := strings.TrimSpace(Body.Title)
	desc := strings.TrimSpace(Body.Description)

	course := domain.Course{
		ID:          uuid.NewString(),
		Title:       title,
		Description: desc,
	}

	if err := handler.db.Create(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "title already exists",
			})
		}

		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "title already exists",
			})
		}

		return httpx.InternalServerError(ctx, err)
	}

	ctx.Location("/courses/" + course.ID)
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"courseId": course.ID,
	})
}
