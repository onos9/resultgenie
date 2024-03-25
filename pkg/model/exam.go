package model

type ExamType struct {
	ID           float64 `json:"id"`
	ActiveStatus float64 `json:"active_status"`
	Title        string  `json:"title"`
	Percentage   float64 `json:"percentage"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	UpdatedBy    float64 `json:"updated_by"`
	SchoolID     float64 `json:"school_id"`
	AcademicID   float64 `json:"academic_id"`
}

type ExamData struct {
	ExamTypes    []ExamType `json:"data"`
	Success bool       `json:"success"`
}
