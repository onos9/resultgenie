package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"math"
	"mime/multipart"
	"os"
	"repot/pkg/chromium"
	"repot/pkg/httpclient"
	"repot/pkg/utils"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	School  School   `json:"school"`
	Student Student  `json:"student"`
	Score   Score    `json:"score"`
	Records []Record `json:"records"`
	Ratings []Rating `json:"ratings"`
	Remark  Remark   `json:"remark"`

	teacherId  float64
	scores     map[string]interface{}
	resultData map[string]interface{}
	client     *httpclient.Client
}

func New(c *httpclient.Client) (*Result, error) {
	return &Result{
		client: c,
	}, nil
}

func (r *Result) Uplaod() error {
	fmt.Println("Uploading result...")
	chrome := chromium.New(r.client)
	data, err := chrome.SendHTML("@generated/index.html")
	if err != nil {
		return err
	}

	date, err := time.Parse("January 2, 2006", r.School.VacationDate)
	if err != nil {
		return err
	}

	timeline := map[string]string{
		"visible_to_student": fmt.Sprintf("%d", 1),
		"student_id":         fmt.Sprintf("%d", int(r.Student.ID)),
		"title":              r.School.Term,
		"date":               date.Format("2006-01-02"),
		"description":        r.School.ResultDesc,
	}

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	for key, value := range timeline {
		_ = w.WriteField(key, value)
	}

	part, err := w.CreateFormFile("document_file", "result.pdf")
	if err != nil {
		return err
	}

	_, err = io.Copy(part, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	r.client.SetHeader("Content-Type", w.FormDataContentType())
	_ = w.Close()

	body, err := r.client.Post("/api/marks-grade", buf)
	if err != nil {
		return err
	}

	data, err = io.ReadAll(body)
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	fmt.Println("Done!")

	return nil
}

func (r *Result) CacheData() error {
	return nil
}

func (r *Result) Render(data interface{}) error {
	d := data.(map[string]interface{})
	var ok bool
	r.resultData, ok = d["student_data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to get result_data data")
	}

	r.scores, ok = d["scores"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to get scores data")
	}

	err := r.processStudentData()
	if err != nil {
		return err
	}

	err = r.processSchoolData()
	if err != nil {
		return err
	}

	err = r.processRecordData()
	if err != nil {
		return err
	}

	err = r.processScoreData()
	if err != nil {
		return err
	}

	err = r.processRatingData()
	if err != nil {
		return err
	}

	file, err := os.Create("generated/index.html")
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(file, "index", *r)
	if err != nil {
		return err
	}

	// if r.Score.Lowest != 0 {
	// err = r.Uplaod()
	// if err != nil {
	// 	return err
	// }
	// }

	return nil
}

func (r *Result) processSchoolData() error {
	if data, ok := r.resultData["school"].(map[string]interface{}); ok {
		_, state := utils.GetLocation(data["address"].(string))
		r.School = School{
			SchoolCity:   "Makurdi",
			SchoolRegion: state,
			ResultDesc:   "TERMLY SUMMARY OF PROGRESS REPORT",
			Term:         r.Student.Term,
			VacationDate: "December 15, 2023",
		}

		if val, ok := data["school_name"].(string); ok {
			r.School.SchoolName = val
		} else {
			return r.Error("failed to get school_name", nil)
		}
	}
	return nil
}

func (r *Result) processStudentData() error {
	var ok bool
	if r.Student.ID, ok = r.resultData["id"].(float64); !ok {
		return fmt.Errorf("failed to get student id")
	}
	if r.Student.FullName, ok = r.resultData["full_name"].(string); !ok {
		return fmt.Errorf("failed to get full_name")
	}
	if r.Student.StudentPhoto, ok = r.resultData["student_photo"].(string); !ok {
		return r.Error("failed to get student_photo", nil)
	}
	if val, ok := r.resultData["class_name"].(string); ok {
		r.Student.ClassName = val
	} else {
		return r.Error("failed to get class_name", nil)
	}
	if val, ok := r.resultData["section_name"].(string); ok {
		r.Student.SectionName = val
	} else {
		return r.Error("failed to get section_name", nil)
	}
	if val, ok := r.resultData["academic"].(map[string]interface{}); ok {
		if r.Student.SessionYear, ok = val["title"].(string); !ok {
			return r.Error("failed to get academic title", nil)
		}
	}

	if val, ok := r.resultData["admission_no"].(float64); ok {
		r.Student.AdminNo = fmt.Sprintf("%04d", int(val))
	} else {
		return r.Error("failed to get admission_no", nil)
	}

	if val, ok := r.resultData["custom_field"].(map[string]interface{}); ok {
		if r.Student.Term, ok = val["Exam Type"].(string); !ok {
			return r.Error("failed to get Exam Type", nil)
		}
		if r.Student.Opened, ok = val["Days School Opened"].(string); !ok {
			return r.Error("failed to get Days School Opened", nil)
		}
		if r.Student.Present, ok = val["Days Present"].(string); !ok {
			return r.Error("failed to get Days Present", nil)
		}
		if r.Student.Absent, ok = val["Days Absent"].(string); !ok {
			return r.Error("failed to get Days Absent", nil)
		}
	}

	if val, ok := r.resultData["category"].(map[string]interface{}); ok {
		if r.Student.Arm, ok = val["category_name"].(string); !ok {
			return r.Error("failed to get category_name", nil)
		}
	}

	return nil
}

