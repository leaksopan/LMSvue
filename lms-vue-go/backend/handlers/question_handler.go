package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lms-vue-go/backend/models"
	"lms-vue-go/backend/repository"
)

// questionRepo adalah repository untuk operasi question
var questionRepo = repository.NewQuestionRepository()

// GetAllQuestions mengembalikan semua soal (tanpa jawaban untuk non-admin)
func GetAllQuestions(c *gin.Context) {
	// Cek apakah user adalah admin dari context yang diset oleh middleware
	userRole, exists := c.Get("userRole")
	isAdmin := exists && userRole == models.RoleAdmin

	// Ambil semua soal dari database
	questions, err := questionRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data soal"})
		return
	}

	// Jika bukan admin, hapus jawaban dari response
	if !isAdmin {
		// Buat salinan questions tanpa jawaban
		questionsWithoutAnswers := make([]models.Question, len(questions))
		for i, q := range questions {
			questionsWithoutAnswers[i] = q
			questionsWithoutAnswers[i].HideAnswer()
		}
		c.JSON(http.StatusOK, gin.H{"data": questionsWithoutAnswers})
		return
	}

	// Jika admin, kirim semua data termasuk jawaban
	c.JSON(http.StatusOK, gin.H{"data": questions})
}

// GetQuestionByID mengembalikan soal berdasarkan ID
func GetQuestionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Validasi ID
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus positif"})
		return
	}

	// Ambil soal dari database
	question, err := questionRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data soal"})
		return
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan"})
		return
	}

	// Cek apakah user adalah admin dari context yang diset oleh middleware
	userRole, exists := c.Get("userRole")
	isAdmin := exists && userRole == models.RoleAdmin
	if !isAdmin {
		// Jika bukan admin, hapus jawaban
		questionCopy := *question
		questionCopy.HideAnswer()
		c.JSON(http.StatusOK, gin.H{"data": questionCopy})
		return
	}

	// Jika admin, kirim semua data termasuk jawaban
	c.JSON(http.StatusOK, gin.H{"data": question})
}

// CreateQuestion menambahkan soal baru
func CreateQuestion(c *gin.Context) {
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Validasi tipe soal
	if question.Type != models.MultipleChoice && question.Type != models.Essay {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe soal tidak valid"})
		return
	}

	// Simpan soal ke database
	err := questionRepo.Create(&question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan soal"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": question})
}

// UpdateQuestion mengupdate soal
func UpdateQuestion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Cari soal berdasarkan ID
	question, err := questionRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data soal"})
		return
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan"})
		return
	}

	// Bind data baru
	var updatedQuestion models.Question
	if err := c.ShouldBindJSON(&updatedQuestion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Update data
	updatedQuestion.ID = uint(id)
	err = questionRepo.Update(&updatedQuestion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate soal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedQuestion})
}

// DeleteQuestion menghapus soal
func DeleteQuestion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Cari soal berdasarkan ID
	question, err := questionRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data soal"})
		return
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan"})
		return
	}

	// Hapus data
	err = questionRepo.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus soal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Soal berhasil dihapus"})
}
