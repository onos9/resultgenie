package result

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math"
	"os"
	"repot/pkg/chromium"
	"repot/pkg/edusms"
	"repot/pkg/model"
	"repot/pkg/utils"
	"strconv"
	"strings"
)

type Result struct {
	School  School   `json:"school"`
	Student Student  `json:"student"`
	Score   Score    `json:"score"`
	Records []Record `json:"records"`
	Ratings []Rating `json:"ratings"`
	Remark  Remark   `json:"remark"`

	teacherId float64
	data      *model.Data
	client    *edusms.Client
	ErrMsg    error
}

func New() (*Result, error) {
	return &Result{
		client: edusms.GetInstance(),
	}, nil
}

func (r *Result) Render(data *model.Data) ([]byte, error) {
	r.data = data
	err := r.processStudentData()
	if err != nil {
		return nil, err
	}

	err = r.processSchoolData()
	if err != nil {
		return nil, err
	}

	err = r.processRecordData()
	if err != nil {
		return nil, err
	}

	err = r.processScoreData()
	if err != nil {
		return nil, err
	}

	err = r.processRatingData()
	if err != nil {
		return nil, err
	}

	file, err := os.Create("generated/index.html")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		return nil, err
	}

	err = tmpl.ExecuteTemplate(file, "index", *r)
	if err != nil {
		return nil, err
	}

	b, err := r.generate()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *Result) generate() ([]byte, error) {
	chrome := chromium.New(r.client.HTTPClient)
	data, err := chrome.SendHTML("@generated/index.html")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Result) processSchoolData() error {
	sch := r.data.Student.School
	month := strings.Split(r.data.Student.CustomField.ExamType, "-")[1]
	_, state := utils.GetLocation(sch.SchoolName)
	r.School = School{
		SchoolCity:   "Makurdi",
		SchoolRegion: state,
		ResultDesc:   "TERMLY SUMMARY OF PROGRESS REPORT",
		Term:         r.Student.Term,
		VacationDate: strings.Trim(month, " "),
	}

	r.School.SchoolName = sch.SchoolName

	return nil
}

func (r *Result) processStudentData() error {
	std := r.data.Student
	r.Student.ID = std.ID
	r.Student.AdminNo = fmt.Sprintf("%04d", int(std.AdmissionNo))
	r.Student.FullName = std.FullName
	r.Student.StudentPhoto = std.StudentPhoto
	r.Student.ClassName = std.ClassName
	r.Student.SectionName = std.SectionName

	feild := std.CustomField
	r.Student.Term = feild.ExamType
	r.Student.Opened = feild.DaysSchoolOpened
	r.Student.Present = feild.DaysPresent
	r.Student.Absent = feild.DaysAbsent

	cat := std.Category
	r.Student.Arm = cat.Category

	return nil
}

func (r *Result) processRecordData() error {
	subs := r.data.Student.Subjects

	for _, sub := range subs {
		marksData := sub.Marks

		score := 0
		marks := make(map[string]float64)
		subject_code := strings.ToUpper(sub.SubjectCode)
		for _, m := range marksData {
			// if m.ExamTitle == nil || m.TotalMarks == 0 {
			// 	return fmt.Errorf("[Error]: %s for %s can not be 0", *m.ExamTitle, sub.SubjectName)
			// }

			marks[*m.ExamTitle] = m.TotalMarks
			if subject_code == "BIBLE" && m.TeacherRemarks != "" {
				r.Remark.Name = "Teacher Remark"
				r.Remark.Comment = m.TeacherRemarks
			}

			score += int(m.TotalMarks)
		}
		grade, color := r.getGrade(score)
		record := Record{
			Score:   score,
			Outcome: grade,
			Grade:   grade,
			Color:   color,
		}

		record.Subject = sub.SubjectName
		var ok bool
		if r.Student.Arm == "GRADERS" {
			record.Mta = marks["FIRST CA"]
			record.Ca = marks["SECOND CA"]
			record.Oral = marks["ORAL"]
			record.Exam = marks["EXAM"]

		} else {
			obj := []model.Objective{}
			err := utils.UnmarshalJson("template/objectives.json", &obj)
			if err != nil {
				return err
			}

			for _, obj := range obj {
				class := strings.ToUpper(obj.ClassName)
				code := strings.ToUpper(obj.SubjectCode)
				if class == r.Student.ClassName && code == subject_code {
					record.Objectives = strings.Split(obj.Text, "|")
					break
				}
			}

			if record.Exam, ok = marks["EXAM"]; !ok {
				return fmt.Errorf("[Error]: EXAM Record for (%s) not found", sub.SubjectName)
			}
		}

		if subject_code == "BIBLE" && r.Remark.Comment == "" {
			return errors.New("[Error]: Teacher Remark not found in BIBLE, Please check your result")
		}

		r.teacherId = sub.TeacherID
		r.Records = append(r.Records, record)
	}

	return nil
}

