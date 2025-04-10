package handlers

import (
	"net/http"
	"strconv"

	"lms-vue-go/backend/models"
	"lms-vue-go/backend/repository"

	"github.com/gin-gonic/gin"
)

// GetStudentAnswers mengembalikan semua jawaban siswa yang sedang login
func GetStudentAnswers(c *gin.Context) {
	// Set CORS headers
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

	// Handle OPTIONS request
	if c.Request.Method == "OPTIONS" {
		c.Status(http.StatusOK)
		return
	}

	// Inisialisasi repository
	studentRepo := repository.NewStudentRepository()
	studentAnswerRepo := repository.NewStudentAnswerRepository()

	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
		return
	}

	// Cari student berdasarkan user ID
	student, err := studentRepo.FindByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil siswa tidak ditemukan"})
		return
	}

	// Ambil semua jawaban siswa
	answers, err := studentAnswerRepo.FindByStudent(student.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil jawaban siswa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": answers})
}

// GetStudentAnswerByQuestion mengembalikan jawaban siswa untuk soal tertentu
func GetStudentAnswerByQuestion(c *gin.Context) {
	// Set CORS headers
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

	// Handle OPTIONS request
	if c.Request.Method == "OPTIONS" {
		c.Status(http.StatusOK)
		return
	}

	// Inisialisasi repository
	studentRepo := repository.NewStudentRepository()
	studentAnswerRepo := repository.NewStudentAnswerRepository()
	questionRepo := repository.NewQuestionRepository()

	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
		return
	}

	// Ambil question ID dari parameter
	questionID, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID soal tidak valid"})
		return
	}

	// Cari student berdasarkan user ID
	student, err := studentRepo.FindByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil siswa tidak ditemukan"})
		return
	}

	// Cek apakah soal ada
	question, err := questionRepo.FindByID(uint(questionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data soal"})
		return
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan"})
		return
	}

	// Ambil jawaban siswa untuk soal tertentu
	answer, err := studentAnswerRepo.FindByStudentAndQuestion(student.ID, uint(questionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil jawaban siswa"})
		return
	}

	if answer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jawaban tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": answer})
}

// SubmitStudentAnswer menyimpan jawaban siswa untuk soal tertentu
func SubmitStudentAnswer(c *gin.Context) {
	// Set CORS headers
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

	// Handle OPTIONS request
	if c.Request.Method == "OPTIONS" {
		c.Status(http.StatusOK)
		return
	}

	// Inisialisasi repository
	studentRepo := repository.NewStudentRepository()
	studentAnswerRepo := repository.NewStudentAnswerRepository()
	questionRepo := repository.NewQuestionRepository()

	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
		return
	}

	// Struktur untuk binding request
	type SubmitAnswerRequest struct {
		QuestionID uint   `json:"question_id" binding:"required"`
		Answer     string `json:"answer" binding:"required"`
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Cari student berdasarkan user ID
	student, err := studentRepo.FindByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil siswa tidak ditemukan"})
		return
	}

	// Cek apakah soal ada
	question, err := questionRepo.FindByID(req.QuestionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data soal"})
		return
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan"})
		return
	}

	// Cek apakah sudah ada jawaban sebelumnya
	existingAnswer, err := studentAnswerRepo.FindByStudentAndQuestion(student.ID, req.QuestionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa jawaban sebelumnya"})
		return
	}

	// Hitung skor otomatis untuk soal pilihan ganda
	var score *int
	if question.Type == models.MultipleChoice {
		// Jika jawaban benar, berikan skor penuh
		if req.Answer == question.Answer {
			fullScore := question.Score
			score = &fullScore
		} else {
			// Jika jawaban salah, berikan skor 0
			zeroScore := 0
			score = &zeroScore
		}
	}

	// Jika sudah ada jawaban, update
	if existingAnswer != nil {
		existingAnswer.Answer = req.Answer
		existingAnswer.Score = score
		err = studentAnswerRepo.Update(existingAnswer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate jawaban"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": existingAnswer, "message": "Jawaban berhasil diupdate"})
		return
	}

	// Jika belum ada jawaban, buat baru
	newAnswer := models.StudentAnswer{
		StudentID:  student.ID,
		QuestionID: req.QuestionID,
		Answer:     req.Answer,
		Score:      score,
	}

	err = studentAnswerRepo.Create(&newAnswer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jawaban"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": newAnswer, "message": "Jawaban berhasil disimpan"})
}

// GradeStudentAnswer memberikan nilai untuk jawaban siswa (hanya untuk guru dan admin)
func GradeStudentAnswer(c *gin.Context) {
	// Inisialisasi repository
	studentAnswerRepo := repository.NewStudentAnswerRepository()

	// Ambil role dari context
	userRole, exists := c.Get("userRole")
	if !exists || (userRole != models.RoleAdmin && userRole != models.RoleTeacher) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin dan guru yang dapat memberikan nilai"})
		return
	}

	// Ambil answer ID dari parameter
	answerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID jawaban tidak valid"})
		return
	}

	// Struktur untuk binding request
	type GradeRequest struct {
		Score int `json:"score" binding:"required"`
	}

	var req GradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Cari jawaban berdasarkan ID
	answer, err := studentAnswerRepo.FindByID(uint(answerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data jawaban"})
		return
	}

	if answer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jawaban tidak ditemukan"})
		return
	}

	// Update skor
	answer.Score = &req.Score
	err = studentAnswerRepo.Update(answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate nilai"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": answer, "message": "Nilai berhasil diupdate"})
}

// GetAllStudentAnswers mengembalikan semua jawaban siswa (hanya untuk admin dan guru)
func GetAllStudentAnswers(c *gin.Context) {
	// Inisialisasi repository
	studentAnswerRepo := repository.NewStudentAnswerRepository()

	// Ambil role dari context
	userRole, exists := c.Get("userRole")
	if !exists || (userRole != models.RoleAdmin && userRole != models.RoleTeacher) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin dan guru yang dapat melihat semua jawaban"})
		return
	}

	// Ambil semua jawaban
	answers, err := studentAnswerRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data jawaban"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": answers})
}
