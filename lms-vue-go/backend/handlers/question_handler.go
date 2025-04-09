package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lms-vue-go/backend/models"
)

// Data dummy untuk soal-soal
var questions = []models.Question{
	{
		ID:       1,
		Type:     models.MultipleChoice,
		Question: "Ibukota Indonesia adalah...",
		Options:  []string{"Jakarta", "Bandung", "Surabaya", "Yogyakarta", "Medan"},
		Answer:   "A",
	},
	{
		ID:       2,
		Type:     models.Essay,
		Question: "Jelaskan mengapa belajar pemrograman penting di era digital?",
	},
	{
		ID:       3,
		Type:     models.MultipleChoice,
		Question: "Bahasa pemrograman yang berjalan di lingkungan browser adalah...",
		Options:  []string{"JavaScript", "Java", "Python", "Go", "C++"},
		Answer:   "A",
	},
	{
		ID:       4,
		Type:     models.MultipleChoice,
		Question: "Mana yang bukan termasuk framework JavaScript?",
		Options:  []string{"Django", "React", "Angular", "Vue", "Svelte"},
		Answer:   "A",
		ImageURL: "https://example.com/frameworks.jpg",
	},
}

// GetAllQuestions mengembalikan semua soal (tanpa jawaban untuk non-admin)
func GetAllQuestions(c *gin.Context) {
	// Cek apakah user adalah admin (implementasikan autentikasi nanti)
	isAdmin := false

	// Jika bukan admin, hapus jawaban dari response
	if !isAdmin {
		// Buat salinan tanpa jawaban
		questionsNoAnswer := make([]models.Question, len(questions))
		for i, q := range questions {
			questionsNoAnswer[i] = q
			questionsNoAnswer[i].Answer = ""
		}
		c.JSON(http.StatusOK, gin.H{"data": questionsNoAnswer})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": questions})
}

// GetQuestionByID mengembalikan soal berdasarkan ID
func GetQuestionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	for _, question := range questions {
		if int(question.ID) == id {
			// Jika bukan admin, hapus jawaban
			isAdmin := false
			if !isAdmin {
				question.Answer = ""
			}
			c.JSON(http.StatusOK, gin.H{"data": question})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan"})
} 