func (r *Result) processRecordData() error {
	objectives, err := utils.UnmarshalJason("template/objectives.json")
	if err != nil {
		return err
	}

	if data, ok := r.resultData["subjects"].([]interface{}); ok {
		for _, rec := range data {
			rec := rec.(map[string]interface{})
			marksData := rec["marks"].([]interface{})

			marks := make(map[string]interface{})
			score := 0
			for _, m := range marksData {
				m := m.(map[string]interface{})
				if title, ok := m["exam_title"].(string); ok {
					marks[title] = m["total_marks"].(float64)
				}
				if comment, ok := m["teacher_remarks"].(string); ok {
					r.Remark.Name = "Teacher Remark"
					r.Remark.Comment = comment
				}
				score += int(m["total_marks"].(float64))
			}
			grade, color := r.getGrade(score)
			record := Record{
				Score:   score,
				Outcome: grade,
				Grade:   grade,
				Color:   color,
			}

			record.Subject, ok = rec["subject_name"].(string)
			if !ok {
				return r.Error("failed to get subject_name", &record.Subject)
			}
			if r.Student.Arm == "GRADERS" {
				record.Mta, ok = marks["FIRST CA"].(float64)
				if !ok {
					return r.Error("failed to get FIRST CA", &record.Subject)
				}
				record.Ca, ok = marks["SECOND CA"].(float64)
				if !ok {
					return r.Error("failed to get SECOND CA", &record.Subject)
				}
				record.Oral, ok = marks["ORAL"].(float64)
				if !ok {
					return r.Error("failed to get ORAL", &record.Subject)
				}
				record.Exam, ok = marks["EXAM"].(float64)
				if !ok {
					return r.Error("failed to get EXAM", &record.Subject)
				}
			} else {
				record.Exam, ok = marks["SCORE"].(float64)
				if !ok {
					return r.Error("failed to get EXAM", &record.Subject)
				}
			}

			r.teacherId, ok = rec["teacher_id"].(float64)
			if !ok {
				return r.Error("failed to get teacher_id", &record.Subject)
			}

			subject_code, ok := rec["subject_code"].(string)
			if !ok {
				return r.Error("failed to get subject_code", &record.Subject)
			}

			for _, obj := range objectives {
				class := strings.ToUpper(obj["class_name"].(string))
				code := strings.ToUpper(obj["subject_code"].(string))
				if class == r.Student.ClassName && code == subject_code {
					record.Objectives = strings.Split(obj["text"].(string), "|")
					break
				}
			}
			r.Records = append(r.Records, record)
		}
	}
	return nil
}

func (r *Result) processScoreData() error {
	r.Score.Total = r.scores["total_score"].(float64)
	r.Score.Average = r.scores["average"].(float64)
	r.Score.Highest = r.scores["highest_average"].(float64)
	r.Score.Lowest = r.scores["lowest_average"].(float64)

	if r.Student.Arm == "EYFS" {
		r.Score.Grading = "EMERGING(0-80) EXPECTED(81-90) EXCEEDING(91-100)"
	} else {
		r.Score.Grading = "A(94-100) B(86-93) C(77-85) D(70-76) E(0-69)"
	}
	return nil
}

func (r *Result) processRatingData() error {
	if data, ok := r.resultData["custom_field"].(map[string]interface{}); ok {
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
	details = append(details, fmt.Sprintf("Student ID: %.0f", r.Student.ID))
	details = append(details, fmt.Sprintf("Error: %s", msg))

	if subject != nil {
		details = append(details, fmt.Sprintf("Subject: %s", *subject))
	}
	err := fmt.Errorf(strings.Join(details, "\n"))
	return err
}
