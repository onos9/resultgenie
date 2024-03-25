package model

type Objective struct {
	ID          float64 `json:"id"`
	ExamID      float64 `json:"exam_id"`
	ClassName   string  `json:"class_name"`
	ClassID     float64 `json:"class_id"`
	SubjectCode string  `json:"subject_code"`
	SubjectID   float64 `json:"subject_id"`
	Text        string  `json:"text"`
}
