package renderer

var paths = []string{
	"template/index.html",
	"template/header.html",
	"template/studentInfo.html",
	"template/record.html",
	"template/score.html",
	"template/rating.html",
	"template/remark.html",
}

var Remarks = []string{"Fair", "Fair", "Good", "very good", "Very Good", "Excellent", "Excellent"}
var Levels = []string{
	"w-0/12",
	"w-1/12",
	"w-2/12",
	"w-3/12",
	"w-4/12",
	"w-5/12",
	"w-6/12",
	"w-7/12",
	"w-8/12",
	"w-9/12",
	"w-10/12",
	"w-11/12",
}

type School struct {
	SchoolName   string
	SchoolCity   string
	SchoolRegion string
	ResultDesc   string
	VacationDate string
	Term         string
}

type Student struct {
	ID           float64
	FullName     string
	StudentPhoto string
	AdminNo      string
	ClassName    string
	SectionName  string
	Term         string
	Arm          string
	SessionYear  string
	Opened       string
	Present      string
	Absent       string
}

type Record struct {
	Subject    string
	Objectives []string
	Outcome    string
	Mta        float64
	Ca         float64
	Oral       float64
	Exam       float64
	Score      int
	Grade      string
	Color      string
}

type Score struct {
	Total   float64
	Average float64
	Highest float64
	Lowest  float64
	Grading string
}

type Rating struct {
	Attribute string
	Rate      int
	Remark    string
	Level     string
}

type Remark struct {
	Name    string
	Comment string
}
