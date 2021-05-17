package storage

import "gorm.io/gorm"

type Degree struct {
	gorm.Model
	Name     string
	Subjects []*Subject `gorm:"many2many:degree_subjects;"`
}

type Subject struct {
	gorm.Model
	Name    string
	Degrees []*Degree `gorm:"many2many:degree_subjects;"`
	Topics  []Topic
}

type Topic struct {
	gorm.Model
	Name      string
	SubjectID int
	Subject   Subject
	MCQs      []MCQ
}

type MCQ struct {
	gorm.Model
	Question    string
	Answer      string
	Explanation string
	OptionOne   string
	OptionTwo   string
	OptionThree string
	OptionFOur  string
	TopicID     int
	Topic       Topic
}
