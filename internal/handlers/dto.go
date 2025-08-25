package handlers

import "github.com/guycanella/api-courses-golang/internal/domain"

type CourseDoc struct {
	ID          string `json:"id" example:"4e70d7c4-5f5b-4f5a-9c9f-0e0b4a7c0d18"`
	Title       string `json:"title" example:"Go para Iniciantes"`
	Description string `json:"description" example:"Curso introdutório de Go"`
	CreatedAt   string `json:"created_at" example:"2025-08-20T15:04:05Z"`
}

type CreateCourseDTO struct {
	Title       string `json:"title" example:"Curso de React"`
	Description string `json:"description" example:"Alguma descrição"`
}

type CourseResponse struct {
	Course domain.Course `json:"course"`
}

type CoursesResponse struct {
	Data  []CourseDoc `json:"data"`
	Page  int         `json:"page"  example:"1"`
	Limit int         `json:"limit" example:"10"`
	Total int64       `json:"total" example:"42"`
}

type CreatedIDResponse struct {
	CourseID string `json:"courseId" example:"4e70d7c4-5f5b-4f5a-9c9f-0e0b4a7c0d18"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

type ValidationErrorResponse struct {
	Errors map[string]string `json:"errors" example:"{\"title\":\"is required\"}"`
}
