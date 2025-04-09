package models

// Student merepresentasikan data siswa
type Student struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Class string `json:"class"`
	Email string `json:"email"`
}

// QuestionType adalah tipe soal (pilihan ganda atau essay)
type QuestionType string

const (
	MultipleChoice QuestionType = "multiple_choice"
	Essay          QuestionType = "essay"
)

// Question merepresentasikan soal ujian/kuis
type Question struct {
	ID       uint         `json:"id"`
	Type     QuestionType `json:"type"`
	Question string       `json:"question"`
	Options  []string     `json:"options,omitempty"` // Hanya digunakan untuk pilihan ganda
	Answer   string       `json:"answer,omitempty"`  // Jawaban benar untuk admin
	ImageURL string       `json:"image_url,omitempty"` // URL gambar jika soal memiliki gambar
}

// StudentAnswer merepresentasikan jawaban siswa terhadap soal
type StudentAnswer struct {
	ID         uint   `json:"id"`
	StudentID  uint   `json:"student_id"`
	QuestionID uint   `json:"question_id"`
	Answer     string `json:"answer"`
	Score      *int   `json:"score,omitempty"` // Nilai, bisa null jika belum dinilai
} 