func (r *Result) processScoreData() error {
	score := r.data.Score
	if score.Lowest.Average == 0 {
		return errors.New("[Error]: Lowest Class Average can not be 0, ID: " + fmt.Sprint(score.Lowest.StudentID))
	}

	r.Score.Total = score.Total
	r.Score.Average = score.Average
	r.Score.Highest = score.Highest.Average
	r.Score.Lowest = score.Lowest.Average

	if r.Student.Arm == "EYFS" {
		r.Score.Grading = "EMERGING(0-80) EXPECTED(81-90) EXCEEDING(91-100)"
	} else {
		r.Score.Grading = "A(94-100) B(86-93) C(77-85) D(70-76) E(0-69)"
	}
	return nil
}

func (r *Result) processRatingData() error {
	var data map[string]interface{}
	inrec, err := json.Marshal(r.data.Student.CustomField)
	if err != nil {
		return err
	}
	err = json.Unmarshal(inrec, &data)
	if err != nil {
		return err
	}

	for key, val := range data {
		if key == "Days School Opened" || key == "Days Present" || key == "Days Absent" || key == "Exam Type" {
			continue
		}

		if val != nil {
			rating := Rating{
				Attribute: key,
			}
			rate, err := strconv.Atoi(val.(string))
			if err != nil {
				return err
			}
			rating.Rate = int((float64(rate) / 5.0) * 100)
			span := math.Floor((float64(rate) / 5.0) * 12)
			if span >= 12 {
				rating.Level = "w-full"
			} else {
				rating.Level = Levels[int(span)]
			}
			rating.Remark = Remarks[rate]
			r.Ratings = append(r.Ratings, rating)
		}
	}

	return nil
}

func (r *Result) getGrade(score int) (string, string) {
	if r.Student.Arm == "EYFS" {
		// EMERGING(0-80) EXPECTED(81-90) EXCEEDING(91-100)
		if score >= 0 && score <= 80 {
			return "EMERGING", "bg-purple-200"
		} else if score >= 81 && score <= 90 {
			return "EXPECTED", "bg-blue-200"
		} else if score >= 91 && score <= 100 {
			return "EXCEEDING", "bg-red-200"
		}

	} else {
		// "A(94-100) B(86-93) C(77-85) D(70-76) E(0-69)"
		if score >= 94 && score <= 100 {
			return "A", "bg-purple-200"
		} else if score >= 86 && score <= 93 {
			return "B", "bg-blue-200"
		} else if score >= 77 && score <= 85 {
			return "C", "bg-yellow-200"
		} else if score >= 70 && score <= 76 {
			return "D", "bg-orange-200"
		} else if score >= 0 && score <= 69 {
			return "E", "bg-red-200"
		}
	}

	return "Outstanding", "bg-red-200"
}

func (r *Result) Error(msg string, subject *string) error {

	var details []string

	details = append(details, fmt.Sprintf("Student: %s", r.Student.FullName))
	details = append(details, fmt.Sprintf("Student No: %s", r.Student.AdminNo))
	details = append(details, fmt.Sprintf("Student ID: %.0f", r.Student.ID))
	details = append(details, fmt.Sprintf("Error: %s", msg))

	if subject != nil {
		details = append(details, fmt.Sprintf("Subject: %s", *subject))
	}
	err := fmt.Errorf(strings.Join(details, "\n"))
	return err
}
