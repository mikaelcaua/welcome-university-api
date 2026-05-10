package dto

type CreateStateRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CreateUniversityRequest struct {
	Name string `json:"name"`
}

type CreateCourseRequest struct {
	Name string `json:"name"`
}

type CreateSubjectRequest struct {
	Name string `json:"name"`
}
