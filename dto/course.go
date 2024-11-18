package dto

type Course struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateCourseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateCourseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
