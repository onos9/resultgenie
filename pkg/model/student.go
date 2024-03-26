package model

type Response struct {
	Data    *Data `json:"data"`
	Success *bool `json:"success"`
}
type Data struct {
	Student *Student `json:"student_data"`
	Score   *Score   `json:"scores"`
}

type Score struct {
	Total   float64      `json:"total_score"`
	Average float64      `json:"average"`
	Lowest  ClassAverage `json:"lowest"`
	Highest ClassAverage `json:"highest"`
}

type ClassAverage struct {
	StudentID float64 `json:"student_id"`
	Average   float64 `json:"class_average"`
}

type Student struct {
	ID              float64 `json:"id"`
	AdmissionNo     float64 `json:"admission_no"`
	AcademicID      float64 `json:"academic_id"`
	FullName        string  `json:"full_name"`
	StudentPhoto    string  `json:"student_photo"`
	ParentID        float64 `json:"parent_id"`
	StudentCategory float64 `json:"student_category_id"`
	ClassID         float64 `json:"class_id"`
	SectionID       float64 `json:"section_id"`
	SchoolID        float64 `json:"school_id"`
	ClassName       string  `json:"class_name"`
	SectionName     string  `json:"section_name"`

	CustomField CustomField `json:"custom_field"`
	Parent      Parent      `json:"parents"`
	School      School      `json:"school"`
	Category    Category    `json:"category"`

	Subjects []Subject  `json:"subjects"`
	Timeline []Timeline `json:"timeline"`
}

type Subject struct {
	ID          float64 `json:"id"`
	SubjectName string  `json:"subject_name"`
	SubjectCode string  `json:"subject_code"`
	Score       float64 `json:"score"`
	TeacherID   float64 `json:"teacher_id"`
	SchoolID    float64 `json:"school_id"`
	ClassID     float64 `json:"class_id"`
	SectionID   float64 `json:"section_id"`
	AcademicID  float64 `json:"academic_id"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	Marks []Mark `json:"marks"`
}
type Mark struct {
	StudentRollNo       float64 `json:"student_roll_no"`
	StudentAddmissionNo float64 `json:"student_addmission_no"`
	TotalMarks          float64 `json:"total_marks"`
	IsAbsent            float64 `json:"is_absent"`
	TeacherRemarks      string  `json:"teacher_remarks" validate:"required"`
	ExamTitle           *string `json:"exam_title"`
	StudentRecordID     float64 `json:"student_record_id"`
	ExamSetupID         float64 `json:"exam_setup_id"`
	ExamTermID          float64 `json:"exam_term_id"`
	StudentID           float64 `json:"student_id"`
	SubjectID           float64 `json:"subject_id"`
	SectionID           float64 `json:"section_id"`
	ClassID             float64 `json:"class_id"`
	SchoolID            float64 `json:"school_id"`
	AcademicID          float64 `json:"academic_id"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Category struct {
	ID         float64 `json:"id"`
	Category   string  `json:"category_name"`
	SchoolID   float64 `json:"school_id"`
	AcademicID float64 `json:"academic_id"`

	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	UpdatedBy float64 `json:"updated_by"`
}

type Parent struct {
	ID             float64 `json:"id"`
	FatherName     string  `json:"fathers_name"`
	MotherName     string  `json:"mothers_name"`
	GuardianEmail  string  `json:"guardians_email"`
	GuardianMobile string  `json:"guardians_mobile"`
	ParentUser     string  `json:"parent_user"`
}

type Timeline struct {
	ID             float64 `json:"id"`
	StaffStudentID float64 `json:"staff_student_id"`
	SchoolID       float64 `json:"school_id"`
	AcademicID     float64 `json:"academic_id"`
	Description    string  `json:"description"`
	File           string  `json:"file"`
	Type           string  `json:"type"`
	Date           string  `json:"date"`
	Visible        float64 `json:"visible_to_student"`

	CreatedBy float64 `json:"created_by"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	UpdatedBy float64 `json:"updated_by"`
}

type Section struct {
	ID          float64 `json:"id"`
	SectionName string  `json:"section_name"`
	Description string  `json:"description"`
	SchoolID    float64 `json:"school_id"`
	AcademicID  float64 `json:"academic_id"`
}

type Academic struct {
	ID                   float64 `json:"id"`
	Year                 string  `json:"year"`
	Title                string  `json:"title"`
	StartingDate         string  `json:"starting_date"`
	EndingDate           string  `json:"ending_date"`
	CopyWithAcademicYear string  `json:"copy_with_academic_year"`
	ActiveStatus         float64 `json:"active_status"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
	CreatedBy            float64 `json:"created_by"`
	UpdatedBy            float64 `json:"updated_by"`
	SchoolID             float64 `json:"school_id"`
}

type School struct {
	ID              float64 `json:"id"`
	SchoolName      string  `json:"school_name"`
	CreatedBy       float64 `json:"created_by"`
	UpdatedAt       string  `json:"updated_at"`
	SchoolCode      string  `json:"school_code"`
	UpdatedBy       float64 `json:"updated_by"`
	IsEmailVerified float64 `json:"is_email_verified"`
	StartingDate    string  `json:"starting_date"`
	EndingDate      string  `json:"ending_date"`
	PackageID       float64 `json:"package_id"`
	PlanType        string  `json:"plan_type"`
	ContactType     string  `json:"contact_type"`
	ActiveStatus    float64 `json:"active_status"`
	IsEnabled       string  `json:"is_enabled"`
	CreatedAt       string  `json:"created_at"`
	AcademicID      float64 `json:"academic_id"`
	SchoolID        float64 `json:"school_id"`
}

type CustomField struct {
	ExamType                  string `json:"Exam Type"`
	AdherentAndIndependent    string `json:"Adherent and Independent"`
	SelfControlAndInteraction string `json:"Self-control and Interaction"`
	FlexibilityAndCreativity  string `json:"Flexibility and Creativity"`
	Meticulous                string `json:"Meticulous"`
	Neatness                  string `json:"Neatness"`
	OverallProgress           string `json:"Overall Progress"`
	DaysSchoolOpened          string `json:"Days School Opened"`
	DaysPresent               string `json:"Days Present"`
	DaysAbsent                string `json:"Days Absent"`
}
