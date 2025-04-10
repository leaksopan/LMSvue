package routes

import (
	"net/http"
	"strconv"
	"time"

	"lms-vue-go/backend/handlers"
	"lms-vue-go/backend/middleware"
	"lms-vue-go/backend/models"
	"lms-vue-go/backend/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter mengatur semua endpoint API
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Tambahkan middleware security headers
	r.Use(middleware.SecurityHeaders())

	// Apply CORS middleware to all routes
	r.Use(middleware.CORSMiddleware())

	// Endpoint untuk health check
	r.GET("/api/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "online",
			"message": "LMS API berjalan dengan baik",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Test endpoint untuk debugging
	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Test endpoint berhasil diakses",
		})
	})

	// Test endpoint for student answers
	r.GET("/api/test-answers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": []map[string]interface{}{
				{
					"id":          1,
					"student_id":  1,
					"question_id": 1,
					"answer":      "A",
					"score":       10,
				},
				{
					"id":          2,
					"student_id":  1,
					"question_id": 2,
					"answer":      "B",
					"score":       5,
				},
			},
		})
	})

	// Public endpoint for student answers
	r.GET("/api/public-answers", func(c *gin.Context) {
		// Return dummy data for now
		c.JSON(http.StatusOK, gin.H{
			"data": []map[string]interface{}{
				{
					"id":            1,
					"student_id":    1,
					"question_id":   1,
					"answer":        "A",
					"score":         10,
					"question_text": "What is the capital of France?",
					"question_type": "multiple_choice",
				},
				{
					"id":            2,
					"student_id":    1,
					"question_id":   2,
					"answer":        "B",
					"score":         5,
					"question_text": "What is the capital of Germany?",
					"question_type": "multiple_choice",
				},
			},
		})
	})

	// Public endpoint for all student answers (admin only)
	r.GET("/api/public-all-answers", func(c *gin.Context) {
		// Initialize repositories
		studentAnswerRepo := repository.NewStudentAnswerRepository()

		// Get all student answers with details
		answers, err := studentAnswerRepo.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student answers"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": answers})
	})

	// Public endpoint for grading student answers
	r.PUT("/api/public-grade-answer/:id", func(c *gin.Context) {
		// Get answer ID from URL parameter
		answerIDStr := c.Param("id")
		answerID, err := strconv.ParseUint(answerIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
			return
		}

		// Parse request body
		var req struct {
			Score int `json:"score" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Initialize repository
		studentAnswerRepo := repository.NewStudentAnswerRepository()

		// Find the answer by ID
		answer, err := studentAnswerRepo.FindByID(uint(answerID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student answer"})
			return
		}

		if answer == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
			return
		}

		// Update the score
		answer.Score = &req.Score
		err = studentAnswerRepo.Update(answer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update answer score"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":    answer,
			"message": "Answer graded successfully",
		})
	})

	// Grup untuk API
	api := r.Group("/api")
	{
		// Routes untuk autentikasi (tidak perlu middleware auth)
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
			// Route untuk mendapatkan data user saat ini (perlu middleware auth)
			auth.GET("/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)
		}

		// Routes untuk siswa (perlu middleware auth)
		students := api.Group("/students", middleware.AuthMiddleware())
		{
			// Semua pengguna dapat melihat daftar siswa
			students.GET("/", handlers.GetAllStudents)
			students.GET("/:id", handlers.GetStudentByID)
			// Endpoint untuk mendapatkan profil siswa sendiri (hanya untuk siswa)
			students.GET("/profile/me", handlers.GetCurrentStudentProfile)
			// Hanya admin dan guru yang dapat mengelola data siswa
			studentAdmin := students.Group("/", middleware.RoleMiddleware(models.RoleAdmin, models.RoleTeacher))
			{
				studentAdmin.POST("/", handlers.CreateStudent)
				studentAdmin.PUT("/:id", handlers.UpdateStudent)
				studentAdmin.DELETE("/:id", handlers.DeleteStudent)
			}
		}

		// Add a public endpoint for student answers with CORS headers
		api.GET("/public/answers/my", func(c *gin.Context) {
			// Set CORS headers
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

			// Get token from Authorization header
			authorization := c.GetHeader("Authorization")
			if authorization == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
				return
			}

			// Extract token from Bearer prefix
			token := ""
			if len(authorization) > 7 && authorization[:7] == "Bearer " {
				token = authorization[7:]
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
				return
			}

			// Validate token and get user ID
			userID, err := middleware.ValidateToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			// Initialize repositories
			studentRepo := repository.NewStudentRepository()
			studentAnswerRepo := repository.NewStudentAnswerRepository()

			// Find student by user ID
			student, err := studentRepo.FindByUserID(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student data"})
				return
			}

			if student == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
				return
			}

			// Get all answers for the student
			answers, err := studentAnswerRepo.FindByStudent(student.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student answers"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"data": answers})
		})

		// Handle OPTIONS requests for public answers
		api.OPTIONS("/public/answers/my", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.Status(http.StatusOK)
		})

		// Add a public endpoint for submitting answers
		api.POST("/public/answers/submit", func(c *gin.Context) {
			// Set CORS headers
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

			// Get token from Authorization header
			authorization := c.GetHeader("Authorization")
			if authorization == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
				return
			}

			// Extract token from Bearer prefix
			token := ""
			if len(authorization) > 7 && authorization[:7] == "Bearer " {
				token = authorization[7:]
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
				return
			}

			// Validate token and get user ID
			userID, err := middleware.ValidateToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			// Bind request body
			type SubmitAnswerRequest struct {
				QuestionID uint   `json:"question_id" binding:"required"`
				Answer     string `json:"answer" binding:"required"`
			}

			var req SubmitAnswerRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
				return
			}

			// Initialize repositories
			studentRepo := repository.NewStudentRepository()
			studentAnswerRepo := repository.NewStudentAnswerRepository()
			questionRepo := repository.NewQuestionRepository()

			// Find student by user ID
			student, err := studentRepo.FindByUserID(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student data"})
				return
			}

			if student == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
				return
			}

			// Check if question exists
			question, err := questionRepo.FindByID(req.QuestionID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch question data"})
				return
			}

			if question == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
				return
			}

			// Check if answer already exists
			existingAnswer, err := studentAnswerRepo.FindByStudentAndQuestion(student.ID, req.QuestionID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check previous answer"})
				return
			}

			// Calculate score for multiple choice questions
			var score *int
			if question.Type == models.MultipleChoice {
				if req.Answer == question.Answer {
					fullScore := question.Score
					score = &fullScore
				} else {
					zeroScore := 0
					score = &zeroScore
				}
			}

			// If answer exists, update it
			if existingAnswer != nil {
				existingAnswer.Answer = req.Answer
				existingAnswer.Score = score
				err = studentAnswerRepo.Update(existingAnswer)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update answer"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"data": existingAnswer, "message": "Answer updated successfully"})
				return
			}

			// Create new answer
			newAnswer := models.StudentAnswer{
				StudentID:  student.ID,
				QuestionID: req.QuestionID,
				Answer:     req.Answer,
				Score:      score,
			}

			err = studentAnswerRepo.Create(&newAnswer)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save answer"})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"data": newAnswer, "message": "Answer submitted successfully"})
		})

		// Handle OPTIONS requests for submitting answers
		api.OPTIONS("/public/answers/submit", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.Status(http.StatusOK)
		})

		// Routes untuk jawaban siswa (perlu middleware auth)
		answers := api.Group("/answers", middleware.AuthMiddleware())
		{
			// Siswa dapat melihat jawaban mereka sendiri
			answers.GET("/my", handlers.GetStudentAnswers)
			answers.GET("/my/question/:questionId", handlers.GetStudentAnswerByQuestion)
			// Siswa dapat mengirimkan jawaban
			answers.POST("/submit", handlers.SubmitStudentAnswer)

			// Hanya admin dan guru yang dapat melihat semua jawaban dan memberikan nilai
			answerAdmin := answers.Group("/", middleware.RoleMiddleware(models.RoleAdmin, models.RoleTeacher))
			{
				answerAdmin.GET("/", handlers.GetAllStudentAnswers)
				answerAdmin.PUT("/:id/grade", handlers.GradeStudentAnswer)
			}
		}

		// Add a public endpoint for questions that doesn't require authentication
		api.GET("/public/questions", func(c *gin.Context) {
			// Set CORS headers
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

			// Initialize repository
			questionRepo := repository.NewQuestionRepository()

			// Get all questions from database
			questions, err := questionRepo.FindAll()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
				return
			}

			// Hide answers for all questions
			for i := range questions {
				questions[i].HideAnswer()
			}

			c.JSON(http.StatusOK, gin.H{"data": questions})
		})

		// Handle OPTIONS requests for public questions
		api.OPTIONS("/public/questions", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
			c.Status(http.StatusOK)
		})

		// Add a public endpoint for students that doesn't require authentication
		api.GET("/public/students", func(c *gin.Context) {
			// Set CORS headers
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

			// Initialize repository
			studentRepo := repository.NewStudentRepository()

			// Get all students from database
			students, err := studentRepo.FindAll()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"data": students})
		})

		// Handle OPTIONS requests for public students
		api.OPTIONS("/public/students", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
			c.Status(http.StatusOK)
		})

		// Routes untuk soal (perlu middleware auth)
		questionsGroup := api.Group("/questions")

		// Handle OPTIONS requests for CORS preflight
		questionsGroup.OPTIONS("/*path", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
			c.Header("Access-Control-Max-Age", "86400") // 24 hours
			c.Status(http.StatusOK)
		})

		questionsGroup.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, // Allow all origins for questions
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length", "Content-Type"},
			AllowCredentials: false, // Set to false when using * for AllowOrigins
			MaxAge:           12 * time.Hour,
		}))

		// Add auth middleware after CORS
		questions := questionsGroup.Group("/", middleware.AuthMiddleware())
		{
			// Semua pengguna dapat melihat soal
			questions.GET("/", handlers.GetAllQuestions)
			questions.GET("/:id", handlers.GetQuestionByID)
			// Hanya admin dan guru yang dapat mengelola soal
			questionAdmin := questions.Group("/", middleware.RoleMiddleware(models.RoleAdmin, models.RoleTeacher))
			{
				questionAdmin.POST("/", handlers.CreateQuestion)
				questionAdmin.PUT("/:id", handlers.UpdateQuestion)
				questionAdmin.DELETE("/:id", handlers.DeleteQuestion)
			}
		}
	}

	return r
}
