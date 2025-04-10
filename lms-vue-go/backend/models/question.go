package models

import "time"

// QuestionType adalah tipe untuk jenis soal
type QuestionType string

const (
	MultipleChoice QuestionType = "multiple_choice"
	Essay          QuestionType = "essay"
)

// Question merepresentasikan soal
type Question struct {
	ID        uint         `json:"id"`
	Type      QuestionType `json:"type"`
	Question  string       `json:"question"`
	Options   []string     `json:"options,omitempty"`
	Answer    string       `json:"answer,omitempty"`
	ImageURL  string       `json:"image_url,omitempty"`
	Score     int          `json:"score"`
	CreatedAt time.Time    `json:"created_at,omitempty"`
	UpdatedAt time.Time    `json:"updated_at,omitempty"`
}

// HideAnswer menghapus jawaban dari soal untuk keamanan
func (q *Question) HideAnswer() {
	q.Answer = ""
}
