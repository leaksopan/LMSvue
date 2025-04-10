package models

import "time"

// Student merepresentasikan data siswa
type Student struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id,omitempty"`
	Name      string    `json:"name"`
	Class     string    `json:"class"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// StudentAnswer merepresentasikan jawaban siswa untuk soal
type StudentAnswer struct {
	ID         uint      `json:"id"`
	StudentID  uint      `json:"student_id"`
	QuestionID uint      `json:"question_id"`
	Answer     string    `json:"answer"`
	Score      *int      `json:"score,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// StudentAnswerWithDetails merepresentasikan jawaban siswa dengan detail siswa dan soal
type StudentAnswerWithDetails struct {
	ID            uint         `json:"id"`
	StudentID     uint         `json:"student_id"`
	QuestionID    uint         `json:"question_id"`
	Answer        string       `json:"answer"`
	Score         *int         `json:"score,omitempty"`
	StudentName   string       `json:"student_name"`
	StudentClass  string       `json:"student_class"`
	UserID        uint         `json:"user_id"`
	QuestionText  string       `json:"question_text"`
	QuestionType  QuestionType `json:"question_type"`
	QuestionScore int          `json:"question_score"`
	CreatedAt     time.Time    `json:"created_at,omitempty"`
	UpdatedAt     time.Time    `json:"updated_at,omitempty"`
}
